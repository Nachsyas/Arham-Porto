package github

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// MaxFileSize limits individual source files fetched from GitHub (256 KiB) per Gate #14.
const MaxFileSize = 256 * 1024

var (
	ErrFileNotFound     = errors.New("file not found on github")
	ErrRateLimited      = errors.New("github api rate limit exceeded")
	ErrFileOversized    = errors.New("file exceeds maximum allowed size (256 KiB)")
	ErrBinaryContent    = errors.New("binary or non-utf8 content rejected")
	ErrPathExcluded     = errors.New("path matches secret or excluded pattern")
	ErrInvalidSHA       = errors.New("invalid commit SHA format: expected 40 hex characters")
	ErrInvalidPath      = errors.New("invalid path: traversal or absolute path not permitted")
	ErrUnauthorizedHost = errors.New("redirect to unauthorized host rejected")
	ErrMissingToken     = errors.New("github token is required")
	ErrUserNotFound     = errors.New("github user not found")
)

var shaRegex = regexp.MustCompile(`^[0-9a-fA-F]{40}$`)

const (
	defaultGitHubAPIOrigin = "https://api.github.com"
	defaultTimeout         = 15 * time.Second
)

// Allowed redirect hosts in production.
var allowedRedirectHosts = map[string]bool{
	"api.github.com":                true,
	"raw.githubusercontent.com":     true,
	"github.com":                    true,
	"codeload.github.com":           true,
	"objects.githubusercontent.com": true,
}

// Client provides hardened, secure access to GitHub contents API.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewClient constructs a production GitHub client pinned strictly to https://api.github.com.
// Users/environment cannot arbitrarily set GITHUB_BASE_URL in production (Gate #12).
func NewClient(token string) *Client {
	return newClientWithBaseURL(defaultGitHubAPIOrigin, token, false)
}

// NewTestClient constructs a client pointing to a custom base URL (e.g. httptest.Server).
// Strictly intended for automated test suites (Gate #12).
func NewTestClient(testBaseURL, token string) *Client {
	return newClientWithBaseURL(testBaseURL, token, true)
}

func newClientWithBaseURL(baseURL, token string, isTest bool) *Client {
	var testHost string
	if isTest {
		if u, err := url.Parse(baseURL); err == nil {
			testHost = u.Host
		}
	}

	httpClient := &http.Client{
		Timeout: defaultTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}

			// Gate #13: Redirect safety - verify target host
			reqHost := req.URL.Host
			if isTest && testHost != "" && reqHost == testHost {
				return nil
			}
			if !allowedRedirectHosts[reqHost] {
				return fmt.Errorf("%w: %s", ErrUnauthorizedHost, reqHost)
			}
			return nil
		},
	}

	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		token:      strings.TrimSpace(token),
		httpClient: httpClient,
	}
}

// ResolveCommitSHA resolves a branch/tag/ref to an immutable 40-character Git SHA-1.
// Validates format strictly (Gate #39).
func (c *Client) ResolveCommitSHA(ctx context.Context, owner, repo, ref string) (string, error) {
	if strings.TrimSpace(ref) == "" {
		ref = "HEAD"
	}

	endpoint := fmt.Sprintf("%s/repos/%s/%s/commits/%s", c.baseURL, url.PathEscape(owner), url.PathEscape(repo), url.PathEscape(ref))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create commit request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.sha")
	req.Header.Set("User-Agent", "Arham-Porto-Indexer/1.0")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("github commit resolution failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("%w: repo %s/%s ref %s", ErrFileNotFound, owner, repo, ref)
	}
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		return "", ErrRateLimited
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d while resolving commit: %s", resp.StatusCode, resp.Status)
	}

	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
	if err != nil {
		return "", fmt.Errorf("failed to read commit response: %w", err)
	}

	sha := strings.TrimSpace(string(bodyBytes))
	if !shaRegex.MatchString(sha) {
		return "", fmt.Errorf("%w: got %q", ErrInvalidSHA, sha)
	}

	return sha, nil
}

// FetchRawFile downloads and validates source file bytes at an immutable commit SHA.
// Enforces max size 256 KiB, UTF-8 validity, and secret exclusions (Gate #14, #33).
func (c *Client) FetchRawFile(ctx context.Context, owner, repo, commitSHA, filePath string) ([]byte, error) {
	// Validate commit SHA (Gate #39)
	if !shaRegex.MatchString(commitSHA) {
		return nil, fmt.Errorf("%w: %s", ErrInvalidSHA, commitSHA)
	}

	// Validate path against traversal
	cleanPath := path.Clean(strings.TrimSpace(filePath))
	if strings.HasPrefix(cleanPath, "../") || strings.HasPrefix(cleanPath, "/") || cleanPath == ".." || cleanPath == "." {
		return nil, fmt.Errorf("%w: %s", ErrInvalidPath, filePath)
	}

	// Defense-in-depth secret & pattern check (Gate #33)
	if excluded, reason := isPathExcluded(cleanPath); excluded {
		return nil, fmt.Errorf("%w (%s): %s", ErrPathExcluded, reason, cleanPath)
	}

	endpoint := fmt.Sprintf("%s/repos/%s/%s/contents/%s?ref=%s",
		c.baseURL,
		url.PathEscape(owner),
		url.PathEscape(repo),
		cleanPath,
		url.QueryEscape(commitSHA),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create contents request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3.raw")
	req.Header.Set("User-Agent", "Arham-Porto-Indexer/1.0")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github fetch request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %s at %s", ErrFileNotFound, cleanPath, commitSHA)
	}
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimited
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d fetching file %s: %s", resp.StatusCode, cleanPath, resp.Status)
	}

	// Bound response body bytes (Gate #14): MaxFileSize + 1 to detect oversize
	limitedReader := io.LimitReader(resp.Body, MaxFileSize+1)
	content, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %w", err)
	}

	if len(content) > MaxFileSize {
		return nil, fmt.Errorf("%w: size %d > %d", ErrFileOversized, len(content), MaxFileSize)
	}

	// Reject binary content
	if isBinaryContent(content) {
		return nil, fmt.Errorf("%w: %s", ErrBinaryContent, cleanPath)
	}

	return content, nil
}

func isBinaryContent(data []byte) bool {
	if bytes.IndexByte(data, 0) != -1 {
		return true
	}
	return !utf8.Valid(data)
}

func isPathExcluded(p string) (bool, string) {
	lower := strings.ToLower(strings.TrimSpace(p))
	secretPatterns := []string{".env", "id_rsa", "id_ecdsa", "id_ed25519"}
	for _, pat := range secretPatterns {
		if lower == pat || strings.HasPrefix(lower, pat+".") || strings.HasSuffix(lower, "/"+pat) || strings.Contains(lower, "/"+pat+".") {
			return true, "secret-bearing pattern match"
		}
	}
	if strings.HasSuffix(lower, ".pem") || strings.HasSuffix(lower, ".key") || strings.HasSuffix(lower, ".crt") {
		return true, "certificate or private key file"
	}
	if strings.HasPrefix(filepath.Base(lower), "credentials.") || strings.HasPrefix(filepath.Base(lower), "service-account") {
		return true, "credential file pattern"
	}
	return false, ""
}

// GenerateCitationURL creates an immutable GitHub URL for a file at a specific commit SHA (Gate #38).
func GenerateCitationURL(owner, repo, commitSHA, filePath string) (string, error) {
	if !shaRegex.MatchString(commitSHA) {
		return "", fmt.Errorf("%w: %s", ErrInvalidSHA, commitSHA)
	}

	cleanPath := path.Clean(strings.TrimSpace(filePath))
	if strings.HasPrefix(cleanPath, "../") || strings.HasPrefix(cleanPath, "/") || cleanPath == ".." || cleanPath == "." {
		return "", fmt.Errorf("%w: %s", ErrInvalidPath, filePath)
	}

	return fmt.Sprintf("https://github.com/%s/%s/blob/%s/%s", owner, repo, commitSHA, cleanPath), nil
}

// FetchUserRepositories retrieves all public repositories for a specified user with pagination.
// Requires valid GITHUB_TOKEN authentication.
func (c *Client) FetchUserRepositories(ctx context.Context, username string) ([]Repository, error) {
	if strings.TrimSpace(c.token) == "" {
		return nil, ErrMissingToken
	}

	cleanUsername := strings.TrimSpace(username)
	if cleanUsername == "" {
		return nil, errors.New("github username cannot be empty")
	}

	var allRepos []Repository
	page := 1
	perPage := 100
	const maxPages = 50

	for page <= maxPages {
		endpoint := fmt.Sprintf("%s/users/%s/repos?per_page=%d&page=%d&type=owner",
			c.baseURL,
			url.PathEscape(cleanUsername),
			perPage,
			page,
		)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create repos request: %w", err)
		}

		req.Header.Set("Accept", "application/vnd.github.v3+json")
		req.Header.Set("User-Agent", "Arham-Porto-Sync/1.0")
		req.Header.Set("Authorization", "Bearer "+c.token)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return nil, fmt.Errorf("github fetch repos failed: %w", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			resp.Body.Close()
			return nil, fmt.Errorf("%w: user %s", ErrUserNotFound, cleanUsername)
		}
		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			return nil, ErrRateLimited
		}
		if resp.StatusCode == http.StatusUnauthorized {
			resp.Body.Close()
			return nil, errors.New("github api authentication failed: invalid or expired token")
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("unexpected status %d fetching repos: %s", resp.StatusCode, resp.Status)
		}

		var pageItems []githubRepoResponse
		decodeErr := json.NewDecoder(resp.Body).Decode(&pageItems)
		linkHeader := resp.Header.Get("Link")
		resp.Body.Close()

		if decodeErr != nil {
			return nil, fmt.Errorf("failed to decode github repos response: %w", decodeErr)
		}

		if len(pageItems) == 0 {
			break
		}

		for _, item := range pageItems {
			allRepos = append(allRepos, item.ToRepository())
		}

		// Check pagination termination:
		// When GitHub Link header is present, continue only if rel="next" is present.
		// When Link header is omitted (e.g. some mock environments), terminate if page is smaller than perPage.
		if linkHeader != "" {
			if !strings.Contains(linkHeader, `rel="next"`) {
				break
			}
		} else {
			if len(pageItems) < perPage {
				break
			}
		}

		page++
	}

	return allRepos, nil
}


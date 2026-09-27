package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGitHubClient(t *testing.T) {
	validSHA := "0123456789abcdef0123456789abcdef01234567"

	// Setup mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/commits/main"):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(validSHA))

		case strings.Contains(r.URL.Path, "/commits/malformed"):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("not-a-40-char-sha"))

		case strings.Contains(r.URL.Path, "/commits/notfound"):
			w.WriteHeader(http.StatusNotFound)

		case strings.Contains(r.URL.Path, "/commits/ratelimited"):
			w.WriteHeader(http.StatusForbidden)

		case strings.Contains(r.URL.Path, "README.md"):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("# Sample Readme Content"))

		case strings.Contains(r.URL.Path, "oversized.txt"):
			w.WriteHeader(http.StatusOK)
			// Write 256 KiB + 10 bytes
			oversized := make([]byte, MaxFileSize+10)
			w.Write(oversized)

		case strings.Contains(r.URL.Path, "binary.dat"):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte{0x00, 0x01, 0x02, 0xFF})

		case strings.Contains(r.URL.Path, "redirect-evil"):
			http.Redirect(w, r, "https://evil-unauthorized-domain.com/evil.txt", http.StatusFound)

		case strings.Contains(r.URL.Path, "timeout"):
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	client := NewTestClient(ts.URL, "test-token")
	ctx := context.Background()

	t.Run("resolve commit SHA success", func(t *testing.T) {
		sha, err := client.ResolveCommitSHA(ctx, "Nachsyas", "EduTrace", "main")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sha != validSHA {
			t.Errorf("expected %s, got %s", validSHA, sha)
		}
	})

	t.Run("resolve commit SHA malformed error", func(t *testing.T) {
		_, err := client.ResolveCommitSHA(ctx, "Nachsyas", "EduTrace", "malformed")
		if err == nil || !errors.Is(err, ErrInvalidSHA) {
			t.Errorf("expected ErrInvalidSHA, got %v", err)
		}
	})

	t.Run("resolve commit SHA 404", func(t *testing.T) {
		_, err := client.ResolveCommitSHA(ctx, "Nachsyas", "EduTrace", "notfound")
		if err == nil || !errors.Is(err, ErrFileNotFound) {
			t.Errorf("expected ErrFileNotFound, got %v", err)
		}
	})

	t.Run("resolve commit SHA rate limited", func(t *testing.T) {
		_, err := client.ResolveCommitSHA(ctx, "Nachsyas", "EduTrace", "ratelimited")
		if err == nil || !errors.Is(err, ErrRateLimited) {
			t.Errorf("expected ErrRateLimited, got %v", err)
		}
	})

	t.Run("fetch raw file success", func(t *testing.T) {
		content, err := client.FetchRawFile(ctx, "Nachsyas", "EduTrace", validSHA, "README.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(string(content), "Sample Readme") {
			t.Errorf("unexpected content: %s", string(content))
		}
	})

	t.Run("fetch raw file invalid SHA rejected", func(t *testing.T) {
		_, err := client.FetchRawFile(ctx, "Nachsyas", "EduTrace", "bad-sha", "README.md")
		if err == nil || !errors.Is(err, ErrInvalidSHA) {
			t.Errorf("expected ErrInvalidSHA, got %v", err)
		}
	})

	t.Run("fetch raw file path traversal rejected", func(t *testing.T) {
		_, err := client.FetchRawFile(ctx, "Nachsyas", "EduTrace", validSHA, "../etc/passwd")
		if err == nil || !errors.Is(err, ErrInvalidPath) {
			t.Errorf("expected ErrInvalidPath, got %v", err)
		}
	})

	t.Run("fetch raw file secret excluded rejected", func(t *testing.T) {
		_, err := client.FetchRawFile(ctx, "Nachsyas", "EduTrace", validSHA, ".env")
		if err == nil || !errors.Is(err, ErrPathExcluded) {
			t.Errorf("expected ErrPathExcluded, got %v", err)
		}
	})

	t.Run("fetch raw file oversized rejected", func(t *testing.T) {
		_, err := client.FetchRawFile(ctx, "Nachsyas", "EduTrace", validSHA, "oversized.txt")
		if err == nil || !errors.Is(err, ErrFileOversized) {
			t.Errorf("expected ErrFileOversized, got %v", err)
		}
	})

	t.Run("fetch raw file binary rejected", func(t *testing.T) {
		_, err := client.FetchRawFile(ctx, "Nachsyas", "EduTrace", validSHA, "binary.dat")
		if err == nil || !errors.Is(err, ErrBinaryContent) {
			t.Errorf("expected ErrBinaryContent, got %v", err)
		}
	})

	t.Run("redirect to unauthorized host rejected", func(t *testing.T) {
		_, err := client.FetchRawFile(ctx, "Nachsyas", "EduTrace", validSHA, "redirect-evil")
		if err == nil || !errors.Is(err, ErrUnauthorizedHost) {
			t.Errorf("expected ErrUnauthorizedHost, got %v", err)
		}
	})

	t.Run("generate citation URL", func(t *testing.T) {
		url, err := GenerateCitationURL("Nachsyas", "EduTrace", validSHA, "contracts/AchievementBadge.sol")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := fmt.Sprintf("https://github.com/Nachsyas/EduTrace/blob/%s/contracts/AchievementBadge.sol", validSHA)
		if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})
}

func TestFetchUserRepositories(t *testing.T) {
	t.Run("successful repository fetch and field mapping", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if auth := r.Header.Get("Authorization"); auth != "Bearer test-secret-token" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if r.URL.Path != "/users/Nachsyas/repos" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{
					"id": 1001,
					"name": "EduTrace",
					"full_name": "Nachsyas/EduTrace",
					"description": "Verifiable credential management",
					"html_url": "https://github.com/Nachsyas/EduTrace",
					"homepage": "https://edutrace.example.com",
					"language": "Go",
					"stargazers_count": 42,
					"forks_count": 8,
					"topics": ["go", "blockchain", "education"]
				},
				{
					"id": 1002,
					"name": "Maritime-AI",
					"full_name": "Nachsyas/Maritime-AI",
					"description": null,
					"html_url": "https://github.com/Nachsyas/Maritime-AI",
					"homepage": null,
					"language": null,
					"stargazers_count": 12,
					"forks_count": 2,
					"topics": null
				}
			]`))
		}))
		defer ts.Close()

		client := NewTestClient(ts.URL, "test-secret-token")
		repos, err := client.FetchUserRepositories(context.Background(), "Nachsyas")
		if err != nil {
			t.Fatalf("unexpected error fetching repos: %v", err)
		}

		if len(repos) != 2 {
			t.Fatalf("expected 2 repositories, got %d", len(repos))
		}

		// Verify repo 1 fields
		r1 := repos[0]
		if r1.GithubID != 1001 || r1.Name != "EduTrace" || r1.FullName != "Nachsyas/EduTrace" {
			t.Errorf("repo 1 identity mismatch: %+v", r1)
		}
		if r1.Description != "Verifiable credential management" || r1.HTMLURL != "https://github.com/Nachsyas/EduTrace" {
			t.Errorf("repo 1 links/desc mismatch: %+v", r1)
		}
		if r1.Homepage != "https://edutrace.example.com" || r1.Language != "Go" || r1.Stars != 42 || r1.Forks != 8 {
			t.Errorf("repo 1 metrics mismatch: %+v", r1)
		}
		if len(r1.Topics) != 3 || r1.Topics[0] != "go" {
			t.Errorf("repo 1 topics mismatch: %+v", r1.Topics)
		}

		// Verify repo 2 null handling
		r2 := repos[1]
		if r2.GithubID != 1002 || r2.Description != "" || r2.Homepage != "" || r2.Language != "" {
			t.Errorf("repo 2 null defaults mismatch: %+v", r2)
		}
		if len(r2.Topics) != 0 {
			t.Errorf("expected empty topics slice, got %+v", r2.Topics)
		}
	})

	t.Run("pagination across multiple pages", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			page := r.URL.Query().Get("page")
			w.Header().Set("Content-Type", "application/json")
			if page == "1" {
				w.Header().Set("Link", fmt.Sprintf(`<%s/users/Nachsyas/repos?page=2&per_page=100>; rel="next"`, "http://"+r.Host))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[
					{"id": 1, "name": "repo-1", "full_name": "Nachsyas/repo-1", "html_url": "https://github.com/Nachsyas/repo-1"},
					{"id": 2, "name": "repo-2", "full_name": "Nachsyas/repo-2", "html_url": "https://github.com/Nachsyas/repo-2"}
				]`))
				return
			}
			if page == "2" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[
					{"id": 3, "name": "repo-3", "full_name": "Nachsyas/repo-3", "html_url": "https://github.com/Nachsyas/repo-3"}
				]`))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
		}))
		defer ts.Close()

		client := NewTestClient(ts.URL, "test-token")
		repos, err := client.FetchUserRepositories(context.Background(), "Nachsyas")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(repos) != 3 {
			t.Fatalf("expected 3 repositories across 2 pages, got %d", len(repos))
		}
		if repos[0].Name != "repo-1" || repos[1].Name != "repo-2" || repos[2].Name != "repo-3" {
			t.Errorf("unexpected repositories order or names: %+v", repos)
		}
	})

	t.Run("API error - missing token", func(t *testing.T) {
		client := NewTestClient("https://api.github.com", "")
		_, err := client.FetchUserRepositories(context.Background(), "Nachsyas")
		if err == nil || !errors.Is(err, ErrMissingToken) {
			t.Errorf("expected ErrMissingToken, got %v", err)
		}
	})

	t.Run("API error - 404 user not found", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer ts.Close()

		client := NewTestClient(ts.URL, "test-token")
		_, err := client.FetchUserRepositories(context.Background(), "non-existent-user")
		if err == nil || !errors.Is(err, ErrUserNotFound) {
			t.Errorf("expected ErrUserNotFound, got %v", err)
		}
	})

	t.Run("API error - rate limited", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}))
		defer ts.Close()

		client := NewTestClient(ts.URL, "test-token")
		_, err := client.FetchUserRepositories(context.Background(), "Nachsyas")
		if err == nil || !errors.Is(err, ErrRateLimited) {
			t.Errorf("expected ErrRateLimited, got %v", err)
		}
	})

	t.Run("API error - server 500 error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		client := NewTestClient(ts.URL, "test-token")
		_, err := client.FetchUserRepositories(context.Background(), "Nachsyas")
		if err == nil || !strings.Contains(err.Error(), "unexpected status 500") {
			t.Errorf("expected status 500 error, got %v", err)
		}
	})

	t.Run("context cancellation terminates pagination", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(50 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
		}))
		defer ts.Close()

		client := NewTestClient(ts.URL, "test-token")
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := client.FetchUserRepositories(ctx, "Nachsyas")
		if err == nil || !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})
}


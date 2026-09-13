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

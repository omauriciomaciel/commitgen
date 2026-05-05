package updatecheck

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLatestVersionReadsGoProxyVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != "commitgen" {
			t.Fatalf("User-Agent = %q, want commitgen", r.Header.Get("User-Agent"))
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"Version":"v1.2.3"}`)); err != nil {
			t.Fatalf("Write() error = %v", err)
		}
	}))
	defer server.Close()

	version, err := LatestVersion(context.Background(), server.Client(), server.URL)
	if err != nil {
		t.Fatalf("LatestVersion() error = %v", err)
	}
	if version != "v1.2.3" {
		t.Fatalf("LatestVersion() = %q, want v1.2.3", version)
	}
}

func TestLatestVersionFallsBackToReleaseTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"tag_name":"v1.2.3"}`)); err != nil {
			t.Fatalf("Write() error = %v", err)
		}
	}))
	defer server.Close()

	version, err := LatestVersion(context.Background(), server.Client(), server.URL)
	if err != nil {
		t.Fatalf("LatestVersion() error = %v", err)
	}
	if version != "v1.2.3" {
		t.Fatalf("LatestVersion() = %q, want v1.2.3", version)
	}
}

func TestLatestVersionErrorsWhenVersionIsMissing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{}`)); err != nil {
			t.Fatalf("Write() error = %v", err)
		}
	}))
	defer server.Close()

	if _, err := LatestVersion(context.Background(), server.Client(), server.URL); err == nil {
		t.Fatal("LatestVersion() error = nil, want error")
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name  string
		left  string
		right string
		want  int
	}{
		{name: "newer patch", left: "v1.1.1", right: "v1.1.0", want: 1},
		{name: "newer minor", left: "1.2.0", right: "1.1.9", want: 1},
		{name: "newer major", left: "2.0.0", right: "1.9.9", want: 1},
		{name: "same with v prefix", left: "v1.1.0", right: "1.1.0", want: 0},
		{name: "older", left: "1.0.9", right: "1.1.0", want: -1},
		{name: "ignores prerelease suffix", left: "v1.1.0-beta.1", right: "1.1.0", want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareVersions(tt.left, tt.right); got != tt.want {
				t.Fatalf("CompareVersions(%q, %q) = %d, want %d", tt.left, tt.right, got, tt.want)
			}
		})
	}
}

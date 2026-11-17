package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// TestUploadRetriesWhileServerBusy checks that `preview upload` waits out a
// busy server: the first attempt is rejected before its body is read, and
// the retry sends the whole archive again with a freshly minted token.
func TestUploadRetriesWhileServerBusy(t *testing.T) {
	archive := bytes.Repeat([]byte("x"), 8<<20)
	path := filepath.Join(t.TempDir(), "fe.tar.gz")
	if err := os.WriteFile(path, archive, 0o644); err != nil {
		t.Fatal(err)
	}

	var mu sync.Mutex
	var tokens []string
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		attempts++
		attempt := attempts
		tokens = append(tokens, r.Header.Get("Authorization"))
		mu.Unlock()
		if attempt == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, `{"error":"busy"}`)
			return
		}
		got, err := io.ReadAll(r.Body)
		if err != nil || !bytes.Equal(got, archive) {
			t.Errorf("retry body: %d bytes (err %v), want the full %d-byte archive", len(got), err, len(archive))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"sha": "abc", "short_sha": "abc", "side": "frontend", "hash": "h", "published": true})
	}))
	defer srv.Close()

	minted := 0
	token := func(context.Context) (string, error) {
		minted++
		return fmt.Sprintf("tok-%d", minted), nil
	}
	var out strings.Builder
	if err := runUpload(context.Background(), srv.URL, &out, "demo", "frontend", "", path, "abc", false, false, token); err != nil {
		t.Fatalf("runUpload: %v\n%s", err, out.String())
	}
	if attempts != 2 || !strings.Contains(out.String(), "server busy") {
		t.Fatalf("attempts = %d, want 2; output:\n%s", attempts, out.String())
	}
	if tokens[0] == tokens[1] {
		t.Fatalf("retry reused token %q; want a fresh one", tokens[1])
	}
}

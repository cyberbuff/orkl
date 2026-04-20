package orkl

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newJSONServer returns a test server that asserts the request matches wantPath
// and wantQuery (url-encoded), then responds with body.
func newJSONServer(t *testing.T, wantPath, wantQuery, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != wantPath {
			t.Errorf("path: got %q, want %q", r.URL.Path, wantPath)
		}
		if got := r.URL.Query().Encode(); got != wantQuery {
			t.Errorf("query: got %q, want %q", got, wantQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, body)
	}))
}

func runCLI(t *testing.T, srv *httptest.Server, args []string) (stdout string, rc int) {
	t.Helper()
	var out bytes.Buffer
	baseURL := ""
	if srv != nil {
		baseURL = srv.URL
	}
	rc = runWithBaseURL(context.Background(), args, &out, io.Discard, baseURL)
	return strings.TrimSpace(out.String()), rc
}

// ---------------------------------------------------------------------------
// library-info
// ---------------------------------------------------------------------------

func TestLibraryInfoCommandPrintsJSON(t *testing.T) {
	srv := newJSONServer(t, "/library/info", "", `{"ok":true}`)
	defer srv.Close()

	got, rc := runCLI(t, srv, []string{"library-info"})
	if rc != 0 {
		t.Fatalf("exit code: got %d, want 0", rc)
	}
	if got != `{"ok":true}` {
		t.Fatalf("output: got %s", got)
	}
}

// ---------------------------------------------------------------------------
// library-version
// ---------------------------------------------------------------------------

func TestLibraryVersionCommand(t *testing.T) {
	srv := newJSONServer(t, "/library/version", "", `{"version":"1"}`)
	defer srv.Close()

	got, rc := runCLI(t, srv, []string{"library-version"})
	if rc != 0 {
		t.Fatalf("exit code: got %d, want 0", rc)
	}
	if got != `{"version":"1"}` {
		t.Fatalf("output: got %s", got)
	}
}

// ---------------------------------------------------------------------------
// library-versions
// ---------------------------------------------------------------------------

func TestLibraryVersionsCommandNoFlags(t *testing.T) {
	srv := newJSONServer(t, "/library/version/entries", "", `{"entries":[]}`)
	defer srv.Close()

	got, rc := runCLI(t, srv, []string{"library-versions"})
	if rc != 0 {
		t.Fatalf("exit code: got %d, want 0", rc)
	}
	if got != `{"entries":[]}` {
		t.Fatalf("output: got %s", got)
	}
}

func TestLibraryVersionsCommandWithFlags(t *testing.T) {
	srv := newJSONServer(t, "/library/version/entries", "limit=5&offset=10&order=asc", `{"entries":[]}`)
	defer srv.Close()

	_, rc := runCLI(t, srv, []string{"library-versions", "--limit", "5", "--offset", "10", "--order", "asc"})
	if rc != 0 {
		t.Fatalf("exit code: got %d, want 0", rc)
	}
}

// ---------------------------------------------------------------------------
// library-entries
// ---------------------------------------------------------------------------

func TestLibraryEntriesCommandNoFlags(t *testing.T) {
	srv := newJSONServer(t, "/library/entries", "", `{"entries":[]}`)
	defer srv.Close()

	got, rc := runCLI(t, srv, []string{"library-entries"})
	if rc != 0 {
		t.Fatalf("exit code: got %d, want 0", rc)
	}
	if got != `{"entries":[]}` {
		t.Fatalf("output: got %s", got)
	}
}

func TestLibraryEntriesCommandWithFlags(t *testing.T) {
	srv := newJSONServer(t, "/library/entries", "limit=10&offset=5&order=desc&order_by=date&origin=web", `{"entries":[]}`)
	defer srv.Close()

	_, rc := runCLI(t, srv, []string{
		"library-entries",
		"--limit", "10", "--offset", "5", "--order-by", "date", "--order", "desc", "--origin", "web",
	})
	if rc != 0 {
		t.Fatalf("exit code: got %d, want 0", rc)
	}
}

// ---------------------------------------------------------------------------
// library-entry
// ---------------------------------------------------------------------------

func TestLibraryEntryCommand(t *testing.T) {
	srv := newJSONServer(t, "/library/entry/abc-123", "", `{"id":"abc-123"}`)
	defer srv.Close()

	got, rc := runCLI(t, srv, []string{"library-entry", "abc-123"})
	if rc != 0 {
		t.Fatalf("exit code: got %d, want 0", rc)
	}
	if got != `{"id":"abc-123"}` {
		t.Fatalf("output: got %s", got)
	}
}

func TestLibraryEntryCommandMissingArg(t *testing.T) {
	_, rc := runCLI(t, nil, []string{"library-entry"})
	if rc != 2 {
		t.Fatalf("exit code: got %d, want 2", rc)
	}
}

// ---------------------------------------------------------------------------
// library-entry-hash
// ---------------------------------------------------------------------------

func TestLibraryEntryHashCommand(t *testing.T) {
	srv := newJSONServer(t, "/library/entry/sha1/deadbeef", "", `{"hash":"deadbeef"}`)
	defer srv.Close()

	got, rc := runCLI(t, srv, []string{"library-entry-hash", "deadbeef"})
	if rc != 0 {
		t.Fatalf("exit code: got %d, want 0", rc)
	}
	if got != `{"hash":"deadbeef"}` {
		t.Fatalf("output: got %s", got)
	}
}

func TestLibraryEntryHashCommandMissingArg(t *testing.T) {
	_, rc := runCLI(t, nil, []string{"library-entry-hash"})
	if rc != 2 {
		t.Fatalf("exit code: got %d, want 2", rc)
	}
}

// ---------------------------------------------------------------------------
// library-search
// ---------------------------------------------------------------------------

func TestLibrarySearchCommandSendsQueryParameters(t *testing.T) {
	srv := newJSONServer(t, "/library/search", `full=true&limit=25&origin=web&query=%22threat+actor%22`, `{"results":[]}`)
	defer srv.Close()

	got, rc := runCLI(t, srv, []string{"library-search", `"threat actor"`, "--full", "--limit", "25", "--origin", "web"})
	if rc != 0 {
		t.Fatalf("exit code: got %d, want 0", rc)
	}
	if got != `{"results":[]}` {
		t.Fatalf("output: got %s", got)
	}
}

func TestLibrarySearchCommandMissingArg(t *testing.T) {
	_, rc := runCLI(t, nil, []string{"library-search"})
	if rc != 2 {
		t.Fatalf("exit code: got %d, want 2", rc)
	}
}

// ---------------------------------------------------------------------------
// source-entries
// ---------------------------------------------------------------------------

func TestSourceEntriesCommand(t *testing.T) {
	srv := newJSONServer(t, "/source/entries", "", `{"entries":[]}`)
	defer srv.Close()

	got, rc := runCLI(t, srv, []string{"source-entries"})
	if rc != 0 {
		t.Fatalf("exit code: got %d, want 0", rc)
	}
	if got != `{"entries":[]}` {
		t.Fatalf("output: got %s", got)
	}
}

// ---------------------------------------------------------------------------
// source-entry
// ---------------------------------------------------------------------------

func TestSourceEntryCommand(t *testing.T) {
	srv := newJSONServer(t, "/source/entry/src-456", "", `{"id":"src-456"}`)
	defer srv.Close()

	got, rc := runCLI(t, srv, []string{"source-entry", "src-456"})
	if rc != 0 {
		t.Fatalf("exit code: got %d, want 0", rc)
	}
	if got != `{"id":"src-456"}` {
		t.Fatalf("output: got %s", got)
	}
}

func TestSourceEntryCommandFullFlag(t *testing.T) {
	srv := newJSONServer(t, "/source/entry/src-456", "full=true", `{"id":"src-456","reports":[]}`)
	defer srv.Close()

	_, rc := runCLI(t, srv, []string{"source-entry", "src-456", "--full"})
	if rc != 0 {
		t.Fatalf("exit code: got %d, want 0", rc)
	}
}

func TestSourceEntryCommandMissingArg(t *testing.T) {
	_, rc := runCLI(t, nil, []string{"source-entry"})
	if rc != 2 {
		t.Fatalf("exit code: got %d, want 2", rc)
	}
}

// ---------------------------------------------------------------------------
// ta-entries
// ---------------------------------------------------------------------------

func TestTAEntriesCommand(t *testing.T) {
	srv := newJSONServer(t, "/ta/entries", "", `{"entries":[]}`)
	defer srv.Close()

	got, rc := runCLI(t, srv, []string{"ta-entries"})
	if rc != 0 {
		t.Fatalf("exit code: got %d, want 0", rc)
	}
	if got != `{"entries":[]}` {
		t.Fatalf("output: got %s", got)
	}
}

// ---------------------------------------------------------------------------
// ta-entry
// ---------------------------------------------------------------------------

func TestTAEntryCommand(t *testing.T) {
	srv := newJSONServer(t, "/ta/entry/ta-789", "", `{"id":"ta-789"}`)
	defer srv.Close()

	got, rc := runCLI(t, srv, []string{"ta-entry", "ta-789"})
	if rc != 0 {
		t.Fatalf("exit code: got %d, want 0", rc)
	}
	if got != `{"id":"ta-789"}` {
		t.Fatalf("output: got %s", got)
	}
}

func TestTAEntryCommandMissingArg(t *testing.T) {
	_, rc := runCLI(t, nil, []string{"ta-entry"})
	if rc != 2 {
		t.Fatalf("exit code: got %d, want 2", rc)
	}
}

// ---------------------------------------------------------------------------
// Error paths
// ---------------------------------------------------------------------------

func TestHTTPErrorReturnsExitCode1(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer srv.Close()

	_, rc := runCLI(t, srv, []string{"library-info"})
	if rc != 1 {
		t.Fatalf("exit code: got %d, want 1", rc)
	}
}

func TestHTTPErrorWithEmptyBodyReturnsExitCode1(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, rc := runCLI(t, srv, []string{"library-info"})
	if rc != 1 {
		t.Fatalf("exit code: got %d, want 1", rc)
	}
}

func TestNetworkErrorReturnsExitCode1(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := srv.URL
	srv.Close() // close before the request

	rc := runWithBaseURL(context.Background(), []string{"library-info"}, &bytes.Buffer{}, io.Discard, url)
	if rc != 1 {
		t.Fatalf("exit code: got %d, want 1", rc)
	}
}

// ---------------------------------------------------------------------------
// CLI argument edge cases
// ---------------------------------------------------------------------------

func TestNoArgsReturnsExitCode2(t *testing.T) {
	rc := RunCLI([]string{}, &bytes.Buffer{}, io.Discard)
	if rc != 2 {
		t.Fatalf("exit code: got %d, want 2", rc)
	}
}

func TestUnknownCommandReturnsExitCode2(t *testing.T) {
	_, rc := runCLI(t, nil, []string{"does-not-exist"})
	if rc != 2 {
		t.Fatalf("exit code: got %d, want 2", rc)
	}
}

// ---------------------------------------------------------------------------
// optionalInt edge case: 0 is treated as "not set"
// ---------------------------------------------------------------------------

func TestOptionalIntZeroOmitsParam(t *testing.T) {
	// Passing --limit 0 should omit the limit parameter entirely.
	srv := newJSONServer(t, "/library/version/entries", "", `{"entries":[]}`)
	defer srv.Close()

	_, rc := runCLI(t, srv, []string{"library-versions", "--limit", "0"})
	if rc != 0 {
		t.Fatalf("exit code: got %d, want 0", rc)
	}
	// The server handler in newJSONServer will call t.Errorf if "limit" appears in the query.
}

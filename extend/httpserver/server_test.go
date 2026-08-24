package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondOK(t *testing.T) {
	w := httptest.NewRecorder()
	respondOK(w, map[string]any{"key": "value"})

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp Response
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 0 {
		t.Errorf("expected code 0, got %d", resp.Code)
	}
	if resp.Msg != "ok" {
		t.Errorf("expected msg 'ok', got '%s'", resp.Msg)
	}
}

func TestRespondErr(t *testing.T) {
	w := httptest.NewRecorder()
	respondErr(w, http.StatusBadRequest, "bad request")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
	var resp Response
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 1 {
		t.Errorf("expected code 1, got %d", resp.Code)
	}
	if resp.Msg != "bad request" {
		t.Errorf("expected msg 'bad request', got '%s'", resp.Msg)
	}
}

func TestParseExchange(t *testing.T) {
	cases := []struct {
		in   string
		want uint8
		err  bool
	}{
		{"sh", 1, false},
		{"sz", 0, false},
		{"bj", 2, false},
		{"SH", 1, false},
		{"xx", 0, true},
	}
	for _, c := range cases {
		ex, err := parseExchange(c.in)
		if c.err {
			if err == nil {
				t.Errorf("parseExchange(%q) expected error", c.in)
			}
			continue
		}
		if err != nil {
			t.Errorf("parseExchange(%q) unexpected error: %v", c.in, err)
		}
		if ex.Uint8() != c.want {
			t.Errorf("parseExchange(%q) = %d, want %d", c.in, ex.Uint8(), c.want)
		}
	}
}

func TestQueryUint16Default(t *testing.T) {
	req := httptest.NewRequest("GET", "/?count=100", nil)
	if got := queryUint16Default(req, "count", 50); got != 100 {
		t.Errorf("expected 100, got %d", got)
	}
	req2 := httptest.NewRequest("GET", "/", nil)
	if got := queryUint16Default(req2, "count", 50); got != 50 {
		t.Errorf("expected default 50, got %d", got)
	}
}

func TestCorsAllowAll(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "https://client.scalar.com")
	w := httptest.NewRecorder()
	cors(next, nil).ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("expected allow all, got %q", got)
	}
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestCorsSpecificOrigins(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	origins := []string{"https://client.scalar.com"}

	// 匹配的 Origin 回显
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "https://client.scalar.com")
	w := httptest.NewRecorder()
	cors(next, origins).ServeHTTP(w, req)
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "https://client.scalar.com" {
		t.Errorf("expected matched origin, got %q", got)
	}

	// 不匹配的 Origin 不加 CORS 头
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.Header.Set("Origin", "https://evil.example.com")
	w2 := httptest.NewRecorder()
	cors(next, origins).ServeHTTP(w2, req2)
	if got := w2.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("expected no allow header, got %q", got)
	}
}

func TestCorsPreflight(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("preflight request should not reach next handler")
	})
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "https://client.scalar.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	w := httptest.NewRecorder()
	cors(next, nil).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Methods"); got != "GET, OPTIONS" {
		t.Errorf("expected 'GET, OPTIONS', got %q", got)
	}
}

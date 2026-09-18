package revproxy_solution

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProxyHandlerSolution(t *testing.T) {
	upstreamServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom-Header") != "ProxyTest" {
			http.Error(w, "missing header", http.StatusBadRequest)
			return
		}
		body, _ := io.ReadAll(r.Body)
		w.Header().Set("X-Upstream-Echo", "Verified")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("Upstream received: " + string(body)))
	}))
	defer upstreamServer.Close()

	proxy := NewProxyHandler(upstreamServer.URL)

	req := httptest.NewRequest(http.MethodPost, "/api/data", strings.NewReader("payload-content"))
	req.Header.Set("X-Custom-Header", "ProxyTest")
	rec := httptest.NewRecorder()

	proxy.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected upstream status 201, got %d", rec.Code)
	}

	if rec.Header().Get("X-Upstream-Echo") != "Verified" {
		t.Errorf("expected header X-Upstream-Echo: Verified, got %q", rec.Header().Get("X-Upstream-Echo"))
	}

	expectedBody := "Upstream received: payload-content"
	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
	}
}

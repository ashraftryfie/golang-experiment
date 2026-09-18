package apigateway

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGatewayRoutesAndMiddlewares(t *testing.T) {
	const validKey = "gatekeeper-super-secret"
	gw := NewGateway(GatewayConfig{APIKey: validKey})

	tests := []struct {
		name           string
		method         string
		path           string
		headers        map[string]string
		wantStatusCode int
		wantHeaderKey  string
		wantBodySubstr string
	}{
		{
			name:           "public health check returns 200 without auth",
			method:         http.MethodGet,
			path:           "/healthz",
			headers:        nil,
			wantStatusCode: http.StatusOK,
			wantHeaderKey:  "X-Request-ID",
			wantBodySubstr: `"status":"ok"`,
		},
		{
			name:           "v1 protected route missing API key returns 401",
			method:         http.MethodGet,
			path:           "/v1/services/orders",
			headers:        nil,
			wantStatusCode: http.StatusUnauthorized,
			wantHeaderKey:  "X-Request-ID",
			wantBodySubstr: `"error":"unauthorized"`,
		},
		{
			name:   "v1 protected route with valid API key returns 200 and path param",
			method: http.MethodGet,
			path:   "/v1/services/orders",
			headers: map[string]string{
				"X-API-Key": validKey,
			},
			wantStatusCode: http.StatusOK,
			wantHeaderKey:  "X-Request-ID",
			wantBodySubstr: `"service":"orders"`,
		},
		{
			name:           "unknown route returns 404",
			method:         http.MethodGet,
			path:           "/unknown/path",
			headers:        nil,
			wantStatusCode: http.StatusNotFound,
			wantHeaderKey:  "X-Request-ID",
			wantBodySubstr: "404 page not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}

			rec := httptest.NewRecorder()
			gw.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatusCode {
				t.Errorf("status code = %d, want %d", rec.Code, tc.wantStatusCode)
			}

			if tc.wantHeaderKey != "" && rec.Header().Get(tc.wantHeaderKey) == "" {
				t.Errorf("expected header %q to be present in response", tc.wantHeaderKey)
			}

			if !containsSubstring(rec.Body.String(), tc.wantBodySubstr) {
				t.Errorf("expected body to contain %q, got %q", tc.wantBodySubstr, rec.Body.String())
			}
		})
	}
}

func containsSubstring(s, sub string) bool {
	return len(s) >= len(sub) && (sub == "" || jsonOrPlainContains(s, sub))
}

func jsonOrPlainContains(s, sub string) bool {
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err == nil {
		raw, _ := json.Marshal(m)
		if stringContains(string(raw), sub) {
			return true
		}
	}
	return stringContains(s, sub)
}

func stringContains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && indexOf(s, sub) >= 0)
}

func indexOf(s, substr string) int {
	n := len(substr)
	if n == 0 {
		return 0
	}
	for i := 0; i+n <= len(s); i++ {
		if s[i:i+n] == substr {
			return i
		}
	}
	return -1
}

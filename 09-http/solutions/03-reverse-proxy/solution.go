package revproxy_solution

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// NewProxyHandler forwards incoming requests to targetBaseURL.
func NewProxyHandler(targetBaseURL string) http.Handler {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	targetBaseURL = strings.TrimRight(targetBaseURL, "/")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		destURL := fmt.Sprintf("%s%s", targetBaseURL, r.URL.RequestURI())

		outReq, err := http.NewRequestWithContext(r.Context(), r.Method, destURL, r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("create outbound request error: %v", err), http.StatusInternalServerError)
			return
		}

		// Copy request headers
		for k, vv := range r.Header {
			for _, v := range vv {
				outReq.Header.Add(k, v)
			}
		}

		resp, err := client.Do(outReq)
		if err != nil {
			http.Error(w, fmt.Sprintf("proxy upstream error: %v", err), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for k, vv := range resp.Header {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}

		// Write status code
		w.WriteHeader(resp.StatusCode)

		// Stream body
		_, _ = io.Copy(w, resp.Body)
	})
}

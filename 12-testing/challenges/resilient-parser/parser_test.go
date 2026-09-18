package kvparser

import (
	"errors"
	"testing"
)

func TestParseKVTable(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "empty string",
			input:   "",
			want:    map[string]string{},
			wantErr: false,
		},
		{
			name:  "standard query string",
			input: "host=127.0.0.1&port=8080",
			want: map[string]string{
				"host": "127.0.0.1",
				"port": "8080",
			},
			wantErr: false,
		},
		{
			name:  "semicolon delimiter and whitespace",
			input: "  user = alice ; role = admin  ",
			want: map[string]string{
				"user": "alice",
				"role": "admin",
			},
			wantErr: false,
		},
		{
			name:  "skips comments and empty pairs",
			input: "host=localhost;&;#comment=ignored&timeout=5s;",
			want: map[string]string{
				"host":    "localhost",
				"timeout": "5s",
			},
			wantErr: false,
		},
		{
			name:    "missing equals delimiter",
			input:   "host=localhost&invalid_token",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty key",
			input:   "=value",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseKV(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ParseKV() err = %v, wantErr = %v", err, tc.wantErr)
			}
			if !tc.wantErr {
				if len(got) != len(tc.want) {
					t.Fatalf("got len %d, want len %d", len(got), len(tc.want))
				}
				for k, v := range tc.want {
					if got[k] != v {
						t.Errorf("key %q: got %q, want %q", k, got[k], v)
					}
				}
			}
		})
	}
}

func BenchmarkParseKV(b *testing.B) {
	sample := "host=10.0.0.1;port=9092&topic=telemetry-stream;client_id=worker-01#comment=ignored&timeout=30s"
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ParseKV(sample)
	}
}

func FuzzParseKV(f *testing.F) {
	f.Add("host=localhost&port=8080")
	f.Add("")
	f.Add(";;;;&&&&&&&")
	f.Add("# full comment line")
	f.Add("a=b;c=d&e=f")
	f.Add("====")

	f.Fuzz(func(t *testing.T, s string) {
		res, err := ParseKV(s)
		// Must never panic
		if err == nil && res == nil {
			t.Errorf("result should not be nil when err is nil")
		}
		if errors.Is(err, ErrMissingDelimiter) {
			// Expected for malformed pairs
		}
	})
}

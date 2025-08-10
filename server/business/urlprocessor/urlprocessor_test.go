package urlprocessor_test

import (
	"strings"
	"testing"

	"github.com/Babatunde50/book-crud/server/business/urlprocessor"
)

func TestURLProcessor_Process(t *testing.T) {
	tests := []struct {
		name      string
		rawURL    string
		operation urlprocessor.Operation
		want      string
		wantErr   bool
	}{
		{
			name:      "canonical example from prompt",
			rawURL:    "https://BYFOOD.com/food-EXPeriences?query=abc/",
			operation: urlprocessor.OpCanonical,
			want:      "https://BYFOOD.com/food-EXPeriences",
		},
		{
			name:      "all operations example from prompt",
			rawURL:    "https://BYFOOD.com/food-EXPeriences?query=abc/",
			operation: urlprocessor.OpAll,
			want:      "https://www.byfood.com/food-experiences",
		},
		{
			name:      "invalid operation type",
			rawURL:    "https://byfood.com/food",
			operation: "invalid",
			wantErr:   true,
		},
		{
			name:      "malformed URL",
			rawURL:    "%%%invalid-url",
			operation: urlprocessor.OpAll,
			wantErr:   true,
		},
	}

	p := urlprocessor.New()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := p.Process(tc.rawURL, string(tc.operation))

			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil for input: %s", tc.rawURL)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tc.want {
				t.Errorf("expected %q, got %q", tc.want, got)
			}

			if strings.TrimSpace(got) == "" {
				t.Errorf("processed URL is unexpectedly empty for: %s", tc.rawURL)
			}
		})
	}
}

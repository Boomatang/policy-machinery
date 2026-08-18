//go:build unit

package controller_test

import (
	"k8s.io/apimachinery/pkg/labels"
	"testing"

	"github.com/kuadrant/policy-machinery/controller"
)

func TestToLabelSelector(t *testing.T) {
	mustParse := func(s string) labels.Selector {
		selector, err := labels.Parse(s)
		if err != nil {
			t.Fatalf("test setup: failed to parse selector %q: %v", s, err)
		}
		return selector
	}

	tests := []struct {
		name  string
		input string
		want  labels.Selector
	}{
		{
			name:  "valid string",
			input: "valid",
			want:  mustParse("valid"),
		},
		{
			name:  "valid equality selector",
			input: "environment=production",
			want:  mustParse("environment=production"),
		},
		{
			name:  "valid inequality selector",
			input: "tier!=frontend",
			want:  mustParse("tier!=frontend"),
		},
		{
			name:  "invalid string falls back to Nothing",
			input: "===",
			want:  labels.Nothing(),
		},
		{
			name:  "empty string matches everything",
			input: "",
			want:  labels.Everything(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := controller.ToLabelSelector(tt.input)
			if got.String() != tt.want.String() {
				t.Errorf("ToLabelSelector(%q) = %q, want %q", tt.input, got.String(), tt.want.String())

			}
		})
	}
}

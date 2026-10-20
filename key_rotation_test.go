package main

import "testing"

func TestRotationPlanRequiresOverlap(t *testing.T) {
	tests := []struct {
		name  string
		grace int
		want  string
	}{
		{"zero downtime", 24, "overlap"},
		{"no overlap", 0, "blocked"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rotationPlan(tt.grace); got != tt.want {
				t.Fatalf("rotationPlan(%d) = %q, want %q", tt.grace, got, tt.want)
			}
		})
	}
}

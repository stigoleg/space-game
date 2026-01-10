package game

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected SemanticVersion
		wantErr  bool
	}{
		{"1.0.0", SemanticVersion{1, 0, 0}, false},
		{"v1.0.0", SemanticVersion{1, 0, 0}, false},
		{"1.2.3", SemanticVersion{1, 2, 3}, false},
		{"v2.5.10", SemanticVersion{2, 5, 10}, false},
		{"10.20.30", SemanticVersion{10, 20, 30}, false},
		{"invalid", SemanticVersion{}, true},
		{"1.2", SemanticVersion{}, true},
		{"1.2.3.4", SemanticVersion{}, true},
		{"a.b.c", SemanticVersion{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseVersion(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseVersion(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ParseVersion(%q) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ParseVersion(%q) = %v, want %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestIsNewerThan(t *testing.T) {
	tests := []struct {
		v1       SemanticVersion
		v2       SemanticVersion
		expected bool
	}{
		{SemanticVersion{1, 0, 0}, SemanticVersion{0, 9, 9}, true},
		{SemanticVersion{1, 1, 0}, SemanticVersion{1, 0, 9}, true},
		{SemanticVersion{1, 0, 1}, SemanticVersion{1, 0, 0}, true},
		{SemanticVersion{1, 0, 0}, SemanticVersion{1, 0, 0}, false},
		{SemanticVersion{1, 0, 0}, SemanticVersion{1, 0, 1}, false},
		{SemanticVersion{1, 0, 0}, SemanticVersion{1, 1, 0}, false},
		{SemanticVersion{1, 0, 0}, SemanticVersion{2, 0, 0}, false},
		{SemanticVersion{2, 0, 0}, SemanticVersion{1, 9, 9}, true},
	}

	for _, tt := range tests {
		t.Run(tt.v1.String()+"_vs_"+tt.v2.String(), func(t *testing.T) {
			result := tt.v1.IsNewerThan(tt.v2)
			if result != tt.expected {
				t.Errorf("%v.IsNewerThan(%v) = %v, want %v", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestVersionString(t *testing.T) {
	v := SemanticVersion{1, 2, 3}
	if v.String() != "1.2.3" {
		t.Errorf("String() = %q, want %q", v.String(), "1.2.3")
	}
	if v.StringWithV() != "v1.2.3" {
		t.Errorf("StringWithV() = %q, want %q", v.StringWithV(), "v1.2.3")
	}
}

func TestVersionEquals(t *testing.T) {
	v1 := SemanticVersion{1, 2, 3}
	v2 := SemanticVersion{1, 2, 3}
	v3 := SemanticVersion{1, 2, 4}

	if !v1.Equals(v2) {
		t.Errorf("%v.Equals(%v) = false, want true", v1, v2)
	}
	if v1.Equals(v3) {
		t.Errorf("%v.Equals(%v) = true, want false", v1, v3)
	}
}

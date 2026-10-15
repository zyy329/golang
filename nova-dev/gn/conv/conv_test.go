package conv

import "testing"

func TestFloatToInt(t *testing.T) {
	testS := []struct {
		id   string
		f    float64
		want int
	}{
		{"", 6.99999999, 7},
		{"", 6.4, 6},
		{"", 6.5, 7},
		{"", 8520095.656, 8520096},
		{"", 1.656, 2},
		{"", -0.4, 0},
		{"", -0.5, -1},
		{"", 0, 0},
		{"", 0.4, 0},
		{"", 0.499999, 0},
		{"", 0.5, 1},
		{"", -1, -1},
		{"", -6.99999999, -7},
		{"", -7.00000001, -7},
		{"", -7.4, -7},
		{"", -7.5, -8},
		{"", -66, -66},
	}

	for _, tt := range testS {
		if got := FloatToInt(tt.f); got != tt.want {
			t.Errorf("FloatToInt(%v) = %v, want %d", tt.f, got, tt.want)
		}
	}
}

package currencyn

import (
	"testing"
)

func TestFromStr(t *testing.T) {
	testS := []struct {
		id       string
		amount   string
		currency string
		want     int64
	}{
		{"Valid IDR", "1000000000000000", "IDR", 1000000000000000},
		{"Valid IDR", "1000", "IDR", 1000},
		{"Valid IDR", "1000.1", "IDR", 1000},
		{"Valid IDR", "1000.666", "IDR", 1000},
		{"Valid IDR -", "-1000", "IDR", -1000},
		{"Valid IDR 0 -", "-1000.1", "IDR", -1000},
		{"Valid IDR 0 -", "-1000.6", "IDR", -1000},
		{"Valid IDR", "1", "IDR", 1},
		{"Valid IDR -", "-1", "IDR", -1},
		{"Valid VND", "1", "VND", 1},
		{"Valid USD", "1000", "USD", 1000000},
		{"Valid USD", "1000.1", "USD", 1000100},
		{"Valid USD", "1000.50", "USD", 1000500},
		{"Valid USD", "1000.501", "USD", 1000501},
		{"Valid USD", "8520095.656", "USD", 8520095656},
		{"Valid USD", "-8520095.656", "USD", -8520095656},
		{"Valid USD", "8520095.6564", "USD", 8520095656},
		{"Valid USD", "-8520095.6564", "USD", -8520095656},
		{"Invalid USD", "8520095.6566", "USD", 8520095656},
		{"Valid USD", "1000.5011", "USD", 1000501},
		{"Valid USD", "1000000000000000", "USD", 1000000000000000000},
		{"Valid USD", "10000000000000.5", "USD", 10000000000000500},
		{"Valid USD -", "-10000000000000.5", "USD", -10000000000000500},
		{"Valid GC", "10000000000000.5", "USD", 10000000000000500},
		{"Valid SC", "10000000000000.5", "USD", 10000000000000500},
		{"Invalid str", "abc", "IDR", 0},
		//{"Invalid currency", "1000.50", "XYZ", 0},
	}

	for _, tt := range testS {
		if got, _ := FromStr(tt.amount, tt.currency); got != tt.want {
			t.Errorf("FromStr(%q, %q) = %d, want %d", tt.amount, tt.currency, got, tt.want)
		}
	}
}

func TestToStr(t *testing.T) {
	testS := []struct {
		id       string
		amount   int64
		currency string
		want     string
	}{
		{"Valid IDR", 100000000000000, "IDR", "100000000000000"},
		{"Valid IDR", 100, "IDR", "100"},
		{"Valid IDR -", -100, "IDR", "-100"},
		{"Valid IDR", 1, "IDR", "1"},
		{"Valid IDR", 1, "IDR", "1"},
		{"Valid VND", 1, "VND", "1"},
		{"Valid USD", 10005000, "USD", "10005"},
		{"Valid USD", 100050, "USD", "100.05"},
		{"Valid USD", 10005001, "USD", "10005.001"},
		{"Valid USD", 100050011, "USD", "100050.011"},
		{"Valid USD", 1000500110000000000, "USD", "1000500110000000"},
		{"Valid USD", 1000500110000011, "USD", "1000500110000.011"},
		{"Valid GC", 1000500110000011, "USD", "1000500110000.011"},
		{"Valid SC", 1000500110000011, "USD", "1000500110000.011"},
		//{"Invalid currency", 100050, "XYZ", "100.05"},
	}

	for _, tt := range testS {
		if got := ToStr(tt.amount, tt.currency); got != tt.want {
			t.Errorf("ToStr(%d, %q) = %q, want %q", tt.amount, tt.currency, got, tt.want)
		}
	}
}

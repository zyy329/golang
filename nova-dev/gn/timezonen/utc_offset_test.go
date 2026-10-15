package timezonen

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestDiff1(t *testing.T) {
	offset := map[string]int{}
	for i := -12; i <= 14; i++ {
		key := fmt.Sprintf("\"UTC%.02d:00\":", i)
		value := i * 60 * 60
		fmt.Printf("%s %d,\n", key, value)
		offset[key] = value
	}
}

func TestDiff2(t *testing.T) {
	offset := map[string]int{}
	ks := []int{-9, -3, 3, 4, 5, 6, 9, 10}
	for _, i := range ks {
		key := fmt.Sprintf("\"UTC%.02d:30\":", i)
		value := i*60*60 + 30*60
		fmt.Printf("%s %d,\n", key, value)
		offset[key] = value
	}
}

func TestOffset(t *testing.T) {
	offsets := GetUtcOffsets()

	for k, v := range offsets {
		// 构造一个包含偏移的时间字符串
		s := "2025-01-01 00:00:00 " + k[len("UTC"):] // “+14:00”
		tm, err := time.Parse("2006-01-02 15:04:05 -07:00", s)
		if err != nil {
			fmt.Println("parse error:", k, err)
			continue
		}

		_, sec := tm.Zone()

		if sec != v {
			t.Errorf("❌ mismatch: %s expect=%d got=%d\n", k, v, sec)
		} else {
			t.Logf("✅ ok: %s = %d\n", k, sec)
		}
	}
}

func TestUtcOffsets(t *testing.T) {
	for k, v := range utcOffset {
		if !strings.HasPrefix(k, "UTC") {
			t.Errorf("invalid key format: %s", k)
			continue
		}

		sign := k[3]
		hm := k[4:]
		parts := strings.Split(hm, ":")
		if len(parts) != 2 {
			t.Errorf("invalid time format: %s", k)
			continue
		}

		hours, err1 := strconv.Atoi(parts[0])
		minutes, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			t.Errorf("invalid hour/minute: %s", k)
			continue
		}

		expected := hours*3600 + minutes*60
		if sign == '-' {
			expected = -expected
		}

		if expected != v {
			t.Errorf("UTC offset mismatch for %s: expected %d, got %d", k, expected, v)
		}
	}
}

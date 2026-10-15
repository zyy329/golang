package timen

import (
	"testing"
	"time"
)

func TestSetZone(t *testing.T) {
	t.Log(time.Local)
	t.Log(time.Now())

	SetZone(7)
	t.Log(time.Local)
	t.Log(time.Now())

	SetZone(6.5)
	t.Log(time.Local)
	t.Log(time.Now())
}

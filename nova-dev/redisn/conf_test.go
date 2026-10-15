package redisn

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestLoadAllConf(t *testing.T) {
	mp, _ := loadAllConf()
	js, _ := json.Marshal(mp)
	t.Logf("conf: %+v", string(js))
}

func TestGetClient(t *testing.T) {
	client := NewClient("common")

	cmd1 := client.Set(context.Background(), "sdcsd", "345345", 1*time.Second)
	t.Logf("cmd: %+v", cmd1)

	cmd := client.Get(context.Background(), "sdcsd")
	t.Logf("cmd: %+v", cmd)
}

package etcdn

import (
	"testing"

	"nova/gn/confkey"
)

func TestLoadEtcd(t *testing.T) {
	//c, err := LoadEtcd(confkey.Common)
	c, err := LoadEtcd(confkey.Common + "sdc")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(c)
}

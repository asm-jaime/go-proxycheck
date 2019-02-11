package pcheck

import (
	"testing"
)

var proxs = []string{
	"118.97.153.250:53281",
	"200.52.144.77:8080",
	"202.51.106.195:8080",
	"93.77.14.13:32410",
	"190.214.13.90:21776",
}

func TestAsynProx(t *testing.T) {
	tlist := ProxyCheck(proxs)
	if len(tlist) < 1 {
		t.Error("all proxy not available, or this proxy list is outdated")
	}
}

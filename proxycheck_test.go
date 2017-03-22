package main

import (
	"fmt"
	"testing"
)

func _TestWriteTProx(t *testing.T) { // {{{
	prox := Prox{}
	// set test default
	prox.file = "data/test_tprox.txt"
	prox.tfile = "data/test_tprox.txt"

	proxies := []string{
		"47.90.75.157:3128",
		"185.117.153.230:3129",
		"213.165.166.210:3128",
		"139.59.102.243:8080",
		"192.99.159.91:8080",
		"83.171.108.210:8081",
		"202.158.52.212:8080",
		"144.217.46.198:3128",
		"195.239.66.166:8081",
		"195.225.172.243:8080",
	}

	prox.tlist = make([]*string, len(proxies))
	for i := range prox.tlist {
		prox.tlist[i] = &proxies[i]
	}

	err := prox.writeTProx()

	if err != nil {
		t.Error("Expected: ", err)
	}

	err = prox.readProx()
	if err != nil {
		t.Error("Expected: ", err)
	}

	if len(prox.list)-1 != len(prox.tlist) {
		fmt.Printf("list: %v", prox.list[10])
		t.Error("\ncount list: ", len(prox.list), ", tlist: ", len(prox.tlist))
	}

} // }}}

func TestSynProx(t *testing.T) {
	prox := Prox{
		tlist: make([]*string, 0),
	}

	proxies := []string{
		"google.com:80",
		"ya.ru:80",
	}

	prox.list = proxies
	prox.synProx()

	if len(prox.tlist) == 0 {
		t.Error("no tlist")
	}

}

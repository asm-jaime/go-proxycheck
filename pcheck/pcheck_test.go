package pcheck

import (
	"fmt"
	"testing"
)

func _TestWriteTProx(t *testing.T) { // {{{
	prox := Prox{}
	// set test default
	prox.File = "../data/test_tprox.txt"
	prox.TFile = "../data/test_tprox.txt"

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

	prox.TList = make([]*string, len(proxies))
	for i := range prox.TList {
		prox.TList[i] = &proxies[i]
	}

	err := prox.writeTProx()

	if err != nil {
		t.Error("Expected: ", err)
	}

	err = prox.readProx()
	if err != nil {
		t.Error("Expected: ", err)
	}

	if len(prox.List)-1 != len(prox.TList) {
		fmt.Printf("list: %v", prox.List[10])
		t.Error("\ncount list: ", len(prox.List), ", tlist: ", len(prox.TList))
	}

} // }}}

func _TestSynProx(t *testing.T) { // {{{
	prox := Prox{
		tlist: make([]*string, 0),
	}

	proxies := []string{
		"google.com:80",
		"ya.ru:80",
	}

	prox.List = proxies
	prox.synProx()

	if len(prox.TList) == 0 {
		t.Error("no tlist")
	}

} // }}}

func TestRequest(t *testing.T) { // {{{
	prox := Prox{}

	proxy := "52.208.118.39:3128"
	req := "http://www.google.com"

	err := prox.oneProx(proxy)
	if err != nil {
		t.Error("proxy not available: ", proxy)
		return
	}
	_, err = prox.req(req)
	if err != nil {
		t.Error("request url not available: ", req)
		return
	}

	res, err := prox.proxyReq(req, proxy)
	if err != nil {
		t.Error("url through proxy not available: ", err)
		return
	}
	fmt.Printf("\neverything is allright, status: %v\n", res.Status)
} // }}}

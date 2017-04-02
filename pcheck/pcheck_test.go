package pcheck

import (
	"fmt"
	"testing"
	"time"
)

// ========== tests
// del prefix '_' for make test available

func _TestRequest(t *testing.T) { // {{{
	prox := Prox{}

	proxy := "52.208.118.39:3128"
	req := "http://www.google.com"

	err := prox.OneProx(proxy)
	if err != nil {
		t.Error("proxy not available: ", proxy)
		return
	}
	_, err = prox.Req(req)
	if err != nil {
		t.Error("request url not available: ", req)
		return
	}

	res, err := prox.ProxyReq(req, proxy)
	if err != nil {
		t.Error("url through proxy not available: ", err)
		return
	}
	fmt.Printf("\neverything is allright, status: %v\n", res.Status)
} // }}}

func _TestSynProx(t *testing.T) { // {{{
	prox := Prox{}
	File := "../data/test_prox.txt"
	TFile := "../data/test_tprox.txt"

	// load proxies from file, (prox.txt as default)
	list, err := prox.ReadProx(File)
	if err != nil {
		t.Error("read proxies: ", err)
		return
	}
	if len(list) < 1 {
		t.Error("empty list of proxy")
	}

	start := time.Now()
	tlist := prox.SynProx(&list)
	elapsed := time.Since(start)
	fmt.Printf("\nsyn: %s\n", elapsed)

	err = prox.WriteTProx(TFile, &tlist)
	if err != nil {
		t.Error("write proxies: ", err)
	}
} // }}}

func TestAsynProx(t *testing.T) { // {{{
	prox := Prox{}
	File := "../data/test_prox.txt"
	TFile := "../data/test_tprox.txt"

	// load proxies from file, (prox.txt as default)
	list, err := prox.ReadProx(File)
	if err != nil {
		t.Error("read proxies: ", err)
		return
	}
	if len(list) < 1 {
		t.Error("empty list of proxy")
	}
	start := time.Now()
	tlist := prox.AsynProx(&list)
	elapsed := time.Since(start)
	fmt.Printf("\nasyn: %s\n", elapsed)

	err = prox.WriteTProx(TFile, &tlist)
	if err != nil {
		t.Error("write proxies: ", err)
	}
} // }}}

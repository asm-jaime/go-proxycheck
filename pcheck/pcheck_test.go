package pcheck

import (
	"fmt"
	"testing"
	"time"
)

// ========== configs

func SetDefault(prox *Prox) { // {{{
	prox.File = "../data/test_prox.txt"
	prox.TFile = "../data/test_tprox.txt"
	prox.Timeout = 1 * time.Second
} // }}}

// ========== tests
// del prefix '_' for make test available

func TestWriteTProx(t *testing.T) { // {{{
	prox := Prox{}
	// set test default
	SetDefault(&prox)

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

	err := prox.WriteTProx()

	if err != nil {
		t.Error("Expected: ", err)
	}

	err = prox.ReadProx()
	if err != nil {
		t.Error("Expected: ", err)
	}

	if len(prox.List)-1 != len(prox.TList) {
		fmt.Printf("list: %v", prox.List[10])
		t.Error("\ncount list: ", len(prox.List), ", tlist: ", len(prox.TList))
	}

} // }}}

func TestSynProx(t *testing.T) { // {{{
	prox := Prox{
		TList: make([]*string, 0),
	}

	proxies := []string{
		"google.com:80",
		"ya.ru:80",
	}

	prox.List = proxies
	prox.SynProx()

	if len(prox.TList) == 0 {
		t.Error("no tlist for synProx (maybe you not have connections?)")
	}

} // }}}

func TestRequest(t *testing.T) { // {{{
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

func TestAsynMainFunc(t *testing.T) { // {{{
	prox := Prox{}
	SetDefault(&prox)

	// load proxies from file, (prox.txt as default)
	err := prox.ReadProx()
	if err != nil {
		t.Error("asyn prox: ", err)
		return
	}

	prox.AsynProx()

	err = prox.WriteTProx()
	if err != nil {
		t.Error("write file: ", err)
	}

	var input string
	fmt.Scanln(&input)
} // }}}

func TestSynMainFunc(t *testing.T) { // {{{
	prox := Prox{}
	SetDefault(&prox)

	// load proxies from file, (prox.txt as default)
	err := prox.ReadProx()
	if err != nil {
		t.Error("read proxies: ", err)
		return
	}

	prox.SynProx()

	err = prox.WriteTProx()
	if err != nil {
		t.Error("write proxies: ", err)
	}

	var input string
	fmt.Scanln(&input)
} // }}}

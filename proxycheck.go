package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type Prox struct {
	file  string
	tfile string
	list  []string
	tlist []*string
	sync.RWMutex
}

func (prox *Prox) setDefault() { // {{{
	prox.file = "data/prox.txt"
	prox.tfile = "data/tprox.txt"

} // }}}

func (prox *Prox) readProx() (err error) { // {{{
	file, err := ioutil.ReadFile(prox.file)
	if err != nil {
		return err
	}
	text := string(file)
	prox.list = strings.Split(text, "\r\n")
	return err
} // }}}

func (prox *Prox) writeTProx() (err error) { // {{{
	fmt.Println("====================")
	fmt.Println("start to write reachable proxy list")
	fmt.Println("====================")

	if len(prox.tlist) < 1 {
		return errors.New("empty list with available proxies")
	}

	_, err = os.Stat(prox.tfile)
	if os.IsExist(err) {
		err = os.Remove(prox.tfile)
		if err != nil {
			return errors.New("delete file failed")
		}
	}

	file, err := os.Create(prox.tfile)
	if err != nil {
		return errors.New("can't create file")
	}
	defer file.Close()

	// write proxy list to file
	for _, proxy := range prox.tlist {
		_, err = file.WriteString(*proxy + "\n")
		if err != nil {
			return errors.New("can't write file")
		}
	}

	// save changes
	err = file.Sync()
	if err != nil {
		return errors.New("can't save changes file")
	}

	fmt.Println("====================")
	fmt.Println("successful complete writing reachable proxy list!")
	fmt.Println("====================")
	return
} // }}}

// ========== syn prox

func (prox *Prox) synProx() { // {{{
	prox.Lock()
	defer prox.Unlock()
	for _, proxy := range prox.list {
		conn, err := net.DialTimeout("tcp", proxy, 1*time.Second)
		fmt.Print(err)
		if err == nil {
			prox.tlist = append(prox.tlist, &proxy)
			conn.Close()
		}
	}
} // }}}

// ========== asyn prox

func (prox *Prox) asynProx() { // {{{
	for i := 0; i < len(prox.list); i++ {
		go prox.Dial(i)
	}
} // }}}

func (prox *Prox) Dial(num int) { // {{{
	prox.Lock()
	defer prox.Unlock()
	conn, err := net.DialTimeout("tcp", prox.list[num], 1*time.Second)
	fmt.Print(err)
	if err == nil {
		defer conn.Close()
		prox.tlist = append(prox.tlist, &prox.list[num])
	}
} // }}}

// ========== request throught proxy

func request(proxy string, req_url string) (response *http.Response, err error) { //{{{
	proxyUrl, err := url.Parse(proxy)
	fmt.Printf("\nurl proxy: %v\n", proxyUrl)
	client := &http.Client{
		Timeout:   time.Second * 1,
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
	}
	response, err = client.Get(req_url) // do request http://example.com
	return response, err
} // }}}

func main() {
	// settings
	prox := Prox{}
	prox.setDefault()

	// load proxies from file, (prox.txt as default)
	err := prox.readProx()
	if err != nil {
		fmt.Printf("\nerr: %v, no proxy-list\n", err)
		return
	}

	prox.asynProx()
	fmt.Print(prox.tlist)

	var input string
	fmt.Scanln(&input)
}

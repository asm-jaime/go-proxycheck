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
	urls  []string
	sync.RWMutex
}

// ========== configs

func (prox *Prox) setDefault() { // {{{
	prox.file = "data/prox.txt"
	prox.tfile = "data/tprox.txt"

} // }}}

// ========== file operations

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
	prox.Lock()
	defer prox.Unlock()

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
	for i := range prox.tlist {
		fmt.Println(*prox.tlist[i])
		_, err = file.WriteString(string(*prox.tlist[i]) + "\r\n")
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

// ========== a prox

func (prox *Prox) oneProx(proxy string) (err error) {
	conn, err := net.DialTimeout("tcp", proxy, 3*time.Second)
	if err == nil {
		defer conn.Close()
	}
	return err
}

// ========== syn prox

func (prox *Prox) synProx() { // {{{
	prox.Lock()
	defer prox.Unlock()
	for i, proxy := range prox.list {
		conn, err := net.DialTimeout("tcp", proxy, 1*time.Second)
		fmt.Print(err)
		if err == nil {
			prox.tlist = append(prox.tlist, new(string))
			prox.tlist[len(prox.tlist)-1] = &prox.list[i]
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

// ========== requests

func (prox *Prox) req(reqURL string) (data string, err error) { // {{{
	res, err := http.Get(reqURL)
	if err != nil {
		return data, err
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return data, err
	}
	data = string(body)
	return data, err
} // }}}

func (prox *Prox) proxyReq(req string, proxy string) (res *http.Response, err error) {
	timeout := time.Duration(1 * time.Second)
	proxyURL, err := url.Parse("http://" + proxy)
	reqURL, err := url.Parse(req)

	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	client := &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}

	res, err = client.Get(reqURL.String())
	return res, err
}

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

	prox.synProx()
	fmt.Printf("\ntlist: %v\n", prox.tlist)

	var input string
	fmt.Scanln(&input)
}

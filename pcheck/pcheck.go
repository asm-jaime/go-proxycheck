package pcheck

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

// Prox structure for processing proxy list
type Prox struct {
	File  string
	TFile string
	List  []string
	TList []*string
	Urls  []string
	sync.RWMutex
}

// ========== configs

func (prox *Prox) SetDefault() { // {{{
	prox.File = "data/prox.txt"
	prox.TFile = "data/tprox.txt"

} // }}}

// ========== file operations

func (prox *Prox) ReadProx() (err error) { // {{{
	file, err := ioutil.ReadFile(prox.File)
	if err != nil {
		return err
	}
	text := string(file)
	prox.List = strings.Split(text, "\r\n")
	return err
} // }}}

func (prox *Prox) WriteTProx() (err error) { // {{{
	prox.Lock()
	defer prox.Unlock()

	fmt.Println("====================")
	fmt.Println("start to write reachable proxy list")
	fmt.Println("====================")

	if len(prox.TList) < 1 {
		return errors.New("empty list with available proxies")
	}

	_, err = os.Stat(prox.TFile)
	if os.IsExist(err) {
		err = os.Remove(prox.TFile)
		if err != nil {
			return errors.New("delete file failed")
		}
	}

	file, err := os.Create(prox.TFile)
	if err != nil {
		return errors.New("can't create file")
	}
	defer file.Close()

	// write proxy list to file
	for i := range prox.TList {
		fmt.Println(*prox.TList[i])
		_, err = file.WriteString(string(*prox.TList[i]) + "\r\n")
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

func (prox *Prox) OneProx(proxy string) (err error) {
	conn, err := net.DialTimeout("tcp", proxy, 3*time.Second)
	if err == nil {
		defer conn.Close()
	}
	return err
}

// ========== syn prox

func (prox *Prox) SynProx() { // {{{
	prox.Lock()
	defer prox.Unlock()
	for i, proxy := range prox.List {
		conn, err := net.DialTimeout("tcp", proxy, 1*time.Second)
		fmt.Print(err)
		if err == nil {
			prox.TList = append(prox.TList, new(string))
			prox.TList[len(prox.TList)-1] = &prox.List[i]
			conn.Close()
		}
	}
} // }}}

// ========== asyn prox

func (prox *Prox) AsynProx() { // {{{
	for i := 0; i < len(prox.List); i++ {
		go prox.Dial(i)
	}
} // }}}

func (prox *Prox) Dial(num int) { // {{{
	prox.Lock()
	defer prox.Unlock()
	conn, err := net.DialTimeout("tcp", prox.List[num], 1*time.Second)
	fmt.Print(err)
	if err == nil {
		defer conn.Close()
		prox.TList = append(prox.TList, &prox.List[num])
	}
} // }}}

// ========== requests

func (prox *Prox) Req(reqURL string) (data string, err error) { // {{{
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

func (prox *Prox) ProxyReq(req string, proxy string) (res *http.Response, err error) {
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

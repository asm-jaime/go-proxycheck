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
	File    string
	TFile   string
	List    []string
	TList   []*string
	Timeout time.Duration
	sync.RWMutex
}

// ========== file operations

func (prox *Prox) ReadProx(rfile string) (list []string, err error) { // {{{
	_, err = os.Stat(rfile)
	if os.IsNotExist(err) {
		return errors.New(rfile + " does not exist")
	}

	file, err := ioutil.ReadFile(rfile)
	if err != nil {
		return list, err
	}

	text := string(file)
	list = strings.Split(text, "\r\n")

	if len(list) < 1 {
		return list, errors.New(rfile + " does not contain any list")
	}

	return list, err
} // }}}

func (prox *Prox) WriteTProx(wfile string, list *[]string) (err error) { // {{{
	fmt.Println("====================")
	fmt.Println("start to write reachable proxy list")
	fmt.Println("====================")

	if len(list) < 1 {
		return errors.New("empty list with available proxies")
	}

	_, err = os.Stat(wfile)
	if os.IsExist(err) {
		err = os.Remove(wfile)
		if err != nil {
			return errors.New("delete file failed")
		}
	}

	file, err := os.Create(wfile)
	if err != nil {
		return errors.New("can't create file")
	}
	defer file.Close()

	// write proxy list to file
	for _, proxy := range *list {
		fmt.Println(proxy)
		_, err = file.WriteString(string(proxy) + "\n")
		if err != nil {
			return errors.New("can't write file")
		}
	}

	// save changes
	err = file.Sync()
	if err != nil {
		return errors.New("can't save changes file")
	}

	fmt.Println("\n====================")
	fmt.Println("successful complete writing reachable proxy list!")
	fmt.Println("====================")
	return
} // }}}

// ========== a prox

func (prox *Prox) OneProx(proxy string) (err error) { // {{{
	conn, err := net.DialTimeout("tcp", proxy, prox.Timeout)
	if err == nil {
		defer conn.Close()
	}
	return err
} // }}}

// ========== syn prox

func (prox *Prox) SynProx() { // {{{
	prox.Lock()
	defer prox.Unlock()
	for i, proxy := range prox.List {
		conn, err := net.DialTimeout("tcp", proxy, prox.Timeout)
		if err == nil {
			fmt.Printf("\n %v available", proxy)
			prox.TList = append(prox.TList, new(string))
			prox.TList[len(prox.TList)-1] = &prox.List[i]
			conn.Close()
		} else {
			fmt.Printf("\n %v not available, err: %v", proxy, err)
		}
	}
} // }}}

// ========== asyn prox

func (prox *Prox) AsynProx() { // {{{
	var wg sync.WaitGroup
	wg.Add(len(prox.List))

	for i := range prox.List {
		go prox.Dial(i, &wg)
	}

	wg.Wait()
} // }}}

func (prox *Prox) Dial(i int, wg *sync.WaitGroup) { // {{{
	conn, err := net.DialTimeout("tcp", prox.List[i], prox.Timeout)
	if err == nil {
		defer conn.Close()

		prox.Lock()
		prox.TList = append(prox.TList, new(string))
		prox.TList[len(prox.TList)-1] = &prox.List[i]
		prox.Unlock()
	} else {
		fmt.Printf("\n %v not available, err: %v", prox.List[i], err)
	}
	wg.Done()
} // }}}

// ========== worker prox

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

func (prox *Prox) ProxyReq(req string, proxy string) (res *http.Response, err error) { // {{{
	timeout := time.Duration(prox.Timeout)
	proxyURL, err := url.Parse("http://" + proxy)
	reqURL, err := url.Parse(req)

	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	client := &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}

	res, err = client.Get(reqURL.String())
	return res, err
} // }}}

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
	List []string
	sync.RWMutex
}

// ========== file operations

func (prox *Prox) ReadProx(rfile string) (list []string, err error) { // {{{
	_, err = os.Stat(rfile)
	if os.IsNotExist(err) {
		return list, errors.New(rfile + " does not exist")
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

func (prox *Prox) WriteTProx(wfile string, tlist *[]string) (err error) { // {{{
	fmt.Println("====================")
	fmt.Println("start to write reachable proxy list")
	fmt.Println("====================")

	if len(*tlist) < 1 {
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
	for _, proxy := range *tlist {
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
	conn, err := net.DialTimeout("tcp", proxy, 1*time.Second)
	if err == nil {
		defer conn.Close()
	}
	return err
} // }}}

// ========== syn prox

func (prox *Prox) SynProx(list *[]string) (tlist []string) { // {{{
	for _, proxy := range *list {
		conn, err := net.DialTimeout("tcp", proxy, 1*time.Second)
		if err == nil {
			defer conn.Close()
			fmt.Printf("\n %v available", proxy)
			tlist = append(tlist, proxy)
		} else {
			fmt.Printf("\n %v not available, err: %v", proxy, err)
		}
	}
	return tlist
} // }}}

// ========== asyn prox

func (prox *Prox) worker(wg *sync.WaitGroup, cs chan string, proxy string) { // {{{
	defer wg.Done()
	conn, err := net.DialTimeout("tcp", proxy, 4*time.Second)
	if err == nil {
		defer conn.Close()
		cs <- proxy
	}
}

func (prox *Prox) monitor(wg *sync.WaitGroup, cs chan string) {
	wg.Wait()
	close(cs)
} // }}}

func (prox *Prox) AsynProx(list *[]string) (tlist []string) { // {{{
	wg := &sync.WaitGroup{}
	cs := make(chan string)
	for _, proxy := range *list {
		wg.Add(1)
		go prox.worker(wg, cs, proxy)
	}
	go prox.monitor(wg, cs)

	for i := range cs {
		tlist = append(tlist, i)
	}
	return tlist
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

func (prox *Prox) ProxyReq(req string, proxy string) (res *http.Response, err error) { // {{{
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
} // }}}

// ========== easy ddos ()
/*
func (prox *Prox) DefProxyGet(req string, proxy string) {
	timeout := time.Duration(1 * time.Second)
	proxyURL, err := url.Parse("http://" + proxy)
	reqURL, err := url.Parse(req)

	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	client := &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}

	_, _ = client.Get(reqURL.String())
}

func (prox *Prox) DefProxyPost(req string, proxy string) {
	timeout := time.Duration(1 * time.Second)
	proxyURL, err := url.Parse("http://" + proxy)
	reqURL, err := url.Parse(req)

	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	client := &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}

	_, _ = client.Post(reqURL.String())
}

func (prox *Prox) worker_ddos(wg *sync.WaitGroup, cs chan string, proxy string) { // {{{
	defer wg.Done()
	conn, err := net.DialTimeout("tcp", proxy, 4*time.Second)
	if err == nil {
		defer conn.Close()
		cs <- proxy
	}
}

func (prox *Prox) monitor_ddos(wg *sync.WaitGroup, cs chan string) {
	wg.Wait()
	close(cs)
} // }}}

func (prox *Prox) AsynDdos(url string, list *[]string) (err error) {

}
*/

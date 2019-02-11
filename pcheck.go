package pcheck

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// ProxyCheck post any proxies and return list of available proxies
func ProxyCheck(list []string) (tlist []string) {
	wg := &sync.WaitGroup{}
	c := make(chan string)

	for _, prox := range list {
		wg.Add(1)
		go func(prox string) {
			conn, err := net.DialTimeout("tcp", prox, 10*time.Second)
			if err == nil {
				defer conn.Close()
				c <- prox
			} else {
				log.Println(err)
			}
			wg.Done()
		}(prox)
	}
	go func() {
		for i := range c {
			tlist = append(tlist, i)
		}
	}()
	wg.Wait()

	return tlist
}

// ProxyReq make a request through a proxy
func ProxyReq(req string, proxy string) (res *http.Response, err error) {
	timeout := time.Duration(2 * time.Second)
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

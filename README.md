# go-proxycheck
  Check proxy list from a file and create file with verifies proxies.

### full install

  * `git clone http://github.com/asm-jaime/go-proxycheck`
  * `go get` inside project
  
### lib only install

  * `go get github.com/asm-jaime/go-proxycheck/pcheck`

### example

```go
package main

import	"github.com/asm-jaime/go-proxycheck/pcheck"

func main() {
  prox := pcheck.Prox{
    File: "data/prox.txt",
    TFile: "data/tprox.txt",
    Timeout = 1 * time.Second,
  }

  err := prox.ReadProx()
  prox.AsynProx()  
  err = prox.WriteTProx()
}
```

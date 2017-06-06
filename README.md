# go-proxycheck
  Check proxy list from a file and create a file with verified proxies.

### full install

  * `git clone http://github.com/asm-jaime/go-proxycheck`
  * `go get` inside project
  
### lib only install

  * `go get github.com/asm-jaime/go-proxycheck/pcheck`

### example (how to use lib)

  * example main.go
  
```go
package main

import(
  "fmt"
  "time"
  "github.com/asm-jaime/go-proxycheck/pcheck"
)

func main() {
  File := "prox.txt"
  TFile := "tprox.txt"
  
  prox := pcheck.Prox{}
  list, err := prox.ReadProx(File)
  if err != nil {
    fmt.Print("\nasyn prox: ", err)
    return
  }
  tlist := prox.AsynProx(&list)
  err = prox.WriteTProx(TFile, &tlist)
  if err != nil {
    fmt.Print("\nwrite file: ", err)
  }
}
```

  * example prox.txt
  
```
47.90.75.157:3128
185.117.153.230:3129
213.165.166.210:3128
139.59.102.243:8080
192.99.159.91:8080
83.171.108.210:8081
```

I'm available to provide any information, suggestion or contribution!

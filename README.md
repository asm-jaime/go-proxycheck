# go-proxycheck
  Check proxy list.

### install

  * `go get github.com/asm-jaime/go-proxycheck/pcheck`

### example (how to use lib)

  * cmd [/example](./example)
  * lib example main.go

```go
package main

import(
  "log"
  "github.com/asm-jaime/go-proxycheck"
)

func main() {
  var proxs = []string{
    "118.97.153.250:53281",
    "200.52.144.77:8080",
    "202.51.106.195:8080",
    "93.77.14.13:32410",
    "190.214.13.90:21776",
  }

  list := pcheck.ProxyCheck(proxs)
  log.Println(list)
}
```

I'm available to provide any information, suggestion or contribution!

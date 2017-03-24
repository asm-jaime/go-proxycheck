package main

import "fmt"

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

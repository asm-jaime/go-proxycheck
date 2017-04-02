package main

import (
	"fmt"
	"os"

	"github.com/asm-jaime/go-proxycheck/pcheck"
	// pcheck "learn.go/ares.projects/go-proxycheck/pcheck"
)

// ========== info

func PrintInfo() { // {{{
	fmt.Println("=== The program is intended: ")
	fmt.Println(" * get a proxy list from a file")
	fmt.Println(" * check this list for available")
	fmt.Println(" * create a file with the available proxies")

	fmt.Println("=== command line params: ")
	fmt.Println("  <file of any proxies> : file with proxy adresses, separated by new line")
	fmt.Println("  default value: data/prox.txt")
	fmt.Println("  file content example:")

	fmt.Println("")
	fmt.Println("    0.0.0.0:8080")
	fmt.Println("    192.168.44.45:80")
	fmt.Println("")

	fmt.Println("  <file of available proxies> : file with available/accessible proxy adresses, separated by new line")
	fmt.Println("  default value: data/tprox.txt")

	fmt.Println("=== start program examples: ")
	fmt.Println("  ./go-proxycheck [<file of any proxies>] [<file of available proxies>]")
	fmt.Println("  go run main.go [<file of any proxies>] [<file of available proxies>]")
} // }}}

// ========== configs

func startChecker(args []string) {
	prox := pcheck.Prox{}
	File := "./data/prox.txt"
	TFile := "./data/tprox.txt"

	if len(args) > 3 ||
		(len(args) > 1 &&
			(args[1] == "--help" || args[1] == "help")) {
		PrintInfo()
		return
	}

	if len(args) > 1 { // set proxy file
		File = args[1]
	}
	if len(args) > 2 { // set available proxy file
		TFile = args[2]
	}

	// info
	fmt.Println("---------------")
	fmt.Println("proxies: ", File)
	fmt.Println("available proxies: ", TFile)
	fmt.Println("---------------")

	// load proxies from file, (prox.txt as default)
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

	var input string
	fmt.Scanln(&input)
}

func main() {
	args := os.Args
	startChecker(args)
}

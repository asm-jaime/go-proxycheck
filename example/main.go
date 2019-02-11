package main

import (
	"bufio"
	"errors"
	"flag"
	"github.com/asm-jaime/go-proxycheck"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func readProx(path string) (list []string, err error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return list, err
	}
	list = strings.Split(string(file), "\n")

	if len(list) < 1 {
		return list, errors.New("list does not exist")
	}
	return list, err
}

func writeProx(path string, list []string) (err error) {
	if len(list) < 1 {
		return errors.New("can not write empty list")
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)

	defer file.Close()
	defer writer.Flush()

	str := strings.Join(list, "\n")

	_, err = writer.WriteString(str)
	return err
}

type flags struct {
	startPtr   *string
	proxFile   *string
	proxFileCh *string
}

func main() {
	fls := &flags{}

	fls.startPtr = flag.String("start", "check", "start cheking proxies")
	fls.proxFile = flag.String("pf", "prox.txt", "file with all proxies")
	fls.proxFileCh = flag.String("pfc", "prox_ch.txt", "file with checked proxies")
	flag.Parse()

	switch *fls.startPtr {
	case "check":
		list, err := readProx(*fls.proxFile)
		if err != nil {
			log.Printf("\nfail read prox: %v", err)
			break
		}
		tlist := pcheck.ProxyCheck(list)
		err = writeProx(*fls.proxFileCh, tlist)
		if err != nil {
			log.Printf("\nwrite file: %v", err)
		}
	default:
		log.Println("cant recognize start command, please, try again")
	}
}

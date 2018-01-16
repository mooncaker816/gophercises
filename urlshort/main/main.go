package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/mooncaker816/gophercises/urlshort"
)

func main() {
	filetype := os.Args[1]
	var filename string
	var myhandler http.HandlerFunc
	switch filetype {
	case "json":
		filename = "urls.json"
		urls, err := getfile(filename)
		if err != nil {
			fmt.Println("get file error: ", err)
		}
		myhandler, err = urlshort.JsonHandler(urls, http.HandlerFunc(hello))
		if err != nil {
			fmt.Println("jsonhandler err: ", err)
		}
	case "yaml":
		filename = "urls.yaml"
		urls, err := getfile(filename)
		if err != nil {
			fmt.Println("get file error: ", err)
		}
		myhandler, err = urlshort.YAMLHandler(urls, http.HandlerFunc(hello))
		if err != nil {
			fmt.Println("yamlhandler err: ", err)
		}
	default:
		filename = ""
		myhandler = http.HandlerFunc(hello)
	}
	http.ListenAndServe("localhost:8000", myhandler)
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func getfile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("open input failed!")
		return nil, err
	}
	defer file.Close()
	urls, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("read input failed!")
		return nil, err
	}
	return urls, nil
}

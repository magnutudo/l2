package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	file, err := os.Create("info.text")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	resp := os.Args[1]
	answ, err := http.Get(resp)
	defer answ.Body.Close()
	io.Copy(file, answ.Body)

}

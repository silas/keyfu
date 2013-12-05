package main

import (
	"fmt"
	"net/http"
	"os"
)

func indexHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "<html><body><h1>KeyFu</h1></body></html>")
}

func handle() {
	http.HandleFunc("/", indexHandler)
}

func serve() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

func main() {
	handle()
	serve()
}

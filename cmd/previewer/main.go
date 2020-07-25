package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	println("hello previewer")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	println("buy buy previewer")
}

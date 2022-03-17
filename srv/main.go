package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	txt := flag.String("t", "", "text")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, *txt)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "10080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

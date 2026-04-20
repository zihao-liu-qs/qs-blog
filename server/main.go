package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zihao-liu-qs/qs-blog/server/api"
)

func main() {
	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	addr := ":8080"
	fmt.Printf("server running at http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

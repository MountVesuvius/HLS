package main

import (
	"fmt"
	"log"
	"net/http"
)

// i can get away with simple mux for rn
func main() {
    log.Println("Starting file-service")

    mux := http.NewServeMux()
    mux.HandleFunc("GET /file/{id}", func (w http.ResponseWriter, req *http.Request) {
        pathval := req.PathValue("id")
        log.Println("Getting ", pathval, "as requested")
        response := fmt.Sprintf("You found file %s", pathval)
       w.Write([]byte(response)) 
    })

    http.ListenAndServe(":8080", mux)
}

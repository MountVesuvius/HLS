package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"server/common/middleware"
	"server/server/db"
	// "server/db"
)

// find an easy way to have this just be loaded in
// possibly env??
const FILE_SERVICE = "http://file-service:8080/"

func GetFile(w http.ResponseWriter, req *http.Request) {
	val := req.PathValue("id")
	path := fmt.Sprintf("%s/file/%s", FILE_SERVICE, val)

	log.Println("Getting file: ", val)

	resp, err := http.Get(path)
	if err != nil {
		log.Println("Error contacting file service", err)
		http.Error(w, "Failed to get file", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	w.Write(body)
}

// won't keep routing in here forever i don't think, not sure yet
func main() {
	database := db.Connect()
	db.Sync(database)
	router := http.NewServeMux()
	router.HandleFunc("GET /file/{id}", GetFile)

	// create versioning
	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))

	applyMiddleware := middleware.CreateChain(
		middleware.EnableCors,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: applyMiddleware(v1),
	}

	log.Println("Starting server")
	server.ListenAndServe()
}

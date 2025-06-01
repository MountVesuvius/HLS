package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// ooo la la so secure
var SECRET_KEY = []byte("ThisIsTemporaryForTesting")

// this is just for examples
func GetFile(w http.ResponseWriter, req *http.Request) {
	pathval := req.PathValue("id")
	log.Println("Getting ", pathval, "as requested")
	response := fmt.Sprintf("You found file %s", pathval)
	w.Write([]byte(response))
}

// generate basic hmac
func generateSignature(path string, expires int64) string {
	data := fmt.Sprintf("%s:%d", path, expires)
	h := hmac.New(sha256.New, SECRET_KEY)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// stitch url
// probably don't need baseURL as an arg
func generatePresignedURL(baseURL, path string, validFor time.Duration) string {
	expires := time.Now().Add(validFor).Unix()
	sig := generateSignature(path, expires)
	return fmt.Sprintf("%s%s?expires=%d&sig=%s", baseURL, path, expires, sig)
}

// check for expiry and hmac match
func verifySignature(req *http.Request) bool {
	path := req.URL.Path
	expireStr := req.URL.Query().Get("expires")
	sig := req.URL.Query().Get("sig")

	expires, err := strconv.ParseInt(expireStr, 10, 64)
	if err != nil || time.Now().Unix() > expires {
		return false
	}

	expectedSig := generateSignature(path, expires)
	return hmac.Equal([]byte(sig), []byte(expectedSig))
}

// test that the sign verification works
func uploadHandler(w http.ResponseWriter, req *http.Request) {
	if !verifySignature(req) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	log.Println("Valid URL")
	// this could be abstracted, not sure if it's even worth it though
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Nice, it worked"))
}

// need to get the presigned url from somewhere
func presignedHandler(w http.ResponseWriter, req *http.Request) {
	// probably should pass the record id from the metadata db to this instead of "file"
	url := generatePresignedURL("http://localhost:8080", "/upload/file", 10*time.Minute)
	fmt.Println("Presigned URL:", url)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(url))
}

func main() {
	router := http.NewServeMux()
	// example route
	router.HandleFunc("GET /file/{id}", GetFile)

	// presigned urls (yes the both get so i can test with my browser quickly)
	router.HandleFunc("GET /upload/", uploadHandler)
	router.HandleFunc("GET /presigned/", presignedHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Starting file-service")
	server.ListenAndServe()
}

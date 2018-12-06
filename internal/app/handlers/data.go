package handlers

import (
	"compress/gzip"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Data serve data files
func Data(rw http.ResponseWriter, r *http.Request) {
	var err error

	// check for malicious path
	path := strings.Replace(r.URL.Path, "/", "", 1)
	if strings.Contains(path, "..") {
		rw.WriteHeader(400)
		rw.Write([]byte{})
		return
	}

	// default paths
	switch path {
	case "data/avatar/default":
		path = "assets/avatar.png"
	}

	// open file
	var file *os.File
	file, err = os.Open(path)
	defer file.Close()

	// set content-type and content-encoding
	rw.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(path))+"; charset=utf-8")
	rw.Header().Set("Content-Encoding", "gzip")

	// security headers
	SecurityHeaders(rw, r)

	// gzip compression
	w, _ := gzip.NewWriterLevel(rw, 2)
	defer w.Close()

	// check for error
	if err != nil {
		// unknown error
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}

	// write file
	_, err = io.Copy(w, file)
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}
}
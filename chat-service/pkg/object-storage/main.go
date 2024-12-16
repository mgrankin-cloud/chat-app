package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var (
	STORAGE_DIR         = "./storage"
	UPLOAD_PREFIX_LEN   = len("/upload/")
	DOWNLOAD_PREFIX_LEN = len("/download/")
)

type Storage struct {
	mu    sync.Mutex
	files map[string][]byte
}

func NewStorage() *Storage {
	return &Storage{
		files: make(map[string][]byte),
	}
}

func (s *Storage) Save(key string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.files[key] = data

	err := os.WriteFile(filepath.Join(STORAGE_DIR, key), data, 0644)
	if err != nil {
		log.Printf("Error creating file: %s: %v", key, err)
		return err
	}

	return nil
}

func (s *Storage) Load(key string) ([]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, exists := s.files[key]
	if exists {
		return data, true
	}

	data, err := os.ReadFile(filepath.Join(STORAGE_DIR, key))
	if err != nil {
		return nil, false
	}

	s.files[key] = data
	return data, true
}

func HandleUpload(w http.ResponseWriter, r *http.Request, s *Storage) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Path[UPLOAD_PREFIX_LEN:]

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading data", http.StatusInternalServerError)
		return
	}

	err = s.Save(key, data)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Obj %s saved successfully", key)
}

func HandleDownload(w http.ResponseWriter, r *http.Request, s *Storage) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Path[DOWNLOAD_PREFIX_LEN:]

	data, exists := s.Load(key)
	if !exists {
		http.Error(w, "Obj not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func HandleList(w http.ResponseWriter, r *http.Request, s *Storage) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	keys := make([]string, 0, len(s.files))
	for key := range s.files {
		keys = append(keys, key)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keys)
}

func main() {
	if _, err := os.Stat(STORAGE_DIR); os.IsNotExist(err) {
		err := os.Mkdir(STORAGE_DIR, 0755)
		if err != nil {
			log.Fatalf("Error creating directory %s: %v", STORAGE_DIR, err)
		}
	}

	storage := NewStorage()

	http.HandleFunc("/upload/", func(w http.ResponseWriter, r *http.Request) {
		HandleUpload(w, r, storage)
	})

	http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		HandleDownload(w, r, storage)
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		HandleList(w, r, storage)
	})

	log.Println("Listen and serve on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

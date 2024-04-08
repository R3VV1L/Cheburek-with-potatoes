package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

// Config представляет конфигурацию для сайта.
type Config struct {
	Dir   string `json:"dir"`
	Path  string `json:"path"`
	Index string `json:"index"`
}

// Configs - срез Config.
type Configs []Config

// FileServerWithPrefix представляет файловый сервер с префиксом.
type FileServerWithPrefix struct {
	Prefix  string
	Handler http.Handler
}

// NewFileServerWithPrefix создает новый файловый сервер с указанным префиксом.
func NewFileServerWithPrefix(prefix string, handler http.Handler) *FileServerWithPrefix {
	return &FileServerWithPrefix{Prefix: prefix, Handler: handler}
}

// LoadConfig загружает конфигурацию из файла или другого источника.
func LoadConfig() (Configs, error) {
	// Загрузка конфигурации из JSON-файла
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	// Распарсивание JSON-данных в структуру Configs
	var configs Configs
	err = json.Unmarshal(data, &configs)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

func main() {
	configs, err := LoadConfig()
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации:", err)
		return
	}

	fsMap := make(map[string]*FileServerWithPrefix)

	for _, config := range configs {
		fs := http.FileServer(http.Dir(config.Dir))
		fsWithPrefix := NewFileServerWithPrefix(config.Path, fs)
		fsMap[config.Path] = fsWithPrefix
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		for k, fs := range fsMap {
			if path == k || path == k+"/" || strings.HasPrefix(path, k+"/"+filepath.Base(fs.Prefix)) {
				fs.Handler.ServeHTTP(w, r)
				return
			}
		}

		http.NotFound(w, r)
	})

	http.ListenAndServe(":80", nil)
}

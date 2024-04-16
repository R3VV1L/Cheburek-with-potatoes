package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// Config представляет конфигурацию для сайта.
type Config struct {
	Dir   string `json:"dir"`   // Директория, где лежит сайт
	Path  string `json:"path"`  // Путь, по которому будет доступен сайт
	Index string `json:"index"` // Название файла-индекса
	Port  int    `json:"port"`  // Номер порта, на котором будет запущен сервер
}

// Configs - срез Config.
type Configs []Config

// LoadConfig загружает конфигурацию из файла или другого источника.
func LoadConfig() (Configs, error) {
	// Получаем текущую рабочую директорию
	currentWorkDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Собираем полный путь до конфига
	configPath := filepath.Join(currentWorkDir, "config.json")

	// Загрузка конфигурации из JSON-файла
	data, err := os.ReadFile(configPath)
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

// StartServers запускает серверы для каждой конфигурации.
func StartServers(configs Configs) {
	// Создаем отдельный обработчик для каждой конфигурации
	for _, config := range configs {
		fs := http.FileServer(http.Dir(config.Dir))
		http.Handle(config.Path, http.StripPrefix(config.Path, fs))

		go func(config Config) {
			fmt.Printf("%s runnig on %d\n", config.Index, config.Port)
			err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
			if err != nil {
				fmt.Printf("Ошибка запуска сервера на порту %d: %v\n", config.Port, err)
			}
		}(config)
	}
}

func main() {
	configs, err := LoadConfig()
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации:", err)
		return
	}

	StartServers(configs)

	select {}
}

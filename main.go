package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Собираем полный путь до конфига
	configPath := filepath.Join(cwd, "config.json")

	// Загрузка конфигурации из JSON-файла
	data, err := ioutil.ReadFile(configPath)
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

	// Создаем отдельный обработчик для каждой конфигурации
	for _, config := range configs {
		fs := http.FileServer(http.Dir(config.Dir))
		http.Handle(config.Path, http.StripPrefix(config.Path, fs))

		go func(config Config) {
			fmt.Printf("Запуск сервера на порту %d...\n", config.Port)
			err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
			if err != nil {
				fmt.Printf("Ошибка запуска сервера на порту %d: %v\n", config.Port, err)
			}
		}(config)
	}

	select {}
}

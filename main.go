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
	Dir      string `json:"dir"`      // Директория, где находится сайт
	Path     string `json:"path"`     // Путь, по которому будет обслуживаться сайт
	Index    string `json:"index"`    // Название файла-индекса
	Port     int    `json:"port"`     // Номер порта, на котором будет слушать сервер
	IndexDir string `json:"indexDir"` // Директория, где находится файл-индекс
}

// Configs - срез Config.
type Configs []Config

// LoadConfig загружает конфигурацию из JSON-файла.
func LoadConfig() (Configs, error) {
	// Получаем текущую рабочую директорию
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Собираем путь к файлу конфигурации
	configPath := filepath.Join(cwd, "config.json")

	// Читаем файл конфигурации
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Распарсиваем JSON-данные в срез Config
	var configs Configs
	err = json.Unmarshal(data, &configs)
	if err != nil {
		return nil, err
	}

	// Устанавливаем поле IndexDir для каждого Config
	for i, config := range configs {
		config.IndexDir = filepath.Join(config.Dir, config.Index)
		if _, err = os.Stat(config.IndexDir); os.IsNotExist(err) {
			return nil, fmt.Errorf("сайт не найден: %s", config.IndexDir)
		}
		configs[i] = config
	}

	return configs, nil
}

// StartServers запускает сервер HTTP для каждой конфигурации.
func StartServers(configs Configs) {
	// Создаем файловый сервер для каждой конфигурации
	for _, config := range configs {
		fileServer := http.FileServer(http.Dir(config.Dir))

		// Устанавливаем обработчик для пути сайта
		handler := http.StripPrefix(config.Path, fileServer)

		// Запускаем HTTP-сервер для сайта
		go func(config Config) {
			fmt.Printf("%s запущен на порту %d\n", config.Index, config.Port)
			err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), handler)
			if err != nil {
				fmt.Printf("Ошибка запуска сервера для %s: %v\n", config.Path, err)
			}
		}(config)
	}
}

func main() {
	configs, err := LoadConfig()
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации:", err)
		os.Exit(1)
	}

	StartServers(configs)

	// Блокируем выполнение программы бесконечно
	select {}
}

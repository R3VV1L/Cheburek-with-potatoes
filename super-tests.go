package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain(t *testing.T) {
	// Создаем mock-сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))
	defer server.Close()

	// Отправляем запрос на mock-сервер
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Запускаем сервер и отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Проверяем статус код и ответ сервера
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидаемый статус код: %d, полученный статус код: %d", http.StatusOK, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "Hello, World!" {
		t.Errorf("Ожидаемый ответ: Hello, World!, полученный ответ: %s", string(body))
	}
}

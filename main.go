package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	defaultPort = "8080"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Обработчик для всех маршрутов
	http.HandleFunc("/", proxyHandler)

	log.Printf("Сервер запущен на порту %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Извлекаем целевой URL из параметра target_url
	targetURL := r.URL.Query().Get("target_url")
	if targetURL == "" {
		http.Error(w, "Отсутствует параметр target_url", http.StatusBadRequest)
		return
	}

	log.Printf("Получен %s запрос на %s, проксируем на %s", r.Method, r.URL.Path, targetURL)

	// 2. Читаем тело исходного запроса (для последующей передачи)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// 3. Создаём новый запрос с тем же методом и телом
	req, err := http.NewRequest(r.Method, targetURL, bytes.NewReader(body))
	if err != nil {
		http.Error(w, "Ошибка создания прокси-запроса", http.StatusInternalServerError)
		return
	}

	// 4. Копируем заголовки исходного запроса (кроме Hop-by-hop)
	copyHeaders(req.Header, r.Header)

	// 5. Выполняем HTTP-запрос с таймаутом
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Ошибка при запросе к %s: %v", targetURL, err)
		http.Error(w, "Ошибка связи с целевым сервером", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 6. Копируем заголовки ответа клиенту
	copyHeaders(w.Header(), resp.Header)

	// 7. Устанавливаем статус ответа
	w.WriteHeader(resp.StatusCode)

	// 8. Копируем тело ответа
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Ошибка при отправке ответа клиенту: %v", err)
	}
}

// copyHeaders копирует все заголовки из src в dst, пропуская hop-by-hop
func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		// Пропускаем заголовки, которые не должны передаваться прокси
		if shouldSkipHeader(k) {
			continue
		}
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// shouldSkipHeader определяет заголовки, которые не должны копироваться при проксировании
func shouldSkipHeader(header string) bool {
	switch header {
	case "Connection", "Keep-Alive", "Proxy-Authenticate", "Proxy-Authorization",
		"Te", "Trailers", "Transfer-Encoding", "Upgrade":
		return true
	default:
		return false
	}
}

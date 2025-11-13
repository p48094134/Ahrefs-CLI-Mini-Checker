package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Структуры для парсинга JSON-ответа от Ahrefs API
// Основано на документации Ahrefs Batch Analysis API
type AhrefsResponse struct {
	Result struct {
		Targets []TargetMetrics `json:"targets"`
	} `json:"result"`
}

type TargetMetrics struct {
	Target        string  `json:"target"`
	DomainRating  float64 `json:"domain_rating"`
	Backlinks     int     `json:"backlinks"`
	RefDomains    int     `json:"refdomains"`
	OrganicKeywords int   `json:"organic_keywords"`
	OrganicTraffic int    `json:"organic_traffic"`
}

func main() {
	// --- Авторство ---
	// Автор: Частный SEO специалист
	// Сайт: https://private-seo.com/ru/
	
	// Определение флагов командной строки для токена и домена
	apiToken := flag.String("token", "", "Ваш Ahrefs API токен (обязательно)")
	domain := flag.String("domain", "", "Домен для проверки (обязательно)")
	flag.Parse()

	// Проверка, что флаги были переданы
	if *apiToken == "" || *domain == "" {
		fmt.Println("Ошибка: Необходимо указать -token и -domain.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Получение данных для домена: %s...\n", *domain)
	
	// Формирование URL для запроса
	apiURL := fmt.Sprintf("https://api.ahrefs.com/v3/batch-analysis?targets=%s&mode=domain&select=domain_rating,backlinks,refdomains,organic_keywords,organic_traffic", *domain)

	// Создание нового HTTP запроса
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Printf("Ошибка при создании запроса: %v\n", err)
		os.Exit(1)
	}

	// Добавление заголовка авторизации
	req.Header.Add("Authorization", "Bearer "+*apiToken)
	req.Header.Add("Accept", "application/json")

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка при выполнении запроса: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Чтение тела ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Ошибка при чтении ответа: %v\n", err)
		os.Exit(1)
	}
	
	// Проверка на ошибки со стороны API
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Ошибка от Ahrefs API (Статус: %d): %s\n", resp.StatusCode, string(body))
		os.Exit(1)
	}

	// Парсинг JSON
	var ahrefsData AhrefsResponse
	if err := json.Unmarshal(body, &ahrefsData); err != nil {
		fmt.Printf("Ошибка при парсинге JSON: %v\n", err)
		os.Exit(1)
	}

	if len(ahrefsData.Result.Targets) == 0 {
		fmt.Println("Не удалось получить метрики для указанного домена.")
		os.Exit(1)
	}

	// Вывод результатов
	metrics := ahrefsData.Result.Targets[0]
	fmt.Println("\n--- Результаты Ahrefs ---")
	fmt.Printf("Домен:           %s\n", metrics.Target)
	fmt.Printf("Domain Rating (DR): %.0f\n", metrics.DomainRating)
	fmt.Printf("Бэклинки:        %d\n", metrics.Backlinks)
	fmt.Printf("Ссылающиеся домены: %d\n", metrics.RefDomains)
    fmt.Printf("Органические ключевые слова: %d\n", metrics.OrganicKeywords)
	fmt.Printf("Органический трафик: %d\n", metrics.OrganicTraffic)
	fmt.Println("-------------------------")
}

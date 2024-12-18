package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "fmt"
)

type Rates struct {
    Rates Key `json:"rates"`
}

type Key struct {
    RUB float64 `json:"RUB"`
}

func getDollarPrice() (Rates, error) {
    url := "https://api.exchangerate-api.com/v4/latest/USD"

    // Запрос
    response, err := http.Get(url)
    if err != nil {
        log.Printf("Error fetching data from %s: %v", url, err)
        return Rates{}, err
    }
    defer response.Body.Close()

    data, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Printf("Error reading response body: %v", err)
        return Rates{}, err
    }

    var rates Rates
    err = json.Unmarshal(data, &rates)
    if err != nil {
        log.Printf("Error unmarshalling JSON data: %v", err)
        return Rates{}, err
    }

    return rates, nil
}

// Обработчик для получения цены
func priceHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*") // Разрешить все источники
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")


    rates, err := getDollarPrice()
    if err != nil {
        http.Error(w, "Failed to get dollar price", http.StatusInternalServerError)
        return
    }

	// курс доллара в JSON
    response := map[string]float64{"RUB": rates.Rates.RUB}
    json.NewEncoder(w).Encode(response)
    fmt.Println(response)
}   

//запуск сервера 
func main() {
    http.HandleFunc("/api/price", priceHandler)
    log.Fatal(http.ListenAndServe("localhost:8080", nil))
}




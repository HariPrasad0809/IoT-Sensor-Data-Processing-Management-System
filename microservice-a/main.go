package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// SensorData struct
type SensorData struct {
	SensorValue float64 `json:"sensor_value"`
	SensorType  string  `json:"sensor_type"`
	ID1         string  `json:"id1"`
	ID2         int     `json:"id2"`
	Timestamp   string  `json:"timestamp"`
}

var (
	sensorType   = "Temperature" // You can change sensor type
	frequencySec = 5             // default 5 seconds
	mutex        = &sync.Mutex{}
)

func generateData() SensorData {
	id1 := string(rune('A' + rand.Intn(26))) // Random capital letter
	id2 := rand.Intn(100)                    // Random int 0-99
	return SensorData{
		SensorValue: rand.Float64() * 100,
		SensorType:  sensorType,
		ID1:         id1,
		ID2:         id2,
		Timestamp:   time.Now().Format(time.RFC3339),
	}
}

func sendData(data SensorData) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling data:", err)
		return
	}

	// Replace with Microservice B URL
	resp, err := http.Post("http://localhost:8081/sensor-data", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error sending data:", err)
		return
	}
	defer resp.Body.Close()
	log.Println("Data sent:", data)
}

func dataLoop() {
	for {
		mutex.Lock()
		delay := time.Duration(frequencySec) * time.Second
		mutex.Unlock()

		data := generateData()
		sendData(data)
		time.Sleep(delay)
	}
}

func changeFrequencyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type Req struct {
		Frequency int `json:"frequency_sec"`
	}

	var req Req
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON"))
		return
	}

	mutex.Lock()
	frequencySec = req.Frequency
	mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Frequency updated"))
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// REST endpoint to change frequency
	http.HandleFunc("/change-frequency", changeFrequencyHandler)

	go func() {
		log.Println("Starting data generation loop...")
		dataLoop()
	}()

	log.Println("Microservice A running on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

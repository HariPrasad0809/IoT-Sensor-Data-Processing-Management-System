package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// SensorData struct
type SensorData struct {
	ID          int     `json:"id"`
	SensorValue float64 `json:"sensor_value"`
	SensorType  string  `json:"sensor_type"`
	ID1         string  `json:"id1"`
	ID2         int     `json:"id2"`
	Timestamp   string  `json:"timestamp"`
}

var db *sql.DB

func connectDB() {
	var err error
	// Update user:password@tcp(host:port)/dbname
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/sensor_db")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create table if not exists
	query := `CREATE TABLE IF NOT EXISTS sensor_data (
		id INT AUTO_INCREMENT PRIMARY KEY,
		sensor_value FLOAT,
		sensor_type VARCHAR(50),
		id1 CHAR(1),
		id2 INT,
		timestamp DATETIME
	)`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to DB and ensured table exists")
}

func insertData(data SensorData) {
	query := "INSERT INTO sensor_data(sensor_value, sensor_type, id1, id2, timestamp) VALUES(?,?,?,?,?)"
	_, err := db.Exec(query, data.SensorValue, data.SensorType, data.ID1, data.ID2, data.Timestamp)
	if err != nil {
		log.Println("Insert error:", err)
	}
}

// POST /sensor-data
func sensorDataHandler(w http.ResponseWriter, r *http.Request) {
	var data SensorData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON"))
		return
	}
	insertData(data)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data inserted"))
}

// GET /sensor-data?id1=A&id2=1&page=1&limit=10
func getDataHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id1 := q.Get("id1")
	id2Str := q.Get("id2")
	pageStr := q.Get("page")
	limitStr := q.Get("limit")

	id2, _ := strconv.Atoi(id2Str)
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := "SELECT id, sensor_value, sensor_type, id1, id2, timestamp FROM sensor_data WHERE id1=? AND id2=? LIMIT ? OFFSET ?"
	rows, err := db.Query(query, id1, id2, limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Query error"))
		return
	}
	defer rows.Close()

	var result []SensorData
	for rows.Next() {
		var d SensorData
		err := rows.Scan(&d.ID, &d.SensorValue, &d.SensorType, &d.ID1, &d.ID2, &d.Timestamp)
		if err != nil {
			log.Println(err)
			continue
		}
		result = append(result, d)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// DELETE /sensor-data?id1=A&id2=1
func deleteDataHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id1 := q.Get("id1")
	id2Str := q.Get("id2")
	id2, _ := strconv.Atoi(id2Str)

	query := "DELETE FROM sensor_data WHERE id1=? AND id2=?"
	_, err := db.Exec(query, id1, id2)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Delete error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data deleted"))
}

func main() {
	connectDB()
	r := mux.NewRouter()

	r.HandleFunc("/sensor-data", sensorDataHandler).Methods("POST")
	r.HandleFunc("/sensor-data", getDataHandler).Methods("GET")
	r.HandleFunc("/sensor-data", deleteDataHandler).Methods("DELETE")

	fmt.Println("Microservice B running on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

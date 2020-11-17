package main

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

var dsn = "root:root@tcp(127.0.0.1:3306)/todolist?charset=utf8mb4&parseTime=True&loc=Local"
var db, _ = gorm.Open("mysql", dsn)



type TodoItemModel struct {
	Id          int `gorm:"primary_key"`
	Description string
	Completed   bool
}

func Healthz(writer http.ResponseWriter, req *http.Request) {
	log.Info("API health is OK")
	writer.Header().Set("Content-Type", "application/json")
	io.WriteString(writer, `{"alive": true}`)
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func main() {
	defer db.Close()

	db.Debug().DropTableIfExists(&TodoItemModel{})
	db.Debug().AutoMigrate(&TodoItemModel{})

	log.Info("Starting API service")
	router := mux.NewRouter()
	router.HandleFunc("/healthz", Healthz).Methods("GET")
	http.ListenAndServe(":8000", router)
}

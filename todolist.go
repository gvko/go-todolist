package main

import (
  "encoding/json"
  "github.com/gorilla/mux"
  log "github.com/sirupsen/logrus"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "io"
  "net/http"
  "time"
)

var conn, _ = mgo.Dial("localhost")
var collection = conn.DB("TutDb").C("ToDo")

type ToDoItem struct {
  Id          bson.ObjectId `bson:"_id,omitempty"`
  Description string
  Done        bool
  Date        time.Time
}

func init() {
  log.SetFormatter(&log.TextFormatter{})
  log.SetReportCaller(true)
}

/*
 * Function to check the health of the service
 */
func Healthz(writer http.ResponseWriter, req *http.Request) {
  log.Info("API health is OK")
  writer.Header().Set("Content-Type", "application/json")
  io.WriteString(writer, `{"alive": true}`)
}

/*
 * Add or update an item in the list
 */
func addItem(writer http.ResponseWriter, req *http.Request) {
  log.Info("Creating a new ToDo item")
  _ = collection.Insert(ToDoItem{
    bson.NewObjectId(),
    req.FormValue("description"),
    false,
    time.Now(),
  })

  result := ToDoItem{}
  _ = collection.Find(bson.M{"description": req.FormValue("description")}).One(&result)
  json.NewEncoder(writer).Encode(result)
}

/*
 * Get an item from the list by ID
 */
func getItemById(writer http.ResponseWriter, req *http.Request) {
  log.Info("Get an item by ID")
  var res ToDoItem

  vars := mux.Vars(req)
  id := vars["id"]

  _ = collection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&res)
  json.NewEncoder(writer).Encode(res)
}

/*
 * Get an item by ID or all items
 */
func getItem(writer http.ResponseWriter, req *http.Request) {
  log.Info("Getting all items")
  var res []ToDoItem

  _ = collection.Find(nil).All(&res)
  json.NewEncoder(writer).Encode(res)
}

func main() {
  conn.SetMode(mgo.Monotonic, true)
  defer conn.Close()

  log.Info("Starting API service at localhost:8000")
  router := mux.NewRouter()

  router.HandleFunc("/healthz", Healthz).Methods("GET")
  router.HandleFunc("/todo", addItem).Methods("POST", "PUT")
  router.HandleFunc("/todo", getItem).Methods("GET")
  router.HandleFunc("/todo/{id}", getItemById).Methods("GET")

  http.ListenAndServe(":8000", router)
}

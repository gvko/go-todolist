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
func Healthz(w http.ResponseWriter, req *http.Request) {
  log.Info("API health is OK")
  w.Header().Set("Content-Type", "application/json")
  io.WriteString(w, `{"alive": true}`)
}

/*
 * Add or update an item in the list
 */
func addItem(w http.ResponseWriter, req *http.Request) {
  log.Info("Creating a new ToDo item")
  _ = collection.Insert(ToDoItem{
    bson.NewObjectId(),
    req.FormValue("description"),
    false,
    time.Now(),
  })

  result := ToDoItem{}
  _ = collection.Find(bson.M{"description": req.FormValue("description")}).One(&result)
  json.NewEncoder(w).Encode(result)
}

/*
 * Get an item from the list by ID
 */
func getItemById(w http.ResponseWriter, req *http.Request) {
  log.Info("Get an item by ID")
  var res ToDoItem

  vars := mux.Vars(req)
  id := vars["id"]

  _ = collection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&res)
  json.NewEncoder(w).Encode(res)
}

/*
 * Get an item by ID or all items
 */
func getItem(w http.ResponseWriter, req *http.Request) {
  log.Info("Getting all items")
  var res []ToDoItem

  _ = collection.Find(nil).All(&res)
  json.NewEncoder(w).Encode(res)
}

/*
 * Patches an item by marking it as Done
 */
func markItemAsDone(w http.ResponseWriter, req *http.Request) {
  log.Info("Updating an item: mark as Done")
  vars := mux.Vars(req)
  id := bson.ObjectIdHex(vars["id"])
  err := collection.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"done": true}})

  if err != nil {
    log.Error("Error: could not update item")
    w.WriteHeader(http.StatusNotFound)
    w.Header().Set("Content-Type", "application/json")
    io.WriteString(w, `{"updated": false, "error": `+err.Error()+`}`)
  } else {
    log.Info("Item updated successfully")
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")
    io.WriteString(w, `{"updated": true}`)
  }
}

/*
 * Delete an item from the list
 */
func deleteItem(w http.ResponseWriter, req *http.Request) {
  log.Info("Deleting an item")
  vars := mux.Vars(req)
  id := vars["id"]
  err := collection.RemoveId(bson.ObjectIdHex(id))

  if err == mgo.ErrNotFound {
    log.Error("Error: could not delete item. Item not found.")
    json.NewEncoder(w).Encode(err.Error())
  } else {
    log.Info("Item deleted successfully")
    io.WriteString(w, "{result: 'OK'}")
  }
}

func main() {
  conn.SetMode(mgo.Monotonic, true)
  defer conn.Close()

  log.Info("Starting API service at localhost:8000")
  router := mux.NewRouter()

  router.HandleFunc("/healthz", Healthz).Methods("GET")
  router.HandleFunc("/todo", addItem).Methods("POST")
  router.HandleFunc("/todo", getItem).Methods("GET")
  router.HandleFunc("/todo/{id}", getItemById).Methods("GET")
  router.HandleFunc("/todo/{id}", markItemAsDone).Methods("PATCH")
  router.HandleFunc("/todo/{id}", deleteItem).Methods("DELETE")

  http.ListenAndServe(":8000", router)
}

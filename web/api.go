package web

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/just1689/pg-gateway/db"
	"github.com/just1689/pg-gateway/query"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

//LB & RB are used to apply them to an array
var LB = []byte("[")
var RB = []byte("]")
var COMMA = []byte(",")

/*
Main API class in which all handle commands are handled
*/

//writer sets the header for what is to be retutrned
func HandleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, PATCH, DELETE, POST, OPTIONS")
}

// checks field, and id values if exists, else throws the error. These are patch only, e.g. update, delete, etc.
func HandlePatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) //route variables for the current request, if any.
	entity := vars["entity"]
	if entity == "" {
		http.Error(w, "You need to supply an entity: /{entity}/{id}", http.StatusBadRequest)
		return
	}
	field := vars["field"]
	if field == "" {
		http.Error(w, "You need to supply an field", http.StatusBadRequest)
		return
	}
	id := vars["id"]
	if id == "" {
		http.Error(w, "You need to supply an id: /{entity}/{id}", http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(r.Body) //reads the body of the http request
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, "Could not read post body", http.StatusBadRequest)
		return
	}

	item := db.Insertable{} //if its insert table, returns raw msg
	err = json.Unmarshal(b, &item)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, "Could not unmarshal item from body", http.StatusBadRequest)
		return
	}

	err = db.Update(entity, field, id, item)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, "Could not update", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

//As above done for delete cmds with the clientReponseWriter and Reader
func HandleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entity := vars["entity"]
	if entity == "" {
		http.Error(w, "You need to supply an entity: /{entity}/{id}", http.StatusBadRequest)
		return
	}
	field := vars["field"]
	if field == "" {
		http.Error(w, "You need to supply an field", http.StatusBadRequest)
		return
	}
	id := vars["id"]
	if id == "" {
		http.Error(w, "You need to supply an id: /{entity}/{id}", http.StatusBadRequest)
		return
	}

	err := db.Delete(entity, field, id)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, "Could not delete", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

//same as handlerDelete cmds. Checks vars and performs the insert operation.
func HandleInsert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entity := vars["entity"]
	if entity == "" {
		http.Error(w, "You need to supply an entity: /{entity}", http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, "Could not read post body", http.StatusBadRequest)
		return
	}

	item := db.Insertable{}
	err = json.Unmarshal(b, &item)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, "Could not unmarshal item from body", http.StatusBadRequest)
		return
	}

	err = db.Insert(entity, item)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, "Could not insert", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func HandleGetMany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	entity := vars["entity"]
	if entity == "" {
		http.Error(w, "You need to supply an entity", http.StatusBadRequest)
		return
	}
	//TODO:check here
	q, err := query.BuildQueryFromURL(r.URL.String()[1:]) //parses the query string for getMany requests
	if err != nil {
		http.Error(w, "Could not build query from http request", http.StatusBadRequest)
		return
	}
	c, err := db.GetByQuery(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "Application/json")
	w.WriteHeader(http.StatusOK)
	rows := 0
	//creates the json array using left and right square brakcets
	w.Write(LB)
	for row := range c {
		rows++
		if rows > 1 {
			w.Write(COMMA)
		}
		w.Write(row)

	}
	w.Write(RB)

}

package app

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/6adfeniks/rest_api_with_db/pkg/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// App is a structure to keep router and database controllers.
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize sets router and connection to database,
// by accepting user, password and dbname as parameters.
func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password,
		dbname)

	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run runs the server with the defined address
func (a *App) Run(addr string) error {
	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		return errors.New("bad address")
	}
	return nil
	//log.Fatal(http.ListenAndServe(addr, a.Router))
}

// initializeRoutes registers new routes with a matcher for the URL path.
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
	a.Router.HandleFunc("/user", a.createUser).Methods("POST")
	a.Router.HandleFunc("/user/{id}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/user/{id}", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/user/{id}", a.deleteUser).Methods("DELETE")
}

// RespondWithError calls custom respond function
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON converts payload to the json format and returns respond to the server about errors
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// getUser returns user if such exists on the server
func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	u := models.User{ID: id}
	if err = u.GetUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			RespondWithError(w, http.StatusNotFound, "User not found")
		default:
			RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondWithJSON(w, http.StatusOK, u)
}

// getUsers returns users if such exists on the server
func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetUsers(a.DB)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, users)
}

// createUser creates new user on the server, returns user if it was created
func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := u.CreateUser(a.DB); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, u)
}

// updateUser updates data about some user on the server
func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var u models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()
	u.ID = id

	if err := u.UpdateUser(a.DB); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, u)
}

// deleteUser deletes user from the server
func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid User Id")
		return
	}

	u := models.User{ID: id}
	if err := u.DeleteUser(a.DB); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

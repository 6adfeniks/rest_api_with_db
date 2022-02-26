package app

import (
	"bytes"
	"fmt"
	"github.com/6adfeniks/rest_api_with_db/internal/config"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	cfg := getConfig()

	a.Initialize(cfg.Database.User, cfg.Database.Password, cfg.Database.Dbname)

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func getConfig() *config.Config {
	cfg, err := config.NewConfig("../../../configs/config2.yml")
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func dropTable() {
	if _, err := a.DB.Exec("DROP table if exists users"); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM users")
	a.DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS users
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    age INT NOT NULL
)`

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestInitialize(t *testing.T) {
	cfg := getConfig()

	b := App{}

	b.Initialize(cfg.Database.User, cfg.Database.Password, cfg.Database.Dbname)
	if a.DB == nil || a.Router == nil{
		t.Errorf("Expected router and database controller, got %v and %v", a.Router, a.DB)
	}
}

func TestRun(t *testing.T){
	b := App{}

	err := b.Run("wow")
	if err.Error() != "bad address"{
		t.Errorf("got: %v, want bad address", err.Error())
	}
}

func addUsers(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		statement := fmt.Sprintf("INSERT INTO users(name, age) VALUES('%s', %d)",
			("User " + strconv.Itoa(i+1)), ((i + 1) * 10))
		a.DB.Exec(statement)
	}
}

//func checkWork(t *testing.T, code int, method, url string) {
//	req, _ := http.NewRequest(method, url, nil)
//	response := executeRequest(req)
//	checkResponseCode(t, code, response.Code)
//}

func TestGetUser(t *testing.T) {
	clearTable()
	addUsers(1)

	req, _ := http.NewRequest("GET", "/user/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/user/w", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	req, _ = http.NewRequest("GET", "/user/23", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

	dropTable()
	req, _ = http.NewRequest("GET", "/user/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusInternalServerError, response.Code)
	ensureTableExists()

}

func TestGetUsers(t *testing.T) {
	clearTable()
	addUsers(2)
	req, _ := http.NewRequest("GET", "/users", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	dropTable()
	req, _ = http.NewRequest("GET", "/users", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusInternalServerError, response.Code)
	ensureTableExists()

}

func TestUpdateUser(t *testing.T) {
	clearTable()
	addUsers(1)

	req, _ := http.NewRequest("GET", "/user/1", nil)
	response := executeRequest(req)

	payload := []byte(`{"name":"test user - updated name","age":21}`)

	req, _ = http.NewRequest("PUT", "/user/1", bytes.NewBuffer(payload))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	payload = []byte(`name":"test user - updated name","age":21}`)

	req, _ = http.NewRequest("PUT", "/user/1", bytes.NewBuffer(payload))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	payload = []byte(`{"name":"test user - updated name","age":21}`)

	req, _ = http.NewRequest("PUT", "/user/w", bytes.NewBuffer(payload))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	payload = []byte(`{"name":"test user - updated name","age":21}`)
	dropTable()
	req, _ = http.NewRequest("PUT", "/user/1", bytes.NewBuffer(payload))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusInternalServerError, response.Code)
	ensureTableExists()

}

func TestCreateUser(t *testing.T) {
	clearTable()

	payload := []byte(`{"name":"test user","age":30}`)

	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	payload = []byte(`name":"test user","age":30}`)

	req, _ = http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	payload = []byte(`{"name":"test user","age":10}`)
	dropTable()
	req, _ = http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusInternalServerError, response.Code)
	ensureTableExists()
}

func TestDeleteUser(t *testing.T) {
	clearTable()
	addUsers(1)

	req, _ := http.NewRequest("DELETE", "/user/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/user/w", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	dropTable()
	req, _ = http.NewRequest("DELETE", "/user/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusInternalServerError, response.Code)
	ensureTableExists()

}


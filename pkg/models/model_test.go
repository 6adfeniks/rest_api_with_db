package models

import (
"database/sql"
"fmt"
"github.com/6adfeniks/rest_api_with_db/internal/config"
_ "github.com/go-sql-driver/mysql"
"log"
"reflect"
"testing"
)

var db *sql.DB

func start(){
	cfg, _ := config.NewConfig("../../configs/config2.yml")

	connectionString := fmt.Sprintf("%s:%s@/%s", cfg.Database.User,
		cfg.Database.Password, cfg.Database.Dbname)
	var err error
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	db.Exec("DELETE FROM users")
	db.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
}

func TestGetUser(t *testing.T) {

	start()
	clearTable()
	var send = &User{ID:1, Name: "Kolyan", Age: 21}
	send.CreateUser(db)

	var tdb = &User{ID: 1}

	tdb.GetUser(db)

	if !(tdb.ID == 1 && tdb.Name == "Kolyan" && tdb.Age == 21) {
		t.Errorf("bad scanning to object from db got %v", tdb)
	}
}

func TestUpdateUser(t *testing.T) {
	var updateData = &User{ID: 1, Name: "Kolyan", Age: 21}
	updateData.UpdateUser(db)

	if !(updateData.Name == "Kolyan" && updateData.Age == 21) {
		t.Errorf("want Kolyan and 21 got %v", updateData)
	}
}

func TestCreateUser(t *testing.T) {
	var tdb = &User{Name: "Anton", Age: 33}
	tdb.CreateUser(db)

	var want int
	db.QueryRow("select last_insert_id()").Scan(&want)

	if tdb.ID != want{
		t.Errorf("want %v got %v",want, tdb.ID)
	}
	tdb.DeleteUser(db)

	var tdb2 = &User{Name: "Anton2etttttttttttttttttttttttttttttttttttttttttttttttttукуаt", Age: 34}
	err := tdb2.CreateUser(db)

	if  err.Error() != "Error 1406: Data too long for column 'name' at row 1" {
		t.Errorf("got: %v, want: ", err.Error())
	}

	//var check bool
	//err = db.QueryRow("select last_insert_id()").Scan(&check)
	//if err.Error() != "Scan error on column index 0, name \"last_insert_id()\": sql/driver: couldn't convert \"39\" into type bool" {
	//	t.Errorf("got: %v, want: Scan error on column index 0, name \"last_insert_id()\": sql/driver: couldn't convert \"39\" into type bool ", err.Error())
	//}
}

func TestGetUsers(t *testing.T) {
	var usersTest = []User{{
		ID: 1,
		Name: "Kolyan",
		Age: 21,
	},
		{
			8,
			"Anton",
			33,
		},
		{
			10,
			"Anton",
			37,
		},
		{
			20,
			"Anton",
			34,
		},
		{
			22,
			"Anton",
			34,
		},
		{
			24,
			"Anton",
			34,
		},
	}

	users, _ := GetUsers(db)

	if reflect.DeepEqual(usersTest, users) {
		t.Errorf("got %v, want: %v", users, usersTest)
	}
	db.Close()
	_, err := GetUsers(db)
	if err.Error() != "sql: database is closed"{
		t.Errorf("got: %v, want: sql: database is closed", err.Error())
	}

}


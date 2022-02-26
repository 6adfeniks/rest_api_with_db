package models

import (
	"database/sql"
)

// User is a model of user in the database
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// GetUser sets parameters to the object of *User if such exists in the database
func (u *User) GetUser(db *sql.DB) error {
	//statement := fmt.Sprintf("SELECT name, age FROM users WHERE id=%d", u.ID)
	//return db.QueryRow(statement).Scan(&u.Name, &u.Age)

	statement := "SELECT name, age FROM users WHERE id = ?"
	return db.QueryRow(statement, u.ID).Scan(&u.Name, &u.Age)
}

// UpdateUser updates data about the user in the database
func (u *User) UpdateUser(db *sql.DB) error {
	//statement := fmt.Sprintf("UPDATE users SET name='%s', age=%d WHERE id=%d", u.Name, u.Age, u.ID)
	//_, err := db.Exec(statement)

	statement := "update users set name = ?, age = ? where id = ?"
	_, err := db.Exec(statement, u.Name, u.Age, u.ID)

	return err
}

// DeleteUser deletes the user from the database
func (u *User) DeleteUser(db *sql.DB) error {
	//statement := fmt.Sprintf("DELETE FROM users WHERE id=%d", u.ID)
	//_, err := db.Exec(statement)

	statement := "delete from users where id = ?"
	_, err := db.Exec(statement, u.ID)

	return err

}

// CreateUser creates new user in the database
func (u *User) CreateUser(db *sql.DB) error {
	//statement := fmt.Sprintf("INSERT INTO users(name, age) VALUES('%s', %d)", u.Name, u.Age)
	//_, err := db.Exec(statement)

	statement := "insert into users(name, age) values(?, ?)"
	_, err := db.Exec(statement, u.Name, u.Age)

	if err != nil {
		return err
	}

	//err = db.QueryRow("select last_insert_id()").Scan(&u.ID)
	//if err != nil {
	//	return err
	//}
	db.QueryRow("select last_insert_id()").Scan(&u.ID)

	return nil
}

// GetUsers returns the slice of the users from the database if such exists
func GetUsers(db *sql.DB) ([]User, error) {

	statement := "select id, name, age from users"
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []User{}

	//for rows.Next() {
	//	var u User
	//	if err = rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
	//		return nil, err
	//	}
	//	users = append(users, u)
	//}
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Age)
		users = append(users, u)
	}

	return users, nil
}

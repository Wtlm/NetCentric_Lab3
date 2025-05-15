package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	// Format: username:password@tcp(host:port)/database_name
	dsn := "root:123456@tcp(127.0.0.1:3306)/user_management"

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("MySQL not responding:", err)
	}

	createTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(100) NOT NULL,
        firstname VARCHAR(100) NOT NULL,
        lastname VARCHAR(100) NOT NULL,
        email VARCHAR(100) NOT NULL,
        avatar TEXT,
        phone VARCHAR(20),
        dob VARCHAR(50),
        country VARCHAR(100),
        city VARCHAR(100),
        street_name VARCHAR(100),
        street_address VARCHAR(200)
    );
    `
	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	fmt.Println("MySQL database connected and users table ready.")
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type UserData struct {
	ID         int    `json:"id"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Status     bool   `json:"status"`
	DateJoined string `json:"datejoined"`
}

// Connect to database
func GetDatabase() (*sql.DB, error) {
	connectionString := "sql5666109:WTIApagXJb@tcp(sql5.freemysqlhosting.net:3306)/sql5666109"

	// Open a database connection
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to the databases!")
	return db, nil
}

// Create User Data
func createUser(c *gin.Context) {
	var newUser UserData

	// Bind JSON data from the request body into newUser
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON data"})
		return
	}

	// Insert the new user into the database
	db, err := GetDatabase()
	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	defer db.Close()

	result, err := db.Exec("INSERT INTO users(username,email,password) VALUES (?,?,?)", newUser.Username, newUser.Email, newUser.Password)
	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{"error": "Failed to insert user into the database"})
		return
	}

	// Get the ID of the newly inserted user
	newUserID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{"error": "Failed to get the ID of the newly inserted user"})
		return
	}

	newUser.ID = int(newUserID)

	c.JSON(201, gin.H{"user": newUser})
}

// Get user date
func getUsersData(sql string) ([]UserData, error) {
	var usersData []UserData
	db, err := GetDatabase()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	// if err != nil {
	// 	log.Fatal(err)
	// 	return nil, err
	// } else {
	rows, err := db.Query(sql)

	if err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		defer rows.Close()

		for rows.Next() {
			var userData UserData
			err := rows.Scan(&userData.ID, &userData.Email)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("ID: %d, Email: %s\n", userData.ID, userData.Email)
			usersData = append(usersData, userData)
		}

		// Check for errors from iterating over rows
		if err = rows.Err(); err != nil {
			log.Fatal(err)
			return nil, err
		}

		return usersData, nil

	}
}

func defaultPage(c *gin.Context) {
	UserDatax, err := getUsersData("SELECT id, email FROM users")

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message":     "Welcome to Secure Api, You are Highly Celebrated!",
			"Error Occur": err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message":           "Welcome to Secure Api, You are Highly Celebrated!",
			"All working finex": UserDatax,
		})
	}
}

func main() {
	r := gin.New()

	//Home Page
	r.GET("/", defaultPage)

	//register users
	r.POST("/v1/reg/create", createUser)
	r.Run()
}

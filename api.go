package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
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
	WalletID   string `json:walletid`
	DateJoined string `json:"datejoined"`
}

type FaintData struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Status     bool   `json:"status"`
	WalletID   string `json:walletid`
	DateJoined string `json:"datejoined"`
}

type AddUser struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Wallet struct {
	ID        int     `json:"id"`
	W_account string  `json:"w_account"`
	Balance   float64 `json:"balance"`
}

// Function to generate a hashcode
func returnString128() (a string) {
	min := 97
	max := 122
	min1 := 65
	max1 := 90
	//a:=""
	for i := 0; i <= 128; i++ {
		choice := rand.Intn(2)
		if choice == 1 {
			a += string(rand.Intn(max-min) + min)
		} else {
			a += string(rand.Intn(max1-min1) + min1)
		}
	}
	return a
}

// Connect to database
func GetDatabase() (*sql.DB, error) {
	connectionString := "admin:beyourself@tcp(securemine.czuybcvcj4h5.us-west-2.rds.amazonaws.com:3306)/securemin"

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

	// fmt.Println("Connected to the databases!")
	return db, nil
}

// Create User Data

// Get user data
func getUsersData() ([]FaintData, error) {
	var faintData []FaintData
	db, err := GetDatabase()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id,username,status,walletid,datejoined FROM users")

	if err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		defer rows.Close()

		for rows.Next() {
			var faintOne FaintData
			err := rows.Scan(&faintOne.ID, &faintOne.Username, &faintOne.Status, &faintOne.WalletID, &faintOne.DateJoined)
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Printf("ID: %d, Email: %s\n", userData.ID, userData.Email)
			faintData = append(faintData, faintOne)
		}

		// Check for errors from iterating over rows
		if err = rows.Err(); err != nil {
			log.Fatal(err)
			return nil, err
		}

		return faintData, nil

	}
}

func getUserData(sql string) ([]UserData, error) {
	var usersData []UserData
	db, err := GetDatabase()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	rows, err := db.Query(sql)

	if err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		defer rows.Close()

		for rows.Next() {
			var userData UserData
			err := rows.Scan(&userData.Email)
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Printf("ID: %d, Email: %s\n", userData.ID, userData.Email)
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

func DefaultPage(c *gin.Context) {
	UserDatax, err := getUsersData()

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

// Function to create wallet
func CreateWallet(c *gin.Context) (id int) {
	//var wallet Wallet
	db, err := GetDatabase()

	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{"Error": "Internal Server Error"})
		return 0
	}

	defer db.Close()
	walletid := returnString128()
	result, err := db.Exec("INSERT INTO wallet(wallet,balance) VALUE(?,?)", walletid, 0.00000002)

	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{"Error": "Unable to create wallet, contact admin"})
		return 0
	}

	walletIds, err := result.LastInsertId()

	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{"Error": "Unable to retrieve Wallet ID"})
		return 0
	}

	walletID := int(walletIds)

	return walletID

}

// Create User Account
func CreateUser(c *gin.Context) {
	var newUser AddUser

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

	condition := "email = '" + newUser.Email + "'"
	userExists, err := getUserData(fmt.Sprintf("SELECT email FROM users WHERE %s", condition))

	if userExists != nil {
		c.JSON(500, gin.H{"Error": "User Record Exists"})
		return
	} else {
		WalletId := CreateWallet(c)

		if WalletId == 0 {
			c.JSON(500, gin.H{"Error": "Internal Server Error"})
			return
		}

		result, err := db.Exec("INSERT INTO users(username,email,walletid,password) VALUES (?,?,?,?)", newUser.Username, newUser.Email, WalletId, newUser.Password)
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

}

func main() {
	r := gin.New()

	//Home Page
	r.GET("/", DefaultPage)

	//register users
	r.POST("/v1/reg/create", CreateUser)
	r.Run(":8000")
}

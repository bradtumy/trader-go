package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

// User from user table in trader-go db
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// Stock from stock table in trader-go db
type Stock struct {
	ID          int     `json:"id"`
	Symbol      string  `json:"symbol"`
	Name        string  `json:"name"`
	Price       float32 `json:"price"`
	TotalShares int     `json:"total_shares"`
}

// Orders from orders table in trader-go db
type Orders struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Symbol   string `json:"symbol"`
	Shares   int    `json:"shares"`
}

// Config struct to hold properties imported from a file at runtime
type Config struct {
	Server struct {
		Port        string `yaml:"port"`
		Host        string `yaml:"host"`
		Certificate string `yaml:"server_cert"`
		Key         string `yaml:"server_key"`
	} `yaml:"server"`
	Database struct {
		Name     string `yaml:"name"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"user"`
		Password string `yaml:"pass"`
	} `yaml:"database"`
}

// db connection global variables
var db *sql.DB

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, "Welcome to the create User!")
	fmt.Println("Endpoint Hit: create user")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	// grab the username from the request
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	username := keyVal["username"]
	fmt.Println("New username is: ", username)

	// prepare insert statement
	_, err = db.Exec("INSERT INTO users (username) VALUES (?)", username)
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "New user was created")
}

func returnAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, "Welcome to the returnAllUsers!")
	fmt.Println("Endpoint Hit: returnAllUsers")

	var users []User

	results, err := db.Query("Select * from users")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var user User
		err = results.Scan(&user.ID, &user.Username)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
}

func returnAllStocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, "Welcome to the returnAllStocks!")
	fmt.Println("Endpoint Hit: returnAllStocks")

	var stocks []Stock

	results, err := db.Query("Select * from stocks")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var stock Stock
		err = results.Scan(&stock.ID, &stock.Symbol, &stock.Name, &stock.Price, &stock.TotalShares)
		if err != nil {
			panic(err.Error())
		}
		stocks = append(stocks, stock)
	}
	json.NewEncoder(w).Encode(stocks)
}

func returnAllAccounts(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the returnAllAccounts!")
	fmt.Println("Endpoint Hit: returnAllAccounts")
}

func returnAllTrades(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the returnAllTrades!")
	fmt.Println("Endpoint Hit: returnAllTrades")
}

func returnAllTransfers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the returnAllTransfers!")
	fmt.Println("Endpoint Hit: returnAllTransfers")
}

// Return all orders in the DB
func returnAllOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, "Welcome to the returnAllOrders!")
	fmt.Println("Endpoint Hit: returnAllOrders")

	var orders []Orders

	results, err := db.Query("select orders.id, users.username, stocks.symbol, shares from orders inner join users on orders.user_id = users.id inner join stocks on orders.stock_id = stocks.id")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var order Orders
		err = results.Scan(&order.ID, &order.Username, &order.Symbol, &order.Shares)
		if err != nil {
			panic(err.Error())
		}
		orders = append(orders, order)
	}
	json.NewEncoder(w).Encode(orders)
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, "Welcome to the Create Orders!")
	fmt.Println("Endpoint Hit: CreateOrders")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	// grab the attributes from the request
	keyVal := make(map[string]int)
	json.Unmarshal(body, &keyVal)
	userid := keyVal["userid"]
	stockid := keyVal["stockid"]
	shares := keyVal["shares"]

	// prepare insert statement
	_, err = db.Exec("INSERT INTO orders (user_id, stock_id, shares) VALUES (?,?,?)", userid, stockid, shares)
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "New order was created")
}

func createStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, "Welcome to the Create Stock!")
	fmt.Println("Endpoint Hit: CreateStock")

	var stock Stock

	// grab the values from the request body
	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		fmt.Println("oh no, that didn't work", err)
		//http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// grab the attributes from the request
	symbol := stock.Symbol
	name := stock.Name
	price := stock.Price
	totalshares := stock.TotalShares

	// prepare insert statement
	_, err = db.Exec("INSERT INTO stocks (symbol, name, price, total_shares) VALUES (?,?,?,?)", symbol, name, price, totalshares)
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "New stock was created")
}

// User must be specified in request and then return users' portfolio (stocks owned and values)
func returnUsersAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []User
	var orders []Orders
	//fmt.Println("Request URI: ", path.Base(r.RequestURI))
	username := path.Base(r.RequestURI)

	// query for userid
	results, err := db.Query("Select * from users where username = ?", username)
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var user User
		err := results.Scan(&user.ID, &user.Username)
		if err != nil {
			panic(err.Error())
		}

		users = append(users, user)
		fmt.Println("Userid: ", users[0].ID)
	}

	// query the orders table with userid and return all rows

	//select orders.id, users.username, stocks.symbol, shares from orders inner join users on orders.user_id = users.id inner join stocks on orders.stock_id = stocks.id
	results, err = db.Query("select orders.id, users.username, stocks.symbol, shares from orders inner join users on orders.user_id = users.id inner join stocks on orders.stock_id = stocks.id where orders.user_id = ?", users[0].ID)

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var order Orders
		err = results.Scan(&order.ID, &order.Username, &order.Symbol, &order.Shares)
		if err != nil {
			panic(err.Error())
		}
		orders = append(orders, order)
	}
	json.NewEncoder(w).Encode(orders)
}
func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(cfg *Config) {
	f, err := os.Open("resources/properties.yml")
	if err != nil {
		fmt.Println("ERROR: I couldn't read the properties file.")
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		fmt.Println("ERROR: I couldn't decode the YAML.")
		processError(err)
	}
}

func handleRequests(cfg Config) {

	router := mux.NewRouter()

	router.HandleFunc("/", homePage)
	router.HandleFunc("/users", returnAllUsers).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/orders", returnAllOrders).Methods("GET")
	router.HandleFunc("/orders", createOrder).Methods("POST")
	router.HandleFunc("/stocks", returnAllStocks).Methods("GET")
	router.HandleFunc("/stocks", createStock).Methods("POST")
	router.HandleFunc("/accounts", returnAllAccounts).Methods("GET")
	router.HandleFunc("/accounts/{i}", returnUsersAccounts).Methods("GET")
	router.HandleFunc("/trade", returnAllTrades)
	router.HandleFunc("/transfer", returnAllTransfers)

	// start the https listener using the signed cert and key
	log.Fatal(http.ListenAndServeTLS(cfg.Server.Port, cfg.Server.Certificate, cfg.Server.Key, router))
}

func main() {

	var cfg Config
	readFile(&cfg)
	var err error

	// set connection variables with values read from the properties file
	port := cfg.Server.Port
	dbConnectString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	fmt.Println("Trader-Go is starting on port: ", port)

	db, err = sql.Open("mysql", dbConnectString)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("MySQL DB Connection Established")
	}
	handleRequests(cfg)

	// close the db when we are done
	defer db.Close()
}

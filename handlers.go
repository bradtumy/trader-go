package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

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
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	username := keyVal["username"]
	fmt.Println("New username is: ", username)
	_, err = db.Exec("INSERT INTO users (username) VALUES (?)", username)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New user was created")
}

func returnAllUsers(w http.ResponseWriter, r *http.Request) {

	if db == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

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

func getOrderHistory(w http.ResponseWriter, r *http.Request) {
	// Authenticate user and get user_id from session or request
	userID := 1 // Placeholder; replace with actual user authentication and retrieval

	rows, err := db.Query("SELECT id, stock_id, shares FROM orders WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []map[string]interface{}
	for rows.Next() {
		var id, stockID, shares int
		if err := rows.Scan(&id, &stockID, &shares); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		order := map[string]interface{}{
			"id":       id,
			"stock_id": stockID,
			"shares":   shares,
		}
		orders = append(orders, order)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	// Replace with your database connection logic
	db, err := sql.Open("mysql", "tradergo:password1@tcp(db:3306)/trader_go")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var username string
	err = db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	userProfile := map[string]string{"id": userID, "username": username}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProfile)
}

func searchStocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryParams := r.URL.Query()
	symbol := queryParams.Get("symbol")
	name := queryParams.Get("name")

	var stocks []Stock
	var rows *sql.Rows
	var err error

	if symbol != "" && name != "" {
		rows, err = db.Query("SELECT * FROM stocks WHERE symbol LIKE ? OR name LIKE ?", "%"+symbol+"%", "%"+name+"%")
	} else if symbol != "" {
		rows, err = db.Query("SELECT * FROM stocks WHERE symbol LIKE ?", "%"+symbol+"%")
	} else if name != "" {
		rows, err = db.Query("SELECT * FROM stocks WHERE name LIKE ?", "%"+name+"%")
	} else {
		http.Error(w, "Please provide either a symbol or name query parameter", http.StatusBadRequest)
		return
	}

	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var stock Stock
		err := rows.Scan(&stock.ID, &stock.Symbol, &stock.Name, &stock.Price, &stock.TotalShares)
		if err != nil {
			panic(err.Error())
		}
		stocks = append(stocks, stock)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stocks)
}

func getKeyMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Example metrics: total number of users, total number of stocks, total number of orders
	var metrics struct {
		TotalUsers  int `json:"total_users"`
		TotalStocks int `json:"total_stocks"`
		TotalOrders int `json:"total_orders"`
	}

	// Get total number of users
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&metrics.TotalUsers)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get total number of stocks
	err = db.QueryRow("SELECT COUNT(*) FROM stocks").Scan(&metrics.TotalStocks)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get total number of orders
	err = db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&metrics.TotalOrders)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(metrics)
}

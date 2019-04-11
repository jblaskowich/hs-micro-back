package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	nats "github.com/nats-io/go-nats"

	_ "github.com/go-sql-driver/mysql"
)

var (
	natsURL  = "demo.nats.io"
	natsPORT = ":4222"
	database *sql.DB
)

type Message struct {
	Title   string
	Content string
	Date    string
}

func getEvents() {
	natsChan := os.Getenv("NATSCHAN")
	if natsChan == "" {
		natsChan = "zjnO12CgNkHD0IsuGd89zA"
	}

	nc, err := nats.Connect(natsURL + natsPORT)
	if err != nil {
		log.Println(err.Error())
	}
	nc.Subscribe(natsChan, func(m *nats.Msg) {
		//fmt.Printf("Received message: %s\n", string(m.Data))
		msg := Message{}
		err := json.Unmarshal(m.Data, &msg)
		if err != nil {
			log.Println(err.Error())
		}
		t := time.Now()
		msg.Date = t.Format("2006-01-02 15:04:05")
		save2Database(msg)
		//log.Printf("Received message, Title: %s, Content: %s, Date: %s\n", msg.Title, msg.Content, msg.Date)
	})
}

func save2Database(m Message) {
	/*
		Prepare your Database by creating the following TABLE:

		CREATE TABLE `post`(
		`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
		`post_title` varchar(64) DEFAULT NULL,
		`post_content` mediumtext,
		`post_date` timestamp NULL DEFAULT NULL,
		PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=latin1;
	*/
	_, err := database.Exec("INSERT INTO post SET post_title=?, post_content=?, post_date=?", m.Title, m.Content, m.Date)
	if err != nil {
		log.Println(err.Error())
	}
}

func main() {
	dbUser := os.Getenv("DBUSER")
	if dbUser == "" {
		dbUser = "user"
	}
	dbPass := os.Getenv("DBPASS")
	if dbPass == "" {
		dbPass = "password"
	}
	dbHost := os.Getenv("DBHOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	dbPort := os.Getenv("DBBASE")
	if dbPort == "" {
		dbPort = ":3306"
	}
	dbBase := os.Getenv("DBBASE")
	if dbBase == "" {
		dbBase = "blowofmouth"
	}
	dbConn := fmt.Sprintf("%s:%s@tcp(%s%s)/%s", dbUser, dbPass, dbHost, dbPort, dbBase)
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Printf("Couldn't connect to %s%s/%s\n"+dbHost, dbPort, dbBase)
		log.Println(err.Error())
	}
	database = db

	getEvents()
	port := os.Getenv("HS-MICRO-BACK")
	if port == "" {
		port = ":9090"
	}
	rtr := mux.NewRouter()
	http.Handle("/", rtr)
	http.ListenAndServe(port, nil)
}

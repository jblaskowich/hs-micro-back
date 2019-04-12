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
	natsURL  = "demo.nats.io"           // Can be superseded by env NATSURL
	natsPort = ":4222"                  // Can be superseded by env NATSPORT
	natsPost = "zjnO12CgNkHD0IsuGd89zA" // Can be superseded by env NATSPOST
	natsGet  = "OWM7pKQNbXd7l75l21kOzA" // Can be superseded by env NATSGET
	database *sql.DB
)

// Message is the representation of a post
type Message struct {
	ID      string
	Title   string
	Content string
	Date    string
}

func watchPost() {
	if os.Getenv("NATSURL") != "" {
		natsURL = os.Getenv("NATSURL")
	}
	if os.Getenv("NATSPORT") != "" {
		natsPort = os.Getenv("NATSPORT")
	}
	if os.Getenv("NATSCHAN") != "" {
		natsPost = os.Getenv("NATSPOST")
	}

	nc, err := nats.Connect("nats://" + natsURL + natsPort)
	if err != nil {
		log.Println(err.Error())
	}
	nc.Subscribe(natsPost, func(m *nats.Msg) {
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

func reqReply() {
	if os.Getenv("NATSURL") != "" {
		natsURL = os.Getenv("NATSURL")
	}
	if os.Getenv("NATSPORT") != "" {
		natsPort = os.Getenv("NATSPORT")
	}
	if os.Getenv("NATSCHAN") != "" {
		natsGet = os.Getenv("NATSGET")
	}

	nc, err := nats.Connect("nats://" + natsURL + natsPort)
	if err != nil {
		log.Println(err.Error())
	}

	/*
		We subscribe to Chan natsGet and we are waiting for a new message
		 the message contains the reference of an inbox, in which we send our content
	*/
	nc.Subscribe(natsGet, func(m *nats.Msg) {
		log.Println("Repl sent on " + m.Reply)
		err := nc.Publish(string(m.Reply), []byte(`[{"ID": 1, "Title": "hello world", "Content": "blablabla"}]`))
		if err != nil {
			log.Println(err.Error())
		}
		nc.Flush()
	})
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

	watchPost()
	reqReply()
	port := os.Getenv("HS-MICRO-BACK")
	if port == "" {
		port = ":9090"
	}
	rtr := mux.NewRouter()
	http.Handle("/", rtr)
	http.ListenAndServe(port, nil)
}

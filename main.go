package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	nats "github.com/nats-io/go-nats"

	_ "github.com/go-sql-driver/mysql"
)

var (
	// NATS
	natsURL  = "demo.nats.io"           // Can be superseded by env NATSURL
	natsPort = ":4222"                  // Can be superseded by env NATSPORT
	natsPost = "zjnO12CgNkHD0IsuGd89zA" // Can be superseded by env NATSPOST
	natsGet  = "OWM7pKQNbXd7l75l21kOzA" // Can be superseded by env NATSGET
	// SQL
	dbUser   = "user"        // Can be superseded by env DBUSER
	dbPass   = "password"    // Can be superseded by env DBPASS
	dbHost   = "127.0.0.1"   // Can be superseded by env DBHOST
	dbPort   = ":3306"       // Can be superseded by env DBPORT
	dbBase   = "blowofmouth" // Can be superseded by env DBBASE
	database *sql.DB
)

// Message is the representation of a post
type Message struct {
	ID      string
	Title   string
	Content string
	Date    string
}

// watchPost capture new posts sent through NATS
func watchPost(url, port, subj string) {

	nc, err := nats.Connect("nats://" + url + port)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Printf("Connect to nats://%s%s on %s\n", url, port, subj)
	}
	nc.Subscribe(subj, func(m *nats.Msg) {
		log.Println("new reccord sent to the database")
		msg := Message{}
		err := json.Unmarshal(m.Data, &msg)
		if err != nil {
			log.Println(err.Error())
		}
		t := time.Now()
		msg.Date = t.Format("2006-01-02 15:04:05")
		err = save2DatabaseSQL(msg)
		if err != nil {
			log.Println(err.Error())
		}
	})
}

// save2DatabaseSQL reccord posts in SQL database (MySQL, MariaDB)
func save2DatabaseSQL(m Message) error {
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
		return err
	}

	return nil
}

// reqReply waits for post request, and return database rows
func reqReply(url, port, subj string) {

	nc, err := nats.Connect("nats://" + url + port)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		log.Printf("Connect to nats://%s%s on %s\n", url, port, subj)
	}

	/*
		We subscribe to Chan natsGet and we are waiting for a new message
		 the message contains the reference of an inbox, in which we send our content
	*/
	nc.Subscribe(subj, func(m *nats.Msg) {
		log.Println("Repl sent on " + m.Reply)
		err := nc.Publish(string(m.Reply), selectPosts())
		//err := nc.Publish(string(m.Reply), []byte(`[{"ID": 1, "Title": "hello world", "Content": "blablabla"}]`))
		if err != nil {
			log.Println(err.Error())
		}
		nc.Flush()
	})
}

// selectPosts is responsible of the queries (SELECT) in the database, and return an array of rows
func selectPosts() []byte {
	messages := []Message{}
	msg, err := database.Query("SELECT id,post_title,post_content,post_date FROM post ORDER BY id DESC;")
	if err != nil {
		log.Println(err.Error())
	}
	defer msg.Close()
	for msg.Next() {
		thisMsg := Message{}
		msg.Scan(&thisMsg.ID, &thisMsg.Title, &thisMsg.Content, &thisMsg.Date)
		messages = append(messages, thisMsg)
	}
	posts, _ := json.Marshal(messages)
	return posts
}

func blockForever() {
	select {}
}

func main() {
	// Database
	if os.Getenv("DBUSER") != "" {
		dbUser = os.Getenv("DBUSER")
	}
	if os.Getenv("DBPASS") != "" {
		dbPass = os.Getenv("DBPASS")
	}
	if os.Getenv("DBHOST") != "" {
		dbHost = os.Getenv("DBHOST")
	}
	if os.Getenv("DBPORT") != "" {
		dbPort = os.Getenv("DBPORT")
	}
	if os.Getenv("DBBASE") != "" {
		dbBase = os.Getenv("DBBASE")
	}

	// NATS
	if os.Getenv("NATSURL") != "" {
		natsURL = os.Getenv("NATSURL")
	}
	if os.Getenv("NATSPORT") != "" {
		natsPort = os.Getenv("NATSPORT")
	}
	if os.Getenv("NATSPOST") != "" {
		natsPost = os.Getenv("NATSPOST")
	}
	if os.Getenv("NATSGET") != "" {
		natsGet = os.Getenv("NATSGET")
	}

	dbConn := fmt.Sprintf("%s:%s@tcp(%s%s)/%s", dbUser, dbPass, dbHost, dbPort, dbBase)
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Printf("Couldn't connect to %s%s/%s\n"+dbHost, dbPort, dbBase)
		log.Fatalln(err.Error())
	} else {
		log.Printf("Connect to %s%s/%s\n", dbHost, dbPort, dbBase)
	}
	database = db

	watchPost(natsURL, natsPort, natsPost)
	reqReply(natsURL, natsPort, natsGet)
	log.Println("service started, waiting for events...")

	blockForever()
}

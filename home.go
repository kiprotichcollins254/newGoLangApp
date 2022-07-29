package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

const(
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "postgres"
	dbname ="fam_survey"
)

type homePageData struct{
	PageTitle string
	PageHeading string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

	db, err := sql.Open("postgres",psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	err = db.Ping()
	if err != nil{
		log.Fatal(err)
	}
	fs:= http.FileServer(http.Dir("static/"))
	http.Handle("/static/",http.StripPrefix("/static/",fs))
    r := mux.NewRouter()
	r.HandleFunc("/books/{title}/page/{page}",func (w http.ResponseWriter,r * http.Request)  {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w,"You've requested the book : %s on page %s . \n" ,title,page)
	})

	r.HandleFunc("/home",func(w http.ResponseWriter, r *http.Request) {

		tmpl := template.Must(template.ParseFiles("static/html/layouts/header.html"))
        data := homePageData{
			PageTitle: "Go App | Home Page",
			PageHeading: "Go To Do Application",
		}
		tmpl.Execute(w,data)
	}).Methods("GET")

	r.HandleFunc("/chat",func(w http.ResponseWriter, r *http.Request) {
		conn,_ := upgrader.Upgrade(w,r,nil)  //error ignore for simplicity

		for {

			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if err = conn.WriteMessage(msgType,msg);err !=nil {
				return
			}
		}
	})

	r.HandleFunc("/chatRoom",func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("static/html/views/websockets.html"))
        data := homePageData{
			PageTitle: "Go App | Home Page",
			PageHeading: "Go To Do Application",
		}
		tmpl.Execute(w,data)
	})

    //  restricting http methods
	// r.HandleFunc("/books/{title}", CreateBook).Methods("POST")
	// r.HandleFunc("/books/{title}", ReadBook).Methods("GET")
	// r.HandleFunc("/books/{title}", UpdateBook).Methods("PUT")
	// r.HandleFunc("/books/{title}", DeleteBook).Methods("DELETE")

	// writing preixes
    // bookrouter := r.PathPrefix("/books").Subrouter()
	// bookrouter.HandleFunc("/", AllBooks)
	// bookrouter.HandleFunc("/{title}", GetBook)

//pg connection 
	r.HandleFunc("/user/{id}",func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"] 
		var (
			username string
			email string
			status bool
		)
		query := "SELECT username,email ,status from users.admins WHERE pk = $1"
		if err := db.QueryRow(query,id).Scan(&username,&email,&status);err != nil {
			log.Fatal(err)
		}

		// fmt.Println("Successfully Connected")
		fmt.Fprintf(w,"User %s email is %s and status is %t \n",username,email,status)
		// fmt.Println(username,email,status)
	}).Methods("POST")

    http.ListenAndServe(":7003", r)
}
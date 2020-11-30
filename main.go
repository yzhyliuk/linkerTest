package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	_ "github.com/lib/pq"
)

//seting connection string to our DB as constant
const connString = "host=localhost port=5432 user=postgres password=080919 dbname=links sslmode=disable"

func main() {
	//adding one general handler func, that will halndle all the requests
	http.HandleFunc("/", GeneralHandler)
	//adding serving static files
	http.Handle("/public", http.FileServer(http.Dir("/public")))
	//starting server
	http.ListenAndServe(":8080", nil)
}

// GeneralHandler : our main handler function. If it recives an GET request - it's cheking http request url and returning home.html
// or redirects to link, if it exists in our DB
func GeneralHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if r.URL.Path == "/" {
			t, _ := template.ParseFiles("public/home.html")
			t.Execute(w, nil)
			return
		}
		targetPage := getLink(r.URL.Path)
		if targetPage != "" {
			http.Redirect(w, r, targetPage, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	if r.Method == "POST" && r.URL.Path == "/addlink" {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		linkURL := buf.String()
		_, err := url.ParseRequestURI(linkURL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid URL format"))
			return
		}
		newURL := addlink(linkURL)
		const domain = "http://localhost:8080/"
		response := domain + newURL
		fmt.Println(response)
		w.Write([]byte(response))
	}
}
func addlink(url string) string {
	db, err := sql.Open("postgres", connString)
	checkErr(err)
	defer db.Close()

	var newRandomUrl string
	for fine := false; fine == false; {
		newRandomUrl = getRandomUrl()
		_, err := db.Exec("INSERT INTO links VALUES($1,$2)", url, newRandomUrl)
		if err == nil {
			fine = true
		}
	}
	return newRandomUrl
}
func getRandomUrl() string {
	const charArray = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	rand.Seed(time.Now().UnixNano())
	newURL := make([]byte, 5)
	for i := range newURL {
		newURL[i] = charArray[rand.Intn(len(charArray))]
	}
	return string(newURL)
}
func getLink(url string) string {
	db, err := sql.Open("postgres", connString)
	checkErr(err)
	defer db.Close()

	url = strings.Trim(url, "/")
	var targetPage string
	err = db.QueryRow("SELECT url FROM links WHERE newurl	 = $1", url).Scan(&targetPage)
	return targetPage
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

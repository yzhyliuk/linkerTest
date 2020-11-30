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
const connString = "host=ec2-54-246-85-151.eu-west-1.compute.amazonaws.com port=5432 user=vkhigwjgxhwbae password=8a7d1118066f00102847c6c5f63592e5d6473379cedbf2d8f02a55b62574a03a dbname=dbdj0dvlvt57dh sslmode=disable"

func initDB()  {
	db, err := sql.Open("postgres",connString)
	checkErr(err)
	defer db.Close()

	_, err = db.Exec("CREATE TABLE links (url VARCHAR(1000), newurl VARCHAR(12) PRIMARY KEY)")
	checkErr(err)
}

func main() {
	//adding one general handler func, that will halndle all the requests
	http.HandleFunc("/", GeneralHandler)
	//adding serving static files
	http.Handle("/public", http.FileServer(http.Dir("/public")))
	//starting server
	http.ListenAndServe(os.Getenv("PORT"), nil)
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
		const domain = "https://link-reducer-test.herokuapp.com/"
		response := domain + newURL
		fmt.Println(response)
		w.Write([]byte(response))
	}
}
func addlink(url string) string {
	db, err := sql.Open("postgres", connString)
	checkErr(err)
	defer db.Close()

	var newRandomURL string
	for fine := false; fine == false; {
		newRandomURL = getRandomURL()
		_, err := db.Exec("INSERT INTO links VALUES($1,$2)", url, newRandomURL)
		if err == nil {
			fine = true
		}
	}
	return newRandomURL
}
func getRandomURL() string {
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

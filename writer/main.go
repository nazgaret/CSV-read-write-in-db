package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

type (
	//PersonStruct for person DB save
	PersonStruct = struct {
		ID           string ` bson:"_id"`
		Name         string ` bson:"name"`
		Email        string ` bson:"email"`
		MobileNumber string ` bson:"mobile_number"`
	}
)

var (
	session *mgo.Session
	err     error
	dbhost  = "localhost"
	r       = strings.NewReplacer("(", "",
		")", "",
		" ", "")
)

func init() {

	session, err = mgo.Dial(dbhost)
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	go func() {
		for range time.Tick(time.Second * 1) {
			сonnectionCheck()
		}
	}()
}

func main() {
	http.HandleFunc("/", routerHandler)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//routerHandler handle a rout
func routerHandler(w http.ResponseWriter, r *http.Request) {
	personSlice := []string{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error ioutil.ReadAll : ", err)
	}
	defer r.Body.Close()
	err = json.Unmarshal(body, &personSlice)
	if err != nil {
		fmt.Println("error marshal body : ", err, string(body))
		return
	}

	saver(personSlice)

	w.Write([]byte("Thank you!"))
}

//saver write to DB
func saver(personSlice []string) {
	person := PersonStruct{
		ID:           personSlice[0],
		Name:         personSlice[1],
		Email:        personSlice[2],
		MobileNumber: personSlice[3],
	}
	person.MobileNumber = "+044" + r.Replace(person.MobileNumber)
	_, err := DB("data").C("persons").UpsertId(person.ID, person)
	if err != nil {
		fmt.Println("ERROR save to DB", err)
	}
}

// DB wrapper for mgo.Session.DB
func DB(dname string) *mgo.Database {
	сonnectionCheck()
	return session.DB(dname)
}

// сonnectionCheck reconect on lose conect
func сonnectionCheck() {
	if err := session.Ping(); err != nil {
		fmt.Println("Lost connection to db!")
		session.Refresh()
		if err := session.Ping(); err == nil {
			fmt.Println("Reconnect to db successful.")
		}
	}
}

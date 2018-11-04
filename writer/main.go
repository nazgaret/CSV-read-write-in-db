package main

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

type (
	personStruct = struct {
		ID           string ` bson:"id"`
		Name         string ` bson:"name"`
		Email        string ` bson:"email"`
		MobileNumber string ` bson:"mobile_number"`
	}
)

const (
	firstFormat = iota + 1
	secondFormat
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
	save := make(chan personStruct)
	for {
		if len(save) > 0 {
			saver(<-save)
		}
	}
}

//saver write to DB
func saver(person personStruct) {
	person.MobileNumber = "+044" + r.Replace(person.MobileNumber)
	_, err := DB("data").C("persons").Upsert(person.ID, person)
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

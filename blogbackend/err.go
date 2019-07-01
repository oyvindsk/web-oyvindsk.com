package main

import (
	"log"
	"net/http"
)

// Error handling functions

func checkAndDie(m string, e error) {
	if e != nil {
		log.Fatal("!! ", m, " : ", e)
	}
}

func checkAndWarn(m string, e error) {
	if e != nil {
		log.Print("! ", m, " : ", e)
	}
}

func checkErr(m string, e error) {
	if e != nil {
		log.Fatal("!! ", m, " : ", e)
	}
}

func checkErrHttp(err error, w http.ResponseWriter) bool {
	if err != nil {
		log.Println("!! ", err)
		http.Error(w, http.StatusText(500), 500)
		return true
	}
	return false
}

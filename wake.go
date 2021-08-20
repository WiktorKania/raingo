package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func WakeMeUpBeforeYouGoGo(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("Woken up on Heroku (probably)")
	rw.Write([]byte("good morning"))
}

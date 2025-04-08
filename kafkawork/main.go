package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello World")
}

func main() {
	router := httprouter.New()
	router.GET("/login", Login)
}

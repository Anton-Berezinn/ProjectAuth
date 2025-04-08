package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func GET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := r.FormValue("user")
	fmt.Println(user)
}

func main() {
	router := httprouter.New()
	router.GET("/", GET)
	if err := http.ListenAndServe("localhost:8080", router); err != nil {
		fmt.Println(err)
		return
	}
}

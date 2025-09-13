package serveProm

import (
	"fmt"
	"net/http"
)

func NewHandler() {
	router := http.NewServeMux()
	router.HandleFunc("/main", AuthMain)
}

type User struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func AuthMain(w http.ResponseWriter, r *http.Request) {
	var user User
	v := r.Body
	fmt.Println(v, user)
	w.Write([]byte("Hello World"))
}

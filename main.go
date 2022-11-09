package main

import (
	"C"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var userData = make(map[string]user)

type user struct {
	Name  string `json:"name"`
	Hobby string `json:"hobby"`
	Age   int8   `json:"age"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "error: %s", err)
		}

		w.Write(body)

		switch r.URL.Query().Get("param") {
		case "world":
			fmt.Fprint(w, "hello world")
		case "cat":
			fmt.Fprint(w, "hello cat")
		default:
			fmt.Fprint(w, "hello everybody")

		}
	})

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "error: %s", err)
		}

		switch r.Method {
		case http.MethodPost:
			oneUser := user{}
			err := json.Unmarshal(body, &oneUser)
			if err != nil {
				fmt.Fprintf(w, "error: %s", err)
				return
			}
			userData[oneUser.Name] = oneUser
			fmt.Fprintf(w, "success, total users: %d", len(userData))
		case http.MethodGet:
			userName := r.URL.Query().Get("user")
			u, ok := userData[userName]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "user %s not found", userName)
				return
			}
			data, err := json.Marshal(u)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "error: %s", err)
			}
			w.Header().Set("Content-type", "application/json")
			w.Write(data)
		}

	})

	http.HandleFunc("/user/changehobby/", func(w http.ResponseWriter, r *http.Request) {
		userName := r.URL.Query().Get("user")
		userHobby := r.URL.Query().Get("hobby")
		u, ok := userData[userName]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "user %s not found", userName)
			return
		}
		u.Hobby = userHobby
		userData[userName] = u
		fmt.Fprintf(w, "success, total users: %d", len(userData))
	})

	fmt.Printf("%v", os.Args[1:])
	port := os.Args[1]
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

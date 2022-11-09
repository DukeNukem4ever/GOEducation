package main

import (
	"C"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type user struct {
	Name  string `json:"name"`
	Hobby string `json:"hobby"`
	Age   int8   `json:"age"`
}

func main() {
	var storage UserStorage
	storage = NewStorage()
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
			storage.Set(oneUser)
			fmt.Fprintf(w, "success, total users: %d", storage.Size())
		case http.MethodGet:
			userName := r.URL.Query().Get("user")
			u, ok := storage.Get(userName)
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
		u, ok := storage.Get(userName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "user %s not found", userName)
			return
		}
		u.Hobby = userHobby
		storage.Set(u)
		fmt.Fprintf(w, "success, total users: %d", storage.Size())
	})

	fmt.Printf("%v", os.Args[1:])
	port := os.Args[1]
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

type UserStorage interface {
	Get(UserName string) (user, bool)
	Set(u user)
	Size() int
}

type MockStorage struct {
}

func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

func (s *MockStorage) Get(UserName string) (user, bool) {
	return user{
		Name:  "Artem",
		Hobby: "Web-Programming",
		Age:   23,
	}, true
}

func (s *MockStorage) Set(u user) {
}

func (s *MockStorage) Size() int {
	return 10
}

type Storage struct {
	data map[string]user
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]user),
	}
}

func (s *Storage) Get(UserName string) (user, bool) {
	u, b := s.data[UserName]
	return u, b
}

func (s *Storage) Set(u user) {
	s.data[u.Name] = u

}

func (s *Storage) Size() int {
	return len(s.data)
}

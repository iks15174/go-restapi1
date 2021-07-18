package myapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//Struct의 변수와 struct tag는 같으면 안된다.
type User struct {
	Id        int       `json: "id"`
	FirstName string    `json: "firstName"`
	LastName  string    `json: "lastName"`
	Email     string    `json: "email"`
	CreatedAt time.Time `json: "createdAt"`
}

var userMap map[int]*User
var lastId int

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Get UserInfo by /users/{id}")
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	user, ok := userMap[id]
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User id:", id)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))

	fmt.Fprint(w, "User Id:", vars["id"])
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)

	/*body to string*/
	// len := r.ContentLength
	// body := make([]byte, len)
	// r.Body.Read(body)
	// fmt.Println("body : ", string(body))

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	lastId++
	user.Id = lastId
	user.CreatedAt = time.Now()
	userMap[user.Id] = user

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
}

func deleteUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	_, ok := userMap[id]
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User Id:", id)
		return
	}
	delete(userMap, id)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted User Id:", id)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	updateUser := new(User)
	err := json.NewDecoder(r.Body).Decode(updateUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	user, ok := userMap[updateUser.Id]
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User Id:", updateUser.Id)
		return
	}
	if updateUser.FirstName != "" {
		user.FirstName = updateUser.FirstName
	}
	if updateUser.LastName != "" {
		user.LastName = updateUser.LastName
	}
	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
}

func NewHandler() http.Handler {
	userMap = make(map[int]*User)
	lastId = 0
	mux := mux.NewRouter()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/users", usersHandler).Methods("GET")
	mux.HandleFunc("/users", createUserHandler).Methods("POST")
	mux.HandleFunc("/users", updateUserHandler).Methods("PUT")
	mux.HandleFunc("/users/{id:[0-9]+}", getUserInfoHandler).Methods("GET")
	mux.HandleFunc("/users/{id:[0-9]+}", deleteUserInfoHandler).Methods("DELETE")

	return mux
}

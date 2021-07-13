package main

import (
	"TuckerYoutube/restapi1/myapp"
	"net/http"
)

func main() {
	http.ListenAndServe(":3000", myapp.NewHandler())
}

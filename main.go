package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SatyendraDhamgaye/mongoDbApi/router"
)

func main() {
	fmt.Println("mongoDB")
	r := router.Router()
	fmt.Println("Server is getting started...")
	log.Fatal(http.ListenAndServe(":4000", r))
	fmt.Println("Listning to port 8000...")

}

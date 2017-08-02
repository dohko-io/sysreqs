package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/dohko-io/sysreqs/mysql"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r.HandleFunc("/package/{name}", PackageHandler).
		Methods("GET").
		Schemes("http", "https")

	r.HandleFunc("/package", AddPackageHandler).
		Methods("POST").
		Schemes("http", "https")

	log.Fatal(http.ListenAndServe(":"+port, r))
}

// PackageHandler returns the dependencies of a given package.
func PackageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repository, err := mysql.NewPackageRepository(os.Getenv("sysreqs.datasource"))
	//
	if err != nil {
		panic(err)
	}
	//
	pkg, err := repository.GetWithName(vars["name"])

	if err != nil {
		panic(err)
	}

	json.Marshal(pkg)
}

// AddPackageHandler handles new packages to be included in the system
func AddPackageHandler(w http.ResponseWriter, r *http.Request) {

}

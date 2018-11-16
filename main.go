package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

/**
This is a Dog Breed comparision API, the only resource right now is a Dog where:
Dog : {
	breed
	avg_weight
	high energy or not?
	color
}

We want GET all, GET single, Create new, Update single, Delete single
*/

type Dog struct {
	gorm.Model
	Breed      string `json:"id,omitempty"`
	AvgWeight  int    `json:"weight,omitempty"`
	HighEnergy bool   `json:"high_energy,omitempty"`
	Color      string `json:"color,omitempty"`
}

// Database access variables - keep these global so the endpoint functions have access to the DB.
var db *gorm.DB
var err error

func GetDogs(w http.ResponseWriter, r *http.Request) {
	var dogs []Dog
	db.Find(&dogs)
	json.NewEncoder(w).Encode(dogs)
}

func GetDog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var dog Dog
	db.First(&dog, vars["id"])
	json.NewEncoder(w).Encode(dog)
}

func CreateDog(w http.ResponseWriter, r *http.Request) {
	var resource Dog
	json.NewDecoder(r.Body).Decode(&resource)
	db.Create(&resource)
	json.NewEncoder(w).Encode(&resource)
}

func DeleteDog(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var resource Dog
	db.First(&resource, params["id"])
	db.Delete(&resource)

	var resources []Dog
	db.Find(&resources)
	json.NewEncoder(w).Encode(&resources)
}

func main() {
	// Connect to the database and populate with Dog Struct
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Dog{})

	// Create
	dog := Dog{Breed: "Shiba Inu", AvgWeight: 15, HighEnergy: true, Color: "cream"}
	db.Create(&dog)

	router := mux.NewRouter()
	// Routes for this api
	router.HandleFunc("/dogs", GetDogs).Methods("GET")
	router.HandleFunc("/dogs/{id}", GetDog).Methods("GET")
	router.HandleFunc("/dogs/{id}", CreateDog).Methods("POST")
	router.HandleFunc("/dogs/{id}", DeleteDog).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

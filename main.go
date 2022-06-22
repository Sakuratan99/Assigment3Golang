package main

import (
	entity "Assigment3Golang/Entity"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var Cuaca entity.Cuaca

func main() {
	go RandomCuaca()
	mux := http.NewServeMux()
	endpoint := http.HandlerFunc(greet)
	mux.Handle("/users", MiddlewareCuaca(endpoint))
	fmt.Println("listening to port 8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}

func MiddlewareCuaca(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func RandomCuaca() {
	for {
		Cuaca.Water = rand.Intn(20)
		Cuaca.Wind = rand.Intn(20)
		Cuaca.StatusWater = "Aman"
		Cuaca.StatusWind = "Aman"
		//water
		if Cuaca.Water <= 5 {
			Cuaca.StatusWater = "Aman"
		} else if Cuaca.Water >= 6 && Cuaca.Water <= 8 {
			Cuaca.StatusWater = "Siaga"
		} else {
			Cuaca.StatusWater = "Bahaya"
		}
		//wind
		if Cuaca.Wind <= 6 {
			Cuaca.StatusWind = "Aman"
		} else if Cuaca.Wind >= 7 && Cuaca.Wind <= 15 {
			Cuaca.StatusWind = "Siaga"
		} else {
			Cuaca.StatusWind = "Bahaya"
		}

		//write json file
		jsonString, _ := json.Marshal(&Cuaca)
		ioutil.WriteFile("status.json", jsonString, os.ModePerm)
		//read from json

		time.Sleep(15 * time.Second)
	}

}

func greet(w http.ResponseWriter, r *http.Request) {
	// msg := "Hello world"
	// fmt.Fprint(w, msg)
	file, _ := ioutil.ReadFile("status.json")
	json.Unmarshal(file, &Cuaca)
	tpl, err := template.ParseFiles("template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	context := entity.Cuaca{
		Water:       Cuaca.Water,
		Wind:        Cuaca.Wind,
		StatusWater: Cuaca.StatusWater,
		StatusWind:  Cuaca.StatusWind,
	}

	tpl.Execute(w, context)
}

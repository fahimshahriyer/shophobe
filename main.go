package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type JsonObject struct {
	Name     string `json:"name"`
	Variance []struct {
		ID struct {
			Display string `json:"display"`
		} `json:"id"`
		Price struct {
			Current struct {
				Price int `json:"price"`
			} `json:"current"`
		} `json:"price"`
		Weight int `json:"weight"`
		Images []struct {
			Image string `json:"image"`
		} `json:"images"`
	} `json:"variance"`
	Category []struct {
		ID struct {
			Display string `json:"display"`
		} `json:"id"`
		Child struct {
			ID struct {
				Display string `json:"display"`
			} `json:"id"`
			Grand struct {
				Display string `json:"display"`
			} `json:"grand"`
		} `json:"child"`
	} `json:"category"`
}

// Product struct which contains product name,weight,price and image src.
type Product struct {
	Name   string
	Weight int
	Price  int
	Image  string
}

var products []Product

func getProducts() {
	root := "./data/"

	files := []string{}
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})

	for _, file := range files {
		file, _ := ioutil.ReadFile(file)
		var productObject JsonObject
		json.Unmarshal([]byte(file), &productObject)

		var product Product

		product = Product{
			Name:   productObject.Name,
			Weight: productObject.Variance[0].Weight,
			Price:  productObject.Variance[0].Price.Current.Price,
			Image:  productObject.Variance[0].Images[0].Image}

		products = append(products, product)
	}
}

func ProductsHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(products))
	t.Execute(w, &products)
}

func ProductsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

func ProductsByVarianceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Variance: %v\n", vars["variance"])
}

func main() {
	getProducts()
	// Router
	router := mux.NewRouter()
	// All Routes.
	router.HandleFunc("/", ProductsHandler)
	router.HandleFunc("/products/category/{category}", ProductsByCategoryHandler)
	router.HandleFunc("/products/variance/{variance}", ProductsByVarianceHandler)

	// Bind to a port and pass our router in
	fmt.Println("Serving on https://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

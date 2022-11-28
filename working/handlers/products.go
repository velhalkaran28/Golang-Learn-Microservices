package handlers

import (
	"context"
	"fmt"
	"kvlearn/data"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Add Products")
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	data.AddProduct(prod)
	p.l.Printf("Prod : %#v", prod)
}

func (p Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	p.l.Println("Handle PUT product ", id)
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	errD := data.UpdateProduct(id, prod)
	if errD == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if errD != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}
		err = prod.Validate()
		if err != nil {
			http.Error(rw, fmt.Sprintf("Error validating the product: %v", err), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}

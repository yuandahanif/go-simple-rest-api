package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type IceCream struct {
	Id      int    `json:"id"`
	Flavour string `json:"flavour"`
}

type JsonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    interface{}
}

var dataIceCream []IceCream

func ListAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	res := JsonResponse{Code: http.StatusOK, Message: "Lihat semua daftar es krim", Data: dataIceCream}
	resJson, err := json.Marshal(res)

	if err != nil {
		fmt.Println("Terjadi kesalahan", err)
		http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError)
		return
	}

	w.Write(resJson)
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	jsonDecoder := json.NewDecoder(r.Body)
	var iceCream IceCream
	err := jsonDecoder.Decode(&iceCream)

	if err != nil {
		fmt.Println("Terjadi kesalahan menambah data", err)
		http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError)
		return
	}

	iceCream.Id = dataIceCream[len(dataIceCream)-1].Id + 1
	dataIceCream = append(dataIceCream, iceCream)

	res := JsonResponse{Code: http.StatusOK, Message: "Berhasil tambah es krim"}
	resJson, _ := json.Marshal(res)

	w.Write(resJson)
}

func Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var notFound bool = true
	iceCreamId := ps.ByName("iceCreamId")
	id, err := strconv.Atoi(iceCreamId)

	if err != nil {
		fmt.Println("Terjadi kesalahan, parameter harus angka", err)
		http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError)
		return
	}

	for index, iceCream := range dataIceCream {
		if id == iceCream.Id {
			notFound = false
			dataIceCream = append(dataIceCream[:index], dataIceCream[index+1:]...)
		}
	}

	if notFound {
		res := JsonResponse{Code: http.StatusNotFound, Message: "Id tidak ditemukan"}
		resJson, _ := json.Marshal(res)
		w.Write(resJson)
		return
	}

	res := JsonResponse{Code: http.StatusOK, Message: "Berhasil hapus es krim"}
	resJson, _ := json.Marshal(res)

	w.Write(resJson)
}

func Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var notFound bool = true
	iceCreamId := ps.ByName("iceCreamId")
	id, err := strconv.Atoi(iceCreamId)

	if err != nil {
		fmt.Println("Terjadi kesalahan, parameter harus angka", err)
		http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError)
		return
	}

	jsonDecoder := json.NewDecoder(r.Body)
	var newIceCream IceCream

	if err := jsonDecoder.Decode(&newIceCream); err != nil {
		fmt.Println("Terjadi kesalahan saat mengubah data", err)
		http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError)
		return
	}

	for index, iceCream := range dataIceCream {
		if id == iceCream.Id {
			notFound = false
			iceCream.Flavour = newIceCream.Flavour

			dataIceCream[index] = iceCream
		}
	}

	if notFound {
		res := JsonResponse{Code: http.StatusNotFound, Message: "Id tidak ditemukan"}
		resJson, _ := json.Marshal(res)
		w.Write(resJson)
		return
	}

	res := JsonResponse{Code: http.StatusOK, Message: "Berhasil ubah es krim"}
	resJson, _ := json.Marshal(res)

	w.Write(resJson)
}

func main() {
	dataIceCream = append(dataIceCream, IceCream{Id: 1, Flavour: "coklat"})
	dataIceCream = append(dataIceCream, IceCream{Id: 2, Flavour: "vanila"})

	router := httprouter.New()

	router.GET("/", ListAll)
	router.POST("/", Create)
	router.DELETE("/:iceCreamId", Delete)
	router.PUT("/:iceCreamId", Update)

	log.Fatal(http.ListenAndServe(":8000", router))

}

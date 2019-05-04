package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Data   interface{}
	Error  error
	Status int
}

type DogHandler func(*http.Request, Kennel) Response

func CreateServer(port string, kennel Kennel) *http.Server {
	router := mux.NewRouter()
	router.Handle("/dogs", DogHttpHandler(GetAllDogs, kennel)).Methods(http.MethodGet)
	router.Handle("/dogs", DogHttpHandler(SaveDog, kennel)).Methods(http.MethodPost)
	router.Handle("/dogs/{id}", DogHttpHandler(GetDog, kennel)).Methods(http.MethodGet)
	router.Handle("/dogs/{id}", DogHttpHandler(UpdateDog, kennel)).Methods(http.MethodPut)
	router.Handle("/dogs/{id}", DogHttpHandler(RemoveDog, kennel)).Methods(http.MethodDelete)
	return &http.Server{Addr: port, Handler: router}
}

func GetDog(r *http.Request, kennel Kennel) Response {
	var dog Dog
	var resp Response
	dog.Id, resp.Error = parseId(r)
	if resp.Error != nil {
		return resp
	}
	resp.Status = http.StatusOK
	ok := kennel.GetOne(&dog)
	if !ok {
		resp.Status = http.StatusNotFound
	}
	resp.Data = dog
	return resp
}

func SaveDog(r *http.Request, kennel Kennel) Response {
	var dog Dog
	var resp Response
	resp.Error = parseBody(r, &dog)
	if resp.Error != nil {
		return resp
	}
	resp.Status = http.StatusOK
	resp.Error = kennel.Save(&dog)
	if resp.Error != nil {
		resp.Status = http.StatusInternalServerError
	}
	resp.Data = dog
	return resp
}

func UpdateDog(r *http.Request, kennel Kennel) Response {
	var dog Dog
	var resp Response
	resp.Error = parseBody(r, &dog)
	if resp.Error != nil {
		return resp
	}
	dog.Id, resp.Error = parseId(r)
	if resp.Error != nil {
		return resp
	}
	resp.Status = http.StatusOK
	resp.Error = kennel.Save(&dog)
	if resp.Error != nil {
		resp.Status = http.StatusInternalServerError
	}
	resp.Data = dog
	return resp
}

func GetAllDogs(r *http.Request, kennel Kennel) Response {
	var resp Response
	dogs := kennel.GetAll()
	resp.Data = dogs
	return resp
}

func RemoveDog(r *http.Request, kennel Kennel) Response {
	var dog Dog
	var resp Response
	dog.Id, resp.Error = parseId(r)
	if resp.Error != nil {
		return resp
	}
	resp.Error = kennel.Remove(&dog)
	if resp.Error != nil {
		resp.Status = http.StatusInternalServerError
	}
	return resp
}

func DogHttpHandler(dogHandler DogHandler, kennel Kennel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := dogHandler(r, kennel)
		if resp.Error != nil {
			WriteError(w, resp)
			return
		}
		WriteSuccess(w, resp)
	}
}

func WriteSuccess(w http.ResponseWriter, resp Response) {
	var jsonData []byte
	resp.Status = http.StatusOK
	jsonData, resp.Error = json.Marshal(resp.Data)
	w.WriteHeader(resp.Status)
	w.Write(jsonData)
}

func WriteError(w http.ResponseWriter, resp Response) {
	w.WriteHeader(resp.Status)
	w.Write([]byte(fmt.Sprintf("error: %#v", resp.Error)))
}

func parseBody(r *http.Request, val interface{}) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, val)
}

func parseId(r *http.Request) (string, error) {
	id := mux.Vars(r)["id"]
	if id == "" {
		return id, fmt.Errorf("dog id is required")
	}
	return id, nil
}
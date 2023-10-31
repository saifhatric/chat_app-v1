package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)




func Routes(addr string){

	app:=mux.NewRouter()
	views :=http.StripPrefix("/",http.FileServer(http.Dir("views/"))) 
	app.Handle("/",views)
	app.HandleFunc("/ws",WsEndpoint)

	log.Fatal(http.ListenAndServe(addr, app))

}
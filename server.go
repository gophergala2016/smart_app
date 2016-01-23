package main

import (
	"./resources"
	"fmt"
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
	"os"
)

func main() {
	wsContainer := restful.NewContainer()

	r := resources.UserResource{}
	r.Register(wsContainer)

	bind := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	log.Printf("listening on %s", bind)

	server := &http.Server{Addr: bind, Handler: wsContainer}
	log.Fatalln(server.ListenAndServe())
}

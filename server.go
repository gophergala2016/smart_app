package main

import (
	"./resources"
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/ahmet/osin-rethinkdb"
	"github.com/emicklei/go-restful"
	"github.com/honeybadger-io/honeybadger-go"
	r "gopkg.in/dancannon/gorethink.v1"
	"log"
	"net/http"
	"os"
)

var (
	server  *osin.Server
	session *r.Session
)

func main() {
	defer honeybadger.Monitor()

	initDb()

	config := osin.NewServerConfig()
	config.ErrorStatusCode = 401

	server = osin.NewServer(config, RethinkDBStorage.New(session))

	wsContainer := restful.NewContainer()

	r := resources.UserResource{}
	r.Register(wsContainer)

	ws := new(restful.WebService)
	ws.Route(ws.POST("/authorize").
		Consumes("application/x-www-form-urlencoded").
		To(authorize))
	wsContainer.Add(ws)

	address := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	log.Printf("listening on %s", address)

	log.Fatalln(http.ListenAndServe(address, honeybadger.Handler(wsContainer)))
}

func initDb() {
	address := fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	session, err := r.Connect(r.ConnectOpts{
		Address:  address,
		Database: os.Getenv("DB_NAME"),
	})
	if err != nil {
		log.Fatalln(err)
	}

	_, err = r.DBCreate(os.Getenv("DB_NAME")).RunWrite(session)
	if err != nil {
		log.Println(err)
	}
}

func authorize(req *restful.Request, res *restful.Response) {
	nr := server.NewResponse()
	defer nr.Close()

	if ar := server.HandleAuthorizeRequest(nr, req.Request); ar != nil {
		if !authenticate(ar, req, res) {
			return
		}

		ar.Authorized = true
		server.FinishAuthorizeRequest(nr, req.Request, ar)
	}

	if nr.IsError && nr.InternalError != nil {
		honeybadger.Notify(nr.InternalError)
	}

	osin.OutputJSON(nr, res.ResponseWriter, req.Request)
}

func authenticate(ar *osin.AuthorizeRequest, req *restful.Request, res *restful.Response) bool {
	r := req.Request

	r.ParseForm()

	if r.Form.Get("username") == "test" && r.Form.Get("password") == "test" {
		return true
	}

	res.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
	res.WriteErrorString(http.StatusUnauthorized, "Not Authorized")

	return false
}

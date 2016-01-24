package main

import (
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/emicklei/go-restful"
	"github.com/jmoiron/sqlx"
	"github.com/ory-am/osin-storage/storage/postgres"
	"log"
	"net/http"
	"os"
)

type Server struct {
	server *osin.Server
	Host   string
	Port   string
}

func (s *Server) Start() {
	config := osin.NewServerConfig()
	config.ErrorStatusCode = 401

	url := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	db, err := sqlx.Open("postgres", url)
	if err != nil {
		log.Fatalln(err.Error())
	}

	storage := postgres.New(db.DB)
	s.server = osin.NewServer(config, storage)

	wsContainer := restful.NewContainer()

	r := UserResource{}
	r.Register(wsContainer, db)

	ws := new(restful.WebService)
	ws.Route(ws.POST("/authorize").
		Consumes("application/x-www-form-urlencoded").
		To(s.authorize))
	wsContainer.Add(ws)

	address := fmt.Sprintf("%s:%s", s.Host, s.Port)
	log.Printf("Listening on %s", address)

	log.Fatalln(http.ListenAndServe(address, wsContainer))
}

func (s *Server) authorize(req *restful.Request, res *restful.Response) {
	nr := s.server.NewResponse()
	defer nr.Close()

	if ar := s.server.HandleAuthorizeRequest(nr, req.Request); ar != nil {
		if !s.authenticate(ar, req, res) {
			return
		}

		ar.Authorized = true
		s.server.FinishAuthorizeRequest(nr, req.Request, ar)
	}

	if nr.IsError && nr.InternalError != nil {
		log.Printf("ERROR: %s", nr.InternalError)
		res.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
	}

	osin.OutputJSON(nr, res.ResponseWriter, req.Request)
}

func (s *Server) authenticate(ar *osin.AuthorizeRequest, req *restful.Request, res *restful.Response) bool {
	r := req.Request

	r.ParseForm()

	if r.Form.Get("username") == "test" && r.Form.Get("password") == "test" {
		return true
	}

	return false
}

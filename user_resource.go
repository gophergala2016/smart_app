package main

import (
	"github.com/emicklei/go-restful"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type User struct {
	Id        uint64      `json:"id"`
	Email     string      `json:"email"`
	Mobile    string      `json:"mobile"`
	Password  string      `json:"-"`
	CreatedAt pq.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt pq.NullTime `db:"updated_at" json:"updated_at"`
	DeletedAt pq.NullTime `db:"deleted_at" json:"-"`
}

type UserResource struct {
	db *sqlx.DB
}

func (u UserResource) Register(c *restful.Container, db *sqlx.DB) {
	u.db = db

	ws := new(restful.WebService)

	ws.Path("/api/v1/users").
		Doc("Manage users").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{id}").
		To(u.findUser).
		Doc("get a user").
		Operation("findUser").
		Param(ws.PathParameter("id", "identifier of the user").
		DataType("string")).
		Writes(User{}))

	ws.Route(ws.PUT("/{id}").
		To(u.updateUser).
		Doc("update a user").
		Operation("updateUser").
		Param(ws.PathParameter("id", "identifier of the user").
		DataType("string")).
		Returns(409, "duplicate id", nil).
		Reads(User{}))

	ws.Route(ws.POST("").
		To(u.createUser).
		Doc("create a user").
		Operation("createUser").
		Reads(User{}))

	ws.Route(ws.DELETE("/{id}").
		To(u.removeUser).
		Doc("delete a user").
		Operation("removeUser").
		Param(ws.PathParameter("id", "identifier of the user").
		DataType("string")))

	c.Add(ws)
}

func (u UserResource) findUser(req *restful.Request, res *restful.Response) {
	id := req.PathParameter("id")

	usr := User{}
	err := u.db.Get(&usr,
		`SELECT
			*
		FROM
			users
		WHERE
			id = $1
		AND
			deleted_at IS NULL`, id)
	if err != nil {
		res.AddHeader("Content-Type", "text/plain")
		res.WriteErrorString(http.StatusNotFound, err.Error())

		return
	}

	res.WriteEntity(usr)
}

func (u *UserResource) createUser(req *restful.Request, res *restful.Response) {
	usr := User{}
	err := req.ReadEntity(&usr)
	if err != nil {
		res.AddHeader("Content-Type", "text/plain")
		res.WriteErrorString(http.StatusInternalServerError, err.Error())

		return
	}

	usr.CreatedAt = pq.NullTime{time.Now(), true}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	if err != nil {
		res.AddHeader("Content-Type", "text/plain")
		res.WriteErrorString(http.StatusInternalServerError, err.Error())

		return
	}

	usr.Password = string(hashedPassword)

	rows, err := u.db.NamedQuery(
		`INSERT INTO
			users (
				email,
				mobile,
				password,
				created_at
			)
		VALUES (
			:email,
	        :mobile,
	        :password,
	        :created_at
	    )
		RETURNING id`, usr)
	if err != nil {
		res.AddHeader("Content-Type", "text/plain")
		res.WriteErrorString(http.StatusInternalServerError, err.Error())

		return
	}

	if rows.Next() {
		rows.Scan(&usr.Id)
	}

	res.WriteHeaderAndEntity(http.StatusCreated, usr)
}

func (u *UserResource) updateUser(req *restful.Request, res *restful.Response) {
	usr := User{}
	err := req.ReadEntity(&usr)
	if err != nil {
		res.AddHeader("Content-Type", "text/plain")
		res.WriteErrorString(http.StatusInternalServerError, err.Error())

		return
	}

	usr.UpdatedAt = pq.NullTime{time.Now(), true}

	_, err = u.db.NamedExec(
		`UPDATE
			users
		SET
			email = :email,
			mobile = :mobile,
			updated_at = :updated_at
		WHERE
			id = :id`, usr)
	if err != nil {
		res.AddHeader("Content-Type", "text/plain")
		res.WriteErrorString(http.StatusInternalServerError, err.Error())

		return
	}

	res.WriteEntity(usr)
}

func (u *UserResource) removeUser(req *restful.Request, res *restful.Response) {
	id := req.PathParameter("id")

	_, err := u.db.Exec(
		`UPDATE
			users
		SET
			deleted_at = $1
		WHERE
			id = $2`, time.Now(), id)
	if err != nil {
		res.AddHeader("Content-Type", "text/plain")
		res.WriteErrorString(http.StatusInternalServerError, err.Error())

		return
	}
}

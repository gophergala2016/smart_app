package resources

import (
	"github.com/emicklei/go-restful"
)

type User struct {
	Id uint64 `json:"id"`
}

type UserResource struct{}

func (u UserResource) Register(container *restful.Container) {
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

	container.Add(ws)
}

func (u UserResource) findUser(request *restful.Request, response *restful.Response) {
	return
}

func (u *UserResource) createUser(request *restful.Request, response *restful.Response) {
	return
}

func (u *UserResource) updateUser(request *restful.Request, response *restful.Response) {
	return
}

func (u *UserResource) removeUser(request *restful.Request, response *restful.Response) {
	return
}

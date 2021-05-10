package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/joyrex2001/kubedock/internal/backend"
	"github.com/joyrex2001/kubedock/internal/model"
	"github.com/joyrex2001/kubedock/internal/server/httputil"
)

// Router is the object that facilitates the kubedock API endpoints.
type Router struct {
	db  *model.Database
	kub backend.Backend
}

// New will instantiate a containerRouter object.
func New(router *gin.Engine, kub backend.Backend) (*Router, error) {
	db, err := model.New()
	if err != nil {
		return nil, err
	}
	cr := &Router{
		db:  db,
		kub: kub,
	}
	cr.initRoutes(router)
	return cr, nil
}

// initRoutes will add all suported routes.
func (cr *Router) initRoutes(router *gin.Engine) {
	// DEBT: "POST /containers/create" and "POST /containers/:id/start" overlap
	// which (currently) means :id that start with a c will 404 by the router;
	// as work around, no container-ids are generated that start with a 'c'...
	router.POST("/containers/create", cr.ContainerCreate)
	router.POST("/containers/:id/start", cr.ContainerStart)

	// DEBT: "POST /networks/create" and "POST /networks/:id/connect"" overlap
	// which (currently) means :id that start with a c will 404 by the router;
	// as work around, no network-ids are generated that start with a 'c'...
	router.POST("/networks/create", cr.NetworksCreate)
	router.POST("/networks/:id/connect", cr.NetworksConnect)
	router.POST("/networks/:id/disconnect", cr.NetworksDisconnect)

	router.GET("/containers/json", cr.ContainerList)
	router.GET("/containers/:id/json", cr.ContainerInfo)
	router.DELETE("/containers/:id", cr.ContainerDelete)
	router.POST("/containers/:id/exec", cr.ContainerExec)
	router.GET("/containers/:id/logs", cr.ContainerLogs)
	router.POST("/exec/:id/start", cr.ExecStart)
	router.GET("/exec/:id/json", cr.ExecInfo)
	router.PUT("/containers/:id/archive", cr.PutArchive)
	router.GET("/networks", cr.NetworksList)
	router.DELETE("/networks/:id", cr.NetworksDelete)
	router.GET("/networks/:id", cr.NetworksInfo)
	router.GET("/images/json", cr.ImageList)
	router.POST("/images/create", cr.ImageCreate)
	router.GET("/images/:image/*json", cr.ImageJSON)
	router.GET("/_ping", cr.Ping)
	router.HEAD("/_ping", cr.Ping)
	router.GET("/info", cr.Info)
	router.GET("/version", cr.Version)
	router.GET("/healthz", cr.Healthz)

	router.POST("/containers/:id/stop", httputil.NotImplemented)
	router.POST("/containers/:id/kill", httputil.NotImplemented)

	// not supported at the moment
	router.GET("/containers/:id/top", httputil.NotImplemented)
	router.GET("/containers/:id/changes", httputil.NotImplemented)
	router.GET("/containers/:id/export", httputil.NotImplemented)
	router.GET("/containers/:id/stats", httputil.NotImplemented)
	router.POST("/containers/:id/resize", httputil.NotImplemented)
	router.POST("/containers/:id/restart", httputil.NotImplemented)
	router.POST("/containers/:id/update", httputil.NotImplemented)
	router.POST("/containers/:id/rename", httputil.NotImplemented)
	router.POST("/containers/:id/pause", httputil.NotImplemented)
	router.POST("/containers/:id/unpause", httputil.NotImplemented)
	router.POST("/containers/:id/attach", httputil.NotImplemented)
	router.GET("/containers/:id/attach/ws", httputil.NotImplemented)
	router.POST("/containers/:id/wait", httputil.NotImplemented)
	router.HEAD("/containers/:id/archive", httputil.NotImplemented)
	router.GET("/containers/:id/archive", httputil.NotImplemented)
	router.POST("/containers/prune", httputil.NotImplemented)
	router.GET("/networks/reaper_default", httputil.NotImplemented)
}

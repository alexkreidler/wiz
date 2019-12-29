//implements a server for the Wiz Processor Server API interface
//default is http
package server

import (
	"github.com/alexkreidler/wiz/api"
	"github.com/gin-gonic/gin"
	"log"
)

type Server struct {
	*gin.Engine
}

func NewServer(server api.ProcessorServer) Server {

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	// Basic GET requests
	router.GET("/processors", func(c *gin.Context) {
		p, err := server.GetAllProcessors()
		if err != nil {
			c.Error(err)
		}
		c.JSON(200, p)
	})
	router.GET("/processors/:procID", func(c *gin.Context) {
		id := c.Param("procID")
		p, err := server.GetProcessor(id)
		if err != nil {
			c.Error(err)
		}
		c.JSON(200, p)
	})
	router.GET("/processors/:procID/runs", func(c *gin.Context) {
		id := c.Param("procID")
		p, err := server.GetRuns(id)
		if err != nil {
			c.Error(err)
		}
		c.JSON(200, p)
	})
	router.GET("/processors/:procID/runs/:runID", func(c *gin.Context) {
		procID := c.Param("procID")
		runID := c.Param("runID")
		p, err := server.GetRun(procID, runID)
		if err != nil {
			c.Error(err)
		}
		c.JSON(200, p)
	})

	// Configuration and manipulation requests
	router.GET("/processors/:procID/runs/:runID/config", func(c *gin.Context) {
		procID := c.Param("procID")
		runID := c.Param("runID")
		p, err := server.GetConfig(procID, runID)
		if err != nil {
			c.Error(err)
		}
		c.JSON(200, p)
	})
	router.POST("/processors/:procID/runs/:runID/config", func(c *gin.Context) {
		procID := c.Param("procID")
		runID := c.Param("runID")
		var data api.Configuration

		err := c.BindJSON(&data)
		if err != nil {
			c.Error(err)
		}

		log.Printf("got config: %#+v", data)
		err = server.Configure(procID, runID, data)
		if err != nil {
			c.Error(err)
		}
		c.Status(200)
	})
	router.GET("/processors/:procID/runs/:runID/data", func(c *gin.Context) {
		procID := c.Param("procID")
		runID := c.Param("runID")
		p, err := server.GetRunData(procID, runID)
		if err != nil {
			c.Error(err)
		}
		c.JSON(200, p)
	})
	router.POST("/processors/:procID/runs/:runID/data", func(c *gin.Context) {
		procID := c.Param("procID")
		runID := c.Param("runID")
		var data api.Data
		err := c.BindJSON(&data)
		//c.MustBindWith()
		if err != nil {
			c.Error(err)
			return
		}

		err = server.AddData(procID, runID, data)
		if err != nil {
			c.Error(err)
		}
		c.Status(200)
	})

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	//router.Run()
	// router.Run(":3000") for a hard coded port
	return Server{router}
}

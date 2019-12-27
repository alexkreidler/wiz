//implements a server for the Wiz Processor Server API interface
//default is http
package server

import (
	"encoding/json"
	"github.com/alexkreidler/wiz/api"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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
		spew.Dump(p)
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
		var data interface{}
		body, err := ioutil.ReadAll(c.Request.Body)

		log.Println("recieved body configuration:", string(body))
		err = json.Unmarshal(body, &data)
		if err != nil {
			c.Error(err)
		}
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
		err := c.BindJSON(data)
		//c.MustBindWith()
		if err != nil {
			c.Error(err)
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

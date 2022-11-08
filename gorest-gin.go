// REST server implemented with Gin, with middleware.
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"example.com/internal/taskstore"
	"github.com/gin-gonic/gin"
)

type taskServer struct {
	store *taskstore.TaskStore
}

func NewTaskServer() *taskServer {
	store := taskstore.New()
	return &taskServer{store: store}
}

func (ts *taskServer) getAllTasksHandler(c *gin.Context) {
	allTasks := ts.store.GetAllTasks()
	c.JSON(http.StatusOK, allTasks)
}

func (ts *taskServer) deleteAllTasksHandler(c *gin.Context) {
	ts.store.DeleteAllTasks()
}

func (ts *taskServer) createTaskHandler(c *gin.Context) {
	type RequestTask struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}

	var rt RequestTask
	if err := c.ShouldBindJSON(&rt); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(rt)
	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)
	c.JSON(http.StatusOK, gin.H{"Id": id})
}

func (ts *taskServer) getTaskHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	task, err := ts.store.GetTask(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, task)
}

func (ts *taskServer) deleteTaskHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if err = ts.store.DeleteTask(id); err != nil {
		c.String(http.StatusNotFound, err.Error())
	}
}

func (ts *taskServer) tagHandler(c *gin.Context) {
	tag := c.Params.ByName("tag")
	tasks := ts.store.GetTasksByTag(tag)
	c.JSON(http.StatusOK, tasks)
}
func (ts *taskServer) helloHandler(c *gin.Context) {
	//fmt.Println("Hello world")
	type HelloWorld struct {
		Text string
	}
	var t = HelloWorld{"Hello World!"}
	//if err := c.ShouldBindJSON(&t); err != nil {
	//	c.String(http.StatusBadRequest, err.Error())
	//	return
	//}
	//id := ts.store.CreateTask(t.Text)
	c.JSON(http.StatusOK, gin.H{"helloworld": t.Text})
}

func (ts *taskServer) authHandler(c *gin.Context) {

	type Provider struct {
		Provider string `json:"provider"`
	}
	var p Provider

	if err := c.ShouldBindJSON(&p); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(p.Provider)

	c.JSON(http.StatusOK, gin.H{"prov": "provider"})

	//var rt RequestTask
	//if err := c.ShouldBindJSON(&rt); err != nil {
	//	c.String(http.StatusBadRequest, err.Error())
	//	return
	//}
	//
	//fmt.Println(rt)
	//id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)
	//c.JSON(http.StatusOK, gin.H{"Id": id})

}

func (ts *taskServer) dueHandler(c *gin.Context) {
	badRequestError := func() {
		c.String(http.StatusBadRequest, "expect /due/<year>/<month>/<day>, got %v", c.FullPath())
	}

	year, err := strconv.Atoi(c.Params.ByName("year"))
	if err != nil {
		badRequestError()
		return
	}

	month, err := strconv.Atoi(c.Params.ByName("month"))
	if err != nil || month < int(time.January) || month > int(time.December) {
		badRequestError()
		return
	}

	day, err := strconv.Atoi(c.Params.ByName("day"))
	if err != nil {
		badRequestError()
		return
	}

	tasks := ts.store.GetTasksByDueDate(year, time.Month(month), day)
	c.JSON(http.StatusOK, tasks)
}

func main() {
	// Set up middleware for logging and panic recovery explicitly.
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	server := NewTaskServer()

	router.POST("/task/", server.createTaskHandler)
	router.GET("/task/", server.getAllTasksHandler)
	router.DELETE("/task/", server.deleteAllTasksHandler)
	router.GET("/task/:id", server.getTaskHandler)
	router.DELETE("/task/:id", server.deleteTaskHandler)
	router.GET("/tag/:tag", server.tagHandler)
	router.GET("/hello", server.helloHandler)
	router.POST("/auth", server.authHandler)
	router.GET("/due/:year/:month/:day", server.dueHandler)

	router.Run("localhost:" + os.Getenv("SERVERPORT"))
}

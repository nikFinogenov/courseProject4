// fetch('http://localhost:4112/auth', {method: 'POST', body: JSON.stringify({provider: 'google'})})
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

func CORS(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:4112")
	c.Writer.Header().Set("Access-Control-Max-Age", "15")
}
func (ts *taskServer) preflight(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	//access-control-allow-headers: authorization,content-type
	//	access-control-allow-methods: POST
	//access-control-allow-origin: *
	c.Header("Access-Control-Allow-Methods", "POST")
	c.Header("Access-Control-Allow-Headers", "authorization, content-type")
	c.JSON(204, "")
}

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

func (ts *taskServer) loginHandler(c *gin.Context) {

	c.Header("Content-type", "text/html charset=utf-8")

	//c.JSON(http.StatusOK, "<form method='POST' action='/auth-fd'><input type='hidden' name='provider' value='google'><button type='submit'>Login</button></form>")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte("<form method='POST' action='/auth-fd'><input type='hidden' name='provider' value='google'><button type='submit'>Login</button></form>"))
}

func (ts *taskServer) authHandler(c *gin.Context) {
	//fmt.Println("qweqwe")
	//c.Header("Access-Control-Allow-Origin", "*")
	//c.Header("Access-Control-Allow-Credentials", "true")
	//c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	//c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")
	type Provider struct {
		Provider string `json:"provider"`
	}
	const client_id = "928253709894-9mjbvcvh2g4ltak6hlmb7kij70dlcdnl.apps.googleusercontent.com"
	const secret = "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"
	var redirect_url string = "http://localhost:" + os.Getenv("SERVERPORT") + "/auth-callback"
	var p Provider
	var url = "https://accounts.google.com/o/oauth2/v2/auth?client_id=" +
		client_id +
		"&response_type=code" +
		"&scope=openid" +
		"&redirect_uri=" +
		redirect_url

	if err := c.ShouldBindJSON(&p); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if p.Provider == "google" {
		CORS(c)
		c.Redirect(302, url)
	}
}
func (ts *taskServer) authHandlerFD(c *gin.Context) {
	//fmt.Println("qweqwe")
	//c.Header("Access-Control-Allow-Origin", "*")
	//c.Header("Access-Control-Allow-Credentials", "true")
	//c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	//c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")
	type Provider struct {
		Provider string `form:"provider"`
	}
	const client_id = "928253709894-9mjbvcvh2g4ltak6hlmb7kij70dlcdnl.apps.googleusercontent.com"
	const secret = "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"
	var redirect_url string = "http://localhost:" + os.Getenv("SERVERPORT") + "/auth-callback"
	var p Provider
	var url = "https://accounts.google.com/o/oauth2/v2/auth?client_id=" +
		client_id +
		"&response_type=code" +
		"&scope=openid" +
		"&redirect_uri=" +
		redirect_url

	if err := c.Bind(&p); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if p.Provider == "google" {
		CORS(c)
		c.Redirect(302, url)
	}
}
func (ts *taskServer) authGETHandler(c *gin.Context) {

	//type Provider struct {
	//	Provider string `json:"provider"`
	//}
	const client_id = "928253709894-9mjbvcvh2g4ltak6hlmb7kij70dlcdnl.apps.googleusercontent.com"
	const secret = "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"
	var redirect_url string = "http://localhost:" + os.Getenv("SERVERPORT") + "/auth-callback"
	//var p Provider
	var url = "https://accounts.google.com/o/oauth2/v2/auth?client_id=" +
		client_id +
		"&response_type=code" +
		"&redirect_uri=" +
		redirect_url +
		"&scope=openid"
	//if err := c.ShouldBindJSON(&p); err != nil {
	//	c.String(http.StatusBadRequest, err.Error())
	//	return
	//}

	c.Redirect(302, url)
	//if p.Provider == "google" {
	//	c.Redirect(302, url)
	//}
}
func (ts *taskServer) authCallbackHandler(c *gin.Context) {
	q := c.Request.URL.Query()
	fmt.Println(q)
	fmt.Println("qweqweqweqweqwe")

	///
	//get jwt from google
	///
	c.JSON(http.StatusOK, "Ya est grut: JWT")
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
	//router.Use(CORS())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	server := NewTaskServer()
	router.OPTIONS("/auth", server.preflight)
	router.GET("/login", server.loginHandler)
	router.POST("/task/", server.createTaskHandler)
	router.GET("/task/", server.getAllTasksHandler)
	router.DELETE("/task/", server.deleteAllTasksHandler)
	router.GET("/task/:id", server.getTaskHandler)
	router.DELETE("/task/:id", server.deleteTaskHandler)
	router.GET("/tag/:tag", server.tagHandler)
	router.GET("/hello", server.helloHandler)
	router.POST("/auth", server.authHandler)
	router.GET("/auth-get", server.authGETHandler)
	router.POST("/auth-fd", server.authHandlerFD)
	router.GET("/auth-callback", server.authCallbackHandler)
	router.GET("/due/:year/:month/:day", server.dueHandler)

	router.Run("localhost:" + os.Getenv("SERVERPORT"))
}

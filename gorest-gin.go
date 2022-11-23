// fetch('http://localhost:4112/auth', {method: 'POST', body: JSON.stringify({provider: 'google'})})
package main

import (
	"bytes"
	"encoding/json"
	"example.com/internal/taskstore"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

const client_id = "928253709894-9mjbvcvh2g4ltak6hlmb7kij70dlcdnl.apps.googleusercontent.com"
const secret = "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"

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

func (ts *taskServer) loginFDHandler(c *gin.Context) {

	c.Header("Content-type", "text/html charset=utf-8")

	//c.JSON(http.StatusOK, "<form method='POST' action='/auth-fd'><input type='hidden' name='provider' value='google'><button type='submit'>Login</button></form>")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte("<form method='POST' action='/auth'><input type='hidden' name='provider' value='google'><button type='submit'>Login</button></form>"))
}

//	func (ts *taskServer) authHandlerFD(c *gin.Context) {
//		type Provider struct {
//			Provider string `json:"provider"`
//		}
//		const client_id = "928253709894-9mjbvcvh2g4ltak6hlmb7kij70dlcdnl.apps.googleusercontent.com"
//		const secret = "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"
//		var redirect_url string = "http://localhost:" + os.Getenv("SERVERPORT") + "/auth-callback"
//		var p Provider
//		var url = "https://accounts.google.com/o/oauth2/v2/auth?client_id=" +
//			client_id +
//			"&response_type=code" +
//			"&scope=openid" +
//			"&redirect_uri=" +
//			redirect_url
//
//		if err := c.ShouldBindJSON(&p); err != nil {
//			c.String(http.StatusBadRequest, err.Error())
//			fmt.Println("BAD request")
//			return
//		}
//		if p.Provider == "google" {
//			fmt.Println("GOOD request")
//			//CORS(c)
//			//c.Redirect(302, url)
//			c.JSON(http.StatusOK, gin.H{"helloworld": t.Text})
//		}
//	}
func (ts *taskServer) authHandler(c *gin.Context) {
	type Provider struct {
		Provider string `form:"provider"`
	}

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
		//CORS(c)
		c.Redirect(302, url)
	}
}

//func (ts *taskServer) authGETHandler(c *gin.Context) {
//
//	//type Provider struct {
//	//	Provider string `json:"provider"`
//	//}
//	const client_id = "928253709894-9mjbvcvh2g4ltak6hlmb7kij70dlcdnl.apps.googleusercontent.com"
//	const secret = "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"
//	var redirect_url string = "http://localhost:" + os.Getenv("SERVERPORT") + "/auth-callback"
//	//var p Provider
//	var url = "https://accounts.google.com/o/oauth2/v2/auth?client_id=" +
//		client_id +
//		"&response_type=code" +
//		"&redirect_uri=" +
//		redirect_url +
//		"&scope=openid"
//	//if err := c.ShouldBindJSON(&p); err != nil {
//	//	c.String(http.StatusBadRequest, err.Error())
//	//	return
//	//}
//
//	c.Redirect(302, url)
//	//if p.Provider == "google" {
//	//	c.Redirect(302, url)
//	//}
//}

func (ts *taskServer) authCallbackGETHandler(c *gin.Context) {
	type GRSP struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   string `json:"expires_in"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
		IdToken     string `json:"id_token"`
	}

	q := c.Request.URL.Query().Get("code")
	//origin := c.Request.URL.Query().Get("origin")
	//fmt.Println(q, "   ", reflect.TypeOf(q))
	url := "https://oauth2.googleapis.com/token"
	values := map[string]string{
		"code":          q,
		"client_id":     client_id,
		"client_secret": secret,
		"grant_type":    "authorization_code",
		"redirect_uri":  "http://localhost:4112/auth-callback"}
	json_data, err := json.Marshal(values)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	var gResp GRSP
	body, err := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &gResp); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	//body =
	//jwt = resp.Body.
	fmt.Println("response Body:", gResp.IdToken)
	//c.Redirect(302, origin + "?token")

}

func (ts *taskServer) authCallbackPOSTHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

func (ts *taskServer) authJsonRedirect(c *gin.Context) {
	type Provider struct {
		Provider string `json:"provider"`
		Origin   string `json:"origin"`
	}
	var p Provider
	var redirect_url string = "http://localhost:" + os.Getenv("SERVERPORT") + "/auth-callback?origin=" + p.Origin

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
		//CORS(c)
		c.JSON(http.StatusOK, gin.H{"link": url})
	}
}

func main() {
	// Set up middleware for logging and panic recovery explicitly.
	router := gin.New()
	//config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"*"}
	//router.Use(cors.New(config))
	router.Use(cors.Default())
	//router.Use(CORS())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	server := NewTaskServer()
	//router.OPTIONS("/auth", server.preflight)//cors
	router.GET("/login", server.loginFDHandler)
	router.GET("/hello", server.helloHandler)

	//router.POST("/auth-fd", server.authHandler)
	//router.GET("/auth-get", server.authGETHandler)
	router.POST("/auth", server.authHandler)
	router.GET("/auth-callback", server.authCallbackGETHandler)
	router.POST("/auth-callback", server.authCallbackPOSTHandler)
	router.POST("/auth/provider-link", server.authJsonRedirect)

	router.Run("localhost:" + os.Getenv("SERVERPORT"))
}

// fetch('http://localhost:4112/auth', {method: 'POST', body: JSON.stringify({provider: 'google'})})
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"example.com/internal/taskstore"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const client_id = "928253709894-9mjbvcvh2g4ltak6hlmb7kij70dlcdnl.apps.googleusercontent.com"
const secret = "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"

type taskServer struct {
	store *taskstore.TaskStore
}

func NewTaskServer() *taskServer {
	store := taskstore.New()
	return &taskServer{store: store}
}
func verifyIdToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken := strings.Replace(c.Request.Header["Authorization"][0], "Bearer ", "", 1)
		_, err := idtoken.Validate(context.Background(), jwtToken, client_id)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
func tokenTypeValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := c.Request.Header["Authorization"]
		if ok {
			jwtToken := c.Request.Header["Authorization"][0]
			if !strings.HasPrefix(jwtToken, "Bearer ") {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
				c.Abort()
				return
			}
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
	}
}
func (ts *taskServer) authHandler(c *gin.Context) {
	type Provider struct {
		Provider string `form:"provider"`
	}

	var redirect_url string = "http://localhost:" + os.Getenv("SERVERPORT") + "/auth-callback"
	var p Provider
	var url = "https://accounts.google.com/o/oauth2/v2/auth?client_id=" +
		client_id +
		"&response_type=code" +
		"&scope=openid email profile phone" +
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

func (ts *taskServer) authCallbackGETHandler(c *gin.Context) {
	type GoogleResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   string `json:"expires_in"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
		IdToken     string `json:"id_token"`
	}

	q := c.Request.URL.Query().Get("code")
	origin := c.Request.URL.Query().Get("origin")
	//fmt.Println(q, "   ", reflect.TypeOf(q))
	url := "https://oauth2.googleapis.com/token"
	values := map[string]string{
		"code":          q,
		"client_id":     client_id,
		"client_secret": secret,
		"grant_type":    "authorization_code",
		"redirect_uri":  "http://localhost:4112/auth-callback?origin=" + origin}
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
	var gResp GoogleResponse
	body, err := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &gResp); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	//body =
	//jwt = resp.Body.
	c.Redirect(302, origin+"?token="+gResp.IdToken)
	//fmt.Println("response Body:", gResp.IdToken)
}

func (ts *taskServer) authCallbackPOSTHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}
func (ts *taskServer) infoHandler(c *gin.Context) {
	jwtToken := c.Request.Header["Authorization"][0]
	jwtToken = strings.Split(jwtToken, " ")[1]
	payload, _ := idtoken.Validate(context.Background(), jwtToken, client_id)
	givenName := strings.TrimSuffix(strings.Replace(fmt.Sprintf("&p",
		payload.Claims["given_name"]),
		"&p%!(EXTRA string=", "", 1), ")")
	currentTime := time.Now()
	c.JSON(http.StatusOK, gin.H{
		"date":       currentTime.Format("01-02-2006"),
		"authorName": givenName,
		"appName":    "auth"})
}
func (ts *taskServer) infoUnsafeHandler(c *gin.Context) {
	currentTime := time.Now()
	c.JSON(http.StatusOK, gin.H{
		"date":       currentTime.Format("01-02-2006"),
		"authorName": "name",
		"appName":    "auth"})
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
func (ts *taskServer) authJsonRedirect(c *gin.Context) {
	type Provider struct {
		Provider string `json:"provider"`
		Origin   string `json:"origin"`
	}
	var p Provider
	if err := c.BindJSON(&p); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	var redirect_url string = "http://localhost:" + os.Getenv("SERVERPORT") + "/auth-callback?origin=" + p.Origin

	var url = "https://accounts.google.com/o/oauth2/v2/auth?client_id=" +
		client_id +
		"&response_type=code" +
		"&scope=openid email profile phone" +
		"&redirect_uri=" +
		redirect_url

	if p.Provider == "google" {
		//CORS(c)
		c.JSON(http.StatusOK, gin.H{"link": url})
	}
}

func main() {
	router := gin.New()
	router.Use(CORSMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	server := NewTaskServer()
	router.GET("/info", tokenTypeValidation(), verifyIdToken(), server.infoHandler)
	router.GET("/infoUnsafe", server.infoUnsafeHandler)
	router.POST("/auth", server.authHandler)
	router.GET("/auth-callback", server.authCallbackGETHandler)
	router.POST("/auth-callback", server.authCallbackPOSTHandler)
	router.POST("/auth/provider-link", server.authJsonRedirect)

	router.Run("localhost:" + os.Getenv("SERVERPORT"))
}

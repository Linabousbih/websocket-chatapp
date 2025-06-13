package main

import (
	"chatapp/auth"
	"chatapp/routes"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

//ServerHTTP handles the HTTP request

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	t.templ.Execute(w, r)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// ✅ Generate and set JWT key
	key := auth.GenerateRandomKey()
	auth.SetJWTKey(key)
	log.Printf("Generated JWT Key: %s\n", key)

	// ✅ Create a new Gin router
	router := gin.Default()

	// ✅ Serve static files
	router.Static("/static", "./static")

	// ✅ HTML routes

	router.GET("/signup", func(c *gin.Context) {
		c.HTML(200, "signup.html", nil)
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})

	router.GET("/", auth.Authenticate(), func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.GET("/chat", auth.Authenticate(), func(c *gin.Context) {
		c.HTML(200, "chat.html", nil)
	})

	// ✅ WebSocket endpoint
	router.GET("/room", func(c *gin.Context) {
		roomName := c.Query("room")
		if roomName == "" {
			c.JSON(400, gin.H{"error": "Room name required"})
			return
		}
		realRoom := getRoom(roomName)
		realRoom.ServeHTTP(c.Writer, c.Request) // fallback to http.Handler
	})

	// ✅ API routes (mounts /api/signup, /api/login, etc.)
	routes.SetupRoutes(router.Group("/api"))

	router.LoadHTMLGlob("templates/*.html")

	// ✅ Start server
	log.Println("Server running at http://localhost:8080")
	router.Run(":8080")
}

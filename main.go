package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	// should be handled in its own routes world, but for now this will work as I learn.

	router := gin.Default()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/**/*")
	router.Static("/static", "static")

	router.GET("/", mainAppFunc)
	router.GET("/db", dbFunc)

	router.Run(":" + port)

}

func mainAppFunc(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/index.tmpl.html", nil)
}

func dbFunc(c *gin.Context) {
	var err error

	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)

	if err != nil {
		fmt.Printf("Error getting your IP: %q", err)
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("IP is %q", ip))

	// create table if it isn't there already.
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS pings (ip text, last_visited timestamp)"); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error storing your IP: %q", err))
		return
	}

	// insert IP
	// if _, err := db.Exec("INSERT INTO pings pings values()"); err != nil {
	// 	c.String(http.StatusInternalServerError, fmt.Sprintf("Error storing your IP: %q", err))
	// 	return
	// }

}

package main

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
)

type Config struct {
	AdminUserID		string	`json:"admin_user_id"`
	AdminPassword	string	`json:"admin_password"`
	AdminEmail		string	`json:"admin_email"`
}

func unmarshalSetting() {

}

func main() {
	// Unmarshal config.json
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	json.Unmarshal(file, &config)

	// gin start
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/login")
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{"title": "Squirrel - Login"})
	})

	r.GET("/admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin.html", gin.H{"title": "Squirrel - Admin"})
	})

	r.GET("/help", func(c *gin.Context) {
		c.HTML(http.StatusOK, "help.html", gin.H{"title": "Squirrel - Help", "adminEmail": config.AdminEmail})
		log.Printf("%s", config.AdminEmail)
	})

	r.Static("/assets", "./assets")

	r.Run(":8080")
}

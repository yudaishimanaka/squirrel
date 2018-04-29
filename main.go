package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	_"github.com/go-sql-driver/mysql"
)

type Config struct {
	AppConfig		App
	DbConfig	Database
}

type App struct {
	AdminUserID   string `json:"admin_user_id"`
	AdminPassword string `json:"admin_password"`
	AdminEmail    string `json:"admin_email"`
}

type Database struct {
	User		string `json:"user"`
	Password	string `json:"password"`
	DbName		string `json:"db_name"`
}

func mysqlConnect(user, password, dbName string) (engine *xorm.Engine, err error) {
	dataSourceName := user+":"+password+"@/"+dbName
	engine, err = xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	return engine, nil
}

func main() {
	// Unmarshal config.json
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	json.Unmarshal(file, &config)

	// MySQL connect
	engine, err := mysqlConnect(config.DbConfig.User, config.DbConfig.Password, config.DbConfig.DbName)
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Close()

	// Gin start
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
		c.HTML(http.StatusOK, "help.html", gin.H{"title": "Squirrel - Help", "adminEmail": config.AppConfig.AdminEmail})
	})

	r.Static("/assets", "./assets")

	r.Run(":8080")
}

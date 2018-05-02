package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type Config struct {
	AppConfig App      `json:"app"`
	DbConfig  Database `json:"db"`
}

type App struct {
	AdminUserID   string `json:"admin_user_id"`
	AdminPassword string `json:"admin_password"`
	AdminEmail    string `json:"admin_email"`
}

type Database struct {
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
}

type User struct {
	UserId		string	`xorm:"not null TEXT"`
	Password	string	`xorm:"not null TEXT"`
	Role		int		`xorm:"not null INT"`
}

type Image struct {
	Id			int
	Path		string
	Uploaded	string
}

func initDatabase(driver, user, password, dbName string, config Config) (e *xorm.Engine, err error) {
	engine, err := xorm.NewEngine(driver, user+":"+password+"@/")
	if err != nil {
		return nil, err
	}

	if _, err := engine.Exec("CREATE DATABASE "+dbName); err != nil {
		log.Printf("Database already exists.")
		return engine, nil
	} else {
		engine.Exec("USE "+dbName)
		engine.CreateTables(User{})
		engine.CreateTables(Image{})
		admin := User{
			UserId: config.AppConfig.AdminUserID,
			Password: config.AppConfig.AdminPassword,
			Role: 1, // Role 0:default, 1:admin
		}
		engine.Insert(admin)
		log.Printf("Success initialize.")

		return engine, nil
	}
}

func main() {
	// Unmarshal config.json
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	json.Unmarshal(file, &config)

	// Init database
	engine, err := initDatabase("mysql", config.DbConfig.User, config.DbConfig.Password, config.DbConfig.DbName, config)
	if err != nil {
		log.Fatal(engine)
	}

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

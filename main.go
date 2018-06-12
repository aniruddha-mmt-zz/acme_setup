package main

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/qor/admin"
	"github.com/qor/assetfs"
	"github.com/qor/qor"
	"github.com/qor/session"
)

// Create a GORM-backend model
type User struct {
	gorm.Model
	Name    string
	Address string
}

// Create another GORM-backend model
type Product struct {
	gorm.Model
	Name        string
	Description string
}

type AdminConfig struct {
	SiteName       string
	DB             *gorm.DB
	SessionManager session.ManagerInterface
	AssetFS        assetfs.Interface
}

func main() {
	DB, _ := gorm.Open("mysql", "root:root@/demo?charset=utf8&parseTime=True&loc=Local")
	DB.CreateTable(&User{}, &Product{})
	DB.AutoMigrate(&User{}, &Product{})

	// Initalization for APIs
	API := admin.New(&qor.Config{DB: DB})
	API.AddResource(&User{})

	// Initalize
	Admin := admin.New(&admin.AdminConfig{DB: DB})

	// Allow to use Admin to manage User, Product
	Admin.AddResource(&User{})
	Admin.AddResource(&Product{})

	// initalize an HTTP request multiplexer
	mux := http.NewServeMux()

	// Mount admin interface to mux
	Admin.MountTo("/admin", mux)

	// Mount API interface to mux
	API.MountTo("/api", mux)

	fmt.Println("Listening on: 9000")
	http.ListenAndServe(":9000", mux)
}

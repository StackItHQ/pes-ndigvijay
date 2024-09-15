package db

import (
	"fmt"
	"log"
	// "os"
	"sync"

	// "github.com/StackItHQ/pes-ndigvijay/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	freedb "github.com/FreeLeh/GoFreeDB"
	"github.com/FreeLeh/GoFreeDB/google/auth"
)

var DBLock sync.Mutex

func InitDB() *gorm.DB {
	DBLock.Lock()
	defer DBLock.Unlock()
	host := "localhost"
	user := "ndv"
	password := "ndv"
	dbname := "mydatabase"
	port := 5432

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&models.User{})
	return db
}

var Store *freedb.GoogleSheetRowStore

func InitSheet() {
	authenticator, err := auth.NewServiceFromFile(
		"../servicetoken.json",
		freedb.FreeDBGoogleAuthScopes,
		auth.ServiceConfig{},
	)
	if err != nil {
		fmt.Println("Error accessing sheet", err)
		return
	}

	Store = freedb.NewGoogleSheetRowStore(
		authenticator,
		"1Bxi8B75FyzeLhCfTnk3QXoX4zMeWNdGa4Q_XbRjuQNw",
		"Sheet1",
		freedb.GoogleSheetRowStoreConfig{
			Columns: []string{
				"id",        
				"name",      
				"email",     
				"password",  
				"age",         
			},
		},
	)
	
}
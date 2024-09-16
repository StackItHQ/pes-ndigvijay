package db

import (
	"fmt"
	"log"
	// "os"
	"sync"
	// "context"
	"github.com/StackItHQ/pes-ndigvijay/backend/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// freedb "github.com/FreeLeh/GoFreeDB"
	// "github.com/FreeLeh/GoFreeDB/google/auth"
	// _ "github.com/lib/pq"
	// "google.golang.org/api/option"
	// "google.golang.org/api/sheets/v4"
)

const (
	credentialsFile = "/Users/digvijaynarayan/Desktop/superjoin/pes-ndigvijay/backend/servicetoken.json"
	spreadsheetId    = "1Bxi8B75FyzeLhCfTnk3QXoX4zMeWNdGa4Q_XbRjuQNw"
	sheetRange       = "Sheet1!A1:D10" 
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

// var Store *freedb.GoogleSheetRowStore

// func InitSheet() {
	// authenticator, err := auth.NewServiceFromFile(
	// 	"../servicetoken.json",
	// 	freedb.FreeDBGoogleAuthScopes,
	// 	auth.ServiceConfig{},
	// )
	// if err != nil {
	// 	fmt.Println("Error accessing sheet", err)
	// 	return
	// }

	// Store = freedb.NewGoogleSheetRowStore(
	// 	authenticator,
	// 	"1Bxi8B75FyzeLhCfTnk3QXoX4zMeWNdGa4Q_XbRjuQNw",
	// 	"Sheet1",
	// 	freedb.GoogleSheetRowStoreConfig{
	// 		Columns: []string{     
	// 			"name",      
	// 			"email",     
	// 			"password",         
	// 		},
	// 	},
	// )
// }


// var SheetsService *sheets.Service

// func InitSheet() {
//     ctx := context.Background()

//     // Initialize Google Sheets API client
//     srv, err := sheets.NewService(ctx, option.WithCredentialsFile(credentialsFile))
//     if err != nil {
//         log.Fatalf("Unable to create Sheets service: %v", err)
//     }
//     SheetsService = srv

//     // Verify the "name", "email", and "password" columns
//     verifyColumns(ctx)
// }

// func verifyColumns(ctx context.Context) {
//     readRange := "Sheet1!A1:D1" // Adjust the range to read the header row
//     resp, err := SheetsService.Spreadsheets.Values.Get(spreadsheetId, readRange).Context(ctx).Do()
//     if err != nil {
//         log.Fatalf("Unable to read sheet: %v", err)
//     }

//     requiredColumns := map[string]bool{
//         "name":     false,
//         "email":    false,
//         "password": false,
//     }

//     if len(resp.Values) > 0 {
//         headers := resp.Values[0]
//         for _, header := range headers {
//             if column, ok := header.(string); ok {
//                 if _, exists := requiredColumns[column]; exists {
//                     requiredColumns[column] = true
//                 }
//             }
//         }

//         for column, present := range requiredColumns {
//             if present {
//                 log.Printf("Column %v is present", column)
//             } else {
//                 log.Printf("Column %v is missing", column)
//             }
//         }
//     } else {
//         log.Println("No headers found")
//     }
// }
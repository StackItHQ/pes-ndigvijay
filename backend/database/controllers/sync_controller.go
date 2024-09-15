package controllers

import (
    "context"
    "log"
    "sync"
    "time"
    "github.com/StackItHQ/pes-ndigvijay/backend/database/models"
    "google.golang.org/api/option"
    "google.golang.org/api/sheets/v4"
    "gorm.io/gorm"
)

const credentialsFile = "/Users/digvijaynarayan/Desktop/superjoin/pes-ndigvijay/backend/servicetoken.json"

// Function to handle syncing Google Sheets and PostgreSQL
func StartSync(gormDB *gorm.DB, stopChan chan struct{}, wg *sync.WaitGroup) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Recovered from panic: %v", r)
        }
    }()

    ctx := context.Background()
    srv, err := sheets.NewService(ctx, option.WithCredentialsFile(credentialsFile))
    if err != nil {
        log.Fatalf("Unable to retrieve Sheets client: %v", err)
    }

    var lastData [][]interface{}

    // PostgreSQL to Sheets operation
    go func() {
        ticker := time.NewTicker(5 * time.Second)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                log.Println("Updating Google Sheet with PostgreSQL data...")
                var users []models.User
                if err := gormDB.Find(&users).Error; err != nil {
                    log.Printf("Error retrieving users from DB: %v", err)
                    continue
                }

                // Prepare data to update Google Sheets
                var sheetData [][]interface{}
                sheetData = append(sheetData, []interface{}{"Name", "Email", "Password"})
                for _, user := range users {
                    sheetData = append(sheetData, []interface{}{user.Name, user.Email, user.Password})
                }

                if err := UpdateGoogleSheet(srv, sheetData); err != nil {
                    log.Printf("Error updating Google Sheet: %v", err)
                } else {
                    log.Println("Successfully updated Google Sheet with PostgreSQL data.")
                }
            case <-stopChan:
                return
            }
        }
    }()

    // Sheets to PostgreSQL operation
    wg.Add(1)
    go func() {
        defer wg.Done()
        ticker := time.NewTicker(5 * time.Second)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                log.Println("Updating PostgreSQL from Google Sheet data...")
                currentData, err := GetSheetData(srv)
                if err != nil {
                    log.Printf("Error retrieving data from Google Sheet: %v", err)
                } else {
                    if !EqualData(lastData, currentData) {
                        log.Println("Sheet data has changed. Updating database...")
                        if err := UpdateDatabase(gormDB, currentData); err != nil {
                            log.Printf("Error updating database: %v", err)
                        }
                        lastData = currentData
                    }
                }
            case <-stopChan:
                return
            }
        }
    }()
}

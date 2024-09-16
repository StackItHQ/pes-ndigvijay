package main

import (
    "context"
    "fmt"
    "log"
    "sync"
    "time"

    "github.com/StackItHQ/pes-ndigvijay/backend/database/db"
    "github.com/StackItHQ/pes-ndigvijay/backend/database/models"
    "github.com/gin-gonic/gin"
    "google.golang.org/api/option"
    "google.golang.org/api/sheets/v4"
    "gorm.io/gorm"
)

const (
	credentialsFile = "/Users/digvijaynarayan/Desktop/superjoin/pes-ndigvijay/backend/servicetoken.json"
	spreadsheetId    = "1Bxi8B75FyzeLhCfTnk3QXoX4zMeWNdGa4Q_XbRjuQNw"
	sheetRange       = "Sheet1!A1:D10" 
)

func main() {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Recovered from panic: %v", r)
        }
    }()

    gormDB := db.InitDB()
    ctx := context.Background()
    srv, err := sheets.NewService(ctx, option.WithCredentialsFile(credentialsFile))
    if err != nil {
        log.Fatalf("Unable to retrieve Sheets client: %v", err)
    }

    var lastData [][]interface{}
    // dataChan := make(chan [][]interface{})
    // errorChan := make(chan error)
    stopChan := make(chan struct{})
    var wg sync.WaitGroup

    // Goroutine for PostgreSQL to Sheets operation
    go func() {
        ticker := time.NewTicker(5 * time.Second)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                // Update Google Sheets from PostgreSQL
                fmt.Println("Updating Google Sheet with PostgreSQL data...")
                err := updateGoogleSheet(srv, gormDB)
                if err != nil {
                    log.Printf("Error updating Google Sheet: %v", err)
                } else {
                    log.Println("Successfully updated Google Sheet with PostgreSQL data.")
                }
            case <-stopChan:
                return
            }
        }
    }()

    // Goroutine for Sheets to PostgreSQL operation
    wg.Add(1)
    go func() {
        defer wg.Done()
        ticker := time.NewTicker(5 * time.Second)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                // Update PostgreSQL from Google Sheets
                fmt.Println("Updating PostgreSQL from Google Sheet data...")
                currentData, err := getSheetData(srv)
                if err != nil {
                    log.Printf("Error retrieving data from Google Sheet: %v", err)
                } else {
                    if !equalData(lastData, currentData) {
                        log.Println("Sheet data has changed. Updating database...")
                        if err := updateDatabase(gormDB, currentData); err != nil {
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

    router := gin.Default()

    go func() {
        if err := router.Run(":8080"); err != nil {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    // Wait for shutdown signal
    for {
        select {
        case <-stopChan:
            fmt.Println("Shutting down...")
            wg.Wait()
            return
        }
    }
}

// Function to get data from Google Sheets
func getSheetData(srv *sheets.Service) ([][]interface{}, error) {
    resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, sheetRange).Context(context.Background()).Do()
    if err != nil {
        return nil, err
    }
    return resp.Values, nil
}

// Function to update PostgreSQL from Google Sheets data
func updateDatabase(db *gorm.DB, data [][]interface{}) error {
    if len(data) == 0 {
        return nil
    }
    for i, row := range data {
        if i == 0 { // Skip header row
            continue
        }
        // Check if the row has the required number of columns
        if len(row) < 3 {
            log.Printf("Row %d does not have enough columns: %v", i, row)
            continue
        }

        name, ok := row[0].(string)
        if !ok {
            log.Printf("Invalid name in row %d: %v", i, row[0])
            continue
        }

        email, ok := row[1].(string)
        if !ok {
            log.Printf("Invalid email in row %d: %v", i, row[1])
            continue
        }

        password, ok := row[2].(string)
        if !ok {
            log.Printf("Invalid password in row %d: %v", i, row[2])
            continue
        }

        user := models.User{Name: name, Email: email, Password: password}
        result := db.Where("email = ?", user.Email).Assign(user).FirstOrCreate(&user)
        if result.Error != nil {
            log.Printf("Error upserting user %v: %v", user.Email, result.Error)
            continue
        }
    }
    return nil
}


// Function to update Google Sheets from PostgreSQL
func updateGoogleSheet(srv *sheets.Service, db *gorm.DB) error {
    var users []models.User
    if err := db.Find(&users).Error; err != nil {
        return err
    }

    var sheetData [][]interface{}
    sheetData = append(sheetData, []interface{}{"Name", "Email", "Password"})
    for _, user := range users {
        sheetData = append(sheetData, []interface{}{user.Name, user.Email, user.Password})
    }

    valueRange := &sheets.ValueRange{
        Values: sheetData,
    }

    _, err := srv.Spreadsheets.Values.Update(spreadsheetId, sheetRange, valueRange).ValueInputOption("RAW").Do()
    return err
}

// Helper function to compare sheet data
func equalData(a, b [][]interface{}) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if len(a[i]) != len(b[i]) {
            return false
        }
        for j := range a[i] {
            if a[i][j] != b[i][j] {
                return false
            }
        }
    }
    return true
}
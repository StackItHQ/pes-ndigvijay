package controllers

import (
    "log"
    "github.com/StackItHQ/pes-ndigvijay/backend/database/models"
    "gorm.io/gorm"
)

// Function to update PostgreSQL from Google Sheets data
func UpdateDatabase(db *gorm.DB, data [][]interface{}) error {
    if len(data) == 0 {
        return nil
    }
    for i, row := range data {
        if i == 0 { // Skip header row
            continue
        }
        if len(row) < 3 { // Ensure there are at least 3 columns
            log.Printf("Row %d does not have enough columns, skipping", i)
            continue
        }

        name, nameOk := row[0].(string)
        email, emailOk := row[1].(string)
        password, passwordOk := row[2].(string)

        if !nameOk || !emailOk || !passwordOk {
            log.Printf("Row %d contains invalid data types, skipping", i)
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



package main

import (
    "fmt"
    "log"
    "sync"

    "github.com/StackItHQ/pes-ndigvijay/backend/database/db"
    "github.com/StackItHQ/pes-ndigvijay/backend/database/controllers"
    "github.com/gin-gonic/gin"
)

func main() {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Recovered from panic: %v", r)
        }
    }()

    gormDB := db.InitDB()
    stopChan := make(chan struct{})
    var wg sync.WaitGroup

    // Start syncing operations
    controllers.StartSync(gormDB, stopChan, &wg)

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

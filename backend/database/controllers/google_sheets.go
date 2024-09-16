package controllers

import (
    "context"
    // "log"
	// "google.golang.org/api/option"
    "google.golang.org/api/sheets/v4"
)

const (
    spreadsheetId = "1Bxi8B75FyzeLhCfTnk3QXoX4zMeWNdGa4Q_XbRjuQNw"
    sheetRange    = "Sheet1!A1:D10"
)

// Function to get data from Google Sheets
func GetSheetData(srv *sheets.Service) ([][]interface{}, error) {
    resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, sheetRange).Context(context.Background()).Do()
    if err != nil {
        return nil, err
    }
    return resp.Values, nil
}

// Function to update Google Sheets from PostgreSQL
func UpdateGoogleSheet(srv *sheets.Service, sheetData [][]interface{}) error {
    valueRange := &sheets.ValueRange{
        Values: sheetData,
    }

    _, err := srv.Spreadsheets.Values.Update(spreadsheetId, sheetRange, valueRange).ValueInputOption("RAW").Do()
    return err
}

// Helper function to compare sheet data
func EqualData(a, b [][]interface{}) bool {
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

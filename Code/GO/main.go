package main

import (
	"fmt"
	"log"

	"golang.org/x/oauth2"
	"google.golang.org/api/sheets/v4"
)

func main() {
    // Create a JWT configurations object for the Google service account
	conf, err := ReadConfig("client_secret.json")

	client := conf.Client(oauth2.NoContext)

	// Create a service object for Google sheets
	srv, err := sheets.New(client)
        if err != nil {
                log.Fatalf("Unable to retrieve Sheets client: %v", err)
        }

        spreadsheetId := "1BSp8k8HODJIXvjcE4ZTOcoDKzgsaqltHccRI4vQVIIY"
        readRange := "Sheet1!A2:E"
        resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
        if err != nil {
                log.Fatalf("Unable to retrieve data from sheet: %v", err)
        }

        if len(resp.Values) == 0 {
                fmt.Println("No data found.")
        } else {
                // fmt.Println("Name, Major:")
                // for _, row := range resp.Values {
                //         // Print columns A and E, which correspond to indices 0 and 4.
                //         fmt.Printf("%s, %s\n", row[0], row[4])
                // }
        }


        ip := "172.31.233.148"
	cpuLoad, err := GetCPULoad(ip)
	if err != nil {
		log.Fatalf("Error getting CPU load: %v", err)
	}

	fmt.Printf("CPU Load: %.2f%%\n", cpuLoad)
}
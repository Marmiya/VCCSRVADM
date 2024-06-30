package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"google.golang.org/api/sheets/v4"
)

func main() {
        // Create a JWT configurations object for the Google service account
        test := true

        /*
          Get the mode from command line arguments
          -1: error
          0: Send Report script to the servers
          1: run
        */
        mod := flag.Int("mode", -1, "mode")
        flag.Parse()

        fmt.Printf("mode: %d\n", *mod)

        conf, err := ReadConfig("client_secret.json")
        if err != nil {
                log.Fatalf("Unable to read client secret file: %v", err)
        }

        client := conf.Client(oauth2.NoContext)

        // Create a service object for Google sheets
        srv, err := sheets.New(client)
        if err != nil {
                log.Fatalf("Unable to retrieve Sheets client: %v", err)
        }

        spreadsheetId := "1BSp8k8HODJIXvjcE4ZTOcoDKzgsaqltHccRI4vQVIIY"
        readRange := "Sheet1!D37:D37"
        resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
        if err != nil {
                log.Fatalf("Unable to retrieve data from sheet: %v", err)
        }

        if len(resp.Values) == 0 {
                fmt.Println("No data found.")
        } else {
                if test {
                        fmt.Println("ip:")
                        for _, row := range resp.Values {
                                fmt.Printf("%s\n", row[0])
                        }
                }
        }

        if *mod == -1 {
                log.Fatalf("Invalid mode")
        } else if *mod == 0 {
                for _, ip := range resp.Values {
                        // There are 2 possible name:{ubuntu, szu} on the server, 
                        // so we will try and find the correct one
                        success := false
                        for _, name := range []string{"ubuntu", "szu"} {
                                err := SendFile(ip[0].(string), name, "../ReportInfo.py")
                                if err == nil {
                                        success = true
                                        break
                                }
                        }
                        if !success {
                                // panic
                                log.Fatalf("Unable to send file to %s", ip[0])
                        }
                }
        }
}


type ResponseData struct {
	Message   string  `json:"message"`
	Timestamp float64 `json:"timestamp"`
}


func getJSONData(ip string) (*ResponseData, error) {
	url := fmt.Sprintf("http://%s:9527/json", ip)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-200 status: %v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var data ResponseData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return &data, nil
}
package main

import (
	"os"
	"net/http"
	"fmt"
	"time"
	"encoding/json"
	"io/ioutil"
	"encoding/csv"
)

//const serverPort = 3333

// Create structs for the data
type People struct {
	People []Person `json:"people"`
	Number int `json:"number"`
	Message string `json:"message"`
}

type Person struct {
	Name string `json:"name"`
	Craft string `json:"craft"`
}

func isLarger(p1 Person, p2 Person) bool {
	// Defines order on for Person constructs
	if p1.Craft > p2.Craft {
		return true
	} else if (p1.Craft == p2.Craft) && (p1.Name > p2.Name) {
		return true
	}
	return false
}

func main() {
	// Get url
	var url = ""
	if len(os.Args) < 2 {
		url = "http://api.open-notify.org/astros.json"
	} else {
		url = os.Args[1]
	}

	// Create client and fetch API
	var client = &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(errors.New("Error fetching API."))
	}
	defer resp.Body.Close()

	// Parse json file into variable people
	byteValue, _ := ioutil.ReadAll(resp.Body)
	var people People
	json.Unmarshal(byteValue, &people)

	// Sort people.People (I used a simple bubble sort)
	var sorted = false
	for !sorted {
		sorted = true
		for i := 0; i < len(people.People)-1; i++ {
			if isLarger(people.People[i], people.People[i+1]) {
				sorted = false
				var temp = people.People[i]
				people.People[i] = people.People[i+1]
				people.People[i+1] = temp
			}
		}
	} 
	
	// Create CSV file
	csvFile, err := os.Create("output.csv")
	if err != nil {
		panic(errors.New("Error creating CSV."))
	}

	// Write data to CSV file
	csvwriter := csv.NewWriter(csvFile)
	header := []string{"Name", "Craft"}
	_ = csvwriter.Write(header)
	for _, person := range people.People {
		data := []string{person.Name, person.Craft}
		_ = csvwriter.Write(data)
	}
	csvwriter.Flush()
	csvFile.Close()

	fmt.Println("Successful")
}
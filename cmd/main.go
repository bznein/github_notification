package main

import (
	"fmt"
	"github.com/andybrewer/mack"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	token := os.Getenv("NOTIFICATION_TOKEN")
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.github.com", nil)
	// ...
	req.Header.Add("Authorization", "token "+token)
	_, err = client.Do(req)
	// ...
	if err != nil {
		log.Fatal(err)
	}
	req, err = http.NewRequest("GET", "https://api.github.com/notifications?all=true", nil)
	req.Header.Add("Authorization", "token "+token)
	response, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(responseData))
	// do stuff
	mack.Notify("Complete")
}

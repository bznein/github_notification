package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/andybrewer/mack"
)

type ApiResponse struct {
	Id string `json:"id"`
}

func main() {
	closeChan := make(chan bool)
	go getNotifications(closeChan)
	<-closeChan
}

func getNotifications(closeChan chan bool) {

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
	req, err = http.NewRequest("GET", "https://api.github.com/notifications", nil)
	req.Header.Add("Authorization", "token "+token)

	for {
		mack.Notify("Excuting")
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
		mack.Notify(string(responseData))
		time.Sleep(time.Second * 60)
	}

}

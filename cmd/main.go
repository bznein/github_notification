package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bznein/github_notification/pkg/notification"
	"github.com/bznein/notipher/pkg/notiphication"
)

func main() {
	closeChan := make(chan bool)
	go getNotifications(closeChan)
	<-closeChan
}

// Max returns the larger of x or y.
func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func getNotifications(closeChan chan bool) {

	token := os.Getenv("NOTIFICATION_TOKEN")
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.github.com", nil)
	req.Header.Add("Authorization", "token "+token)
	_, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	pollTime := 60
	etag := ""
	req, err = http.NewRequest("GET", "https://api.github.com/notifications", nil)
	req.Header.Add("Authorization", "token "+token)

	for {
		req.Header.Set("If-None-Match", etag)
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		etag = response.Header.Get("ETag")

		status := response.StatusCode
		switch status {
		case 200:
		case 304:
			timeT, _ := strconv.Atoi(response.Header.Get("X-Poll-Interval"))
			pollTime = max(pollTime, timeT)
			time.Sleep(time.Second * time.Duration(pollTime))
			continue
		default:
			fmt.Fprintf(os.Stderr, "Error code not 200 or 304: %d", status)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		res := notification.NotificationResponse{}
		err = json.Unmarshal(responseData, &res)
		if err != nil {
			log.Fatal(err)
		}

		newReq, err := http.NewRequest("GET", res[0].Subject.Url, nil)
		newReq.Header.Add("Authorization", "token "+token)
		newResponse, err := client.Do(newReq)
		issue := notification.Issue{}

		newIssue, err := ioutil.ReadAll(newResponse.Body)
		err = json.Unmarshal(newIssue, &issue)
		if err != nil {
			log.Fatal(err)
		}
		for _, notification := range res {
			n := notiphication.Notiphication{}
			n.AppIcon = "./assets/GitHub-Mark-32px.png"
			n.Title = notification.Repository.Description
			n.Subtitle = notification.Subject.Title
			n.Link = issue.HtmlUrl
			n.DropdownLabel = "Remind me"
			actions := notiphication.Actions{}
			n.Actions = actions
			actions["5 Minutes"] = func() { go resendNotification(n, time.Minute*5) }
			actions["10 Minutes"] = func() { go resendNotification(n, time.Minute*10) }
			actions["15 Minutes"] = func() { go resendNotification(n, time.Minute*15) }
			n.AsyncPush()
		}

		time.Sleep(time.Second * time.Duration(pollTime))
	}

}

func resendNotification(n notiphication.Notiphication, sleep time.Duration) {
	time.Sleep(sleep)
	n.AsyncPush()
}

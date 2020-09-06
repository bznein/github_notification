package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"time"

	"github.com/bznein/github_notification/pkg/notification"
	"github.com/bznein/notipher/pkg/notiphication"
)

type Configuration struct {
	GithubToken    string   `json:"github_token"`
	RetryInterval1 int      `json:"retry_interval_1"`
	RetryInterval2 int      `json:"retry_interval_2"`
	RetryInterval3 int      `json:"retry_interval_3"`
	IgnoreList     []string `json:"ignore_list"`
}

func main() {
	closeChan := make(chan bool)
	config := getConfig()
	go getNotifications(closeChan, config)
	<-closeChan
}

func getDefaultNotification() Configuration {
	token, ok := os.LookupEnv("NOTIFICATION_TOKEN")
	if !ok {
		log.Fatal("Errors in loading custom config and can't find NOTIFICATION_TOKEN, aborting")
	}
	return Configuration{GithubToken: token, RetryInterval1: 5, RetryInterval2: 10, RetryInterval3: 15}
}

func getConfig() Configuration {
	user, err := user.Current()
	if err != nil {
		return getDefaultNotification()
	}
	jsonFile, err := os.Open(user.HomeDir + "/.config/github_notifications/config.json")
	if err != nil {
		return getDefaultNotification()
	}
	defer jsonFile.Close()
	configFilesBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return getDefaultNotification()
	}
	config := Configuration{}
	err = json.Unmarshal(configFilesBytes, &config)
	if err != nil {
		return getDefaultNotification()
	}
	return config
}

// Max returns the larger of x or y.
func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func contains(container []string, elem string) bool {
	for _, e := range container {
		if e == elem {
			return true
		}
	}
	return false
}

func getNotifications(closeChan chan bool, config Configuration) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.github.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "token "+config.GithubToken)
	_, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	pollTime := 60
	etag := ""
	req, err = http.NewRequest("GET", "https://api.github.com/notifications", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "token "+config.GithubToken)

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
		issue := notification.Issue{}
		err = json.Unmarshal(responseData, &res)
		if err != nil {
			log.Fatal(err)
		}

		for _, notification := range res {
			if contains(config.IgnoreList, notification.Reason) {
				continue
			}
			n := notiphication.Notiphication{}
			if notification.Subject.Url != "" {
				newReq, err := http.NewRequest("GET", notification.Subject.Url, nil)
				newReq.Header.Add("Authorization", "token "+config.GithubToken)
				newResponse, err := client.Do(newReq)

				newIssue, err := ioutil.ReadAll(newResponse.Body)
				err = json.Unmarshal(newIssue, &issue)
				if err != nil {
					log.Fatal(err)
				}
				n.Link = issue.HtmlUrl
			}
			n.AppIcon = "./assets/GitHub-Mark-32px.png"
			n.Title = notification.Repository.Description
			n.Subtitle = notification.Subject.Title
			n.DropdownLabel = "Options"
			actions := notiphication.Actions{}
			n.Actions = actions
			actions["5 Minutes"] = func() { go resendNotification(n, time.Minute*time.Duration(config.RetryInterval1)) }
			actions["10 Minutes"] = func() { go resendNotification(n, time.Minute*time.Duration(config.RetryInterval2)) }
			actions["15 Minutes"] = func() { go resendNotification(n, time.Minute*time.Duration(config.RetryInterval3)) }
			actions["Ignore this type of notification"] = func() { ignoreNotification(notification.Reason) }
			n.AsyncPush()
		}

		time.Sleep(time.Second * time.Duration(pollTime))
	}

}

func ignoreNotification(reason string) {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	configPath := user.HomeDir + "/.config/github_notifications/config.json"
	jsonFile, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	configFilesBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	config := Configuration{}
	err = json.Unmarshal(configFilesBytes, &config)
	if err != nil {
		log.Fatal(err)
	}
	config.IgnoreList = append(config.IgnoreList, reason)
	editedConfig, err := json.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(configPath, editedConfig, 0644)
}

func resendNotification(n notiphication.Notiphication, sleep time.Duration) {
	time.Sleep(sleep)
	n.AsyncPush()
}

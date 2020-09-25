package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"
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
	LogDir         string   `json:"log_dir"`
}

var logger *log.Logger

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
	user, err := user.Current()
	if err != nil {
		log.Fatal("Tried to load default config but can't read user homedir for default logDir")
	}
	return Configuration{GithubToken: token, RetryInterval1: 5, RetryInterval2: 10, RetryInterval3: 15, IgnoreList: []string{}, LogDir: user.HomeDir + "/gn_logs"}
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
	return mergeWithDefault(config)
}

func mergeWithDefault(override Configuration) Configuration {
	defaultC := getDefaultNotification()
	if override.GithubToken != "" {
		defaultC.GithubToken = override.GithubToken
	}
	if override.RetryInterval1 > 0 {
		defaultC.RetryInterval1 = override.RetryInterval1
	}
	if override.RetryInterval2 > 0 {
		defaultC.RetryInterval2 = override.RetryInterval2
	}
	if override.RetryInterval3 > 0 {
		defaultC.RetryInterval3 = override.RetryInterval3
	}
	// No merging as the default is, for now, empty
	defaultC.IgnoreList = override.IgnoreList

	if override.LogDir != "" {
		defaultC.LogDir = override.LogDir
	}
	return defaultC
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
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile(strings.Replace(config.LogDir, "~", user.HomeDir, 1)+"/"+time.Now().Format("20060102"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if logger == nil {
		logger = log.New(f, "", log.LstdFlags)
	}
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
			logger.Printf("Error in executing request: %s", err)
			continue
		}

		etag = response.Header.Get("ETag")

		status := response.StatusCode
		switch status {
		case 200:
		case 304:
			timeT, _ := strconv.Atoi(response.Header.Get("X-Poll-Interval"))
			logger.Printf("Received 304 response\n")
			pollTime = max(pollTime, timeT)
			time.Sleep(time.Second * time.Duration(pollTime))
			continue
		default:
			logger.Printf("ERROR: code not 200 or 304, but was %d", status)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logger.Printf("Error in reading response body: %s", err)
			continue
		}
		res := notification.NotificationResponse{}
		issue := notification.Issue{}
		err = json.Unmarshal(responseData, &res)
		if err != nil {
			logger.Printf("Error in unmarshaling response: %s", err)
			continue
		}

		logger.Printf("Received code 200 with the following body: %+v", res)
		for _, notification := range res {
			if contains(config.IgnoreList, notification.Reason) {
				logger.Printf("Ignoring %s because is contained in ignorelist %+v", notification.Reason, config.IgnoreList)
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
					logger.Printf("Error in unmarshaling issue: %s", err)
					continue
				}
				n.Link = issue.HtmlUrl
				logger.Printf("Setting notification link as %s", n.Link)
			}
			n.AppIcon = "./assets/GitHub-Mark-32px.png"
			n.Title = notification.Repository.Description
			n.Subtitle = notification.Subject.Title
			n.DropdownLabel = "Options"
			actions := notiphication.Actions{}
			n.Actions = actions
			actions["Remind in 5 Minutes"] = func() {
				logger.Printf("First remind option clicked")
				go resendNotification(n, time.Minute*time.Duration(config.RetryInterval1))
			}
			actions["Remind in 10 Minutes"] = func() {
				logger.Printf("Second remind option clicked")
				go resendNotification(n, time.Minute*time.Duration(config.RetryInterval2))
			}
			actions["Remind in 15 Minutes"] = func() {
				logger.Printf("Third remind option clicked")
				go resendNotification(n, time.Minute*time.Duration(config.RetryInterval3))
			}
			actions["Ignore this type of notification"] = func() {
				logger.Printf("Ignoring notification: %+v", notification.Reason)
				ignoreNotification(notification.Reason)
			}
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
	logger.Printf("%s time passed, resending notification", sleep)
	n.AsyncPush()
}

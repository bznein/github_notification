package notification

import (
	"fmt"
	"log"
	"os/exec"
)

type Notification struct {
	Title   string
	Link    string
	Actions []string
}

func (n Notification) Push() string {

	fmt.Printf("Command: %+v\n", "alerter -message "+n.Title+" -actions "+n.Actions[0])
	cmd := exec.Command("alerter", "-message", n.Title, "-actions", n.Actions[0])
	fmt.Println(cmd)
	response, err := cmd.Output()
	if err != nil {

		log.Fatal(err)
	}
	return string(response)
}

//./alerter -message "Deploy now on UAT ?" -actions Now,"Later today","Tomorrow" -dropdownLabel "When ?"

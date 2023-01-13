package alert

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func SendWLEDPulse() error {
	//First get the current WLED settings
	currentWled := getWLEDSettings()

	workingDir, _ := os.Getwd()
	fmt.Printf(workingDir)

	appRoot := os.Getenv("APP_ROOT")
	fullPath := appRoot + "/config/wled_yellow_alert_post.json"
	configFile, err := os.Open(fullPath)
	if err != nil {
		return err
	}

	byteValue, _ := ioutil.ReadAll(configFile)
	alertPulseString := string(byteValue)

	alertLength, err := strconv.Atoi(os.Getenv("ALERT_LENGTH_SECONDS"))
	if err != nil {
		fmt.Printf("ALERT_LENGTH_SECONDS environment variable not found, or not an integer. Defaulting to 5 seconds.\n")
		alertLength = 5
	}
	sendWLEDCommand(alertPulseString)
	//Wait a number of seconds and return the lights to their prior state.
	time.Sleep(time.Duration(alertLength) * time.Second)
	sendWLEDCommand(currentWled)
	return nil
}

func sendWLEDCommand(bodyString string) {
	base_url := os.Getenv("WLED_IP")
	url := base_url + "/json/state"

	var jsonprep = []byte(bodyString)

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonprep))
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()
}

func getWLEDSettings() string {
	base_url := os.Getenv("WLED_IP")
	url := base_url + "/json/state"

	response, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	responseString := string(b)

	return responseString
}

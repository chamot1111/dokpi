package dokpi

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
)

func sendCmd(cmd, appName string) (status []appStatus, err error) {
	resp, err := http.Get(listenBaseURL + cmd + "&name=" + html.EscapeString(appName))
	if err != nil {
		err = fmt.Errorf("Connection problem: %s", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var msgReceived msg
	err = json.Unmarshal(body, &msgReceived)
	if err != nil {
		err = fmt.Errorf("Body format error: %s", err)
		return
	}

	status = msgReceived.appsStatus

	if !msgReceived.cmdSuccess {
		err = fmt.Errorf("Command fail: %s", msgReceived.errorMsg)
		return
	}
	return
}

func sendStart(appName string) {
	sendCmd(listenStart, appName)
}

func sendStop(appName string) {
	sendCmd(listenStop, appName)
}

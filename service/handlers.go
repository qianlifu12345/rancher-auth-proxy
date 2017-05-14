package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/rancher-auth-filter-service/manager"
)

//RequestData is for the JSON output
type RequestData struct {
	Headers map[string][]string `json:"headers,omitempty"`
	// Body    map[string]interface{} `json:"body,omitempty"`
	EnvID string `json:"envID,omitempty"`
}

//AuthorizeData is for the JSON output
type AuthorizeData struct {
	Message string `json:"message,omitempty"`
}

//MessageData is for the JSON output
type MessageData struct {
	Data []interface{} `json:"data,omitempty"`
}

//ProxyError structure contains the error resource definition
type ProxyError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

//ValidationHandler is a handler for cookie token and returns the request headers and accountid and projectid
func ValidationHandler(w http.ResponseWriter, r *http.Request) {

	reqestData := RequestData{}
	input, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Info(err)
		logrus.Infof("Cannot extract the request")
		w.WriteHeader(400)
		return
	}
	praseCookieErr := json.Unmarshal(input, &reqestData)
	if praseCookieErr != nil {
		logrus.Info(praseCookieErr)
		logrus.Infof("Cannot parse the request.")
		w.WriteHeader(400)
		return
	}
	envid := reqestData.EnvID

	var cookie []string
	if len(reqestData.Headers["Cookie"]) >= 1 {
		cookie = reqestData.Headers["Cookie"]
	} else {
		logrus.Infof("No Cookie found.")
		w.WriteHeader(http.StatusOK)
		return
	}
	var cookieString string
	if len(cookie) >= 1 {

		for i := range cookie {
			if strings.Contains(cookie[i], "token") {
				cookieString = cookie[i]
			}
		}

	} else {
		logrus.Infof("No token found in cookie.")
		w.WriteHeader(http.StatusOK)
		return
	}

	tokens := strings.Split(cookieString, ";")
	tokenValue := ""
	if len(tokens) >= 1 {
		for i := range tokens {
			if strings.Contains(tokens[i], "token") {
				if len(strings.Split(tokens[i], "=")) > 1 {
					tokenValue = strings.Split(tokens[i], "=")[1]
				}
			}

		}
	} else {
		logrus.Errorf("No token found")
		ReturnHTTPError(w, r, 200, fmt.Sprintf("No token found"))
		return
	}
	if tokenValue == "" {
		logrus.Errorf("No token found")
		ReturnHTTPError(w, r, 200, fmt.Sprintf("No token found"))
		return
	}

	//check if the token value is empty or not
	if tokenValue != "" {
		logrus.Infof("token:" + tokenValue)
		logrus.Infof("envid:" + envid)
		projectID, accountID := "", ""
		if envid != "" {
			projectID, accountID, err = getAccountAndProject(manager.URL, envid, tokenValue)
			if accountID == "Unauthorized" {
				logrus.Errorf("Unauthorized")
				ReturnHTTPError(w, r, 401, fmt.Sprintf("Unauthorized"))
				return
			}

			if accountID == "Forbidden" {
				logrus.Errorf("Forbidden")
				ReturnHTTPError(w, r, 403, fmt.Sprintf("Forbidden"))
				return
			}
			if err != nil {
				logrus.Errorf("Error getting the accountid and projectid: %v", err)
				ReturnHTTPError(w, r, 404, fmt.Sprintf("Error getting the accountid and projectid : %v", err))
				return
			}
		} else {
			accountID, err = getAccountID(manager.URL, tokenValue)
			if accountID == "Unauthorized" {
				logrus.Errorf("Unauthorized")
				ReturnHTTPError(w, r, 401, fmt.Sprintf("Unauthorized"))
				return
			}
			if err != nil {
				logrus.Errorf("Error getting the accountid : %v", err)
				ReturnHTTPError(w, r, 404, fmt.Sprintf("Error getting the accountid : %v", err))
				return
			}
		}

		//construct the responseBody
		var headerBody = make(map[string][]string)
		// var Body = make(map[string]interface{})

		requestHeader := reqestData.Headers
		for k, v := range requestHeader {
			headerBody[k] = v
		}

		headerBody["X-API-Account-Id"] = []string{accountID}
		if projectID != "" {
			headerBody["X-API-Project-Id"] = []string{projectID}
		}
		// var responseBody RequestData
		reqestData.Headers = headerBody
		// responseBody.Body = Body
		//convert the map to JSON format
		if responseBodyString, err := json.Marshal(reqestData); err != nil {
			logrus.Info(err)
			w.WriteHeader(500)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(responseBodyString)
		}

	}
}

//get the projectID and accountID from rancher API
func getAccountAndProject(host string, envid string, token string) (string, string, error) {

	client := &http.Client{}
	requestURL := host + "v2-beta/projects/" + envid + "/accounts"
	fmt.Println(requestURL)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		logrus.Infof("Cannot connect to the rancher server. Please check the rancher server URL")
		return "", "", err
	}
	cookie := http.Cookie{Name: "token", Value: token}
	req.AddCookie(&cookie)
	resp, err := client.Do(req)
	if err != nil {
		logrus.Infof("Cannot connect to the rancher server. Please check the rancher server URL")
		return "", "", err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Infof("Cannot read the reponse body")
		return "", "", err
	}
	authMessage := AuthorizeData{}
	err = json.Unmarshal(bodyText, &authMessage)
	fmt.Println(string(bodyText))
	if err != nil {
		logrus.Infof("Cannot extract authorization JSON")
		return "", "", err
	}
	if authMessage.Message == "Unauthorized" {
		logrus.Infof("Unauthorized token")
		err := errors.New("Unauthorized token")
		return "Unauthorized", "Unauthorized", err
	}

	porjectid := resp.Header.Get("X-Api-Account-Id")
	userid := resp.Header.Get("X-Api-User-Id")
	if porjectid == "" || userid == "" {
		logrus.Infof("Cannot get porjectid or userid")
		err := errors.New("Forbidden")
		return "Forbidden", "Forbidden", err

	}
	if porjectid == userid {
		logrus.Infof("Cannot valid project id")
		err := errors.New("Cannot valid project id")
		return "", "", err

	}

	return porjectid, userid, nil
}

//get the accountID from rancher API
func getAccountID(host string, token string) (string, error) {

	client := &http.Client{}
	requestURL := host + "v2-beta/accounts"
	req, err := http.NewRequest("GET", requestURL, nil)
	cookie := http.Cookie{Name: "token", Value: token}
	req.AddCookie(&cookie)
	resp, err := client.Do(req)
	if err != nil {
		logrus.Infof("Cannot connect to the rancher server. Please check the rancher server URL")
		return "", err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	authMessage := AuthorizeData{}
	err = json.Unmarshal(bodyText, &authMessage)
	// fmt.Println(string(bodyText))
	if err != nil {
		logrus.Infof("Unmarshal token fail")
		return "", err
	}
	if authMessage.Message == "Unauthorized" {
		logrus.Infof("Unauthorized token")
		err := errors.New("Unauthorized token")
		return "Unauthorized", err
	}
	messageData := MessageData{}
	err = json.Unmarshal(bodyText, &messageData)
	if err != nil {
		logrus.Infof("Cannot extract accounts JSON")
		err := errors.New("Cannot extract accounts JSON")
		return "", err
	}
	result := ""
	//get id from the data
	for i := 0; i < len(messageData.Data); i++ {

		idData, suc := messageData.Data[i].(map[string]interface{})
		if suc {
			id, suc := idData["id"].(string)
			kind, namesuc := idData["kind"].(string)
			if suc && namesuc {
				//if the token belongs to admin, only return the admin token
				if kind == "admin" {
					return id, nil
				}
			} else {
				logrus.Infof("Cannot extract accounts JSON")
				err := errors.New("Cannot extract accounts JSON")
				return "", err
			}
			result = id

		}

	}

	return result, nil

}

//ReturnHTTPError handles sending out CatalogError response
func ReturnHTTPError(w http.ResponseWriter, r *http.Request, httpStatus int, errorMessage string) {
	svcError := ProxyError{
		Status:  strconv.Itoa(httpStatus),
		Message: errorMessage,
	}
	writeError(w, svcError)
}

func writeError(w http.ResponseWriter, svcError ProxyError) {
	status, err := strconv.Atoi(svcError.Status)
	if err != nil {
		logrus.Errorf("Error writing error response %v", err)
		w.Write([]byte(svcError.Message))
		return
	}
	w.WriteHeader(status)

	jsonStr, err := json.Marshal(svcError)
	if err != nil {
		logrus.Errorf("Error writing error response %v", err)
		w.Write([]byte(svcError.Message))
		return
	}
	w.Write([]byte(jsonStr))
}

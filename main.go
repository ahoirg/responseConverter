//main is a package that contains the logic and tests of the responseConverter.
package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const successStatusCode = 200

var maxGoroutines int64 = 10

//main is first method that is called when the consol application runs. It handles the parameters entered by the user.
func main() {
	var parameters []string = os.Args[1:]
	CreateRequest(parameters)
}

// CreteRequest is the public method that creates a request with the entered parameters.
// here it is checked whether the program can be called with the correct parameters.
func CreateRequest(parameters []string) {
	isValidParameters := true

	if parameters[0] == "-parallel" {
		if len(parameters) < 3 {
			fmt.Println("Missing parameter entered.")
			fmt.Println("Example Request: -parallel 3 www.gooogle.com")
			isValidParameters = false
		}

		numOfGoroutines, err := strconv.ParseInt(parameters[1], 10, 64)
		if err != nil {
			fmt.Println("If you are using '-parallel', you must enter this command as the first parameter and enter the number of goroutines after the '-parallel' command.")
			isValidParameters = false
		} else if numOfGoroutines < 1 {
			fmt.Println("The number of goroutines cannot be less than 1.")
			isValidParameters = false
		}

		parameters = parameters[2:]
		maxGoroutines = numOfGoroutines
	}

	if isValidParameters && len(parameters) > 0 {
		createUrlAndPrintResponse(parameters)
	}
}

// createUrlAndPrintResponse is the method by which goroutines are created.
// within this method, createUrl(), getResponse(), convert() methods are called respectively. It prints valid results to the console.
func createUrlAndPrintResponse(parameters []string) {
	var ch chan struct{} = make(chan struct{}, maxGoroutines) //sets how many goroutines will be used at the same time.
	var wg sync.WaitGroup

	for i := 0; i < len(parameters); i++ {
		wg.Add(1)
		ch <- struct{}{} // blocks the use of goroutines above the limit.

		go func(parameter string) { // each parameter will be processed in different goroutines.
			url, err := createUrl(parameter)
			if err == nil {
				response, err := getResponse(url)
				if err == nil {
					md5text := convert(response)
					fmt.Println(url + " " + md5text)
				}
			}

			<-ch
			wg.Done()
		}(parameters[i])
	}

	wg.Wait() // blocks the main goroutine termination before other goroutines run.
}

//createUrl adds schema to the url If the url doesn't have any schema. If it's not a valid url, an error is returned.
func createUrl(item string) (result string, err error) {
	if !strings.Contains(item, "://") {
		item = "http://" + item
	}

	u, err := url.Parse(item)
	if err != nil || u.Host == "" || strings.HasPrefix(u.Host, ".") || strings.Contains(u.Host, "..") || strings.Count(":", u.Host) > 1 {
		return item, fmt.Errorf("isNotValid")
	}

	if u.Scheme == "" || !(u.Scheme == "fmt" || u.Scheme == "http" || u.Scheme == "https") {
		return item, fmt.Errorf("isNotValid")
	}

	if !strings.Contains(u.Host, ":") && !strings.Contains(u.Host, ".") {
		return item, fmt.Errorf("isNotValid")
	}

	if u.Port() != "" {
		_, portValidation := strconv.ParseInt(u.Port(), 10, 64)
		if portValidation != nil {
			return item, fmt.Errorf("isNotValid")
		}
	}

	return item, err
}

//getResponse creates a request and returns its response as a string.
func getResponse(url string) (string, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != successStatusCode {
		return "", err
	}

	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	response := string(responseData)
	return response, err
}

//convert converts string value to mv5 hash and returns as a string.
func convert(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

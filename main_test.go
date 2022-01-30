//main is a package that contains the logic and tests of the urlValidator.
package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type TestCase struct {
	url      string
	expected string
}

var validTestUrls = []TestCase{
	/* The hash values ​​here are converted from http://www.md5.cz/ */
	{"https://github.com/ahoirg", "https://github.com/ahoirg 9b01cd6a3ffffc17c492768e2349d58e"},
	{"http://github.com/ahoirg", "http://github.com/ahoirg 9b01cd6a3ffffc17c492768e2349d58e"},
	{"https://www.github.com/ahoirg", "https://www.github.com/ahoirg 427090706ed244c10e72f0b4c4fe3139"},
	{"http://www.github.com/ahoirg", "http://www.github.com/ahoirg 397aa820cea01bd703c01ac2c5318023"},
	{"github.com/ahoirg/exercismorg-go-answers", "http://github.com/ahoirg/exercismorg-go-answers 3be737c066791c0c5be459ae5970a9a6"},
	{"www.github.com/ahoirg/exercismorg-go-answers", "http://www.github.com/ahoirg/exercismorg-go-answers 28bb26efb61ae5a9f210e54256f3ff50"},
}

var InvalidTestUrls = []TestCase{
	{"htp://example.co.uk:1030/software/index.html", ""},
	{"example..com", ""},
	{"://www.example.co.uk:1030/software/index.html", ""},
	{"https://:www.example.co.uk/software/index.html", ""},
	{"http://.example.co.uk:1030/software/index.html", ""},
	{"https://www.examspleqwe.com/aaaaas", ""},
	{"https://www.thereIsNoDomainName.com", ""},
	{"https://www.oppsWrongWay.com", ""},
}

var pdfTestUrls = []TestCase{
	{"google.com", "http://google.com 8ff1c478ccca08cca025b028f68b352f"},
	{"adjust.com", "http://adjust.com 6b2560b9a5262571258cc173248b7492"},
	{"yandex.com", "http://yandex.com 4baab01ff9ff0f793bf423aeef539c9d"},
	{"facebook.com", "http://facebook.com ccae5ffa91c4936aef3efd5091a43f3e"},
	{"twitter.com", "http://twitter.com 857efe81a54c8b5c2241846ac4f08e66"},
	{"reddit.com/r/funny", "http://reddit.com/r/funny ff3b2b7dcd9e716ca0adcbd208061c9a"},
	{"reddit.com/r/notfunny", "http://reddit.com/r/notfunny ff3b2b7dcd9e716ca0adcbd208061c9a"},
	{"yahoo.com", "http://yahoo.com e2d50a30b7bfbfda097d72e32578c6a6"},
	{"baroquemusiclibrary.com", "http://baroquemusiclibrary.com 8e5138a0111364f08b10d37ed3371b11"},
}

func TestCreateRequest(t *testing.T) {
	testcases := [][]TestCase{validTestUrls, pdfTestUrls, InvalidTestUrls}

	for k := 0; k < 3; k++ {
		validTestCases := testcases[0]
		for i := 0; i < len(validTestCases); i++ {
			r, w, _ := os.Pipe()
			os.Stdout = w
			CreateRequest([]string{validTestCases[i].url})
			w.Close()

			out, _ := ioutil.ReadAll(r)
			temp := strings.Split(string(out), "\n")
			result := strings.Split(temp[0], " ")

			expecteds := strings.Split(validTestCases[i].expected, " ")
			if result[0] != expecteds[0] {
				t.Errorf("Expected:%s,Result:%s,Index:%v", expecteds[0], result[0], i)
				break
			}
			/* 	//The responses change therefore the md5 values ​​change. I couldn't find how to test this part.
			if result[1] != expecteds[1] {
				t.Errorf("Expected:%s,Result:%s,Index:%v", expecteds[1], result[1], i)
				break
			}
			*/
		}
	}
}

func TestCreateRequestWithParallel(t *testing.T) {
	urls := make([]string, len(pdfTestUrls)+2)
	urls[0] = "-parallel"
	urls[1] = "3"
	for i := 2; i < len(pdfTestUrls)+2; i++ {
		urls[i] = pdfTestUrls[i-2].url
	}

	r, w, _ := os.Pipe()
	os.Stdout = w
	CreateRequest(urls)
	w.Close()

	out, _ := ioutil.ReadAll(r)
	temp := strings.Split(string(out), "\n")

	if len(temp)-1 != len(pdfTestUrls) {
		t.Errorf("resultCount:%v , expectedCount%v", len(temp), len(pdfTestUrls))
	}
}

func TestCreateRequestWithMissNumOfGoroutines(t *testing.T) {
	missEnterGorutinesNum := []string{"-parallel", "http://google.com", "2"}

	r, w, _ := os.Pipe()
	os.Stdout = w
	CreateRequest(missEnterGorutinesNum)
	w.Close()

	out, _ := ioutil.ReadAll(r)
	temp := strings.Split(string(out), "\n")

	if temp[0] != "If you are using '-parallel', you must enter this command as the first parameter and enter the number of goroutines after the '-parallel' command." {
		t.Errorf("An incorrect request was created. Error message should be received.")
	}

}
func TestCreateRequestWithInvalidParallel(t *testing.T) {
	invalidRun := []string{"-parallel", "http://google.com", "http://google.com"}
	r, w, _ := os.Pipe()
	os.Stdout = w
	CreateRequest(invalidRun)
	w.Close()

	out, _ := ioutil.ReadAll(r)
	temp := strings.Split(string(out), "\n")

	if temp[0] != "If you are using '-parallel', you must enter this command as the first parameter and enter the number of goroutines after the '-parallel' command." {
		t.Errorf("An incorrect request was created. Error message should be received.")
	}

}

func TestCreateRequestWithInvalidRequest(t *testing.T) {
	inValidRequest := []string{"-parallel", "2"}

	r, w, _ := os.Pipe()
	os.Stdout = w
	CreateRequest(inValidRequest)
	w.Close()

	out, _ := ioutil.ReadAll(r)
	temp := strings.Split(string(out), "\n")

	if temp[0] != "Missing parameter entered." {
		t.Errorf("An incorrect request was created. Error message should be received.")
	}

}

func TestCreateRequestWithZeroGoRoutine(t *testing.T) {
	urls := []string{"-parallel", "0", "http://google.com"}

	r, w, _ := os.Pipe()
	os.Stdout = w
	CreateRequest(urls)
	w.Close()

	out, _ := ioutil.ReadAll(r)
	temp := strings.Split(string(out), "\n")

	if temp[0] != "The number of goroutines cannot be less than 1." {
		t.Errorf("Wrong error message")
	}
}

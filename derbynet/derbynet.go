package derbynet

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var fullUrl = "http://192.168.0.236/action.php"
var client http.Client

// init will create a new client with a cookie jar,
// which will consequently be used in all POST operations
func init() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	client = http.Client{
		Jar: jar,
	}
}

// GetCookie will post to the derby net URL, and log in as the
// Timer user, returning the cookie from the response
func GetCookie() {
	resp, err := client.PostForm(fullUrl, url.Values{
		"action":   {"role.login"},
		"name":     {"Timer"},
		"password": {""}})

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s\n", string(body))

	for _, c := range resp.Cookies() {
		log.Printf("cookie received: %s=%s\n", c.Name, c.Value)
	}
}

func timerMessage(msg string, params url.Values) string {
	if params == nil {
		params = make(url.Values)
	}
	
	params.Set("action", "timer-message")
	params.Set("message", msg)

	resp, err := client.PostForm(fullUrl, params)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

func Heartbeat() {
	log.Printf("heartbeat returned: %s\n", timerMessage("HEARTBEAT", nil))
}

func Hello() {
	log.Printf("hello returned: %s\n", timerMessage("HELLO", nil))
}

func Identified() {
	params := make(url.Values)
	params.Set("lane_count", "4")
	params.Set("timer", "github.com/Nescient/gpio-timer")
	params.Set("human", "GPIO Timer")
	log.Printf("identified returned: %s\n", timerMessage("IDENTIFIED", params))
}


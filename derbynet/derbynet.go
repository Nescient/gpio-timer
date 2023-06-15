// derbynet is a package to send and receive HTTP messages with the derbynet server.
// see https://derbynet.org/ and https://derbynet.org/builds/docs/Developers-%20Timer%20Messages.pdf
package derbynet

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// fullUrl is the complete URL to the derbynet action page
var fullUrl = "http://192.168.0.236/action.php"

// client is our saved "connected" client (that has the cookie)
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

// timerMessage is a helper function to send timer messages
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

// Heartbeat sends the heartbeat message to the server
func Heartbeat() {
	log.Printf("heartbeat returned: %s\n", timerMessage("HEARTBEAT", nil))
}

// Hello sends the hello message to the server
func Hello() {
	log.Printf("hello returned: %s\n", timerMessage("HELLO", nil))
}

// Identified sends the identified message to the server with the
// provided identification string (probably the git rev)
func Identified(ident string) {
	params := make(url.Values)
	params.Set("lane_count", "4")
	params.Set("timer", "github.com/Nescient/gpio-timer")
	params.Set("human", "GPIO Timer")
	params.Set("ident", ident)
	log.Printf("identified returned: %s\n", timerMessage("IDENTIFIED", params))
}


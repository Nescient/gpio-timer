package derbynet

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var baseUrl = "localhost"
var actionUrl = "/action.php"
var fullUrl = "http://192.168.0.236/action.php"
var client http.Client

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
func GetCookie() string {
	response, err := client.PostForm(fullUrl, url.Values{
		"action":   {"role.login"},
		"name":     {"Timer"},
		"password": {""}})

	//okay, moving on...
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s\n", string(body))

	//c := response.GetCookie("PHPSESSID")

	for _, c := range response.Cookies() {
		log.Printf("name: %s, value %s\n", c.Name, c.Value)
		if c.Name == "PHPSESSID" {
			return c.Name + "=" + c.Value
		}
	}

	return string(body)
}

func Hello(cookie string) {
	response, err := client.PostForm(fullUrl, url.Values{
		"action":  {"timer-message"},
		"message": {"HELLO"}})

	//okay, moving on...
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s\n", string(body))
}

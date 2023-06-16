// derbynet is a package to send and receive HTTP messages with the derbynet server.
// see https://derbynet.org/ and https://derbynet.org/builds/docs/Developers-%20Timer%20Messages.pdf
package derbynet

import (
	"github.com/antchfx/xmlquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	// "github.com/antchfx/xpath"
	"encoding/xml"
	"strconv"
)

// fullUrl is the complete URL to the derbynet action page
var fullUrl = "http://192.168.0.236/action.php"

// client is our saved "connected" client (that has the cookie)
var client http.Client

// roundID indicates which race round is currently active
var roundID = 0

// heat indicates which race heat is currently active
var heat = 0

// HeatReady is the XML structure that indicates a race is
// ready to start.
type HeatReady struct {
	XMLName xml.Name `xml:"heat-ready"`
	// lane-mask: The attribute value is a decimal representation of a bit mask showing which lanes
	// will be occupied for the heat. E.g., lane-mask=“14” means lane 1 (20) will be empty but cars
	// will be in lanes 2 (21), 3 (22), and 4 (23).
	laneMask int `xml:"lane-mask,attr"`
	// class: This attribute contains a human-readable name for the current racing class (group name)
	class string `xml:"class,attr"`
	// round: The ordinal round number for the class (e.g., 1 for the first round, 2 for the second, etc.).
	round int `xml:"round,attr"`
	// roundid: The internal integer identifier of the round. Roundids are unique across all rounds in
	// the database. (Thus, two different classes may each have a first round, but those rounds will
	// have different roundids.)
	roundID int `xml:"roundid,attr"`
	// heat: The number of the current heat within the current round, with 1 being the first heat
	heat int `xml:"heat,attr"`
}

type ActionResponse struct {
	XMLName xml.Name `xml:"action-response"`
	hr HeatReady `xml:"heat-ready"`
}

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

// processResponse parses the XML response and attempts to
// respond with requested information or set local variables
func processResponse(msg string) {
	log.Printf(msg)
	doc, err := xmlquery.Parse(strings.NewReader(msg))
	if err != nil {
		log.Fatal(err)
	}
	root := xmlquery.FindOne(doc, "//action-response")

	// log.Println("--------------")
	// log.Println(root.InnerText())
	// log.Println("--------------")

	if q := root.SelectElement("//remote-log"); q != nil {
		
		log.Println("--------------")
		log.Println(q.Data)
		log.Println(q.Attr)
		// log.Println(q.Attr)
		for _, a := range q.Attr {
			log.Println(a)
			log.Println(a.Name.Local)
			log.Println(a.Value)
		}
		log.Println("--------------")
	}	


	if q := root.SelectElement("//query"); q != nil {
		Flags()
	}
	if l := xmlquery.FindOne(doc, "//action-response/remote-log@send"); l != nil {
		SendLogs(l.InnerText() == "true")
	}
	for i, n := range xmlquery.Find(doc, "//action-response/failure@code") {
		log.Printf("Error %d: %s\n", i, n.InnerText())
	}
	if q := root.SelectElement("//heat-ready"); q != nil {
		log.Println("got heat ready:")
		// <heat-ready lane-mask= class= round= roundid= heat=/>
		// var h HeatReady
		var a ActionResponse
		// if err := xml.Unmarshal([]byte(msg), &h); err != nil {
		// 	log.Fatal(err)
		// }
		if err := xml.Unmarshal([]byte(msg), &a); err != nil {
			log.Fatal(err)
		}
		log.Println(a)
		roundID = a.hr.roundID
		heat = a.hr.heat
		log.Println(a.hr.class)
		log.Printf("First Got round %d and heat %d\n", roundID, heat)
		for _, a := range q.Attr {
			if a.Name.Local == "roundid" {
				roundID, err = strconv.Atoi(a.Value)
				if err != nil{
					log.Fatal(err)
				}
			} else if a.Name.Local == "heat" {
				heat, err = strconv.Atoi(a.Value)
				if err != nil{
					log.Fatal(err)
				}
			}
		}
		log.Printf("SEcond Got round %d and heat %d\n", roundID, heat)
	}
}

// Heartbeat sends the heartbeat message to the server
func Heartbeat() string {
	msg := timerMessage("HEARTBEAT", nil)
	log.Printf("heartbeat returned: %s\n", msg)
	return msg
}

// Hello sends the hello message to the server
func Hello() {
	msg := timerMessage("HELLO", nil)
	processResponse(msg)
}

// Identified sends the identified message to the server with the
// provided identification string (probably the git rev)
func Identified(ident string) {
	params := make(url.Values)
	params.Set("lane_count", "4")
	params.Set("timer", "github.com/Nescient/gpio-timer")
	params.Set("human", "GPIO Timer")
	params.Set("ident", ident)
	processResponse(timerMessage("IDENTIFIED", params))
}

// Flag sends the complete set of command-line flags, detected serial ports,
// and available timer classes to the server, encoded as additional parameters.
func Flags() {
	params := make(url.Values)
	// params.Set("flag-{flagname}", "{type}:{value}")
	// params.Set("desc-{flagname}", "{description}")
	params.Set("ports", "")
	params.Set("device-github.com Nescient gpio-timer", "GPIO Timer")
	processResponse(timerMessage("FLAGS", params))
}

// Started sends a message to the server indicating that the start gate has opened
func Started() {
	processResponse(timerMessage("STARTED", nil))
}

// Finished sends a message to the server to report results.  It is accompanied by
// additional parameters indicating the time and place of each lane that had a reported
// result. E.g., the message body might be,
// action=timer-message&message=FINISHED&roundid=7&heat=4
// &lane1=3.21&place1=2&lane2=3.33&place2=3&lane3=3.14&place3=1
// Here derby-timer.jar is confirming that the results it's reporting are for the fourth heat of the racing
// round whose identifier is 7, that the car in lane 1 has a time of 3.21 seconds and came in second in this
// heat, etc.
// The server may reject the reported results if e.g. they don't correspond to the currently-running heat
// (i.e., if the heat and roundid values aren't what were expected).
func Finished() {
	params := make(url.Values)
	params.Set("roundid", "")
	params.Set("heat", "")
	params.Set("lane1", "")
	params.Set("place1", "")
	params.Set("lane2", "")
	params.Set("place2", "")
	params.Set("lane3", "")
	params.Set("place3", "")
	params.Set("lane4", "")
	params.Set("place4", "")
	processResponse(timerMessage("FINISHED", params))
}

// SendLogs enables and disables sending log information to the server.
// If remote logging is active, derby-timer.jar makes POST requests to the post-timer-log.php URL. The
// POST request body has content-type text/plain type. The server appends the request body text to the
// logged text captured on the server.
func SendLogs(en bool) {
	if en {
		log.Printf("sending logs...")
	} else {
		log.Printf("not sending logs...")
	}
}

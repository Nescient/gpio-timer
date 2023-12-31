// derbynet is a package to send and receive HTTP messages with the derbynet server.
// see https://derbynet.org/ and https://derbynet.org/builds/docs/Developers-%20Timer%20Messages.pdf
package derbynet

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/antchfx/xmlquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// fullUrl is the complete URL to the derbynet action page
var fullUrl = "http://192.168.0.236/action.php"

// logUrl is the complete URL to the derbynet logging page
var logUrl = "http://192.168.0.236/post-timer-log.php"

// HeatReady is the XML structure that indicates a race is
// ready to start.
type HeatReady struct {
	XMLName xml.Name `xml:"heat-ready"`
	// lane-mask: The attribute value is a decimal representation of a bit mask showing which lanes
	// will be occupied for the heat. E.g., lane-mask=“14” means lane 1 (20) will be empty but cars
	// will be in lanes 2 (21), 3 (22), and 4 (23).
	LaneMask int `xml:"lane-mask,attr"`
	// class: This attribute contains a human-readable name for the current racing class (group name)
	Class string `xml:"class,attr"`
	// round: The ordinal round number for the class (e.g., 1 for the first round, 2 for the second, etc.).
	Round int `xml:"round,attr"`
	// roundid: The internal integer identifier of the round. Roundids are unique across all rounds in
	// the database. (Thus, two different classes may each have a first round, but those rounds will
	// have different roundids.)
	RoundID int `xml:"roundid,attr"`
	// heat: The number of the current heat within the current round, with 1 being the first heat
	Heat int `xml:"heat,attr"`
}

// RemoteLog is the XML structure for the logging node
// When true, the timer should post to the post-timer-log.php URL
type RemoteLog struct {
	XMLName xml.Name `xml:"remote-log"`
	Send    bool     `xml:"send,attr"`
}

// ActionResponse is the XML structure that represents the response
// from the derbynet server
type ActionResponse struct {
	XMLName xml.Name  `xml:"action-response"`
	Heat    HeatReady `xml:"heat-ready"`
	Log     RemoteLog `xml:"remote-log"`
}

// DerbyNet holds most of the useful objects to communicate with
// the derbynet server
type DerbyNet struct {
	// client is our saved "connected" client (that has the cookie)
	client http.Client

	// sendLogs indicates that the derbynet server has asked for logs
	sendLogs bool

	// i thought golang was thread safe.  it is not
	lock sync.Mutex

	// the last received HeatReady response
	heat HeatReady

	// an efficient way to wait for a heat
	waitHeat chan int
}

// init will create a new client with a cookie jar,
// which will consequently be used in all POST operations
func (this *DerbyNet) Initialize() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	this.client = http.Client{
		Jar: jar,
	}
	this.sendLogs = false
	this.waitHeat = make(chan int, 1)
}

// GetCookie will post to the derby net URL, and log in as the
// Timer user, returning the cookie from the response
func (this *DerbyNet) GetCookie() {
	this.lock.Lock()
	resp, err := this.client.PostForm(fullUrl, url.Values{
		"action":   {"role.login"},
		"name":     {"Timer"},
		"password": {""}})
	this.lock.Unlock()

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
func (this *DerbyNet) timerMessage(msg string, params url.Values) string {
	if params == nil {
		params = make(url.Values)
	}

	params.Set("action", "timer-message")
	params.Set("message", msg)

	this.lock.Lock()
	resp, err := this.client.PostForm(fullUrl, params)
	this.lock.Unlock()
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body.")
		log.Fatal(err)
	}

	return string(body)
}

// processResponse parses the XML response and attempts to
// respond with requested information or set local variables
func (this *DerbyNet) processResponse(msg string) {
	var respMsg ActionResponse
	if err := xml.Unmarshal([]byte(msg), &respMsg); err != nil {
		log.Fatal(err)
	}

	this.SendLogs(respMsg.Log.Send)
	if this.sendLogs {
		log.Println("VVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVV")
		log.Printf("%s", msg)
		log.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	}

	doc, err := xmlquery.Parse(strings.NewReader(msg))
	if err != nil {
		log.Fatal(err)
	}
	root := xmlquery.FindOne(doc, "//action-response")
	if q := root.SelectElement("//heat-ready"); q != nil {
		this.lock.Lock()
		this.heat = respMsg.Heat
		log.Printf("Heat %d is ready.\n", this.heat.Heat)
		this.waitHeat <- this.heat.Heat
		this.lock.Unlock()
	}

	if q := root.SelectElement("//query"); q != nil {
		this.Flags()
	}

	for i, n := range xmlquery.Find(doc, "//action-response/failure@code") {
		log.Printf("[ERROR%d] %s: %s\n", i, n.Attr[0].Value, n.InnerText())
	}

	// unused messages
	if q := root.SelectElement("//success"); q != nil {
		log.Println("last command successful")
	}
	if q := root.SelectElement("//abort"); q != nil {
		this.lock.Lock()
		this.heat.Heat = 0
		this.waitHeat <- this.heat.Heat
		this.lock.Unlock()
	}
	if q := root.SelectElement("//remote-start"); q != nil {
		log.Println("remote-start message ignored")
	}
	if q := root.SelectElement("//assign-flag"); q != nil {
		log.Println("assign-flag message ignored")
	}
	if q := root.SelectElement("//assign-port"); q != nil {
		log.Println("assign-port message ignored")
	}
	if q := root.SelectElement("//assign-device"); q != nil {
		log.Println("assign-device message ignored")
	}
}

// Heartbeat sends the heartbeat message to the server
func (this *DerbyNet) Heartbeat() {
	msg := this.timerMessage("HEARTBEAT", nil)
	this.processResponse(msg)
}

// Hello sends the hello message to the server
func (this *DerbyNet) Hello() {
	msg := this.timerMessage("HELLO", nil)
	this.processResponse(msg)
}

// Identified sends the identified message to the server with the
// provided identification string (probably the git rev)
func (this *DerbyNet) Identified(ident string) {
	params := make(url.Values)
	params.Set("lane_count", "4")
	params.Set("timer", "github.com/Nescient/gpio-timer")
	params.Set("human", "GPIO Timer")
	params.Set("ident", ident)
	this.processResponse(this.timerMessage("IDENTIFIED", params))
}

// Flag sends the complete set of command-line flags, detected serial ports,
// and available timer classes to the server, encoded as additional parameters.
func (this *DerbyNet) Flags() {
	params := make(url.Values)
	// params.Set("flag-{flagname}", "{type}:{value}")
	// params.Set("desc-{flagname}", "{description}")
	params.Set("ports", "")
	params.Set("device-github.com Nescient gpio-timer", "GPIO Timer")
	this.processResponse(this.timerMessage("FLAGS", params))
}

// WaitForHeat waits until the next heat has been started by the server
// it uses the heat.Heat (unique ID) to check if a heat is valid
func (this *DerbyNet) WaitForHeat() bool {
	select {
	case <-this.waitHeat:
		return this.heat.Heat > 0
	case <-time.After(10 * time.Second):
		return false
	}
}

// Started sends a message to the server indicating that the start gate has opened
func (this *DerbyNet) Started() {
	this.processResponse(this.timerMessage("STARTED", nil))
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
func (this *DerbyNet) Finished(lane1 float64, lane2 float64, lane3 float64, lane4 float64) {
	params := make(url.Values)
	params.Set("roundid", strconv.Itoa(this.heat.RoundID))
	params.Set("heat", strconv.Itoa(this.heat.Heat))
	if lane1 != 0.0 {
		params.Set("lane1", fmt.Sprintf("%.5f", lane1))
	}
	// params.Set("place1", "")
	if lane2 != 0.0 {
		params.Set("lane2", fmt.Sprintf("%.5f", lane2))
	}
	// params.Set("place2", "")
	if lane3 != 0.0 {
		params.Set("lane3", fmt.Sprintf("%.5f", lane3))
	}
	// params.Set("place3", "")
	if lane4 != 0.0 {
		params.Set("lane4", fmt.Sprintf("%.5f", lane4))
	}
	// params.Set("place4", "")
	this.processResponse(this.timerMessage("FINISHED", params))
}

// logPost is a pointless struct to allow me to create
// a custom write function for io.Writer
type logPost struct {
	client *DerbyNet
}

// Write implements a post call for the derbynet logging
// response is <success>XXXX bytes</success>
func (l logPost) Write(p []byte) (n int, err error) {
	size := len(p)
	reader := bytes.NewReader(p)
	l.client.lock.Lock()
	resp, err := l.client.client.Post(logUrl, "text/plain", reader)
	l.client.lock.Unlock()
	if err != nil {
		fmt.Println(err)
		return size - reader.Len(), err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return size, err
	}

	fmt.Println(string(body))

	doc, err := xmlquery.Parse(strings.NewReader(string(body)))
	if err != nil {
		fmt.Println(err)
		return size, err
	}
	root := xmlquery.FindOne(doc, "//success")
	if root != nil {
		i := size
		fmt.Sscanf(root.InnerText(), "%d bytes", &i)
		return i, nil
	}

	return size, nil
}

// SendLogs enables and disables sending log information to the server.
// If remote logging is active, derby-timer.jar makes POST requests to the post-timer-log.php URL. The
// POST request body has content-type text/plain type. The server appends the request body text to the
// logged text captured on the server.
func (this *DerbyNet) SendLogs(en bool) {
	if this.sendLogs != en {
		this.sendLogs = en
		if en {
			log.Printf("sending logs...")
			log.SetOutput(&logPost{})
		} else {
			log.SetOutput(os.Stdout)
			log.Printf("not sending logs...")
		}
	}
}

// Terminate indicates that this timer is terminating
func (this *DerbyNet) Terminate() {
	this.lock.Lock()
	this.heat.Heat = 0
	this.waitHeat <- this.heat.Heat
	this.lock.Unlock()
	params := make(url.Values)
	params.Set("detectable", "0")
	params.Set("error", "GPIO Timer is terminating.  Sorry!")
	this.timerMessage("MALFUNCTION", params)
}

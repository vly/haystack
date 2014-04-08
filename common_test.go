package haystack

import (
	"testing"
	"time"
)

// TestToJSON verifies that Message JSON encoding is working as expected
func TestToJSON(t *testing.T) {
	temp := &Message{
		"1234", // uid
		"test.com/testpage.html", // ref site
		time.Now().String(),      // timestamp
		"Pageview",               // event type
		"{'a': 1}"}               // data blob
	tempJson, _ := temp.ToJSON()
	temp2 := new(Message)
	if ok := temp2.FromJSON(tempJson); ok {
		tempJson2, _ := temp2.ToJSON()
		if string(tempJson) == string(tempJson2) {
			return
		}
	}
	t.Fatalf("ToJSON test failed on TestToJSON\n")
}

// TestPreFlight verifies preflight function is stopping oversized
// messages from being transmitted.
func TestPreFlight(t *testing.T) {
	temp := &Message{"1234", "test.com/testpage.html", time.Now().String(),
		"Pageview", "{'a': 1}"}
	if _, ok := temp.ToJSON(); !ok {
		t.Fatalf("PreFlight test failed to pass good message")
	}

	// Lets try an oversized message
	temp = &Message{"1234", "test.com/testpage.html", time.Now().String(),
		"Pageview", "{'a': 1}"}
	temp.Data = string(make([]byte, MAX_MSG_SIZE))
	if _, ok := temp.ToJSON(); ok {
		t.Fatalf("PreFlight test failed to stop bad message")
	}
}

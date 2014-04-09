package haystack

import (
	"log"
	"testing"
	"time"
)

const (
	tempStreamName string = "PRODUCER_TEST"
)

// initialise connection to Kinesis instance
func TestKinesisInit(t *testing.T) {
	if _, ok := InitKinesis(); !ok {
		log.Println("Failed to init Kinesis connection.")
		t.Fail()
	}
}

// check if stream exists
func TestCheckStream(t *testing.T) {
	if k, ok := InitKinesis(); ok {
		if test := CheckStream(k, "RANDOM_STREAM_NAME"); test {
			t.Fail()
		}
	}
}

// open or create a kinesis stream
func TestStreamInit(t *testing.T) {
	if k, ok := InitKinesis(); ok {
		// create a new stream
		if ok := CreateStream(k, tempStreamName, 1); !ok {
			t.Fail()
		}
		// creating for the second time should fail
		if ok := CreateStream(k, tempStreamName, 1); ok {
			t.Fail()
		}
	}
}

// put message into stream
func TestSendMessage(t *testing.T) {
	ch := make(chan bool)
	tempMessage := &Message{
		"1234", // uid
		"test.com/testpage.html", // ref site
		time.Now().String(),      // timestamp
		"Pageview",               // event type
		"{'a': 1}"}
	if k, ok := InitKinesis(); ok {

		go SendMessage(k, tempStreamName, tempMessage, ch)
	}

	// Hang around for a response
	select {
	case resp := <-ch:
		if !resp {
			log.Println("Message failed to send")
			t.Fail()
		}
		return
	case <-time.After(10 * time.Second):
		log.Println("Message send timeout")
		t.Fail()
		return
	}

}

// delete stream
func TestDeleteStream(t *testing.T) {
	if k, ok := InitKinesis(); ok {
		if ok := DeleteStream(k, tempStreamName); ok {
			return
		}
	}

	log.Println("Failed to delete stream")
	t.Fail()
}

package haystack

import (
	"fmt"
	"log"
	"testing"
	"time"
)

const (
	StreamName string = "CONSUMER_TEST"
)

// TestGetMessage gets a specific record
func TestPrint(t *testing.T) {
	if k, ok := InitKinesis(); ok {
		if test := CheckStream(k, StreamName); !test {
			if ok := CreateStream(k, StreamName, 1); !ok {
				t.Fail()
			}
		}
		ch := make(chan bool)
		for i := 0; i < 5; i++ {
			tempMessage := &Message{
				"1234", // uid
				"test.com/testpage.html", // ref site
				time.Now().String(),      // timestamp
				"Pageview",               // event type
				fmt.Sprintf("{'a': %d}", i)}
			go SendMessage(k, StreamName, tempMessage, ch)
		}

		// Hang around for a response
		select {
		case resp := <-ch:
			if !resp {
				log.Println("Message failed to send")
				t.Fail()
			}

		case <-time.After(10 * time.Second):
			log.Println("Message send timeout")
			t.Fail()
			return
		}

		if data, ok := GetMessages(k, StreamName, "shardId-000000000000"); ok {
			PrintRecords(data)

		} else {
			log.Println("Failed to retrieve messages from stream")
			t.Fail()
		}
	}
}

// TestGetMessages attempts to retrieve all messages in the test stream
func TestGetMessages(t *testing.T) {
	if k, ok := InitKinesis(); ok {
		if _, ok := GetMessages(k, StreamName, "shardId-000000000000"); !ok {
			log.Println("Failed to retrieve messages from stream")
			t.Fail()
		}
	}
	if k, ok := InitKinesis(); ok {
		if ok := DeleteStream(k, StreamName); ok {
			return
		}
	}
}

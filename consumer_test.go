package haystack

import (
	"log"
	"testing"
)

// TestGetMessages attempts to retrieve all messages in the test stream
func TestGetMessages(t *testing.T) {
	if k, ok := InitKinesis(); ok {
		if _, ok := GetMessages(k, tempStreamName); !ok {
			log.Println("Failed to retrieve messages from stream")
			t.Fail()

		}
	}
}

// TestGetMessage gets a specific record
func TestGetMessage(t *testing.T) {}

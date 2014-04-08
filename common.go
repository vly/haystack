package haystack

import (
	"encoding/json"
	"log"
)

const (
	MAX_MSG_SIZE uint32 = (50 << (10 * 1))
)

// Message represents the clickstream
// data to be passed to a stream.
type Message struct {
	Uid       string // user ID
	Ref       string // referral site
	Timestamp string // timestamp
	EventType string // type of event
	Data      string // data blob
}

// ToJSON marshalls the Message struct into a JSON blob
func (m *Message) ToJSON() ([]byte, bool) {
	out, err := json.Marshal(m)
	if err != nil {
		log.Println("Failed to marshal the message")
	} else if ok := m.preFlight(out); ok {
		return out, true
	}
	return out, false
}

// FromJSON (re)contructs a Message struct from a JSON blob
func (m *Message) FromJSON(data []byte) bool {
	if err := json.Unmarshal(data, &m); err != nil {
		return false
	}
	return true
}

// preFlight verifies if data size is under the 50KB limit
// set by AWS Kinesis.
func (m *Message) preFlight(data []byte) bool {
	if uint32(len(data)) <= MAX_MSG_SIZE {
		return true
	}
	log.Println("Message too large")
	return false
}

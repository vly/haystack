package haystack

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	kinesis "github.com/sendgridlabs/go-kinesis"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Base struct for Haystack producer
type Producer struct {
	Conn *kinesis.Kinesis
}

// LogFile passes logging messages to a flatfile
func LogFile(message string) {
	f, err := os.OpenFile("haystack_producer_messages.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(message)
}

// generate a dictionary from passed params
func GenDict(data []string) map[string]string {
	r_map := make(map[string]string)
	for i, _ := range data {
		x := strings.Split(data[i], "_")
		answer := ""
		if len(x) == 2 {
			answer = x[1]
		}
		r_map[x[0]] = answer
	}
	return r_map
}

// PassQuery sends the message to a given stream
func PassQuery(stream string, message string) (ok bool) {
	return true
}

// Authentication key struct
type Keys struct {
	AccessKey string
	SecretKey string
}

// SetEnv stores keys as environment vars
func (k *Keys) SetEnv() {
	err := os.Setenv("AWS_ACCESS_KEY", k.AccessKey)
	err = os.Setenv("AWS_SECRET_KEY", k.SecretKey)
	if err != nil {
		log.Fatalln("Failed to set environment vars")
	}
}

// loadTokens imports auth keys from rootkey.csv (AWS keyset)
func loadTokens(key *Keys) bool {
	if tokens, err := ioutil.ReadFile("rootkey.csv"); err == nil {
		temp := strings.Split(string(tokens), "\n")
		if len(temp) == 2 {
			key.AccessKey = strings.SplitAfter(temp[0], "=")[1]
			key.SecretKey = strings.SplitAfter(temp[1], "=")[1]
			return true
		}
	}
	log.Println("Error loading tokens from rootkey.csv")
	return false
}

// InitKinesis initialises a connection to the Kinesis service
func InitKinesis() (*kinesis.Kinesis, bool) {
	k := new(Keys)
	if ok := loadTokens(k); ok {
		ksis := kinesis.New(k.AccessKey, k.SecretKey)
		return ksis, true
	}
	return nil, false
}

// CheckStream verifies whether a stream already exists
func CheckStream(k *kinesis.Kinesis, streamname string) bool {
	args := kinesis.NewArgs()
	args.Add("StreamName", streamname)
	temp, err := k.ListStreams(args)
	if err != nil {
		log.Println("Could not obtain list of streams.")
	} else {

		for _, b := range temp.StreamNames {
			if b == streamname {
				return true
			}
		}
	}
	return false
}

// CreateStream creates a new stream
func CreateStream(k *kinesis.Kinesis, tempStreamName string, shards int) bool {
	if err := k.CreateStream(tempStreamName, shards); err != nil {
		return false
	}

	// wait for Stream ready state
	timeout := make(chan bool, 30)
	resp := &kinesis.DescribeStreamResp{}
	log.Printf("Waiting for stream to be created")
	for {
		args := kinesis.NewArgs()
		args.Add("StreamName", tempStreamName)
		resp, _ = k.DescribeStream(args)
		log.Printf(".")

		if resp.StreamDescription.StreamStatus != "ACTIVE" {
			time.Sleep(4 * time.Second)
			timeout <- true
		} else {
			break
		}
	}

	return true
}

// DeleteStream deletes an existing stream
func DeleteStream(k *kinesis.Kinesis, streamName string) bool {
	if err := k.DeleteStream(streamName); err == nil {
		return true
	}
	return false
}

// SendMessage generates a new JSON blob and sends to a stream
func SendMessage(k *kinesis.Kinesis, streamName string, msg *Message, comms chan bool) {
	args := kinesis.NewArgs()
	args.Add("StreamName", streamName)
	args.Add("PartitionKey", fmt.Sprintf("partitionKey-%d", 1))
	if data, ok := msg.ToJSON(); ok {
		args.AddData(data)
		if _, err := k.PutRecord(args); err == nil {
			comms <- true
			return
		}
	}
	log.Println("Failed to send message")
	comms <- false
	return
}

// ServerInit launches the Producer endpoint
// ...need to figure out how to do graceful shutdown.
func ServerInit(quit chan bool) {
	// Init main handler of RESTful endpoints
	m := martini.Classic()

	// Init Kinesis comms object
	streamName := "testStream"
	k, ok := InitKinesis()
	if !ok {
		log.Fatalln("Could not init connection to Kinesis server")
	}
	// check if stream exists, create if not
	if test := CheckStream(k, streamName); !test {
		if ok := CreateStream(k, streamName, 1); !ok {
			log.Fatalln("Failed to create stream.")
		}
	}
	// make a buffered comms channel
	ch := make(chan bool, 100)

	// Process message passed via GET request
	m.Get("/log", func(res http.ResponseWriter, r *http.Request) (int, string) {
		block := strings.SplitAfter(r.RequestURI, "/log?")[1]
		var data []string
		if strings.Contains(block, "%7C") {
			data = strings.Split(strings.SplitAfter(r.RequestURI, "/log?")[1], "&")
		} else {
			data = strings.Split(strings.SplitAfter(r.RequestURI, "/log?")[1], "&")
		}
		for _, b := range data {
			log.Println("d: " + b)
		}

		tempMessage := &Message{
			data[0],                  // uid
			"test.com/testpage.html", // ref site
			time.Now().String(),      // timestamp
			"Pageview",               // event type
			fmt.Sprintf("{'data': %s}", data[1])}

		go SendMessage(k, streamName, tempMessage, ch)

		res.Header().Set("Content-Type", "application/json")
		return 200, "{'status': 'ok'}"
	})

	m.Patch("/", func() {
		// update something
	})

	// Process message passed via POST request
	m.Post("/log", binding.Json(Message{}), binding.ErrorHandler,
		func(msg Message, params martini.Params, r *http.Request, res http.ResponseWriter) (int, string) {

			msg.Timestamp = time.Now().String()
			res.Header().Set("Content-Type", "application/json")
			go SendMessage(k, streamName, &msg, ch)

			if out, ok := msg.ToJSON(); ok {
				return 200, string(out)
			}

			return 501, ""
		})

	m.Put("/", func() {
		// replace something
	})

	m.Delete("/", func() {
		// destroy something
	})

	m.Options("/", func() {
		// http options
	})

	m.NotFound(func() string {
		// handle 404
		log.Printf("%s\n", "Yep...")
		return "Something went wrong."
	})

	//m.Run()
	err := http.ListenAndServe(":5004", m)
	if err != nil {
		log.Fatal(err)
	}

	// quit if told to
	// have a feeling this won't work
	select {
	case <-quit:

		return
	}
}

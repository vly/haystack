package haystack

import (
	"github.com/codegangsta/martini"
	kinesis "github.com/sendgridlabs/go-kinesis"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// LogFile passes logging messages to a flatfile
func LogFile(message string) {
	f, err := os.OpenFile("ping_messages.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(message)
}

// LogStream is the dogfood function that
func LogStream(message string) {
	// something
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

type Keys struct {
	AccessKey string
	SecretKey string
}

func (k *Keys) SetEnv() {
	err := os.Setenv("AWS_ACCESS_KEY", k.AccessKey)
	err = os.Setenv("AWS_SECRET_KEY", k.SecretKey)
	if err != nil {
		log.Fatalln("Failed to set environment vars")
	}
}

func loadTokens(key *Keys) {
	log.Println("Loading tokens")
	if tokens, err := ioutil.ReadFile("rootkey.csv"); err == nil {
		temp := strings.Split(string(tokens), "\n")
		if len(temp) == 2 {
			key.AccessKey = strings.SplitAfter(temp[0], "=")[1]
			key.AccessKey = key.AccessKey[:len(key.AccessKey)-1]
			key.SecretKey = strings.SplitAfter(temp[1], "=")[1]
			return
		}
	}
	log.Fatalln("Error loading tokens from rootkey.csv")
	return
}

func InitKinesis() (*kinesis.Kinesis, bool) {
	k := new(Keys)
	loadTokens(k)
	ksis := kinesis.New(k.AccessKey, k.SecretKey)

	return ksis, true
}

func CheckStreams(k *kinesis.Kinesis, streamname string) bool {
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

func Init() {
	// Init main handler of RESTful endpoints
	m := martini.Classic()

	// Init Kinesis comms object
	k, ok := InitKinesis()
	if !ok {
		log.Fatalln("Could not init connect")
	}

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

		// check if stream exists, create new one if it doesn't
		if ok := CheckStreams(k, data[0]); !ok {
			if err := k.CreateStream(data[0], 1); err != nil {
				log.Fatalln("Error creating new stream.")
			}
		}

		res.Header().Set("Content-Type", "application/json")
		return 200, "{}"
	})

	m.Patch("/", func() {
		// update something
	})

	// Process message passed via POST request
	m.Post("/log", func(r *http.Request, res http.ResponseWriter) (int, string) {
		if err := r.ParseForm(); err != nil {
			log.Printf("%s", "nothing posted")
		}
		data_values := make(map[string]string)
		for a, b := range r.Form {
			data_values[a] = b[0]
		}
		res.Header().Set("Content-Type", "application/json")
		// SubmitTicket(data_values)
		return 200, "{}"
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
	err := http.ListenAndServe(":5001", m)
	if err != nil {
		log.Fatal(err)
	}
}

Haystack
========

A minimal clickstream service written in Go.
Utilises AWS Kinesis as the cloud Kafka equivalent.

Package includes the producer service, consumer, and front-end examples to get you started.

Pull requests welcome!

[ ![Codeship Status for vly/haystack](https://www.codeship.io/projects/f0185e10-a132-0131-a163-46e13967588c/status?branch=master)](https://www.codeship.io/projects/18279)

### AWS Authentication
In order to authenticate with your Kinesis instance, create a new key, download as a csv (rootkey.csv), and place in application root.

### Sample Producer
` curl -X POST -d '{"cid": "12345", "type": "pageview", "ref": "http://localhost/test.html", "context": ""}' http://localhost:5004/log`

### Dependencies

- Using Sendgridlab's kinesis library: [https://github.com/sendgridlabs/go-kinesis](https://github.com/sendgridlabs/go-kinesis)
- Martini for producer endpoint: [https://github.com/codegangsta/martini](https://github.com/codegangsta/martini)

### Releases
- 0.0.1 initialisation

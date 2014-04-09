package haystack

import (
	"fmt"
	kinesis "github.com/sendgridlabs/go-kinesis"
	"log"
)

func GetShardIterator(k *kinesis.Kinesis, streamName string, shardId string) string {
	args := kinesis.NewArgs()
	args.Add("StreamName", streamName)
	args.Add("ShardId", shardId)
	args.Add("ShardIteratorType", "TRIM_HORIZON")
	a, _ := k.GetShardIterator(args)
	return a.ShardIterator
}

func GetMessages(k *kinesis.Kinesis, streamName string) (out [][]byte, ok bool) {
	shardIterator := GetShardIterator(k, streamName, "shardId-000000000000")
	args := kinesis.NewArgs()
	args.Add("ShardIterator", shardIterator) // only one shard at the moment
	if resp, err := k.GetRecords(args); err == nil {
		fmt.Println(resp)

		for {
			args = kinesis.NewArgs()
			args.Add("ShardIterator", shardIterator)
			recordResp, err := k.GetRecords(args)
			if len(recordResp.Records) > 0 {
				for _, d := range recordResp.Records {
					res, err := d.GetData()
					if err != nil {
						LogFile(err.Error())
					}
					fmt.Printf("GetRecords Data: %v, err: %v\n", string(res), err)
					out = append(out, res)
				}
			} else if recordResp.NextShardIterator == "" || shardIterator == recordResp.NextShardIterator || err != nil {
				LogFile(fmt.Sprintf("GetRecords ERROR: %v\n", err))
				break
			}

			shardIterator = recordResp.NextShardIterator
		}
	} else {
		log.Println(err.Error())
	}
	return out, true
}

// StoreRecords is an example of a local store
func StoreRecords(records []byte) bool {
	return true
}

// PrintRecords
func PrintRecords(records [][]byte) {
	for a, b := range records {
		log.Printf("%d: %s\n", a, string(b))
	}
}

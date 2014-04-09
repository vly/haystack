package haystack

import (
	"fmt"
	kinesis "github.com/sendgridlabs/go-kinesis"
)

func GetMessages(k *kinesis.Kinesis, streamName string) ([]byte, bool) {
	args := kinesis.NewArgs()
	args.Add("StreamName", streamName)
	args.Add("ShardId", "000000000000")
	args.Add("ShardIteratorType", "TRIM_HORIZON")
	if resp, err := k.GetRecords(args); err == nil {
		fmt.Println(resp)
		resp10, _ := k.GetShardIterator(args)
		shardIterator := resp10.ShardIterator
		for {
			args = kinesis.NewArgs()
			args.Add("ShardIterator", shardIterator)
			resp11, err := k.GetRecords(args)

			if len(resp11.Records) > 0 {
				fmt.Printf("GetRecords Data BEGIN\n")
				for _, d := range resp11.Records {
					res, err := d.GetData()
					fmt.Printf("GetRecords Data: %v, err: %v\n", string(res), err)
				}
				fmt.Printf("GetRecords Data END\n")
			} else if resp11.NextShardIterator == "" || shardIterator == resp11.NextShardIterator || err != nil {
				fmt.Printf("GetRecords ERROR: %v\n", err)
				break
			}

			shardIterator = resp11.NextShardIterator
		}
	}
	return []byte(""), true
}

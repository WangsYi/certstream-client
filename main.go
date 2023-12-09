package main

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/WangsYi/certstream-go"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("example")

func main() {
	// 捕获全局异常
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	// The false flag specifies that we want heartbeat messages.
	stream, errStream := certstream.CertStreamEventStream(false, "ws://127.0.0.1:8080/full-stream")
	for {
		select {
		case jq := <-stream:
			messageType, err := jq.String("message_type")

			if err != nil {
				log.Fatal("Error decoding jq string")
			}

			log.Info("Message type -> ", messageType)
			log.Info("recv: ", jq)
			data, err := jq.Interface("data")
			if err != nil {
				log.Fatal("Error decoding jq string")
			}
			dataRaw, err := json.Marshal(data)
			if err != nil {
				log.Fatal("Error decoding jq string")
			}
			dataStr := string(dataRaw)

			// 用时间戳为文件名创建文件，并保存data字符串到文件
			tm := time.Now().Unix()
			dir := "./data/" + strconv.FormatInt(tm/1000, 10)
			_, err = os.Stat(dir)
			if err != nil {
				err = os.MkdirAll(dir, os.ModePerm)
				if err != nil {
					log.Fatal(err)
				}
			}

			f, err := os.Create(dir + "/" + strconv.FormatInt(time.Now().Unix(), 10) + ".txt")
			if err != nil {
				log.Fatal(err)
			}
			_, err = f.WriteString(dataStr)
			if err != nil {
				err = f.Close()
				if err != nil {
					log.Fatal(err)
				}
				log.Fatal(err)
			}
			f.Sync()
			err = f.Close()
			if err != nil {
				log.Fatal(err)
			}

		case err := <-errStream:
			log.Error(err)
		}
	}
}

package activemq

import (
	stomp "github.com/go-stomp/stomp"
)


type produceData struct {
	Payload    string `json:"payload"`
	Headers map[string]string `json:"headers"`
}

type AmqpClient struct {
	Username string
	Password string
	Address  string
	Channels []string

	Custom     bool
	Processor  func(inp string) string
	Subscriber func(inp string) string

	conn   *stomp.Conn
	sub    map[string]*stomp.Subscription
	msgIn  []byte
	msgOut []byte
}

type Process struct {
	IAmqpClient
}

type IAmqpClient interface {
	Connect()
	LogMessageBeforeProcess(msg, peid string, isManual bool) (rsp string)
	LogMessageAfterProcess(msg, peid string , isManual bool) (rsp string)
	LogResponse(msg, source string, isManual bool) (rsp string)
	LogResponseError(msg, source string, isManual bool) (rsp string)
	LogProcessError(msg, source string, isManual bool) (rsp string)
	
	Process(inp1 string) (out2 string)
	Send(msg string) string
}

package activemq

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	stomp "github.com/go-stomp/stomp"
)

// Connect to activeMQ
var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{
	// stomp.ConnOpt.Login("userid", "userpassword"),
	// stomp.ConnOpt.Host("localhost"),
	stomp.ConnOpt.HeartBeat(360*time.Second, 360*time.Second),
	stomp.ConnOpt.HeartBeatError(360 * time.Second),
}

func (am *AmqpClient) Connect() {

	fmt.Println("\n2).um trying to connect to the activemq")

	options = append(options, stomp.ConnOpt.Login(am.Username, am.Password))
	conn, err := stomp.Dial("tcp", am.Address, options...)
	if err != nil {
		os.Exit(1)
	}

	am.conn = conn
	am.sub = make(map[string]*stomp.Subscription)

	if am.Custom {
		if am.Processor == nil {
			os.Exit(1)
		}
		if am.Subscriber == nil {
			os.Exit(1)
		}
	}

	// starting all channels
	go am.listner()
}

func (am *AmqpClient) listner() {
	fmt.Println("\n3).um listening to the activemq")

	go func() {
		channelName := am.Channels[0]
		sub, _ := am.conn.Subscribe(channelName,
			stomp.AckAuto,
		// stomp.SubscribeOpt.Header("", ""),
		)
		am.sub[channelName] = sub

		fmt.Println("\n3.1).um subscribing to the IN ")

		for {
			// msg := am.subscribeIn(sub, channelName)
			msg,peid := am.subscribeIn(sub, channelName)
			if msg != nil {	
				fmt.Println("\nReplay services called. peid : " , peid)
				// fmt.Println("\nReplay services called. MSG  : " , string())
				
				var p = new(produceData)
				err := json.Unmarshal(am.msgIn, &p)
				// fmt.Println("\nUmarshalled msg : ", len(p.Payload) , "\n unmarshalled headers : " , len(p.Headers))
				if err != nil {
					am.processAndSend(string(am.msgIn), peid, true)
					am.msgIn = nil
				}

			}
		}
	}()

	go func() {
		channelName := am.Channels[1]

		sub, _ := am.conn.Subscribe(channelName, stomp.AckAuto)
		am.sub[channelName] = sub

		fmt.Println("\n3.1).um subscribing to the OUT ")

		for {
			msg ,peid := am.subscribeOut(sub, channelName)
			if msg != nil {

				var p = new(produceData)
				err := json.Unmarshal(am.msgOut, &p)
				if err != nil {
					am.onSend(string(am.msgIn),peid,true)
				}

			}
		}
	}()
}

func (am *AmqpClient) subscribeIn(s *stomp.Subscription, channel string) (msg interface{}, peid string) {

	m := <-s.C

	peid = string(m.Header.Get("CH_PEID"))
	if m.Body != nil {
		am.msgIn = m.Body
		msg = m.Body
		return
	}
	return
}

func (am *AmqpClient) subscribeOut(s *stomp.Subscription, channel string) (msg interface{}, peid string) {

	m := <-s.C
	peid = string(m.Header.Get("CH_PEID"))
	if m.Body != nil {
		am.msgOut = m.Body
		msg = m.Body
		return
	}
	return
}

func (am *AmqpClient) produce(channelName, message string, source map[string]string) (err error) {
	var p = new(produceData)
	p.Payload = message
	p.Headers = source
	b, _ := json.Marshal(p)
	// bytes := []byte(source)
	// b := json.Unmarshal(bytes,&p)
	fmt.Println("\nByte MSg : " , string(b))
retry:
	serr := am.conn.Send(
		channelName,  // destination
		"text/plains", // content-type
		[]byte(b),
		// stomp.SendOpt.Receipt,
		// stomp.SendOpt.Header("", ""),
	) // body
	if serr != nil {
		if serr == stomp.ErrAlreadyClosed {
			am.Connect()
			goto retry
		}

		err = serr
		return
	}

	return nil
}

func (am *AmqpClient) processAndSend(msg, peid string,isManual bool) (out4 string) {

	fmt.Println("\n4). um in ProcessAndSend func: received msg : " , msg)
	//do spme process or conversion
	out2 := am.Process(msg)
	fmt.Println("OUT2", out2)
	// logType=INFO, direction=OUT
	out3 := am.LogMessageAfterProcess(out2, peid, isManual)
	fmt.Println("OUT3", out3)
	
	out4 = am.onSend(out2,peid,isManual)
	return
}

func (am *AmqpClient) onSend(msg, peid string, isManual bool) (out4 string) {

	fmt.Println("\n4). um in onSend func: received msg : " , msg , "\n PEID :", peid)
	
	//come here
	out4 = am.Send(msg)
	am.LogResponse(out4,peid,isManual)

	if(len(out4) == 0){
		am.LogResponseError("Error From endpoint", peid, isManual)
	}
	return
}

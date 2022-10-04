package activemq

import (
	// "encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/susinda/constants"
)

// func passToMQ(am *AmqpClient, msg, source ){
// 	am.produce("codeobe/log", msg, source)
// }
var logBeforeProcessHeaders = make(map[string]string)
var logAfterProcessHeaders = make(map[string]string)
var logResponseHeaders = make(map[string]string)
var logProcessErrorHeaders = make(map[string]string)
var logResponseErrorHeaders = make(map[string]string)

func (am *AmqpClient) LogMessageBeforeProcess(msg, peid string, isManual bool) (rsp string) {

	fmt.Println("\nLogBeforeProcess")
	getHeaderLogBeforeProcess(msg,peid,isManual)
	fmt.Println("\n logBeforeProcessHeaders ",logBeforeProcessHeaders  , "\nmsg: ", msg )
	am.produce("codeobe/log", msg, logBeforeProcessHeaders )
	rsp = msg
	return
}

func (am *AmqpClient) LogMessageAfterProcess(msg, peid string, isManual bool) (rsp string) {

	fmt.Println("\nLogMessageAfterProcess")
	getHeaderLogAfterProcess(msg,peid,isManual)
	fmt.Println("\n LogMessageAfterProcess ",logAfterProcessHeaders  , "\nmsg: ", msg )
	am.produce("codeobe/log", msg, logAfterProcessHeaders )
	rsp = msg
	return

}

func (am *AmqpClient) LogProcessError(msg, peid string, isManual bool) (rsp string) {

	fmt.Println("\n LogProcessError")
	getHeaderLogProcessError(msg,peid,isManual)
	fmt.Println("\n LogProcessError ",logProcessErrorHeaders  , "\nmsg: ", msg )
	am.produce("codeobe/log", msg, logProcessErrorHeaders )
	rsp = msg
	return

}


func (am *AmqpClient) LogResponse(msg, peid string, isManual bool) (rsp string) {

	fmt.Println("\nLogResponse")
	getHeaderLogResponse(msg,peid,isManual)
	fmt.Println("\n LogResponse ",logResponseHeaders  , "\nmsg: ", msg )
	am.produce("codeobe/log", msg, logResponseHeaders )
	rsp = msg
	return

}


func (am *AmqpClient) LogResponseError(msg, peid string, isManual bool) (rsp string) {

	fmt.Println("\n LogResponseError")
	getHeaderLogResponseError(msg,peid,isManual)
	fmt.Println("\n LogResponse ",logResponseErrorHeaders  , "\nmsg: ", msg )
	am.produce("codeobe/log", msg, logResponseErrorHeaders )
	rsp = msg
	return

}


func (am *AmqpClient) Process(inp1 string ) (out2 string) {

	fmt.Println("\nProcess")
	// if am.Custom {
	// 	return am.Processor(inp1)
	// }

	//Default Logic
	return strings.ToUpper(inp1)
}

func (am *AmqpClient) Send(inp2 string) (out4 string) {

	fmt.Println("\n5). Um in Send Function")
	// if am.Custom {
	// 	return am.Subscriber(inp2)
	// }

	//Default Logic
	return "CALLED FROM BACKEND API " + inp2
}

func Close(am *AmqpClient) {
	defer am.conn.Disconnect()
}


const(

	DOMAIN_KEY = "CH_DOMAINKEY"
	PEID = "CH_PEID"

	LOG_TYPE = "CH_LogType"
	LOG_TYPE_INFO = "INFO"
	LOG_TYPE_ERROR = "ERROR"

	LOG_DIRECTION = "CH_LogDirection"
	LOG_DIRECTION_S = "S"
	LOG_DIRECTION_T = "T"
	LOG_DIRECTION_TR = "T_R"
	LOG_DIRECTION_SR = "S_R"

	EXECUTION_MODE = "CH_ExecutionMode"
	EXECUTION_MODE_LISTENER = "ListenerExecution"
	EXECUTION_MODE_MANUAL = "ManualRetry"
	EXECUTION_MODE_SCHEDULED = "ScheduledExecution";

	EXECUTION_USER = "CH_ExecutionUser"
	EXECUTION_USER_SYSTEM = "System"

	ENDPOINT_TYPE = "CH_EndpointType";
	ENDPOINT_TYPE_FTP = "FTP"
	ENDPOINT_TYPE_MAIL = "MAIL"
	ENDPOINT_TYPE_HTTP = "HTTP"
	// ENDPOINT_TYPE_JMS string

	ENDPOINT_TYPE_JMS = "JMS"
	ENDPOINT_TYPE_SOLACE = "SOLACE"

	// ENDPOINT_META string
	ENDPOINT_META = "CH_EndpointMeta"

	JMS_DESTINATION = "CH_jms_destination"
	JMS_HOSTPORT = "CH_jms_hostPort";
	JMS_USER = "CH_jms_user"

	SOURCE_APP = "CH_SourceApp"
	TARGET_APP = "CH_TargetApp"

	PROCESS_NAME = "CH_ProcessName"
	PROCESS_GROUP = "CH_ProcessGroup"

)

func getHeaderLogBeforeProcess(req string,peid string,isManual bool){
	fmt.Println("getHeaderLogBeforeProcess")

	logBeforeProcessHeaders[DOMAIN_KEY] = getDomkey()
	logBeforeProcessHeaders[PEID] = peid
	logBeforeProcessHeaders[LOG_TYPE] = LOG_TYPE_INFO
	logBeforeProcessHeaders[LOG_DIRECTION] = LOG_DIRECTION_S
	logBeforeProcessHeaders[EXECUTION_MODE] = EXECUTION_MODE_LISTENER
	logBeforeProcessHeaders[EXECUTION_USER] = EXECUTION_USER_SYSTEM
	logBeforeProcessHeaders[SOURCE_APP] = "TestSrcApp"
	logBeforeProcessHeaders[TARGET_APP] = "TestTrgtApp"
	logBeforeProcessHeaders[PROCESS_NAME] = "emp-service"
	logBeforeProcessHeaders[PROCESS_GROUP] = "emp-service"
	logBeforeProcessHeaders[ENDPOINT_TYPE] = ENDPOINT_TYPE_JMS
	logBeforeProcessHeaders[ENDPOINT_META] = getInputJMSMetadata()
	
}

func getHeaderLogAfterProcess(req string,peid string,isManual bool){
	fmt.Println("getHeaderLogBeforeProcess")

	logAfterProcessHeaders[DOMAIN_KEY] = getDomkey()
	logAfterProcessHeaders[PEID] = peid
	logAfterProcessHeaders[LOG_TYPE] = LOG_TYPE_INFO
	logAfterProcessHeaders[LOG_DIRECTION] = LOG_DIRECTION_T
	logAfterProcessHeaders[EXECUTION_MODE] = EXECUTION_MODE_LISTENER
	logAfterProcessHeaders[EXECUTION_USER] = EXECUTION_USER_SYSTEM
	logAfterProcessHeaders[SOURCE_APP] = "TestSrcApp"
	logAfterProcessHeaders[TARGET_APP] = "TestTrgtApp"
	logAfterProcessHeaders[PROCESS_NAME] = "emp-service"
	logAfterProcessHeaders[PROCESS_GROUP] = "emp-service"
	logAfterProcessHeaders[ENDPOINT_TYPE] = ENDPOINT_TYPE_JMS
	logAfterProcessHeaders[ENDPOINT_META] = getOutputJMSMetadata()
	
}

func getHeaderLogResponse(req string,peid string,isManual bool){
	fmt.Println("getHeaderLogBeforeProcess")

	logResponseHeaders[DOMAIN_KEY] = getDomkey()
	logResponseHeaders[PEID] = peid
	logResponseHeaders[LOG_TYPE] = LOG_TYPE_INFO
	logResponseHeaders[LOG_DIRECTION] = LOG_DIRECTION_TR
	logResponseHeaders[EXECUTION_MODE] = EXECUTION_MODE_LISTENER
	logResponseHeaders[EXECUTION_USER] = EXECUTION_USER_SYSTEM
	logResponseHeaders[SOURCE_APP] = "TestSrcApp"
	logResponseHeaders[TARGET_APP] = "TestTrgtApp"
	logResponseHeaders[PROCESS_NAME] = "emp-service"
	logResponseHeaders[PROCESS_GROUP] = "emp-service"
	logResponseHeaders[ENDPOINT_TYPE] = ENDPOINT_TYPE_JMS
	logResponseHeaders[ENDPOINT_META] = getOutputJMSMetadata()
	
}

func getHeaderLogProcessError(req string,peid string,isManual bool){
	fmt.Println("getHeaderLogBeforeProcess")

	logProcessErrorHeaders[DOMAIN_KEY] = getDomkey()
	logProcessErrorHeaders[PEID] = peid
	logProcessErrorHeaders[LOG_TYPE] = LOG_TYPE_ERROR
	logProcessErrorHeaders[LOG_DIRECTION] = LOG_DIRECTION_S
	logProcessErrorHeaders[EXECUTION_MODE] = EXECUTION_MODE_LISTENER
	logProcessErrorHeaders[EXECUTION_USER] = EXECUTION_USER_SYSTEM
	logProcessErrorHeaders[SOURCE_APP] = "TestSrcApp"
	logProcessErrorHeaders[TARGET_APP] = "TestTrgtApp"
	logProcessErrorHeaders[PROCESS_NAME] = "emp-service"
	logProcessErrorHeaders[PROCESS_GROUP] = "emp-service"
	logProcessErrorHeaders[ENDPOINT_TYPE] = ENDPOINT_TYPE_JMS
	logProcessErrorHeaders[ENDPOINT_META] = getInputJMSMetadata()
	
}

func getHeaderLogResponseError(req string,peid string,isManual bool){
	fmt.Println("getHeaderLogBeforeProcess")

	logResponseErrorHeaders[DOMAIN_KEY] = getDomkey()
	logResponseErrorHeaders[PEID] = peid
	logResponseErrorHeaders[LOG_TYPE] = LOG_TYPE_ERROR
	logResponseErrorHeaders[LOG_DIRECTION] = LOG_DIRECTION_S
	logResponseErrorHeaders[EXECUTION_MODE] = EXECUTION_MODE_LISTENER
	logResponseErrorHeaders[EXECUTION_USER] = EXECUTION_USER_SYSTEM
	logResponseErrorHeaders[SOURCE_APP] = "TestSrcApp"
	logResponseErrorHeaders[TARGET_APP] = "TestTrgtApp"
	logResponseErrorHeaders[PROCESS_NAME] = "emp-service"
	logResponseErrorHeaders[PROCESS_GROUP] = "emp-service"
	logResponseErrorHeaders[ENDPOINT_TYPE] = ENDPOINT_TYPE_JMS
	logResponseErrorHeaders[ENDPOINT_META] = getInputJMSMetadata()
	
}


func getInputJMSMetadata() string{
	// meta string
	meta  := JMS_DESTINATION + "=" + viper.GetString(constants.DATABASE + "." + constants.CHANNELONE) + "," + JMS_HOSTPORT + "=" + "stomp://" + viper.GetString(constants.DATABASE + "." + constants.ADDRESS) + ":" + viper.GetString(constants.DATABASE + "." + constants.PORT) + "," +  JMS_USER + "=" + viper.GetString(constants.DATABASE + "." + constants.USER)

	fmt.Println("\ngetInputJMSMetadata  : " , meta)

	return meta

}

func getOutputJMSMetadata() string{
	// meta string
	meta  := JMS_DESTINATION + "=" + viper.GetString(constants.DATABASE + "." + constants.CHANNELTWO) + "," + JMS_HOSTPORT + "=" + "stomp://" + viper.GetString(constants.DATABASE + "." + constants.ADDRESS) + ":" + viper.GetString(constants.DATABASE + "." + constants.PORT) + "," +  JMS_USER + "=" + viper.GetString(constants.DATABASE + "." + constants.USER)

	fmt.Println("\ngetOutputJMSMetadata  : " , meta)

	return meta

}


func getDomkey() string{

	domkey := viper.GetString(constants.DATABASE + "." + constants.DOMAINKEY)
	return domkey
}


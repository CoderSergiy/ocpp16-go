/*	==========================================================================
	OCPP 1.6 Protocol
	Filename: server.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/ocpp16-go/
	Purpose: Server implementation for the ocpp example
			 Includes gorutine and charger structs
	=============================================================================
*/

package main

import (
    "github.com/julienschmidt/httprouter"
    "net/http"
	"time"
	"github.com/gorilla/websocket"
	"github.com/CoderSergiy/ocpp16-go/example"
	"github.com/CoderSergiy/ocpp16-go/core"
	"github.com/CoderSergiy/golib/logging"
	"github.com/CoderSergiy/golib/timelib"
	"github.com/CoderSergiy/golib/tools"
)

const (
	configFilePath 	string = "/tmp/configs.json"
	logFilesPath 	string = "/tmp/logs/server"
)

var (
	log logging.Log
	ServerConfigs example.Configs
)


/****************************************************************************************
 *	Struct 	: GoRoutineSet
 * 
 * 	Purpose : Handles goroutine parameters
 *
*****************************************************************************************/
type GoRoutineSet struct {
	WriteChannel	chan string
	Active 			bool
	DateCreated 	int64
}

/****************************************************************************************
 *
 * Function : GoRoutineSetConstructor (Constructor)
 *
 *  Purpose : Creates a new instance of the GoRoutineSet
 *
 *	  Input : Nothing
 *
 *	Return : GoRoutineSet object
*/
func GoRoutineSetConstructor() GoRoutineSet {
    gSet := GoRoutineSet{}
    gSet.Active = true
    gSet.DateCreated = int64(time.Now().Unix())
	gSet.WriteChannel = make (chan string)
    return (gSet)
}




/****************************************************************************************
 *
 * Function : main
 *
 *  Purpose : Main method to start server
 *
 *   Return : Nothing
*/
func main() {
	log = logging.LogConstructor(logFilesPath, true)
	log.Info_Log("Server started.")

	// Set server configurations from file
	configs, configErr := example.SetConfigsFromFile(configFilePath)
	if configErr != nil {
		log.Error_Log("Cannot set configs from file '%v' with error '%v'", configFilePath, configErr)
		return
	}
	ServerConfigs = configs
	log.Info_Log("Set configs from file '%s'", configFilePath)
	log.Info_Log("Uploaded '%v' chargers configurations", len(configs.Chargers))
	log.Info_Log("Chargers '%v'", configs.MaxQueueSize)

	// Define http router
	router := httprouter.New()
	// Set router for the ocpp V1.6 connection in the json format
    router.GET("/ocppj/1.6/:chargerName", wsChargerHandler)
	// Start server
	log.Error_Log("Server fata errorr: '%v'", http.ListenAndServe(":8080", router))
}


/****************************************************************************************
 *
 * Function : wsChargerHandler
 *
 *  Purpose : Handler
 *
 *    Input : w http.ResponseWriter - http response
 *			  r *http.Request - http request object
 *			  ps httprouter.Params - router parameter
 *
 *   Return : Nothing
*/
func wsChargerHandler (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tm := timelib.EventTimerConstructor()
    log.Info_Log("Handle income request from Host '%v' and Path '%v'", r.URL.Host, r.URL.Path)

	chargerName := ps.ByName("chargerName")
	log.Info_Log("HTTP is connected. Charger name '%v'", chargerName)

	// Get Charger from the Configs
	chargerObj, err := ServerConfigs.GetChargerObj(chargerName)
	if err != nil {
		// There is no charger with specifed name in the configs
		http.Error(w, "Not allowed to connect for specified charger", http.StatusBadRequest)
		log.Error_Log("Not allowed to connect for specified charger '%v' with error '%v'", chargerName, err)
		return
	}

	// Check if charger already has connection with server
	if chargerObj.WebSocketConnected {
		// Requested charger is connected to server already
		log.Info_Log("Charger '%v' is connected already", chargerName)
		return
	}

	log.Info_Log("Charger '%v' is exists and websocket connection is not established yet. Will try now", chargerName)

	//Convert http request to WebSocket
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		log.Error_Log("Could not open websocket connection for charger name '%v' with error '%v'", chargerName, err)
		return
	}

	log.Info_Log("Websocket is connected. Charger name '%v'", chargerName)

	// Create log instance with file name "server.{chargerName}"
	chargerLog := logging.LogConstructor(logFilesPath + "." + chargerName, true)
	//log.Info_Log("Logging in file '%v'", chargerLog.fileName)

	ocppHandlers := example.OCPPHandlersConstructor()

	// Authorise request
	chargerObj.AuthConnection = ocppHandlers.Authorisation(chargerName, r)
	// Store remote IP
	chargerObj.InboundIP = r.RemoteAddr
	// Set Charger's WebSocket flag as connected
	chargerObj.WebSocketConnected = true
	// Add charger details to ocppHandlers
	ocppHandlers.Charger = chargerObj

	ocppHandlers.Log = chargerLog
	// Define Gorutine set
	goroutineSet := GoRoutineSetConstructor()

	// Start Read and Write gorutine
	go logReaderRD(conn, chargerName, &goroutineSet, &ocppHandlers, &chargerLog)
	go logReaderWR(conn, chargerName, &goroutineSet, &chargerLog)

	log.Info_Log("OCPPRequestHandler is finished in %v", tm.PrintTimerString())
}

/****************************************************************************************
 *
 * Function : logReaderRD
 *
 *  Purpose : Goroutine method to read from connected websocket
 *
 *    Input : conn *websocket.Conn - websocket connection pointer
 *			  chargerName string - current charger name
 *			  gs *GoRoutineSet - goroutine struct for current charger
 *			  ocppHandlers *example.OCPPHandlers - defined ocppHandlers for the charger
 *			  chargerLog *logging.Log - pointer to the charger log file
 *
 *   Return : Nothing
 *
*/

func logReaderRD(conn *websocket.Conn, chargerName string, gs *GoRoutineSet, ocppHandlers *example.OCPPHandlers, chargerLog *logging.Log) {
	defer conn.Close()
	chargerLog.Info_Log("[%v] Start RD goroutine for charger '%v'", tools.GetGoID(), chargerName)

	// Define OCPP Handler Class
	centralSystem := core.CentralSystemHandlerConstructor(ocppHandlers)

	for {
		// Read websocket message
		_, message, readingSocketError := conn.ReadMessage()

		if readingSocketError != nil {
			chargerLog.Error_Log("[%v] Client is disconnected with error: '%v'", tools.GetGoID(), readingSocketError)
			gs.Active = false
			gs.WriteChannel <- "stop"
			break
		}

		chargerLog.Info_Log("[%v] Received '%v'", tools.GetGoID(), string(message))

		// Call OCPP message handler
		response, responseErr, socketStatus := centralSystem.HandleIncomeMessage(string(message))
		gs.Active = socketStatus

		if responseErr != nil {
			chargerLog.Error_Log("Response error: '%v'", responseErr)
		}

		// If handler generated callResult message - send it to the charger
		if response != "" {
			gs.WriteChannel <- response
		}

	}

	// At this point webSocket needs to be closed
	// Clear memory allocated for ocppHandler
	ocppHandlers = nil

	chargerLog.Info_Log("[%v] Reading goroutine is finished", tools.GetGoID())
}

/****************************************************************************************
 *
 * Function : logReaderWR
 *
 *  Purpose : Goroutine method to write to connected websocket
 *
 *    Input : conn *websocket.Conn - websocket connection pointer
 *			  chargerName string - current charger name
 *			  gs *GoRoutineSet - goroutine struct for current charger
 *			  chargerLog *logging.Log - pointer to the charger log file
 *
 *   Return : Nothing
 *
*/

func logReaderWR(conn *websocket.Conn, chargerName string, gs *GoRoutineSet, chargerLog *logging.Log) {
	defer conn.Close()
    chargerLog.Info_Log("[%v] Start WR gorutine for charger '%v'", tools.GetGoID(), chargerName)

	for {
		// Wait for the message from channel
		message := <- gs.WriteChannel

		// Check if gorutine must to be finished
		if gs.Active == false {
			log.Info_Log("[%v] Close writing goroutine", tools.GetGoID())
			break
		}

		//Send response to the charger
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		    chargerLog.Error_Log("Send error: '%v'", err)
		    return
		}

		chargerLog.Info_Log("Sent to charger '%v'", message)
	}
	// At this point webSocket needs to be closed
	// Clear memory allocated for goroutineSet
	close(gs.WriteChannel)
	gs = nil

	chargerLog.Info_Log("[%v] Writing goroutine is finished", tools.GetGoID())
}

/*	==========================================================================
	OCPP 1.6 Protocol
	Filename: server.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/ocpp16-go/
	Purpose: Server implementation for the ocpp example
			 Includes gorutine and charger structs

	Handlers supported by server:
		1. messageStatusHandler
		2. triggerActionHandler
		3. wsChargerHandler
	=============================================================================
*/

package main

import (
	"encoding/json"
	"github.com/CoderSergiy/golib/logging"
	"github.com/CoderSergiy/golib/timelib"
	"github.com/CoderSergiy/golib/tools"
	"github.com/CoderSergiy/ocpp16-go/core"
	"github.com/CoderSergiy/ocpp16-go/example"
	"github.com/CoderSergiy/ocpp16-go/messages"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	configFilePath string = "/tmp/configs.json"
	logFilesPath   string = "/tmp/logs/server"
)

var (
	log           logging.Log
	ServerConfigs example.Configs
	MQueue        example.SimpleMessageQueue
)

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
	log.Info_Log("Max queue size is %v", configs.MaxQueueSize)

	// Init message queue
	MQueue = example.SimpleMessageQueueConstructor()

	// Define http router
	router := httprouter.New()
	// Handle clients API requests
	router.GET("/message/:messageReference/status", messageStatusHandler)
	router.POST("/command/:chargerName/triggeraction/:action", triggerActionHandler)
	// Set router for the ocpp V1.6 connection in the json format
	router.GET("/ocppj/1.6/:chargerName", wsChargerHandler)
	// Start server
	log.Error_Log("Server fata errorr: '%v'", http.ListenAndServe(":8080", router))
}

/****************************************************************************************
 *
 * Function : messageStatusHandler
 *
 *  Purpose : Handles client request to get message status
 *
 *    Input : w http.ResponseWriter - http response
 *            r *http.Request - http request object
 *            ps httprouter.Params - router parameter
 *
 *   Return : Nothing
 */
func messageStatusHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tm := timelib.EventTimerConstructor()
	log.Info_Log("Handle income messageStatusHandler request from Host '%v' and Path '%v'", r.URL.Host, r.URL.Path)

	reference := ps.ByName("messageReference")

	if reference == "" {
		log.Error_Log("Reference parameter is empty")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Info_Log("Reference is '%v'", reference)

	message, success := MQueue.GetMessage(reference)
	if !success {
		http.Error(w, string(example.CreateFailResponse("Message is not exist")), http.StatusOK)
		log.Error_Log("Message is not exists in the queue with uniqueID: '%v'. Finished in '%v'", reference, tm.PrintTimerString())
		return
	}

	jsonResult, err := json.Marshal(message)
	if err != nil {
		http.Error(w, string(example.CreateFailResponse("Internal Application Error")), http.StatusOK)
		log.Error_Log("Cannot marshal message '%v'. Finished in '%v'", reference, tm.PrintTimerString())
		return
	}

	// Send response in json format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResult)

	log.Info_Log("messageStatusHandler is finished in %v", tm.PrintTimerString())
}

/****************************************************************************************
 *
 * Function : triggerActionHandler
 *
 *  Purpose : Handle clients request to generate triggerAction request to the charger
 *
 *    Input : w http.ResponseWriter - http response
 *			  r *http.Request - http request object
 *			  ps httprouter.Params - router parameter
 *
 *   Return : Nothing
 */
func triggerActionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tm := timelib.EventTimerConstructor()
	log.Info_Log("Handle income TriggerAction request from Host '%v' and Path '%v'", r.URL.Host, r.URL.Path)

	chargerName := ps.ByName("chargerName")
	action := ps.ByName("action")

	if chargerName == "" || action == "" {
		log.Error_Log("One of the parameters are empty charger '%v' action '%v'", chargerName, action)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Info_Log("Charger name '%v' and action '%v'", chargerName, action)

	// Get Charger from the Configs
	chargerObj, err := ServerConfigs.GetChargerObj(chargerName)
	if err != nil {
		// There is no charger with specified name in the configs
		log.Error_Log("Not allowed to connect for specified charger '%v' with error '%v'", chargerName, err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Sanitize the TriggerMessage type from the request
	if !core.SanitizeTriggerMessageType(action) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Error_Log("TriggerMessage type '%v' is not supported", action)
		return
	}

	id := uuid.New()
	// Generate Call request payload for the TriggerMessage
	triggerMessagePayload := core.CreateTriggerMessageRequestPayload(core.TriggerMessageType(action), 0)
	// Generate Call request to the charger
	callMessageRequest := messages.CreateCallMessage(
		id.String(),
		core.ACTION_TRIGGERMESSAGE,
		triggerMessagePayload.GetPayload(),
	)

	// Convert Call message to string
	callMessageString, messageErr := callMessageRequest.ToString()
	if messageErr != nil {
		http.Error(w, string(example.CreateFailResponse("Internal Error")), http.StatusOK)
		log.Error_Log("Error to generate callMessage: '%v'. Finished in '%v'", messageErr, tm.PrintTimerString())
		return
	}

	// Create message for the queue
	queueMessage := example.Message{Action: action, Sent: callMessageString, Status: example.MESSAGE_TYPE_NEW, Received: ""}
	// Add message to the queue
	addingErr := MQueue.Add(callMessageRequest.UniqueID, queueMessage)
	if addingErr != nil {
		http.Error(w, string(example.CreateFailResponse("Internal Error")), http.StatusOK)
		log.Error_Log("Error to add message to the queue: '%v'. Finished in '%v'", addingErr, tm.PrintTimerString())
		return
	}

	log.Info_Log(callMessageString)
	(*chargerObj).WriteChannel <- callMessageRequest.UniqueID

	// Send response in json format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(example.CreateSuccessResponse(callMessageRequest.UniqueID))

	log.Info_Log("triggerActionHandler is finished in %v", tm.PrintTimerString())
}

/****************************************************************************************
 *
 * Function : wsChargerHandler
 *
 *  Purpose : Handler the requests from charger.
 *			  Socket is upgrading to WebSocket further
 *
 *    Input : w http.ResponseWriter - http response
 *			  r *http.Request - http request object
 *			  ps httprouter.Params - router parameter
 *
 *   Return : Nothing
 */
func wsChargerHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tm := timelib.EventTimerConstructor()
	log.Info_Log("Handle income wsCharger request from Host '%v' and Path '%v'", r.URL.Host, r.URL.Path)

	chargerName := ps.ByName("chargerName")
	log.Info_Log("HTTP is connected. Charger name '%v'", chargerName)

	// Get Charger from the Configs
	chargerObj, err := ServerConfigs.GetChargerObj(chargerName)
	if err != nil {
		// There is no charger with specified name in the configs
		http.Error(w, "Not allowed to connect for specified charger", http.StatusBadRequest)
		log.Error_Log("Not allowed to connect for specified charger '%v' with error '%v'", chargerName, err)
		return
	}

	// Check if charger already has connection with server
	if (*chargerObj).WebSocketConnected {
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
	chargerLog := logging.LogConstructor(logFilesPath+"."+chargerName, true)

	ocppHandlers := example.OCPPHandlersConstructor()
	ocppHandlers.Log = chargerLog
	ocppHandlers.MQueue = &MQueue

	// Store remote IP
	(*chargerObj).InboundIP = r.RemoteAddr
	// Set Charger's WebSocket flag as connected
	(*chargerObj).WebSocketConnected = true
	// Add charger details to ocppHandlers
	ocppHandlers.Charger = chargerObj

	// Authorise request
	(*chargerObj).AuthConnection = false
	ocppHandlers.Authorisation(chargerName, r)
	// Socket activity flag
	isSocketActive := true

	// Start Read and Write gorutine
	go logReaderRD(conn, chargerName, &isSocketActive, chargerObj, &ocppHandlers, &chargerLog)
	go logReaderWR(conn, chargerName, &isSocketActive, chargerObj, &chargerLog)

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
 *			  chargerObj *example.Charger - pointer on the charger obj
 *			  isSocketActive *bool - Socket activity flag
 *			  ocppHandlers *example.OCPPHandlers - defined ocppHandlers for the charger
 *			  chargerLog *logging.Log - pointer to the charger log file
 *
 *   Return : Nothing
 *
 */

func logReaderRD(conn *websocket.Conn, chargerName string, isSocketActive *bool, chargerObj *example.Charger, ocppHandlers *example.OCPPHandlers, chargerLog *logging.Log) {
	defer conn.Close()
	chargerLog.Info_Log("[%v] Start RD goroutine for charger '%v'", tools.GetGoID(), chargerName)

	// Define OCPP Handler Class
	centralSystem := core.CentralSystemHandlerConstructor(ocppHandlers)

	for {
		// Read websocket message
		_, rawMessage, readingSocketError := conn.ReadMessage()

		if readingSocketError != nil {
			chargerLog.Error_Log("[%v] Client is disconnected with error: '%v'", tools.GetGoID(), readingSocketError)
			*isSocketActive = false
			(*chargerObj).WebSocketConnected = false
			//gs.WriteChannel <- "stop"
			break
		}

		chargerLog.Info_Log("[%v] Received '%v'", tools.GetGoID(), string(rawMessage))

		// Add arrived rawMessage to the queue
		_, uniqueID, err := messages.GetMessageTypeFromRaw(string(rawMessage))
		if err != nil {
			chargerLog.Error_Log("[%v] Cannot get uniqueid from : '%v'", tools.GetGoID(), err)
			continue
		}
		qMessage := example.Message{Action: "", Received: string(rawMessage), Status: example.MESSAGE_TYPE_RECEIVED, Sent: ""}
		// Add message to the queue
		addingErr := MQueue.Add(uniqueID, qMessage)
		if addingErr != nil {
			log.Error_Log("Error to add message to the queue: '%v'", addingErr)
		}

		// Call OCPP message handler
		response, responseErr, socketStatus := centralSystem.HandleIncomeMessage(string(rawMessage))
		*isSocketActive = socketStatus

		if responseErr != nil {
			chargerLog.Error_Log("Response error: '%v'", responseErr)
		}

		// If handler generated callResult message - send it to the charger
		if response != "" {
			(*chargerObj).WriteChannel <- uniqueID
		}

	}

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
 *			  isSocketActive *bool - Socket activity flag
 *			  chargerObj *example.Charger - pointer on the charger obj
 *			  chargerLog *logging.Log - pointer to the charger log file
 *
 *   Return : Nothing
 *
 */

func logReaderWR(conn *websocket.Conn, chargerName string, isSocketActive *bool, chargerObj *example.Charger, chargerLog *logging.Log) {
	defer conn.Close()
	chargerLog.Info_Log("[%v] Start WR gorutine for charger '%v'", tools.GetGoID(), chargerName)

	for {
		// Wait for the message from channel
		uniqueID := <-(*chargerObj).WriteChannel
		//log.Info_Log("[%v] logReaderWR. UniqueuID '%v'", tools.GetGoID(), uniqueID)

		// Check if gorutine must to be finished
		if *isSocketActive == false {
			log.Info_Log("[%v] Close writing goroutine", tools.GetGoID())
			break
		}

		// Get message from the queue
		qMessage, _ := MQueue.GetMessage(uniqueID)

		//Send response to the charger
		if err := conn.WriteMessage(websocket.TextMessage, []byte(qMessage.Sent)); err != nil {
			chargerLog.Error_Log("Send error: '%v'", err)
			return
		}

		// update status in the queue
		qMessage.Status = example.MESSAGE_TYPE_COMPLETED
		MQueue.UpdateByUniqueID(uniqueID, qMessage)

		chargerLog.Info_Log("Sent to charger '%v'", qMessage.Sent)
	}
	// At this point webSocket needs to be closed
	// Clear memory allocated for goroutineSet
	//close((*chargerObj).WriteChannel)
	//gs = nil

	chargerLog.Info_Log("[%v] Writing goroutine is finished", tools.GetGoID())
}

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
	"github.com/CoderSergiy/golib/logging"
	"github.com/CoderSergiy/golib/timelib"
	"github.com/CoderSergiy/golib/tools"
	"github.com/CoderSergiy/ocpp16-go/core"
	"github.com/CoderSergiy/ocpp16-go/example"
	"github.com/CoderSergiy/ocpp16-go/messages"
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
	router.GET("/message/:messageReference/status", messageStatusAPIHandler)
	router.GET("/charger/:chargerName/status", chargerStatusAPIHandler)
	router.POST("/command/:chargerName/triggeraction/:action", triggerActionAPIHandler)
	// Set router for the ocpp V1.6 (json) connection
	router.GET("/ocppj/1.6/:chargerName", wsChargerHandler)
	// Start server
	log.Error_Log("Server fata errorr: '%v'", http.ListenAndServe(":8080", router))
}

/****************************************************************************************
 *
 * Function : messageStatusAPIHandler
 *
 *  Purpose : Handles client request to get message status
 *
 *    Input : w http.ResponseWriter - http response
 *            r *http.Request - http request object
 *            ps httprouter.Params - router parameter
 *
 *   Return : Nothing
 */
func messageStatusAPIHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tm := timelib.EventTimerConstructor()
	log.Info_Log("Handle income messageStatusAPIHandler request from Host '%v' and Path '%v'", r.URL.Host, r.URL.Path)
	// Get Message from queue
	example.GetMessageStatusAPI(ps.ByName("messageReference"), &MQueue, &log, w)
	log.Info_Log("messageStatusAPIHandler is finished in %v", tm.PrintTimerString())
}

/****************************************************************************************
 *
 * Function : chargerStatusAPIHandler
 *
 *  Purpose : Handles client request to get message status
 *
 *    Input : w http.ResponseWriter - http response
 *            r *http.Request - http request object
 *            ps httprouter.Params - router parameter
 *
 *   Return : Nothing
 */
func chargerStatusAPIHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tm := timelib.EventTimerConstructor()
	log.Info_Log("Handle income chargerStatusAPIHandler request from Host '%v' and Path '%v'", r.URL.Host, r.URL.Path)
	// Get Message from queue
	example.GetChargerStatusAPI(ps.ByName("chargerName"), &ServerConfigs, &log, w)
	log.Info_Log("chargerStatusAPIHandler is finished in %v", tm.PrintTimerString())
}

/****************************************************************************************
 *
 * Function : triggerActionAPIHandler
 *
 *  Purpose : Handle clients request to generate triggerAction request to the charger
 *
 *    Input : w http.ResponseWriter - http response
 *			  r *http.Request - http request object
 *			  ps httprouter.Params - router parameter
 *
 *   Return : Nothing
 */
func triggerActionAPIHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tm := timelib.EventTimerConstructor()
	log.Info_Log("Handle income TriggerActionAPI request from Host '%v' and Path '%v'", r.URL.Host, r.URL.Path)
	// Call Trigger Action API
	example.TriggerActionAPI(&ServerConfigs, &MQueue, &log, ps, w)
	log.Info_Log("triggerActionAPIHandler is finished in %v", tm.PrintTimerString())
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
	log.Info_Log("[%v] HTTP is connected", chargerName)

	// Get Charger from the Configs
	chargerObj, err := ServerConfigs.GetChargerObj(chargerName)
	if err != nil || chargerObj == nil {
		// There is no charger with specified name in the configs
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Error_Log("[%v] There is no charger with specified name with an error '%v'", chargerName, err)
		return
	}

	// Check if charger already has connection with server
	if chargerObj.WebSocketConnected {
		// Requested charger is connected to server already
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Info_Log("[%v] Charger is connected already", chargerName)
		return
	}

	log.Info_Log("[%v] Charger is exists and websocket connection is not established yet. Will try now", chargerName)

	//Convert http request to WebSocket
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Error_Log("[%v] Could not open websocket connection with error '%v'", chargerName, err)
		return
	}

	log.Info_Log("[%v] Connection is upgraded to Websocket type", chargerName)

	// Create log instance with file name "server.{chargerName}"
	chargerLog := logging.LogConstructor(logFilesPath+"."+chargerName, true)
	// Create OCPP Hadlers
	ocppHandlers := example.OCPPHandlersConstructor()

	// Update charger object
	chargerObj.InboundIP = r.RemoteAddr                                    // Store remote IP
	chargerObj.WebSocketConnected = true                                   // Set Charger's WebSocket flag as connected
	chargerObj.AuthConnection = ocppHandlers.Authorisation(chargerName, r) // Authorise request

	// Update ocppHandlers object
	ocppHandlers.Log = chargerLog     // Add log
	ocppHandlers.MQueue = &MQueue     // Add pointer to the Message queue
	ocppHandlers.Charger = chargerObj // Add charger details to ocppHandlers

	// Define socket activity flag
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

		if *isSocketActive == false {
			break
		}

		// Read websocket message
		_, rawMessage, readingSocketError := conn.ReadMessage()

		if readingSocketError != nil {
			chargerLog.Error_Log("[%v] Client is disconnected with error: '%v'", tools.GetGoID(), readingSocketError)
			*isSocketActive = false
			chargerObj.WebSocketConnected = false
			chargerObj.WriteChannel <- "wakeup" // Send string to wakeup write gorutine
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
			log.Error_Log("[%v] Error to add message to the queue: '%v'", tools.GetGoID(), addingErr)
		}

		// Call OCPP message handler
		response, responseErr, socketStatus := centralSystem.HandleIncomeMessage(string(rawMessage))
		*isSocketActive = socketStatus

		if responseErr != nil {
			chargerLog.Error_Log("[%v] Response error: '%v'", tools.GetGoID(), responseErr)
		}

		// If handler generated callResult message - send it to the charger
		if response != "" {
			chargerObj.WriteChannel <- uniqueID
		}
	}

	// Clear Charger parameters
	chargerObj.Disconnected()

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
		uniqueID := <-chargerObj.WriteChannel
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
			chargerLog.Error_Log("[%v] Send error: '%v'", tools.GetGoID(), err)
			return
		}

		// update status in the queue
		qMessage.Status = example.MESSAGE_TYPE_COMPLETED
		MQueue.UpdateByUniqueID(uniqueID, qMessage)

		chargerLog.Info_Log("[%v] Sent to charger '%v'", tools.GetGoID(), qMessage.Sent)
	}

	// Clear Charger parameters
	chargerObj.Disconnected()

	chargerLog.Info_Log("[%v] Writing goroutine is finished", tools.GetGoID())
}

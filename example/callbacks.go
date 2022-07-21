/*	==========================================================================
	OCPP 1.6 Protocol
	Filename: callbacks.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/ocpp16-go/example
	Purpose: All callback methods to handle ocpp messages
			 In this file using simplequeue for demo purpose
	=============================================================================
*/

package example

import (
	"time"
	"net/http"
	"github.com/CoderSergiy/golib/logging"
	"github.com/CoderSergiy/ocpp16-go/core"
	"github.com/CoderSergiy/ocpp16-go/messages"
)

const (
    WEBSOCKET_KEEP_OPEN		bool = true
	WEBSOCKET_CLOSE			bool = false
)

/****************************************************************************************
 *	Struct 	: OCPPHandlers
 * 
 * 	Purpose : Handles struct for each connected charger
 *
*****************************************************************************************/
type OCPPHandlers struct {
	Charger		Charger				// Charger struct which connected to the server
	Log 		logging.Log			// Pointer to the log
	txQueue 	SimpleMessageQueue 	// For example queue will be here
}

/****************************************************************************************
 *
 * Function : OCPPHandlersConstructor (Constructor)
 *
 *  Purpose : Creates a new instance of the OCPPHandlers
 *
 *	  Input : Nothing
 *
 *	Return : OCPPHandlers object
*/
func OCPPHandlersConstructor() OCPPHandlers {
	ocppHandlers := OCPPHandlers{}
	ocppHandlers.init()
	return ocppHandlers
}

/****************************************************************************************
 *
 * Function : OCPPHandlers::init
 *
 *  Purpose : Initiate variables of the OCPPHandlers structure
 *
 *	  Input : Nothing
 *
 *	 Return : Nothing
*/
func (cs *OCPPHandlers) init () {
	cs.txQueue.init() 	// Clear queues
}

/****************************************************************************************
 *
 * Function : OCPPHandlers::finaliseReqHandler
 *
 *  Purpose : Using to finalise Call massage
 *
 *    Input : callMessage messages.CallMessage - original Call message
 *			  responseMessage messages.Message - response message
 *			  socketStatus bool - false - when connection to charger needs to be closed
 *
 *   Return : string - response message in string format
 *			  error - if happened, nil otherwise
 *			  bool - false - when connection to charger needs to be closed, otherwise true
 *
*/
func (cs *OCPPHandlers) finaliseReqHandler (callMessage messages.CallMessage, responseMessage messages.Message, socketStatus bool) (string, error, bool) {
	
	// Convert response message to string format
	messageStr, err := responseMessage.ToString()
	if err != nil {
		return "", err, socketStatus
	}

	return messageStr, nil, socketStatus
}

/****************************************************************************************
 *
 * Function : OCPPHandlers::finaliseRespHandler
 *
 *  Purpose : Using to finalise Response message
 *			  Delete message from the queue
 *
 *    Input : uniqueID string - message's unique ID
 *			  socketStatus bool - false - when connection to charger needs to be closed
 *
 *   Return : error - if happened, nil otherwise
 *			  bool - false - when connection to charger needs to be closed, otherwise true
 *
*/
func (cs *OCPPHandlers) finaliseRespHandler (uniqueID string, socketStatus bool) (error, bool) {
	
	// Before end the handler delete message from txQueue
	err := cs.txQueue.DeleteByUniqueID(uniqueID)

	if err == nil {
		cs.Log.Info_Log("%v", cs.txQueue.printStatus())
	}

	return err, socketStatus
}

/****************************************************************************************
 *
 * Function : OCPPHandlers::GetActionHandler
 *
 *  Purpose : Get Action of the message from queue by unique ID
 *
 *    Input : uniqueID string - message's unique ID
 *
 *   Return : message's action in string format 
 *
*/
func (cs *OCPPHandlers) GetActionHandler (uniqueID string) string {

	// Get action from tx queue by UniqueID
	message, success := cs.txQueue.GetMessage(uniqueID)

	if success {
		// Message exists - return action
		return message.Action
	}

	// Message is not in the queue
	return ""
}





/****************************************************************************************
 *
 * Function : OCPPHandlers::Authorisation
 *
 * Purpose : Using to Authorise charger before allow websocket connection
 *
 *   Input : chargerName string - charger name to be validated
 *           request *http.Request - http request object 
 *
 *  Return : true - when charger is authorised, otherwise false
 *
*/
func (cs *OCPPHandlers) Authorisation (chargerName string, request *http.Request) (bool) {

	cs.Log.Info_Log("Auth request from URL '%s'", request.RequestURI)
	cs.Log.Info_Log("Header is '%s'", request.Header["Authorisation"])

	cs.Charger.AuthConnection = true
	cs.Log.Info_Log("[%v] Authorisation is '%v'", chargerName, cs.Charger.AuthConnection)

	return true
}



/* Define Call Handlers =========================================================================================
=================================================================================================================
*/

/****************************************************************************************
 *
 * Function : OCPPHandlers::BootNotificationRequestHandler
 *
 *  Purpose : Handle BootNotificationRequest
 *
 *    Input : callMessage messages.CallMessage - original Call message
 *
 *   Return : string - response message in string format
 *			  error - if happened, nil otherwise
 *			  bool - false - when connection to charger needs to be closed, otherwise true
 *
*/
func (cs *OCPPHandlers) BootNotificationRequestHandler (callMessage messages.CallMessage) (string, error, bool) {
	cs.Log.Info_Log("[%v] BootNotificationRequest Action", callMessage.UniqueID)

	// Define registration status for the response
	status := core.RegistrationStatusPending
	if cs.Charger.AuthConnection {
		status = core.RegistrationStatusAccepted
	}
	//Create payload
	bootNotificationResp := core.BootNotificationResponse{
		Status: status,
		HeartbeatInterval: 300,
		CurrentTime: time.Now().Format("2006-01-02 15:04:05.000"),
	}
	// Create CallResult message
	callMessageResponse := messages.CallResultMessageWithParam (
		callMessage.UniqueID,
		bootNotificationResp.GetPayload(),
	)

	return cs.finaliseReqHandler(callMessage, &callMessageResponse, WEBSOCKET_KEEP_OPEN)
}

/****************************************************************************************
 *
 * Function : OCPPHandlers::AuthorizeRequestHandler
 *
 *  Purpose : Handle AuthorizeRequest
 *
 *    Input : callMessage messages.CallMessage - original Call message
 *
 *   Return : string - response message in string format
 *			  error - if happened, nil otherwise
 *			  bool - false - when connection to charger needs to be closed, otherwise true
 *
*/
func (cs *OCPPHandlers) AuthorizeRequestHandler (callMessage messages.CallMessage) (string, error, bool) {
	cs.Log.Info_Log("[%v] AuthorizeRequest Action", callMessage.UniqueID)

	// Check Auth flag
	if cs.Charger.AuthConnection == false {
		// Charger is not authorised
		callErrMess := messages.CallErrorMessageConstructor()
		messageStr, err := callErrMess.ToString()
		return messageStr, err, WEBSOCKET_KEEP_OPEN
	}
	// Create ErrorResult message
	callMessageResponse := messages.CallErrorMessage{}

	return cs.finaliseReqHandler(callMessage, &callMessageResponse, WEBSOCKET_KEEP_OPEN)
}

/****************************************************************************************
 *
 * Function : OCPPHandlers::HeartbeatRequestHandler
 *
 *  Purpose : Handle HeartbeatRequest
 *
 *    Input : callMessage messages.CallMessage - original Call message
 *
 *   Return : string - response message in string format
 *			  error - if happened, nil otherwise
 *			  bool - false - when connection to charger needs to be closed, otherwise true
 *
*/
func (cs *OCPPHandlers) HeartbeatRequestHandler (callMessage messages.CallMessage) (string, error, bool) {
	cs.Log.Info_Log("[%v] HeartbeatRequest Action", callMessage.UniqueID)

	// Create payload of the heartbeat response
	heartBeatResponse := core.HeartBeatResponse{}
	heartBeatResponse.Init()

	// Create CallResultMessage
	callMessageResponse := messages.CallResultMessageWithParam (
		callMessage.UniqueID,
		heartBeatResponse.GetPayload(),
	)

	return cs.finaliseReqHandler(callMessage, &callMessageResponse, WEBSOCKET_KEEP_OPEN)
}





/* Define Response Handlers ===================================================================================
===============================================================================================================
*/

/****************************************************************************************
 *
 * Function : OCPPHandlers::AuthorizeResponseHandler
 *
 *  Purpose : Handle AuthorizeResponse
 *
 *    Input : callMessage messages.CallMessage - original Call message
 *
 *   Return : string - response message in string format
 *			  error - if happened, nil otherwise
 *			  bool - false - when connection to charger needs to be closed, otherwise true
 *
*/
func (cs *OCPPHandlers) AuthorizeResponseHandler (callResultMessage messages.CallResultMessage) (error, bool) {
	cs.Log.Info_Log("[%v] AuthorizeResponse Action", callResultMessage.UniqueID)

	// Perform Auth Routine

	return cs.finaliseRespHandler(callResultMessage.UniqueID, WEBSOCKET_KEEP_OPEN)
}





/* Define Error Handler ==============================================================================
======================================================================================================
*/

/****************************************************************************************
 *
 * Function : OCPPHandlers::OCPPErrorHandler
 *
 *  Purpose : Handle OCPP Error
 *
 *    Input : callMessage messages.CallMessage - original Call message
 *
 *   Return : string - response message in string format
 *			  error - if happened, nil otherwise
 *			  bool - false - when connection to charger needs to be closed, otherwise true
 *
*/
func (cs *OCPPHandlers) OCPPErrorHandler (callErrortMessage messages.CallErrorMessage) (error, bool) {
	cs.Log.Info_Log("[%v] OCPPErrorHandler", callErrortMessage.UniqueID)

	// Handle OCPP error

	return cs.finaliseRespHandler(callErrortMessage.UniqueID, WEBSOCKET_KEEP_OPEN)
}
/*	==========================================================================
	OCPP 1.6 Protocol
	Filename: handlers.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/ocpp16-go/core
	Purpose: File keep files to handle incoming OCPP requests
	=============================================================================
*/

package core

import (
	"reflect"
	"encoding/json"
	"errors"
	"github.com/CoderSergiy/ocpp16-go/messages"
)

const (
    REQUEST_TYPE_HANDLER        string = "RequestHandler"
    RESPONSE_TYPE_HANDLER		string = "ResponseHandler"
    ERROR_TYPE_HANDLER          string = "OCPPErrorHandler"

	GET_ACTION_HANDLER			string = "GetActionHandler"
)

type CallHandlers interface {}

/****************************************************************************************
 *	Struct 	: CallbackRoutine
 * 
 * 	Purpose :   Object stores callback handlers.
 *				Using different struct to isolate callback methods from RequestHandler
 *
*****************************************************************************************/
type CallbackRoutine struct {
	CallHandlers
}

/****************************************************************************************
 *
 * Function : CallbackRoutine::getHandler
 *
 *  Purpose : Get handler struct by provided name
 *
 *	 Return : Reflect value
*/
func (callbackRoutine *CallbackRoutine) getHandler (methodName string) reflect.Value {
	return (reflect.ValueOf(callbackRoutine.CallHandlers).MethodByName(methodName))
}




/****************************************************************************************
 *	Struct 	: RequestHandler
 * 
 * 	Purpose : Object handles the OCPP Message request
 *
*****************************************************************************************/
type RequestHandler struct {
	APIhadlers 			CallbackRoutine

}

/****************************************************************************************
 *
 * Function : CentralSystemHandlerConstructor (Constructor)
 *
 *  Purpose : Creates a new instance of the RequestHandler
 *
 *	  Input : callbackRoutines interface{} - routines to handle OCPP requests
 *
 *	Return : RequestHandler object
*/
func CentralSystemHandlerConstructor(callbackRoutines interface{}) RequestHandler {
	rh := RequestHandler{}
	rh.APIhadlers.CallHandlers = callbackRoutines
	return rh
}

/****************************************************************************************
 *
 * Function : RequestHandler::callRequestHandler
 *
 *  Purpose : Call method by provided name in CallbackRoutine
 *
 *	  Input : data interface{} - data to pass into handler
 *			  handlerMethodName string - name of the method
 *
 *	Return : string - response
 *			 error - if happened, nil otherwise
 *			 bool - true if needs to keep websocket open, false otherwise
*/
func (requestHandler *RequestHandler) callRequestHandler(data interface{}, handlerMethodName string) (string, error, bool) {
	
	methodCall := requestHandler.APIhadlers.getHandler(handlerMethodName)
	if !methodCall.IsValid() {
		return "", errors.New("Cannot find CallRequest handler"), true
	}

	// Call Request method from APIhadlers and return value
	in := make([]reflect.Value, methodCall.Type().NumIn())
	in[0] = reflect.ValueOf(data)
	response := methodCall.Call(in)

	return response[0].String(), errors.New(response[1].String()), response[2].Bool()
}

/****************************************************************************************
 *
 * Function : RequestHandler::callResponseHandler
 *
 *  Purpose : Call method by provided name in CallbackRoutine
 *
 *	  Input : data interface{} - data to pass into handler
 *			  handlerMethodName string - name of the method
 *
 *	Return : string - response
 *			 error - if happened, nil otherwise
 *			 bool - true if needs to keep websocket open, false otherwise
*/
func (requestHandler *RequestHandler) callResponseHandler(data interface{}, handlerMethodName string) (string, error, bool) {

	methodCall := requestHandler.APIhadlers.getHandler(handlerMethodName)
	if !methodCall.IsValid() {
		return "", errors.New("Cannot find CallResponse handler"), true
	}

	// Call Response method from APIhadlers and return value
	in := make([]reflect.Value, methodCall.Type().NumIn())
	in[0] = reflect.ValueOf(data)
	response := methodCall.Call(in)
	return "", errors.New(response[0].String()), response[1].Bool()
}

/****************************************************************************************
 *
 * Function : RequestHandler::getActionByUniqueID
 *
 *  Purpose : Call method to obtain action by provided uniqueID of the message
 *
 *	  Input : uniqueID string - message uniqueID
 *
 *	 Return : string - Action name
*/
func (requestHandler *RequestHandler) getActionByUniqueID (uniqueID string) string {
	methodCall := requestHandler.APIhadlers.getHandler(GET_ACTION_HANDLER)
	if !methodCall.IsValid() {
		return ""
	}

	// Call method to get action by uniqueID
	in := make([]reflect.Value, methodCall.Type().NumIn())
	in[0] = reflect.ValueOf(uniqueID)
	response := methodCall.Call(in)
	return response[0].String()
}

/****************************************************************************************
 *
 * Function : RequestHandler::getMessageType
 *
 *  Purpose : Get message type from raw OCPP message
 *
 *   Return : int - message type
 *			  error when cannot unmarshal message, nil otherwise
 *			  
*/
func (requestHandler *RequestHandler) getMessageType (rawMessage string) (int, error) {
	var typeID int
	parametersArray := []interface{}{&typeID}
	if err := json.Unmarshal([]byte(rawMessage), &parametersArray); err != nil {
		return 0, err
	}

	return typeID, nil
}

/****************************************************************************************
 *
 * Function : RequestHandler::HandleIncomeMessage
 *
 *  Purpose : Parse, unmarshal and validate request body from request
 *
 *    Input : rawMessage string - raw message to parse and validate
 *
 *   Return : string - response
 *			  error - if happened, nil otherwise
 *			  bool - true if needs to keep websocket open, false otherwise
*/
func (requestHandler *RequestHandler) HandleIncomeMessage(rawMessage string) (string, error, bool) {

	// First check for income message
	if rawMessage == "" {
		return "", errors.New("Body of the request is empty"), true
	}

	// Get message type from the raw text
	messageType, errMessageType := requestHandler.getMessageType(rawMessage)
	if errMessageType != nil {
		return "", errMessageType, true
	}

	// Handle Call message
	if messageType == int(messages.MESSAGE_TYPE_CALL) {
		// Create CallMessage obj from raw message
		callMessageObj := messages.CreateCallMessageCreator(rawMessage)

		return requestHandler.callRequestHandler(
			callMessageObj,
			callMessageObj.Action + REQUEST_TYPE_HANDLER)
	}

	// Handle Call Result message
	if messageType == int(messages.MESSAGE_TYPE_CALL_RESULT) {
		// Create CallResultMessage obj from raw message
		callResultObj := messages.CallResultMessageCreator(rawMessage)
		
		// To call correct CallResult handler we need action.
		// Action is not exist in CallResult message.
		// We are getting it from sent messages queue
		action := requestHandler.getActionByUniqueID(callResultObj.UniqueID)

		// Get Action from the queue by uniqueID, as call result has not included one
		return requestHandler.callResponseHandler(
			callResultObj,
			action + RESPONSE_TYPE_HANDLER)
	}

	// Handle Call Error message
	if messageType ==int(messages.MESSAGE_TYPE_CALL_ERROR) {
		// Create CallErrorMessage obj from raw message
		callErrorObj := messages.CallErrorMessageCreator(rawMessage)
		return requestHandler.callResponseHandler(
			callErrorObj,
			ERROR_TYPE_HANDLER)
	}

	//Error_Log("Unknown type of the request %d", TypeID)
	return "", errors.New("Request Handler is not found"), true
}






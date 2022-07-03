/*	==========================================================================
	OCPP 1.6 Protocol
	Filename: call.go
	Owner: Sergiy Safronov
	Purpose: File describes all routinse for the Call object of the OCPP protocol
	=============================================================================
*/

package messages

import (
	"encoding/json"
)


type MessageType int

const (
    MESSAGE_TYPE_CALL           MessageType = 2
)

/****************************************************************************************
 *	Struct 	: CallMessage
 * 
 * 	Purpose : Object handles the Call message structure
 *
*****************************************************************************************/
type CallMessage struct {
	UniqueID		string
	Action			string 
	Payload			map[string]interface{}
	Signature 		string
}

/****************************************************************************************
 *
 * Function : CallMessageConstructor (Constructor)
 *
 *  Purpose : Creates a new instance of the CallMessage object
 *
 *	Return : CallMessage object
*/
func CallMessageConstructor () CallMessage {
	//Define CallMessage object
	callMessageObj := CallMessage{}
	callMessageObj.init()
	return callMessageObj
}

/****************************************************************************************
 *
 * Function : CallMessage::createCallMessage (Constructor)
 *
 *  Purpose : Creates a new instance of the CallMessage object
 *
 *	Return : nil is cannot unmarshal message, otherwise CallMessage object
*/
func CreateCallMessageConstructor (rawMessage string) CallMessage {
	callMessageObj := CallMessage{}

	// Load JSON from string and check JSON structure
	if err := callMessageObj.unpackMessage(rawMessage); err != nil {
		return nil
	}

	if callMessageObj.MessageTypeID != MESSAGE_TYPE_CALL {
		return nil
	}

	// Validate 

	return callMessageObj
}

/****************************************************************************************
 *
 * Function : CallMessage::createCallMessage (Constructor)
 *
 *  Purpose : Creates a new instance of the CallMessage object
 *
 *	Return : nil is cannot unmarshal message, otherwise CallMessage object
*/
func CreateCallMessageConstructor (uniqueID string, action string, payload map[string]interface{}) CallMessage {
	callMessageObj := CallMessage{}
	callMessageObj.UniqueID = uniqueID
	callMessageObj.Action = action
	callMessageObj.Payload = payload

	return callMessageObj
}

/****************************************************************************************
 *
 * Function : CallMessage::init
 *
 *  Purpose : Initiate variables of the CallMessage structure
 *
 *	Return : Nothing
*/
func (callMessage *CallMessage) init() {
	callMessage.UniqueID = ""
	callMessage.Action = ""
	callMessage.Payload = make(map[string]interface{})
	callMessage.Signature = ""
}

/****************************************************************************************
 *
 * Function : CallMessage::addPayload
 *
 *  Purpose : Add data to the payload
 *
 *	Return : Nothing
*/
func (callMessage *CallMessage) addPayload (key string, value string) {

	callMessage.Payload[key] = value
}

/****************************************************************************************
 *
 * Function : CallMessage::getMessageType
 *
 *  Purpose : Return Call message type
 *
 *	Return : MessageType
*/
func (callMessage CallMessage) getMessageType () MessageType {
	return MESSAGE_TYPE_CALL
}

/****************************************************************************************
 *
 * Function : CallMessage::unpackMessage
 *
 *  Purpose : unmarshal request body to CallMessage struct
 *
 *	Return : error when cannot unmarshal message, otherwise nil
*/
func (callMessage *CallMessage) unpackMessage (raw_message string) error {

	parametersArray := []interface{}{
		&callMessage.MessageTypeID,
		&callMessage.UniqueID,
		&callMessage.Action,
		&callMessage.Payload,
		&callMessage.Signature
	}

	if err := json.Unmarshal([]byte(raw_message), &parametersArray); err != nil {
		return err
	}

	return nil
}

/****************************************************************************************
 *
 * Function : CallMessage::CreateCallMessage
 *
 *  Purpose : Initiate variables of the CallMessage structure
 *
 *	Return : Nothing
*/
func (callMessage *CallMessage) CreateCallMessage () (string,error) {
	var parametersArray []interface{}
	if cm.Signature == "" {
		parametersArray = []interface{}{
			&cm.MessageTypeID,
			&cm.UniqueID,
			&cm.Action,
			&cm.Payload
		}
	} else {
		parametersArray = []interface{}{
			&cm.MessageTypeID,
			&cm.UniqueID,
			&cm.Action,
			&cm.Payload,
			&cm.Signature
		}
	}

	jsonResult, err := json.Marshal(parametersArray)
	if err != nil {
		return ""
	}

	return string(jsonResult)
}
/*	==========================================================================
	OCPP 1.6 Protocol
	Filename: call.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/ocpp16-go/messages
	Purpose: File describes all routinse for the Call object of the OCPP protocol
	=============================================================================
*/

package messages

import (
	"encoding/json"
	"errors"
)

const (
	MESSAGE_TYPE_CALL MessageType = 2
)

/****************************************************************************************
 *	Struct 	: CallMessage
 *
 * 	Purpose : Object handles the Call message structure
 *
*****************************************************************************************/
type CallMessage struct {
	UniqueID  string
	Action    string
	Payload   map[string]interface{}
	Signature string
}

/****************************************************************************************
 *
 * Function : CallMessageConstructor (Constructor)
 *
 *  Purpose : Creates a new instance of the CallMessage object
 *
 *	Return : CallMessage object
 */
func CallMessageConstructor() CallMessage {
	//Define CallMessage object
	callMessageObj := CallMessage{}
	callMessageObj.init()
	return callMessageObj
}

/****************************************************************************************
 *
 * Function : CreateCallMessageCreator (Constructor)
 *
 *  Purpose : Creates a new instance of the CallMessage object using raw message
 *
 *    Input : rawMessage string - raw message to parse and validate
 *
 *	 Return : CallMessage object
 */
func CreateCallMessageCreator(rawMessage string) CallMessage {
	callMessageObj := CallMessage{}

	// Load JSON from string and check JSON structure
	if err := callMessageObj.unpackMessage(rawMessage); err != nil {
		return CallMessage{}
	}

	// Validate

	return callMessageObj
}

/****************************************************************************************
 *
 * Function : CallMessage::CreateCallMessageWithParam (Constructor)
 *
 *  Purpose : Creates a new instance of the CallMessage object with provided parameters
 *
 *	 Return : CallMessage object
 */
func CreateCallMessageWithParam(uniqueID string, action string, payload map[string]interface{}) CallMessage {
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
 * 	 Return : Nothing
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
 *    Input : 	key string - parameter name
 *				value string - parameter value
 *
 *	 Return : Nothing
 */
func (callMessage *CallMessage) addPayload(key string, value string) {

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
func (callMessage CallMessage) getMessageType() MessageType {
	return MESSAGE_TYPE_CALL
}

/****************************************************************************************
 *
 * Function : CallMessage::unpackMessage
 *
 *  Purpose : unmarshal message to create CallMessage struct
 *
 *    Input : 	rawMessage string - raw message to parse and validate
 *
 *	 Return : error when cannot unmarshal message, otherwise nil
 */
func (callMessage *CallMessage) unpackMessage(rawMessage string) error {
	var messageTypeID int
	parametersArray := []interface{}{
		&messageTypeID,
		&callMessage.UniqueID,
		&callMessage.Action,
		&callMessage.Payload,
		&callMessage.Signature,
	}

	if err := json.Unmarshal([]byte(rawMessage), &parametersArray); err != nil {
		return err
	}

	if messageTypeID != int(MESSAGE_TYPE_CALL) {
		return errors.New("Message type ID is not matched CallMessage")
	}

	return nil
}

/****************************************************************************************
 *
 * Function : CallMessage::ToString
 *
 *  Purpose : Convert CallMessage struct to string  message
 *
 *	 Return : string
 *			  error if happened, nil otherwise
 */
func (callMessage *CallMessage) ToString() (string, error) {
	messageType := int(MESSAGE_TYPE_CALL)
	var parametersArray []interface{}
	if callMessage.Signature == "" {
		parametersArray = []interface{}{
			&messageType,
			&callMessage.UniqueID,
			&callMessage.Action,
			&callMessage.Payload,
		}
	} else {
		parametersArray = []interface{}{
			&messageType,
			&callMessage.UniqueID,
			&callMessage.Action,
			&callMessage.Payload,
			&callMessage.Signature,
		}
	}

	jsonResult, err := json.Marshal(parametersArray)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}

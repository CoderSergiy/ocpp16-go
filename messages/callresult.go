/*	==========================================================================
	OCPP 1.6 Protocol
	Filename: callresult.go
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

const MESSAGE_TYPE_CALL_RESULT MessageType = 3

/****************************************************************************************
 *	Struct 	: CallResultMessage
 *
 * 	Purpose : Object handles the Call Result message structure
 *
*****************************************************************************************/
type CallResultMessage struct {
	UniqueID  string
	Payload   map[string]interface{}
	Signature string
}

/****************************************************************************************
 *
 * Function : CallResultMessageConstructor (Constructor)
 *
 *  Purpose : Creates a new instance of the CallResultMessage object
 *
 *	 Return : CallResultMessage
 */

func CallResultMessageConstructor() CallResultMessage {
	callResultObj := CallResultMessage{}
	callResultObj.init()
	return callResultObj
}

/****************************************************************************************
 *
 * Function : CreateCallResultMessage (Constructor)
 *
 *  Purpose : Creates a new instance of the CallResultMessage from raw message
 *
 *    Input : uniqueID string - id to create message
 *			  payload map[string]interface{} - payload of the message
 *
 *	 Return : CallResultMessage
 */
func CreateCallResultMessage(uniqueID string, payload map[string]interface{}) CallResultMessage {
	callResultObj := CallResultMessage{}
	callResultObj.UniqueID = uniqueID
	callResultObj.Payload = payload
	callResultObj.Signature = ""
	return callResultObj
}

/****************************************************************************************
 *
 * Function : CallResultMessage::init
 *
 *  Purpose : Initiate variables of the CallResultMessage structure
 *
 *	Return : Nothing
 */
func (callResultMessage *CallResultMessage) init() {
	callResultMessage.UniqueID = ""
	callResultMessage.Payload = make(map[string]interface{})
	callResultMessage.Signature = ""
}

/****************************************************************************************
 *
 * Function : CallResultMessageCreator (Constructor)
 *
 *  Purpose : Creates a new instance of the CallResultMessage from raw message
 *
 *    Input : rawMessage string - raw message to parse and validate
 *
 *	 Return : CallResultMessage
 */
func CallResultMessageCreator(rawMessage string) CallResultMessage {
	callResultObj := CallResultMessage{}

	// Load JSON from string and check JSON structure
	if err := callResultObj.unpackMessage(rawMessage); err != nil {
		return CallResultMessage{}
	}

	// Validate

	return callResultObj
}

/****************************************************************************************
 *
 * Function : CallResultMessage::getMessageType
 *
 *  Purpose : Return Call Result message type
 *
 *	Return : MessageType
 */
func (crm CallResultMessage) getMessageType() MessageType {
	return MESSAGE_TYPE_CALL_RESULT
}

/****************************************************************************************
 *
 * Function : CallResultMessage::unpackMessage
 *
 *  Purpose : Parse, unmarshal and validate request body from request
 *
 *    Input : rawMessage string - raw message to parse and validate
 *
 *	 Return : error when cannot unmarshal message, otherwise nil
 */
func (callResultMessage *CallResultMessage) unpackMessage(rawMessage string) error {
	var messageTypeID int
	parametersArray := []interface{}{
		&messageTypeID,
		&callResultMessage.UniqueID,
		&callResultMessage.Payload,
		&callResultMessage.Signature,
	}

	if err := json.Unmarshal([]byte(rawMessage), &parametersArray); err != nil {
		return err
	}

	if messageTypeID != int(MESSAGE_TYPE_CALL_RESULT) {
		return errors.New("Message type ID is not matched CallResultMessage")
	}

	return nil
}

/****************************************************************************************
 *
 * Function : CallResultMessage::ToString
 *
 *  Purpose : Convert CallResultMessage struct to string  message
 *
 *    Input : Nothing
 *
 *	 Return : string
 *			  error if happened, nil otherwise
 */
func (callResultMessage *CallResultMessage) ToString() (string, error) {
	var parametersArray []interface{}

	if callResultMessage.Signature == "" {
		parametersArray = []interface{}{
			int(MESSAGE_TYPE_CALL_RESULT),
			&callResultMessage.UniqueID,
			&callResultMessage.Payload,
		}
	} else {
		parametersArray = []interface{}{
			int(MESSAGE_TYPE_CALL_RESULT),
			&callResultMessage.UniqueID,
			&callResultMessage.Payload,
			&callResultMessage.Signature,
		}
	}

	jsonResult, err := json.Marshal(parametersArray)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}

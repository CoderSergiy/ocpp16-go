/*	==========================================================================
	OCPP 1.6 Protocol
	Filename: callerror.go
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

const MESSAGE_TYPE_CALL_ERROR MessageType = 4

/****************************************************************************************
 *	Struct 	: CallErrorMessage
 *
 * 	Purpose : Object handles the Call Error message structure
 *
*****************************************************************************************/
type CallErrorMessage struct {
	UniqueID         string
	ErrorCode        string
	ErrorDescription string
	ErrorDetails     map[string]interface{}
}

/****************************************************************************************
 *
 * Function : CallErrorMessageConstructor (Constructor)
 *
 *  Purpose : Creates a new instance of the CallErrorMessage object
 *
 *	Return : CallErrorMessage
 */
func CallErrorMessageConstructor() CallErrorMessage {
	callErrorObj := CallErrorMessage{}
	callErrorObj.init()
	return callErrorObj
}

/****************************************************************************************
 *
 * Function : CallErrorMessageCreator (Constructor)
 *
 *  Purpose : Creates a new instance of the CallErrorMessage object
 *
 *    Input : rawMessage string - raw message to parse and validate
 *
 *	 Return : CallErrorMessage
 */
func CallErrorMessageCreator(rawMessage string) CallErrorMessage {
	callErrorObj := CallErrorMessage{}

	// Load JSON from string and check JSON structure
	if err := callErrorObj.unpackMessage(rawMessage); err != nil {
		return CallErrorMessage{}
	}

	// Validate

	return callErrorObj
}

/****************************************************************************************
 *
 * Function : CallErrorMessage::init
 *
 *  Purpose : Initiate variables of the CallErrorMessage structure
 *
 *	 Return : Nothing
 */
func (cem *CallErrorMessage) init() {
	cem.UniqueID = ""
	cem.ErrorCode = ""
	cem.ErrorDescription = ""
	cem.ErrorDetails = make(map[string]interface{})
}

/****************************************************************************************
 *
 * Function : CallError::CreateCallError (Constructor)
 *
 *  Purpose : Creates a new instance of the CallMessage object
 *
 *	Return : nil is cannot unmarshal message, otherwise CallMessage object
 */
/*
func CreateCallError (uniqueID string, errorCode string, errorDescription string, errorDetails map[string]interface{}) string {
	ce := CallErrorConstructor()
	ce.UniqueID = uniqueID
	ce.ErrorCode = errorCode
	ce.ErrorDescription = errorDescription
	ce.ErrorDetails = errorDetails

	return ce.createCallErrorMessage()
}
*/
/****************************************************************************************
 *
 * Function : CallErrorMessage::getMessageType
 *
 *  Purpose : Return Call Error message type
 *
 *	 Return : MessageType
 */
func (callErrorMessage CallErrorMessage) getMessageType() MessageType {
	return MESSAGE_TYPE_CALL_ERROR
}

/****************************************************************************************
 *
 * Function : CallErrorMessage::unpackMessage
 *
 *  Purpose : Unmarshal message from raw format to CallMessage struct
 *
 *    Input : rawMessage string - raw message to parse and validate
 *
 *	 Return : error when cannot unmarshal message, otherwise nil
 */
func (callErrorMessage *CallErrorMessage) unpackMessage(raw_message string) error {
	var messageTypeID int
	tmp := []interface{}{
		&messageTypeID,
		&callErrorMessage.UniqueID,
		&callErrorMessage.ErrorCode,
		&callErrorMessage.ErrorDescription,
		&callErrorMessage.ErrorDetails,
	}

	if err := json.Unmarshal([]byte(raw_message), &tmp); err != nil {
		return err
	}

	if messageTypeID != int(MESSAGE_TYPE_CALL_ERROR) {
		return errors.New("Message type ID is not matched CallErrorMessage")
	}

	return nil
}

/****************************************************************************************
 *
 * Function : CallErrorMessage::ToString
 *
 *  Purpose : Convert CallErrorMessage struct to string  message
 *
 *    Input : Nothing
 *
 *	 Return : string
 *			  error if happened, nil otherwise
 */
func (callErrorMessage *CallErrorMessage) ToString() (string, error) {

	messageTypeID := MESSAGE_TYPE_CALL_ERROR
	var parametersArray []interface{}
	if len(callErrorMessage.ErrorDetails) > 0 {
		parametersArray = []interface{}{
			&messageTypeID,
			&callErrorMessage.UniqueID,
			&callErrorMessage.ErrorCode,
			&callErrorMessage.ErrorDescription,
			&callErrorMessage.ErrorDetails,
		}
	} else {
		parametersArray = []interface{}{
			&messageTypeID,
			&callErrorMessage.UniqueID,
			&callErrorMessage.ErrorCode,
			&callErrorMessage.ErrorDescription,
		}
	}

	jsonResult, err := json.Marshal(parametersArray)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}
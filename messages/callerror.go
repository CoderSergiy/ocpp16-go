/*	==========================================================================
	OCPP 1.6 Protocol
	Filename: callerror.go
	Owner: Sergiy Safronov
	Purpose: File describes all routinse for the Call object of the OCPP protocol
	=============================================================================
*/

package messages

const	MESSAGE_TYPE_CALL_ERROR    MessageType = 4

/****************************************************************************************
 *	Struct 	: CallErrorMessage
 * 
 * 	Purpose : Object handles the Call Error message structure
 *
*****************************************************************************************/
type CallErrorMessage struct {
	UniqueID			string
	ErrorCode			string
	ErrorDescription	string
	ErrorDetails		map[string]interface{}
}

/****************************************************************************************
 *
 * Function : CallErrorConstructor (Constructor)
 *
 *  Purpose : Creates a new instance of the CallResult object
 *
 *	Return : CallErrorMessage
*/
func CallErrorMessageConstructor () CallErrorMessage {
	callErrorObj := CallErrorMessage{}
	callErrorObj.init()
	return callErrorObj
}

/****************************************************************************************
 *
 * Function : CallErrorMessageConstructor (Constructor)
 *
 *  Purpose : Creates a new instance of the CallErrorMessage object
 *
 *	Return : CallErrorMessage
*/
func CallErrorMessageConstructor (rawMessage string) CallErrorMessage {
	callErrorObj := CallErrorMessage{}

	// Load JSON from string and check JSON structure
	if err := callErrorObj.unpackMessage(rawMessage); err != nil {
		return nil
	}

	if callErrorObj.MessageTypeID != MESSAGE_TYPE_CALL_ERROR {
		return nil
	}

	// Validate 

	return callErrorObj
}

/****************************************************************************************
 *
 * Function : CallErrorMessage::init
 *
 *  Purpose : Creates a new instance of the CallResult object
 *
 *	Return : CallErrorMessage object
*/
func (cem *CallErrorMessage) init () {
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
 *	Return : MessageType
*/
func (cem CallErrorMessage) getMessageType () MessageType {
	return MESSAGE_TYPE_CALL_ERROR
}

/****************************************************************************************
 *
 * Function : CallErrorMessage::unpackMessage
 *
 *  Purpose : unmarshal message from raw format to CallMessage struct
 *
 *	Return : error when cannot unmarshal message, otherwise nil
*/
func (callErrorMessage *CallErrorMessage) unpackMessage (raw_message string) error {

	tmp := []interface{}{
		&callErrorMessage.MessageTypeID,
		&callErrorMessage.UniqueID,
		&callErrorMessage.Action,
		&callErrorMessage.Payload,
		&callErrorMessage.Signature
	}

	if err := json.Unmarshal([]byte(raw_message), &tmp); err != nil {
		return err
	}

	return nil
}

/****************************************************************************************
 *
 * Function : CallError::createCallErrorMessage
 *
 *  Purpose : Initiate variables of the CallMessage structure
 *
 *	Return : Nothing
*/
func (callErrorMessage *CallErrorMessage) createCallErrorMessage () (string,error) {
	var parametersArray []interface{}
	if len(callErrorMessage.ErrorDetails) > 0 {
		parametersArray = []interface{}{
			&callErrorMessage.MessageTypeID,
			&callErrorMessage.UniqueID,
			&callErrorMessage.ErrorCode,
			&callErrorMessage.ErrorDescription,
			&callErrorMessage.ErrorDetails
		}
	} else {
		parametersArray = []interface{}{
			&callErrorMessage.MessageTypeID,
			&callErrorMessage.UniqueID,
			&callErrorMessage.ErrorCode,
			&callErrorMessage.ErrorDescription}
	}

	jsonResult, err := json.Marshal(parametersArray)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}
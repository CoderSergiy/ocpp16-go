/*	==========================================================================
	OCPP 1.6 Protocol
	Filename: callresult.go
	Owner: Sergiy Safronov
	Purpose: File describes all routinse for the Call object of the OCPP protocol
	=============================================================================
*/

package messages

const	MESSAGE_TYPE_CALL_RESULT    MessageType = 3

/****************************************************************************************
 *	Struct 	: CallResultMessage
 * 
 * 	Purpose : Object handles the Call Result message structure
 *
*****************************************************************************************/
type CallResultMessage struct {
	UniqueID		string
	Payload			map[string]interface{}
	Signature		string
}

/****************************************************************************************
 *
 * Function : CallResultMessage::CallResultConstructor (Constructor)
 *
 *  Purpose : Creates a new instance of the CallResult message object
 *
 *	Return : CallResultMessage
*/
func CallResultMessageConstructor () CallResultMessage {
	callResultObj := CallResultMessage{}
	callResultObj.init()
	return callResultObj
}

/****************************************************************************************
 *
 * Function : CallResultMessageConstructor (Constructor)
 *
 *  Purpose : Creates a new instance of the CallResultMessage object
 *
 *	Return : CallResultMessage
*/
func CallResultMessageConstructor (rawMessage string) CallResultMessage {
	callResultObj := CallResultMessage{}

	// Load JSON from string and check JSON structure
	if err := callResultObj.unpackMessage(rawMessage); err != nil {
		return nil
	}

	if callResultObj.MessageTypeID != MESSAGE_TYPE_CALL_RESULT {
		return nil
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
func (crm CallResultMessage) getMessageType () MessageType {
	return MESSAGE_TYPE_CALL_RESULT
}

/****************************************************************************************
 *
 * Function : CallResultMessage::unpackMessage
 *
 *  Purpose : Parse, unmarshal and validate request body from request
 *
 *	Return : error when cannot unmarshal message, otherwise nil
*/
func (crm *CallResultMessage) unpackMessage (raw_message string) error {

	parametersArray := []interface{}{
		&crm.MessageTypeID,
		&crm.UniqueID,
		&crm.Payload,
		&crm.Signature
	}

	if err := json.Unmarshal([]byte(raw_message), &parametersArray); err != nil {
		return err
	}

	return nil
}
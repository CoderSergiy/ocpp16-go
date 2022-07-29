/*	==========================================================================
	OCPP 1.6 Protocol
	Filename: trigger_message.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/ocpp16-go/core
	Purpose: Describe all methods to work with BootNotification OCPP message
	=============================================================================
*/

package core

type TriggerMessageStatus string
type TriggerMessageType string

const (
	// TriggerMessageStatuses
	TriggerMessageStatusAccepted       TriggerMessageStatus = "Accepted"
	TriggerMessageStatusRejected       TriggerMessageStatus = "Rejected"
	TriggerMessageStatusNotImplemented TriggerMessageStatus = "NotImplemented"
	// TriggerMessageTypes
	TriggerMessageTypeBootNotification              TriggerMessageType = "BootNotification"
	TriggerMessageTypeDiagnosticsStatusNotification TriggerMessageType = "DiagnosticsStatusNotification"
	TriggerMessageTypeFirmwareStatusNotification    TriggerMessageType = "FirmwareStatusNotification"
	TriggerMessageTypeHeartbeat                     TriggerMessageType = "Heartbeat"
	TriggerMessageTypeMeterValues                   TriggerMessageType = "MeterValues"
	TriggerMessageTypeStatusNotification            TriggerMessageType = "StatusNotification"

	ACTION_TRIGGERMESSAGE string = "TriggerMessage"
)

/****************************************************************************************
 *	Struct 	: TriggerMessageRequestPayload
 *
 * 	Purpose : Handles parameters and methods for current struct
 *
*****************************************************************************************/
type TriggerMessageRequestPayload struct {
	requestedMessage TriggerMessageType //`json:"requestedMessage"`
	connectorId      int                //`json:"connectorId"`
}

/****************************************************************************************
 *
 * Function : CreateTriggerMessageRequestPayload (Constructor)
 *
 *  Purpose : Creates a new instance of the TriggerMessageRequestPayload object with specified values
 *
 *    Input : reqMessageType TriggerMessageType - type of the Trigger MEssage request
 *			  connectorId int - connector of the Charge Point
 *
 *	 Return : TriggerMessageRequestPayload object
 */
func CreateTriggerMessageRequestPayload(reqMessageType TriggerMessageType, connectorId int) TriggerMessageRequestPayload {
	triggerMessageRequestPayload := TriggerMessageRequestPayload{}

	triggerMessageRequestPayload.requestedMessage = reqMessageType
	triggerMessageRequestPayload.connectorId = connectorId

	return triggerMessageRequestPayload
}

/****************************************************************************************
 *
 * Function : TriggerMessageRequestPayload::GetPayload
 *
 *  Purpose : Generate payload using TriggerMessageRequestPayload struct
 *
 *	  Input : Nothing
 *
 *	 Return : map[string]interface{} - map of the payloads values
 */
func (triggerMessageRequestPayload *TriggerMessageRequestPayload) GetPayload() map[string]interface{} {

	payload := make(map[string]interface{})
	payload["requestedMessage"] = string(triggerMessageRequestPayload.requestedMessage)
	if triggerMessageRequestPayload.connectorId > 0 {
		payload["connectorId"] = triggerMessageRequestPayload.connectorId
	}

	return payload
}

/****************************************************************************************
 *	Struct 	: TriggerMessageResponsePayload
 *
 * 	Purpose : Handles parameters and methods for current struct
 *
*****************************************************************************************/
type TriggerMessageResponsePayload struct {
	status TriggerMessageStatus //`json:"requestedMessage"`
}

/****************************************************************************************
 *
 * Function : SanitizeTriggerMessageType
 *
 *  Purpose : Sanitize TriggerMessage type - if it matches API list
 *
 *	  Input : triggerMessageType string - type in string format to sanitize
 *
 *	 Return : true - when type is matches, otherwise false
 */
func SanitizeTriggerMessageType(triggerMessageType string) bool {

	switch triggerMessageType {
	case string(TriggerMessageTypeBootNotification),
		string(TriggerMessageTypeDiagnosticsStatusNotification),
		string(TriggerMessageTypeFirmwareStatusNotification),
		string(TriggerMessageTypeHeartbeat),
		string(TriggerMessageTypeMeterValues),
		string(TriggerMessageTypeStatusNotification):
		return true
	}

	return false
}

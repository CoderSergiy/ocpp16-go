/*	==========================================================================
	OCPP 1.6 Protocol
	Filename: message.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/ocpp16-go/messages
	Purpose: Interface to all message's structs
	=============================================================================
*/

package messages

import (
	"encoding/json"
)

type MessageType int

/****************************************************************************************
 *	Interface : Message
 *
 * 	  Purpose : Interface for Call, CallResult and CallError
 *
*****************************************************************************************/
type Message interface {
	getMessageType() MessageType
	ToString() (string, error)
}

/****************************************************************************************
 *
 * Function : GetMessageTypeFromRaw
 *
 *  Purpose : Get message type and uniqueid from raw OCPP message
 *
 *    Input : rawMessage string - raw OCPP message
 *
 *   Return : string - uniqueID
 *            int - message type
 *            error when cannot unmarshal message, nil otherwise
 *
 */
func GetMessageTypeFromRaw(rawMessage string) (int, string, error) {
	var uniqueID string
	var typeID int
	parametersArray := []interface{}{&typeID, &uniqueID}
	if err := json.Unmarshal([]byte(rawMessage), &parametersArray); err != nil {
		return 0, "", err
	}

	return typeID, uniqueID, nil
}

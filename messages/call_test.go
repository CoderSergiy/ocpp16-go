/*	==========================================================================
	OCPP 1.6 Protocol
	Filename: call_test.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/ocpp16-go/messages
	Purpose: File with test cases for CallMessage
	=============================================================================
*/

package messages

import (
	"fmt"
	"testing"
)

/****************************************************************************************
 *
 * Function : TestCallStruct
 *
 *  Purpose : Method to test initiation of the CallMessage struct
 *
 *   Return : Nothing
 */
func TestCallStruct(t *testing.T) {

	callMessageObj := CallMessageConstructor()

	if callMessageObj.getMessageType() != 2 {
		t.Error(fmt.Printf("Wrong CallMessage type: '%v' instead of 2", callMessageObj.getMessageType()))
	}
}

/****************************************************************************************
 *
 * Function : TestIncomeMessage
 *
 *  Purpose : Method to test income Call message parsing and validation
 *
 *   Return : Nothing
 */
func TestIncomeMessage(t *testing.T) {

	callMessageObj := CreateCallMessageCreator("[2,\"A123.234\",\"BootNotification\",{\"chargePointModel\":\"SingleSocketCharger\",\"chargePointVendor\":\"VendorX\"}]")

	if callMessageObj.getMessageType() != 2 {
		t.Error(fmt.Printf("Wrong CallMessage type: '%v' instead of 2", callMessageObj.getMessageType()))
	}

	if callMessageObj.UniqueID != "A123.234" {
		t.Error(fmt.Printf("Wrong UniqueID : '%v' instead of 'A123.234'", callMessageObj.UniqueID))
	}

	if callMessageObj.Action != "BootNotification" {
		t.Error(fmt.Printf("Wrong Action : '%v' instead of 'BootNotification'", callMessageObj.UniqueID))
	}

	if callMessageObj.Signature != "" {
		t.Error("Signature is not empty")
	}
}

/****************************************************************************************
 *
 * Function : TestOutgoingMessage
 *
 *  Purpose : Method to test generating Call message
 *
 *   Return : Nothing
 */
func TestOutgoingMessage(t *testing.T) {

	payload := make(map[string]interface{})

	// Create CallResult message
	callMessageResponse := CreateCallMessageWithParam(
		"29591-56097986-1",
		"BootNotification",
		payload,
	)
	callMessageResponse.addPayload("chargePointModel", "SingleSocketCharger")
	callMessageResponse.addPayload("chargePointVendor", "VendorX")

	generatedMessage, err := callMessageResponse.ToString()
	if err != nil {
		t.Error(fmt.Printf("Error when generating message '%v'", err))
		return
	}

	callMessageText := "[2,\"29591-56097986-1\",\"BootNotification\",{\"chargePointModel\":\"SingleSocketCharger\",\"chargePointVendor\":\"VendorX\"}]"

	if generatedMessage != callMessageText {
		t.Error(fmt.Printf("Generated message '%v' is not matched expected '%v'", generatedMessage, callMessageText))
	}
}

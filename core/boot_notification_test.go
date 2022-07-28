/*	==========================================================================
	OCPP 1.6 Protocol
	Filename: boot_notification_test.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/ocpp16-go/core
	Purpose: File with test cases for CallMessage
	=============================================================================
*/

package core

import (
	"fmt"
	"github.com/CoderSergiy/ocpp16-go/messages"
	"testing"
	"time"
)

/****************************************************************************************
 *
 * Function : TestBootNotificationResponse
 *
 *  Purpose : Test BootNotification response, including payload
 *
 *   Return : Nothing
 */
func TestBootNotificationResponse(t *testing.T) {

	uniqueID := "jnfow234-345mkregm"

	//Create payload
	bootNotificationRespPayload := BootNotificationResponsePayload{
		status:            RegistrationStatusPending,
		heartbeatInterval: 300,
		currentTime:       time.Now().Format("2006-01-02 15:04:05.000"),
	}

	bootNotificationResp := messages.CreateCallResultMessage(uniqueID, bootNotificationRespPayload.GetPayload())

	messageStr, err := bootNotificationResp.ToString()
	if err != nil {
		t.Error(fmt.Printf("Error when generating message '%v'", err))
		return
	}

	t.Log(fmt.Printf("Success '%v'", messageStr))
}

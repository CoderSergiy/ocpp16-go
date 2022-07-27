/*	==========================================================================
	OCPP 1.6 Protocol
	Filename: boot_notification.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/ocpp16-go/core
	Purpose: Describe all methods to work with BootNotification OCPP message
	=============================================================================
*/

package core

type RegistrationStatus string

const (
	RegistrationStatusAccepted RegistrationStatus = "Accepted"
	RegistrationStatusPending  RegistrationStatus = "Pending"
	RegistrationStatusRejected RegistrationStatus = "Rejected"

	ACTION_BOOTNOTIFICATION string = "BootNotification"
)

/****************************************************************************************
 *	Struct 	: BootNotificationResponse
 *
 * 	Purpose :   Object stores callback handlers.
 *				Using different struct to isolate callback methods from RequestHandler
 *
*****************************************************************************************/
type BootNotificationResponse struct {
	CurrentTime       string             //`json:"currentTime"`
	HeartbeatInterval int                //`json:"heartbeatInterval"`
	Status            RegistrationStatus //`json:"status"`
}

/****************************************************************************************
 *
 * Function : bootNotificationResponse::GetPayload
 *
 *  Purpose : Generate payload from BootNotificationResponse struct
 *
 *	  Input : Nothing
 *
 *	 Return : map[string]interface{} - map of the
 */
func (bootNotificationResponse *BootNotificationResponse) GetPayload() map[string]interface{} {

	payload := make(map[string]interface{})
	payload["status"] = string(bootNotificationResponse.Status)
	payload["currentTime"] = bootNotificationResponse.CurrentTime
	payload["heartbeatInterval"] = bootNotificationResponse.HeartbeatInterval

	return payload

}

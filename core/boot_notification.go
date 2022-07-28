/*	==========================================================================
	OCPP 1.6 Protocol
	Filename: boot_notification.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/ocpp16-go/core
	Purpose: Describe all methods to work with BootNotification OCPP message
	=============================================================================
*/

package core

import (
	"time"
)

type RegistrationStatus string

const (
	RegistrationStatusAccepted RegistrationStatus = "Accepted"
	RegistrationStatusPending  RegistrationStatus = "Pending"
	RegistrationStatusRejected RegistrationStatus = "Rejected"

	ACTION_BOOTNOTIFICATION string = "BootNotification"
)

/****************************************************************************************
 *	Struct 	: BootNotificationResponsePayload
 *
 * 	Purpose :   Object stores callback handlers.
 *				Using different struct to isolate callback methods from RequestHandler
 *
*****************************************************************************************/
type BootNotificationResponsePayload struct {
	currentTime       string             //`json:"currentTime"`
	heartbeatInterval int                //`json:"heartbeatInterval"`
	status            RegistrationStatus //`json:"status"`
}

/****************************************************************************************
 *
 * Function : CreateBootNotificationResponsePayload (Constructor)
 *
 *  Purpose : Creates a new instance of the BootNotificationResponsePayload object with default values
 *
 *    Input : status RegistrationStatus - new status
 *
 *	 Return : BootNotificationResponsePayload object
 */
func CreateBootNotificationResponsePayload(status RegistrationStatus) BootNotificationResponsePayload {
	bootNotificationRespPayload := BootNotificationResponsePayload{}

	bootNotificationRespPayload.status = status
	bootNotificationRespPayload.heartbeatInterval = 300
	bootNotificationRespPayload.currentTime = time.Now().Format("2006-01-02 15:04:05.000")

	return bootNotificationRespPayload
}

/****************************************************************************************
 *
 * Function : BootNotificationResponsePayload::GetPayload
 *
 *  Purpose : Generate payload from BootNotificationResponsePayload struct
 *
 *	  Input : Nothing
 *
 *	 Return : map[string]interface{} - map of the
 */
func (bootNotificationResPayload *BootNotificationResponsePayload) GetPayload() map[string]interface{} {

	payload := make(map[string]interface{})
	payload["status"] = string(bootNotificationResPayload.status)
	payload["currentTime"] = bootNotificationResPayload.currentTime
	payload["heartbeatInterval"] = bootNotificationResPayload.heartbeatInterval

	return payload
}

/****************************************************************************************
 *	Struct 	: BootNotificationRequestPayload
 *
 * 	Purpose :   Object stores callback handlers.
 *				Using different struct to isolate callback methods from RequestHandler
 *
*****************************************************************************************/
type BootNotificationRequestPayload struct {
	chargeBoxSerialNumber   string //`json:"chargeBoxSerialNumber"`
	chargePointModel        string //`json:"chargePointModel"`
	chargePointSerialNumber string //`json:"chargePointSerialNumber"`
	chargePointVendor       string //`json:"chargePointVendor"`
	firmwareVersion         string //`json:"firmwareVersion"`
	iccid                   string //`json:"iccid"`
	imsi                    string //`json:"imsi"`
	meterSerialNumber       string //`json:"meterSerialNumber"`
}

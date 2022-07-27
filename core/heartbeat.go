/*	==========================================================================
	OCPP 1.6 Protocol
	Filename: heartbeat.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/ocpp16-go/core
	Purpose: Describe all methods to work with HeartBeat OCPP message
	=============================================================================
*/

package core

import (
	"time"
)

const (
	ACTION_HEARTBEAT string = "Heartbeat"
)

/****************************************************************************************
 *	Struct 	: HeartBeatResponse
 *
 * 	Purpose : Structure to work with HeartBeat response
 *
*****************************************************************************************/
type HeartBeatResponse struct {
	CurrentTime string //`json:"currentTime"`
}

/****************************************************************************************
 *
 * Function : HeartBeatResponse::Init
 *
 *  Purpose : Initiate variables of the HeartBeatResponse structure
 *
 *	  Input : Nothing
 *
 *	 Return : Nothing
 */
func (heartBeatResponse *HeartBeatResponse) Init() {
	heartBeatResponse.CurrentTime = time.Now().Format("2006-01-02 15:04:05.000")
}

/****************************************************************************************
 *
 * Function : HeartBeatResponse::GetPayload
 *
 *  Purpose : Generate payload from HeartBeatResponse struct
 *
 *	  Input : Nothing
 *
 *	 Return : map[string]interface{} - map of the
 */
func (heartBeatResponse *HeartBeatResponse) GetPayload() map[string]interface{} {

	payload := make(map[string]interface{})
	payload["currentTime"] = heartBeatResponse.CurrentTime

	return payload

}

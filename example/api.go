/*	==========================================================================
	OCPP 1.6 Protocol
	Filename: api.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/ocpp16-go/example
	Purpose: API to control the server by clients
			 File includes APIs:
				- messageStatusHandler
				- triggerActionHandler
				- chargerStatusHandler
	=============================================================================
*/

package example

import (
	"encoding/json"
	"github.com/CoderSergiy/golib/logging"
	"github.com/CoderSergiy/ocpp16-go/core"
	"github.com/CoderSergiy/ocpp16-go/messages"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

/****************************************************************************************
 *
 * Function : GetMessageStatusAPI
 *
 *  Purpose : Get message from the queue and send to the client by http
 *
 *    Input : chargerName string - charger name
 *            serverConfigs *Configs - pointer to the chargers arrays
 *            log *logging.Log - pointer to the log
 *            w http.ResponseWriter - http response
 *
 *   Return : Nothing
 */
func GetChargerStatusAPI(chargerName string, serverConfigs *Configs, log *logging.Log, w http.ResponseWriter) {
	log.Info_Log("GetChargerStatusAPI")

	if chargerName == "" {
		log.Error_Log("chargerName parameter is empty")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Info_Log("ChargerName is '%v'", chargerName)

	// Get Charger from the Configs
	chargerObj, err := serverConfigs.GetChargerObj(chargerName)
	if err != nil || chargerObj == nil {
		// There is no charger with specified name in the configs
		log.Error_Log("GetChargerObj for '%v' returns error '%v'", chargerName, err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	jsonResult, err := json.Marshal(chargerObj)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Error_Log("[%s] Cannot marshal charger obj", chargerName)
		return
	}

	// Send response in json format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResult)
}

/****************************************************************************************
 *
 * Function : GetMessageStatusAPI
 *
 *  Purpose : Get message from the queue and send to the client by http
 *
 *    Input : reference string - unique reference
 *            MQueue *SimpleMessageQueue - pointer to the Message Queue
 *            log *logging.Log - pointer to the log
 *            w http.ResponseWriter - http response
 *
 *   Return : Nothing
 */
func GetMessageStatusAPI(reference string, MQueue *SimpleMessageQueue, log *logging.Log, w http.ResponseWriter) {
	log.Info_Log("[%s] GetMessageStatus", reference)

	if reference == "" {
		log.Error_Log("Reference parameter is empty")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Info_Log("Reference is '%v'", reference)

	message, success := MQueue.GetMessage(reference)
	if !success {
		http.Error(w, string(CreateFailResponse("Message is not exist")), http.StatusOK)
		log.Error_Log("[%s] Message is not exists in the queue with uniqueID", reference)
		return
	}

	jsonResult, err := json.Marshal(message)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Error_Log("[%s] Cannot marshal message", reference)
		return
	}

	// Send response in json format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResult)
}

/****************************************************************************************
 *
 * Function : TriggerActionAPI
 *
 *  Purpose : Handles TriggerAction API request
 *
 *    Input : serverConfigs *Configs - pointer to the chargers arrays
 *            MQueue *SimpleMessageQueue - pointer to the Message Queue
 *            log *logging.Log - pointer to the log
 *            ps httprouter.Params - router parameters
 *            w http.ResponseWriter - http response
 *
 *   Return : Nothing
 */
func TriggerActionAPI(serverConfigs *Configs, MQueue *SimpleMessageQueue, log *logging.Log, ps httprouter.Params, w http.ResponseWriter) {
	log.Info_Log("TriggerActionAPI")

	chargerName := ps.ByName("chargerName")
	action := ps.ByName("action")

	if chargerName == "" || action == "" {
		log.Error_Log("One of the parameters are empty charger '%v' action '%v'", chargerName, action)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Info_Log("Charger name '%v' and action '%v'", chargerName, action)

	// Get Charger from the Configs
	chargerObj, err := serverConfigs.GetChargerObj(chargerName)
	if err != nil || chargerObj == nil {
		// There is no charger with specified name in the configs
		log.Error_Log("GetChargerObj for '%v' returns error '%v'", chargerName, err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Sanitize the TriggerMessage type from the request
	if !core.SanitizeTriggerMessageType(action) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Error_Log("TriggerMessage type '%v' is not supported", action)
		return
	}

	id := uuid.New()
	// Generate Call request payload for the TriggerMessage
	triggerMessagePayload := core.CreateTriggerMessageRequestPayload(core.TriggerMessageType(action), 0)
	// Generate Call request to the charger
	callMessageRequest := messages.CreateCallMessage(
		id.String(),
		core.ACTION_TRIGGERMESSAGE,
		triggerMessagePayload.GetPayload(),
	)

	// Convert Call message to string
	callMessageString, messageErr := callMessageRequest.ToString()
	if messageErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Error_Log("Error to generate callMessage: '%v'", messageErr)
		return
	}

	// Create message for the queue
	queueMessage := Message{Action: action, Sent: callMessageString, Status: MESSAGE_TYPE_NEW, Received: ""}
	// Add message to the queue
	addingErr := MQueue.Add(callMessageRequest.UniqueID, queueMessage)
	if addingErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Error_Log("Error to add message to the queue: '%v'", addingErr)
		return
	}

	// Send to write goroutine message's uniqueid
	(*chargerObj).WriteChannel <- callMessageRequest.UniqueID

	// Send response in json format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(CreateSuccessResponse(callMessageRequest.UniqueID))
}

/****************************************************************************************
 *	Struct 	: APIResponse
 *
 * 	Purpose : Object handles API response parameters
 *
*****************************************************************************************/
type APIResponse struct {
	Status      string `json:"status"`
	Reference   string `json:"reference"`
	Description string `json:"description"`
}

/****************************************************************************************
 *
 * Function : APIResponse::toBytes
 *
 *  Purpose : Convert APIResponse to json format
 *
 *	  Input : Nothing
 *
 *	 Return : []byte - json format of the APIResponse
 */
func (apiResponse *APIResponse) toBytes() []byte {
	jsonResult, err := json.Marshal(apiResponse)
	if err != nil {
		return []byte("Internal Application Error")
	}

	return jsonResult
}

/****************************************************************************************
 *
 * Function : CreateSuccessResponse
 *
 *  Purpose : Create Successfull API response in json format
 *
 *	  Input : reference string - reference for the response
 *
 *	 Return : []byte - json format of the APIResponse
 */
func CreateSuccessResponse(reference string) []byte {
	apiResponse := APIResponse{}

	apiResponse.Status = "success"
	apiResponse.Description = ""
	apiResponse.Reference = reference

	return apiResponse.toBytes()
}

/****************************************************************************************
 *
 * Function : CreateFailResponse
 *
 *  Purpose : Create Failed API response in json format
 *
 *	  Input : description string - description of the failed status
 *
 *	 Return : string - json format of the APIResponse
 */
func CreateFailResponse(description string) string {
	apiResponse := APIResponse{}

	apiResponse.Status = "fail"
	apiResponse.Description = description
	apiResponse.Reference = ""

	return string(apiResponse.toBytes())
}

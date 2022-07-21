/*	==========================================================================
	OCPP 1.6 Protocol
	Filename: simplequeue.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/ocpp16-go/example
	Purpose: All necassary methods to organise messages' queue
			 Current structures using for a demo purpose only
	=============================================================================
*/

package example

import (
	"fmt"
	"io/ioutil"
	"sync"
	"errors"
	"encoding/json"
)


type QueueMessageType int

const (
    MESSAGE_TYPE_SEND           QueueMessageType = 1
    MESSAGE_TYPE_SENT           QueueMessageType = 2
    MESSAGE_TYPE_RECEIVED       QueueMessageType = 3
)

/****************************************************************************************
 *	Struct 	: Message
 * 
 * 	Purpose : Struct handles message's parameters
 *
*****************************************************************************************/
type Message struct {
	Action		string	// Action of the message
	Content		string	// Message content
}

/****************************************************************************************
 *	Struct 	: SimpleMessageQueue
 * 
 * 	Purpose : Struct handles messages queue routines
 *
*****************************************************************************************/
type SimpleMessageQueue struct {
	MaxSize			int
	MessageQueue	map[string]Message
	queueMux 		sync.Mutex
}

/****************************************************************************************
 *
 * Function : SimpleMessageQueue::init
 *
 *  Purpose : Initiate variables of the SimpleMessageQueue structure
 *
 *	  Input : Nothing
 *
 *	 Return : Nothing
*/
func (queue *SimpleMessageQueue) init () {
	queue.MaxSize = 10
	queue.MessageQueue = make(map[string]Message)
}

/****************************************************************************************
 *
 * Function : SimpleMessageQueue::Add
 *
 *  Purpose : Add message to the queue
 *
 *    Input : uniqueID string - id of the message
 *			  message Message - raw message to parse and validate
 *
 *   Return : error - if happened, nil otherwise
 *
*/
func (queue *SimpleMessageQueue) Add (uniqueID string, message Message) error {
	// Lock the queue before any changes
	queue.queueMux.Lock()
	defer queue.queueMux.Unlock()

	// Check if not reached the max size
	if len(queue.MessageQueue) > queue.MaxSize {
		return errors.New("Reached max queue size")
	}

	// Add message to the queue
	queue.MessageQueue[uniqueID] = message

	return nil
}

/****************************************************************************************
 *
 * Function : SimpleMessageQueue::DeleteByUniqueID
 *
 *  Purpose : Delete message from the queue by unique id
 *
 *    Input : uniqueID string - id of the message
 *
 *   Return : error - if happened, nil otherwise
 *
*/
func (queue *SimpleMessageQueue) DeleteByUniqueID (uniqueID string) error {
	// Lock the queue before any changes
	queue.queueMux.Lock()
	defer queue.queueMux.Unlock()

	// If unique id exists in the queue - delete it
	if _, isKeyPresent := queue.MessageQueue[uniqueID]; isKeyPresent {
		delete(queue.MessageQueue, uniqueID)
		return nil
	}

	// Otherwise return an error
	return errors.New("DeleteByUniqueID. Message with pointed uniqueID is not exists")
}

/****************************************************************************************
 *
 * Function : SimpleMessageQueue::GetMessage
 *
 *  Purpose : Get message from the queue by unique id
 *
 *    Input : uniqueID string - id of the message
 *
 *   Return : Message
 *			  bool - true when message exists, false otherwise
 *
*/
func (queue *SimpleMessageQueue) GetMessage (uniqueID string) (Message, bool) {
	// Check if uniqueID is exists in the queue
	if message, isKeyPresent := queue.MessageQueue[uniqueID]; isKeyPresent {
		// Return message
		return message, true
	}

	return Message{}, false
}

/****************************************************************************************
 *
 * Function : SimpleMessageQueue::printStatus
 *
 *  Purpose : Print parameters of the queue
 *
 *    Input : Nothing
 *
 *   Return : string - Parameters of the queue in the text format
 *
*/
func (queue *SimpleMessageQueue) printStatus () string {
	return fmt.Sprintf("Size of the queue is %v where max size set to %v", len(queue.MessageQueue), queue.MaxSize)
}







/****************************************************************************************
 *	Struct 	: Charger
 * 
 * 	Purpose : Struct handles charger parameters in the gorutines
 *
*****************************************************************************************/
type Charger struct {
	AuthToken			string
	HeartBeatInterval	int
	AuthConnection		bool
	WebSocketConnected	bool
	InboundIP			string
}

/****************************************************************************************
 *	Struct 	: Configs
 * 
 * 	Purpose : Object handles configurations from the file
 *
*****************************************************************************************/
type Configs struct {
	Chargers			map[string]Charger	`json:"Chargers"`
	MaxQueueSize		int					`json:"MaxQueueSize"`
}

/****************************************************************************************
 *
 * Function : ServerConfigsConstructor (Constructor)
 *
 *  Purpose : Creates a new instance of the Configs
 *
 *	  Input : Nothing
 *
 *	Return : Configs object
*/
func ServerConfigsConstructor() Configs {
	configs := Configs{}
	configs.init()
	return configs
}

/****************************************************************************************
 *
 * Function : Configs::init
 *
 *  Purpose : Initiate variables of the Configs structure
 *
 *	  Input : Nothing
 *
 *	 Return : Nothing
*/
func (conf *Configs) init() {
	conf.Chargers = make(map[string]Charger)
	conf.MaxQueueSize = 10
}

/****************************************************************************************
 *
 * Function : Configs::GetChargerObj
 *
 *  Purpose : Get charger from the configuration structure
 *
 *	  Input : chargerName string - Name of the charger in the queue
 *
 *	 Return : Charger - charger object
 * 			  error - error if happened
*/
func (conf *Configs) GetChargerObj (chargerName string) (Charger, error) {

	if charger, isKeyPresent := conf.Chargers[chargerName]; isKeyPresent {
		// Requested charger is exists in the configs
		return charger, nil
	}

	// Charger details is not exists in the configs
	return Charger{}, errors.New("Charger is not exists in configs")

}



/****************************************************************************************
 *	Struct 	: ChargerFromFile
 * 
 * 	Purpose : Struct handles charger's parameters from config file
 *
*****************************************************************************************/
type ChargerFromFile struct {
	Name				string 	`json:"Name"`
	Authorization		string 	`json:"Authorization"`
	HeartBeatInterval	int 	`json:"HeartBeatInterval"`
}

/****************************************************************************************
 *	Struct 	: FileConfigs
 * 
 * 	Purpose : Struct handles configurations from file
 *
*****************************************************************************************/
type FileConfigs struct {
	Chargers			[]ChargerFromFile	`json:"Chargers"`
	MaxQueueSize		int					`json:"MaxQueueSize"`
}

/****************************************************************************************
 *
 * Function : SetConfigsFromFile (Constructor)
 *
 *  Purpose : Set configs struct from the file
 *
 *	  Input : fileName string - filename with settings for the test server
 *
 *	 Return : Configs - Configs object
 * 			  error - error if happened
*/
func SetConfigsFromFile(fileName string) (Configs, error) {
	configs := Configs{}
	configs.init()

	// Check if filename is empty
	if fileName == "" {
		return configs, errors.New("Filename is empty")
	}
	
	// Read file context to the buffer
    fileContentBytes, fileError := ioutil.ReadFile(fileName)
    if fileError != nil {
        return configs, fileError
    }

	// Unmarshal the content of the file to the Configs struct
	conf := FileConfigs{}
	if err := json.Unmarshal(fileContentBytes, &conf); err != nil {
		return configs, err
	}

	configs.MaxQueueSize = conf.MaxQueueSize

	for _, charger := range conf.Chargers {
		chargerConf := Charger{}
		chargerConf.AuthToken = charger.Authorization
		chargerConf.HeartBeatInterval = charger.HeartBeatInterval

		configs.Chargers[charger.Name] = chargerConf
	}

	return configs, nil
}



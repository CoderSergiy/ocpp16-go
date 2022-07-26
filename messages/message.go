/*	==========================================================================
	OCPP 1.6 Protocol
	Filename: message.go
	Owner: Sergiy Safronov
	Source : github.com/CoderSergiy/ocpp16-go/messages
	Purpose: Interface to all message's structs
	=============================================================================
*/

package messages

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
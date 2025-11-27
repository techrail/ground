// Package stdresponse contains some structures which can be used to send some standard responses
package stdresponse

// StatusMsgResponse is for reporting the status of a service
type StatusMsgResponse struct {
	Status      string         `json:"status"`
	CurrUtcTime string         `json:"currUtcTime"`
	ServiceName string         `json:"serviceName"`
	Details     map[string]any `json:"details,omitempty"`
}

// SingleMessage is for responding with a single string message
type SingleMessage struct {
	Message string `json:"message"`
}

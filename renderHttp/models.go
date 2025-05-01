package render

import "encoding/json"

type jsonResponseSuccess struct {
	OperationalLog []string    `json:"operationalLog,omitempty"`
	StackTrace     []string    `json:"stackTrace,omitempty"`
	Data           interface{} `json:"data"`
}

// String just gets the json representation or a error string
// The ideal thing to do would be to not use this method to encode the response. Instead, we should always use the
// render methods to send the success json response
func (e jsonResponseSuccess) String() string {
	// String representation of the Error Response. Can only be JSON
	successResponseJson, err := json.Marshal(e)
	if err != nil {
		return "E#1MZHO4 - JSON Encode failed"
	}

	return string(successResponseJson)
}

// ==============================================================
type jsonResponseFailure struct {
	Code           string   `json:"code"`
	Message        string   `json:"message"`
	DevMsg         string   `json:"devMsg,omitempty"`
	StackTrace     []string `json:"stackTrace,omitempty"`
	OperationalLog []string `json:"operationalLog,omitempty"`
}

// String just gets the json representation or a error string
// The ideal thing to do would be to not use this method to encode the response. Instead, we should always use the
// render methods to send the failure json response
func (e jsonResponseFailure) String() string {
	// String representation of the Error Response. Can only be JSON
	successResponseJson, err := json.Marshal(e)
	if err != nil {
		return "E#1N19DN - JSON Encode failed"
	}

	return string(successResponseJson)
}

// ==============================================================

// SingleMessageResponse is for sending a single message response to the client.
// Useful when just a single `200 OK` or `201 CREATED` would be ok but you still want to send a message to the client
// about what happened. e.g. "The blog post was created" or "The upload was successful" etc.
type SingleMessageResponse struct {
	Message string `json:"message"`
}

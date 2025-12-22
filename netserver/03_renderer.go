// Package netserver allows you to get a net/http server up and running with some helper methods available
// right from the start
package netserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/techrail/ground/constants/customCtxKey"
	"github.com/techrail/ground/constants/customHeaders"
	"github.com/techrail/ground/constants/httpheaders"
	"github.com/techrail/ground/logger"
	"github.com/techrail/ground/typs"
	"github.com/techrail/ground/typs/appError"
)

type Renderer struct{}

func (r *Renderer) JsonWithFailureUsingErrorType(w http.ResponseWriter, rq *http.Request, errTy appError.Typ) {
	errId := typs.GetRandomAlphaString(32)
	if errTy.IsBlankNetworkError() {
		logger.Println(fmt.Sprintf("E#2R0L3T: ErrID: %v, Error: %v @@@@@ DevMsg: %v", errId, errTy, errTy.DevMsg))
		r.JsonWithFailure(w, rq, 500, "2R0L3T", "Internal error. Error logged with ID "+errId, errTy.DevMsg)
	}
	r.JsonWithFailure(w, rq, errTy.HttpResponseCode, errTy.Code, errTy.Message, errTy.DevMsg)
}

func (r *Renderer) GetReqCtxValueAsString(rq *http.Request, key string) string {
	val, ok := rq.Context().Value(key).(string)
	if !ok {
		return ""
	}
	return val
}

func (r *Renderer) JsonWithFailure(w http.ResponseWriter, rq *http.Request, httpCode int, errorCode string, errorMessage string, devMessage string) {
	addFixedHeaders(w)
	w.Header().Set(httpheaders.ContentType, "application/json; charset=utf-8")
	w.Header().Set(customHeaders.RequestId, r.GetReqCtxValueAsString(rq, customCtxKey.RequestId))

	if httpCode > 199 && httpCode < 300 {
		errMsg := fmt.Sprintf("E#2R0MLR: Failure Renderer Called with non-failure HTTP code: %v", httpCode)
		logger.Println(errMsg)
	}

	// MARKER: Capturing stack trace
	// REWRITE THIS PART AND ENABLE THE STACK TRACE
	// ========
	// var stackTrace []byte
	// stackTrace = nil
	// captureStackTrace := r.GetReqCtxValueAsString(rq, (customCtxKey.StackTraceRequested)
	// if captureStackTrace != nil {
	// 	// Option was set in the context. Try to check its value
	// 	valBool, ok := captureStackTrace.(bool)
	// 	if ok && valBool {
	// 		// Stacktrace was requested
	// 		stackTrace = debug.Stack()
	// 	}
	// }
	// stackTraceStr := string(stackTrace)
	// stackTraceStr = strings.ReplaceAll(stackTraceStr, "\t", "    ")
	// stackTraceStrLines := strings.Split(stackTraceStr, "\n")
	// if len(stackTrace) == 0 {
	// 	// To ensure that if a blank stack trace is sent by the runtime, it is discarded
	// 	stackTraceStrLines = nil
	// }
	// // MARKER: Capturing operational log
	// var opLog []string
	// respondWithOplog := ctx.UserValue(customCtxKey.OpLogRequested)
	// if respondWithOplog != nil {
	// 	// Option was set in the context. Try to check its value
	// 	valBool, ok := respondWithOplog.(bool)
	//
	// 	if ok {
	// 		if valBool == true {
	// 			res := ctx.UserValue(customCtxKey.CtxOperationLogContent)
	// 			if res == nil {
	// 				errMsg := "E#1MZFU2 - Unexpected nil value found"
	// 				logger.Println(errMsg)
	// 				opLog = []string{
	// 					errMsg,
	// 				}
	// 			} else {
	// 				// Try to assert
	// 				if oprLog, typeAsserted := res.([]string); !typeAsserted {
	// 					errMsg := "E#1MZFUD - Incorrect data format"
	// 					logger.Println(errMsg)
	// 					opLog = []string{
	// 						errMsg,
	// 					}
	// 				} else {
	// 					opLog = oprLog
	// 				}
	// 			}
	// 		}
	// 	} else {
	// 		opLog = []string{
	// 			"E#1MZH45 - Value against user key was not in expected data type",
	// 		}
	// 	}
	// }
	//
	// // MARKER: Checking if DevMessage is allowed or not
	// devMsg := ""
	// if val, ok := ctx.UserValue(customCtxKey.DevMsgAllowedInFailure).(bool); ok != false {
	// 	// Key contains a valid boolean value. Check if the value is true
	// 	if val {
	// 		// devMessage can be sent
	// 		devMsg = devMessage
	// 	}
	// }
	//
	// logger.Println(
	// 	fmt.Sprintf("%v - %v [::DevMsg::]-> %v", errorCode, errorMessage, devMsg))
	//
	devMsg := ""
	stackTraceStrLines := []string{}
	opLog := []string{}

	resp := jsonResponseFailure{
		Code:           errorCode,
		Message:        errorMessage,
		DevMsg:         devMsg,             // NOTE: Fix this
		StackTrace:     stackTraceStrLines, // NOTE: Fix this
		OperationalLog: opLog,              // NOTE: Fix this

	}.String()
	//
	// ========
	w.WriteHeader(httpCode)
	_, _ = w.Write([]byte(resp))
}

// JsonStringWithSuccess is supposed to set the response code and string body value in the context response
// The parameter `jsonBody` is supposed to be a valid json.
func (r *Renderer) JsonStringWithSuccess(w http.ResponseWriter, rq *http.Request, httpCode int, jsonBody string) {
	addFixedHeaders(w)
	w.Header().Set(httpheaders.ContentType, "application/json; charset=utf-8")

	reqId := r.GetReqCtxValueAsString(rq, customCtxKey.RequestId)
	if reqId != "" {
		w.Header().Set(customHeaders.RequestId, reqId)
	}

	// Log for human error
	if httpCode < 200 || httpCode > 299 {
		logger.Debug(fmt.Sprintf("2R2DAC: Success Renderer Called with non-success HTTP code: %v", httpCode))
	}

	// The JSON body should be a valid JSON string
	_, err := json.Marshal(jsonBody)
	if err != nil {
		// It's not. We should send it as a single message response
		// but first log the error with context, so we know what kind of request caused the failure.
		// TODO: Make the following work. Enable contextual logging
		// logger.LogWithContext(ctx, "E#1MZEM0 - Non-JSON string sent to JsonStringWithSuccess")
		// Todo: Make the function below and enable it
		r.JsonStructWithSuccess(w, rq, httpCode, SingleMessageResponse{
			Message: jsonBody,
		})
	}

	// Create the JSON string which confirms with the rest of the structure
	responseBody := fmt.Sprintf(`{"data":%v}`, jsonBody)
	w.WriteHeader(httpCode)
	w.Write([]byte(responseBody))
}

// JsonBytesWithSuccess is supposed to set the response code and []byte body value in the context response
func (r *Renderer) JsonBytesWithSuccess(w http.ResponseWriter, rq *http.Request, httpCode int, jsonBody []byte) {
	addFixedHeaders(w)
	w.Header().Set(httpheaders.ContentType, "application/json; charset=utf-8")

	reqId := r.GetReqCtxValueAsString(rq, customCtxKey.RequestId)
	if reqId != "" {
		w.Header().Set(customHeaders.RequestId, reqId)
	}

	// Log for human error
	if httpCode < 200 || httpCode > 299 {
		errMsg := fmt.Sprintf("E#1MZEZ9 - Success Renderer Called with non-success HTTP code: %v", httpCode)
		logger.Println(errMsg)
	}

	w.WriteHeader(httpCode)
	w.Write([]byte(jsonBody))
}

func (r *Renderer) JsonStructWithSuccess(w http.ResponseWriter, rq *http.Request, httpCode int, structToMarshal interface{}) {
	addFixedHeaders(w)
	w.Header().Set(httpheaders.ContentType, "application/json; charset=utf-8")

	reqId := r.GetReqCtxValueAsString(rq, customCtxKey.RequestId)
	if reqId != "" {
		w.Header().Set(customHeaders.RequestId, reqId)
	}

	// Log for human error
	if httpCode < 200 || httpCode > 299 {
		errMsg := fmt.Sprintf("E#2R2MJ9: Success Renderer Called with non-success HTTP code: %v", httpCode)
		logger.Println(errMsg)
	}

	// MARKER: Capturing stack trace
	// var stackTrace []byte
	// stackTrace = nil
	// captureStackTrace := ctx.UserValue(customCtxKey.StackTraceRequested)
	// if captureStackTrace != nil {
	// 	// Option was set in the context. Try to check its value
	// 	valBool, ok := captureStackTrace.(bool)
	// 	if ok && valBool {
	// 		// Stacktrace was requested
	// 		stackTrace = debug.Stack()
	// 	}
	// }
	// stackTraceStr := string(stackTrace)
	// stackTraceStr = strings.ReplaceAll(stackTraceStr, "\t", "    ")
	// stackTraceStrLines := strings.Split(stackTraceStr, "\n")
	// if len(stackTrace) == 0 {
	// 	// To ensure that if a blank stack trace is sent by the runtime, it is discarded
	// 	stackTraceStrLines = nil
	// }
	//
	// MARKER: Capturing operational log
	// var opLog []string
	// respondWithOplog := ctx.UserValue(customCtxKey.OpLogRequested)
	// if respondWithOplog != nil {
	// 	// Option was set in the context. Try to check its value
	// 	valBool, ok := respondWithOplog.(bool)
	//
	// 	if ok {
	// 		if valBool == true {
	// 			res := ctx.UserValue(customCtxKey.CtxOperationLogContent)
	// 			if res == nil {
	// 				errMsg := "E#1MZFU2 - Unexpected nil value found"
	// 				logger.Println(errMsg)
	// 				opLog = []string{
	// 					errMsg,
	// 				}
	// 			} else {
	// 				// Try to assert
	// 				if oprLog, typeAsserted := res.([]string); !typeAsserted {
	// 					errMsg := "E#1MZFUD - Incorrect data format"
	// 					logger.Println(errMsg)
	// 					opLog = []string{
	// 						errMsg,
	// 					}
	// 				} else {
	// 					opLog = oprLog
	// 				}
	// 			}
	// 		}
	// 	} else {
	// 		opLog = []string{
	// 			"E#1MZH45 - Value against user key was not in expected data type",
	// 		}
	// 	}
	// }
	//

	opLog := []string{}
	stackTraceStrLines := []string{}
	successResponse := jsonResponseSuccess{
		OperationalLog: opLog,
		StackTrace:     stackTraceStrLines,
		Data:           structToMarshal,
	}

	// Attempt to marshal
	jsonBytes, err := json.Marshal(successResponse)
	if err != nil {
		// Something went wrong when trying to marshal
		errMsg := fmt.Sprintf("E#1MZIZS - JSON Marshalling failed: %v", err)
		logger.Println(errMsg)
		// Send error response
		// TODO: build this function
		r.JsonWithFailure(w, rq, http.StatusInternalServerError, "E#2R2MMB", "JSON Marshalling failed", errMsg)
		return
	}

	r.JsonBytesWithSuccess(w, rq, httpCode, jsonBytes)
}

// ==== The types that are reused in the responses ====

type jsonResponseSuccess struct {
	OperationalLog []string `json:"operationalLog,omitempty"`
	StackTrace     []string `json:"stackTrace,omitempty"`
	Data           any      `json:"data"`
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

func addFixedHeaders(w http.ResponseWriter) {
	w.Header().Add("Server", "Apache 2.4")
	w.Header().Add("X-Powered-By", "PHP/7.2.12")
}

// File ends here

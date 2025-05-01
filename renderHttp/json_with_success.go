package render

import (
	"encoding/json"
	"fmt"
	"github.com/techrail/ground/constants/customCtxKey"
	"github.com/techrail/ground/constants/customHeaders"
	"github.com/techrail/ground/logger"
	"github.com/valyala/fasthttp"
	"runtime/debug"
	"strings"
)

// JsonStringWithSuccess is supposed to set the response code and string body value in the context response
// The parameter `jsonBody` is supposed to be a valid json.
func JsonStringWithSuccess(ctx *fasthttp.RequestCtx, httpCode int, jsonBody string) {
	addFixedHeaders(ctx)
	ctx.Response.Header.Set(fasthttp.HeaderContentType, "application/json; charset=utf-8")
	if ctx.UserValue(customCtxKey.RequestId) != nil {
		ctx.Response.Header.Set(customHeaders.RequestId, fmt.Sprintf("%v", ctx.UserValue(customCtxKey.RequestId)))
	}

	// Log for human error
	if httpCode < 200 || httpCode > 299 {
		logger.Debug(fmt.Sprintf("1MUKL3 - Success Renderer Called with non-success HTTP code: %v", httpCode))
	}

	// The JSON body should be a valid JSON string
	_, err := json.Marshal(jsonBody)
	if err != nil {
		// It's not. We should send it as a single message response
		// but first log the error with context, so we know what kind of request caused the failure.
		// TODO: Make the following work. Enable contextual logging
		logger.LogWithContext(ctx, "E#1MZEM0 - Non-JSON string sent to JsonStringWithSuccess")
		// Todo: Make the function below and enable it
		JsonStructWithSuccess(ctx, httpCode, SingleMessageResponse{
			Message: jsonBody,
		})
	}

	// Create the JSON string which confirms with the rest of the structure
	responseBody := fmt.Sprintf(`{"data":%v}`, jsonBody)

	ctx.SetStatusCode(httpCode)
	ctx.Response.SetBodyString(responseBody)
}

// JsonBytesWithSuccess is supposed to set the response code and []byte body value in the context response
func JsonBytesWithSuccess(ctx *fasthttp.RequestCtx, httpCode int, jsonBody []byte) {
	addFixedHeaders(ctx)
	ctx.Response.Header.Set(fasthttp.HeaderXContentTypeOptions, "application/json; charset=utf-8")
	ctx.Response.Header.Set(customHeaders.RequestId, fmt.Sprintf("%v", ctx.UserValue(customCtxKey.RequestId)))

	// Log for human error
	if httpCode < 200 || httpCode > 299 {
		errMsg := fmt.Sprintf("E#1MZEZ9 - Success Renderer Called with non-success HTTP code: %v", httpCode)
		logger.Println(errMsg)
	}

	ctx.SetStatusCode(httpCode)
	ctx.Response.SetBody(jsonBody)
}

func JsonStructWithSuccess(ctx *fasthttp.RequestCtx, httpCode int, structToMarshal interface{}) {
	addFixedHeaders(ctx)
	ctx.Response.Header.Set(fasthttp.HeaderContentType, "application/json; charset=utf-8")
	ctx.Response.Header.Set(customHeaders.RequestId, fmt.Sprintf("%v", ctx.UserValue(customCtxKey.RequestId)))

	// Log for human error
	if httpCode < 200 || httpCode > 299 {
		errMsg := fmt.Sprintf("E#1MZFH5 - Success Renderer Called with non-success HTTP code: %v", httpCode)
		logger.Println(errMsg)
	}

	// MARKER: Capturing stack trace
	var stackTrace []byte
	stackTrace = nil
	captureStackTrace := ctx.UserValue(customCtxKey.StackTraceRequested)
	if captureStackTrace != nil {
		// Option was set in the context. Try to check its value
		valBool, ok := captureStackTrace.(bool)
		if ok && valBool {
			// Stacktrace was requested
			stackTrace = debug.Stack()
		}
	}
	stackTraceStr := string(stackTrace)
	stackTraceStr = strings.ReplaceAll(stackTraceStr, "\t", "    ")
	stackTraceStrLines := strings.Split(stackTraceStr, "\n")
	if len(stackTrace) == 0 {
		// To ensure that if a blank stack trace is sent by the runtime, it is discarded
		stackTraceStrLines = nil
	}

	// MARKER: Capturing operational log
	var opLog []string
	respondWithOplog := ctx.UserValue(customCtxKey.OpLogRequested)
	if respondWithOplog != nil {
		// Option was set in the context. Try to check its value
		valBool, ok := respondWithOplog.(bool)

		if ok {
			if valBool == true {
				res := ctx.UserValue(customCtxKey.CtxOperationLogContent)
				if res == nil {
					errMsg := "E#1MZFU2 - Unexpected nil value found"
					logger.Println(errMsg)
					opLog = []string{
						errMsg,
					}
				} else {
					// Try to assert
					if oprLog, typeAsserted := res.([]string); !typeAsserted {
						errMsg := "E#1MZFUD - Incorrect data format"
						logger.Println(errMsg)
						opLog = []string{
							errMsg,
						}
					} else {
						opLog = oprLog
					}
				}
			}
		} else {
			opLog = []string{
				"E#1MZH45 - Value against user key was not in expected data type",
			}
		}
	}

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
		JsonWithFailure(ctx, fasthttp.StatusInternalServerError, "E#1MZJ36", "JSON Marshalling failed", errMsg)
		return
	}

	JsonBytesWithSuccess(ctx, httpCode, jsonBytes)
}

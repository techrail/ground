package render

import (
	"fmt"
	"github.com/techrail/ground/constants/customCtxKey"
	"github.com/techrail/ground/constants/customHeaders"
	"github.com/techrail/ground/logger"
	types "github.com/techrail/ground/typs"
	"github.com/techrail/ground/typs/appError"
	"github.com/valyala/fasthttp"
	"runtime/debug"
	"strings"
)

func JsonWithFailureUsingErrorType(ctx *fasthttp.RequestCtx, errTy appError.Typ) {
	errId := types.GetRandomAlphaString(32)
	if errTy.IsBlankNetworkError() {
		logger.Println(fmt.Sprintf("E#1MZJCN - ErrID: %v, Error: %v @@@@@ DevMsg: %v", errId, errTy, errTy.DevMsg))
		JsonWithFailure(ctx, 500, "1MZJDY", "Internal error. Error logged with ID "+errId, errTy.DevMsg)
	}
	JsonWithFailure(ctx, errTy.HttpResponseCode, errTy.Code, errTy.Message, errTy.DevMsg)
}

// JsonWithFailure is supposed to set a failure response code and other details
func JsonWithFailure(ctx *fasthttp.RequestCtx, httpCode int, errorCode string, errorMessage string, devMessage string) {
	addFixedHeaders(ctx)
	ctx.Response.Header.Set(fasthttp.HeaderContentType, "application/json; charset=utf-8")
	ctx.Response.Header.Set(customHeaders.RequestId, fmt.Sprintf("%v", ctx.UserValue(customHeaders.RequestId)))

	if httpCode > 199 && httpCode < 300 {
		errMsg := fmt.Sprintf("E#1N17JF - Failure Renderer Called with non-failure HTTP code: %v", httpCode)
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

	// MARKER: Checking if DevMessage is allowed or not
	devMsg := ""
	if val, ok := ctx.UserValue(customCtxKey.DevMsgAllowedInFailure).(bool); ok != false {
		// Key contains a valid boolean value. Check if the value is true
		if val {
			// devMessage can be sent
			devMsg = devMessage
		}
	}

	logger.Println(
		fmt.Sprintf("%v - %v [::DevMsg::]-> %v", errorCode, errorMessage, devMsg))

	resp := jsonResponseFailure{
		Code:           errorCode,
		Message:        errorMessage,
		DevMsg:         devMsg,
		StackTrace:     stackTraceStrLines,
		OperationalLog: opLog,
	}.String()

	ctx.SetStatusCode(httpCode)
	ctx.Response.SetBodyString(resp)
}

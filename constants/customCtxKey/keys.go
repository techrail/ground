package customCtxKey

// This file contains the list of keys in the request context against which certain values will be stored in the context

const RequestId = "ctx_RequestId"
const RandomValue = "ctx_RandomValue"

// OpLogRequested tells if the operational log was requested or not
const OpLogRequested = "ctx_OpLogRequested"

// CtxOperationLogContent would contain the log lines that are being emitted while the request is being served
const CtxOperationLogContent = "ctx_OperationLogContent"

// StackTraceRequested indicates if the stack trace was requested or not
const StackTraceRequested = "ctx_StackTraceRequested"

// DevMsgAllowedInFailure tells whether the dev msg was expected by client in case of an error
const DevMsgAllowedInFailure = "ctx_DevMsgAllowedInFailure"

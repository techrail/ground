package appError

import (
	"fmt"

	"github.com/techrail/ground/constants"
	"github.com/techrail/ground/constants/errCode"
)

type Typ struct {
	Level            Level
	Code             string // LMID
	Message          string // The actual error message
	HttpResponseCode int    // In case we are trying to use this type for returning a network error
	DevMsg           string // Message meant for developers only (usually makes sense against a network error)
	WrappedError     *Typ   // Any wrapped errors that we want to embed in this error
	ExtraData        string // When we need to pass more values (errors are values (as said by Rob Pike))
}

var BlankError Typ

func init() {
	BlankError = Typ{
		Level:   Unknown,
		Code:    errCode.BlankErrorCode,
		Message: constants.EmptyString,
	}
}

func (e Typ) Wrap(wrappedErr Typ) Typ {
	e.WrappedError = &wrappedErr
	return e
}

func (e Typ) Error() string {
	return e.String()
}

func (e Typ) String() string {
	retVal := fmt.Sprintf("%s#%s: %s", e.Level.ShortStr(), e.Code, e.Message)
	if e.WrappedError != nil {
		retVal = retVal + "\n  [Wraps error ==>]\n" + e.WrappedError.String()
	}
	return retVal
}

func (e Typ) IsBlank() bool {
	if e.Level == Unknown &&
		e.Code == errCode.BlankErrorCode &&
		e.Message == constants.EmptyString &&
		e.DevMsg == constants.EmptyString {
		return true
	}
	return false
}

func (e Typ) IsNotBlank() bool {
	return !e.IsBlank()
}

func (e Typ) WrapsErrorCode(errCode string) bool {
	if e.WrappedError != nil {
		return e.WrappedError.WrapsErrorCode(errCode)
	}

	if e.Code == errCode {
		return true
	}
	return false
}

func (e Typ) IsBlankNetworkError() bool {
	if e.IsBlank() && e.HttpResponseCode == constants.EmptyInt {
		return true
	}
	return false
}

func (e Typ) IsNotBlankNetworkError() bool {
	return !e.IsBlankNetworkError()
}

func (e Typ) WrapsErrorLevel(errLvl Level, checkCurrErr bool) bool {
	if checkCurrErr && e.Level == errLvl {
		return true
	}

	err := &e
	for {
		err = e.WrappedError
		if err == nil {
			// We have come to the last error
			break
		}
		if err.Level == errLvl {
			return true
		}
	}
	return false
}

func NewError(errLevel Level, code string, msg string, wrappedError ...Typ) Typ {
	var wErr *Typ

	if len(wrappedError) > 0 && wrappedError[0].IsNotBlank() {
		wErr = &wrappedError[0]
	}

	return Typ{
		Level:        errLevel,
		Code:         code,
		Message:      msg,
		WrappedError: wErr,
	}
}

func NewNetworkError(httpResponseCode int, errLvl Level, code string, msg string, devmsg string, wrappedError ...Typ) Typ {
	var wErr *Typ

	if httpResponseCode < 100 || httpResponseCode > 599 {
		return Typ{
			Level:            Alert,
			Code:             "1DJ1A2",
			Message:          "Internal Error",
			HttpResponseCode: 500,
			DevMsg:           fmt.Sprintf("Function was supplied an invalid value: %v", httpResponseCode),
			WrappedError: &Typ{
				Level:            errLvl,
				Code:             code,
				Message:          msg,
				WrappedError:     wErr,
				HttpResponseCode: httpResponseCode,
				DevMsg:           devmsg,
			},
		}
	}

	if len(wrappedError) > 0 {
		wErr = &wrappedError[0]
	}

	return Typ{
		Level:            errLvl,
		Code:             code,
		Message:          msg,
		WrappedError:     wErr,
		HttpResponseCode: httpResponseCode,
		DevMsg:           devmsg,
	}
}

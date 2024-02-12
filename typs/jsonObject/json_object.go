package jsonObject

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/techrail/ground/constants/errCode"
	"github.com/techrail/ground/typs"
	"github.com/techrail/ground/typs/appError"
)

// IMPORTANT: Please don't  try to replace the usage of this type in the DB model fields with something you think is
// 	better (like string, or a type from another package) without actually testing the full set of effects.
// 	IF YOU HAVE AN IDEA TO IMPROVE THE IMPLEMENTATION, PLEASE READ THE COMMENT BELOW FIRST.

// NOTE: You might wonder why this type was created?
// 	The reason was - the "more_data" fields in our database is a JSON/JSONB type. When we get it from the DB,
// 	the driver returns it as a string and if we try to send that value in a response, it goes out as a string
// 	Now, we know what happens when a JSON is sent to clients as a string - all sorts of escape characters
// 	are added to it (the \" and \\\n characters show up). This causes major headache to the clients and is a pain
// 	for the eyes to look at. REST API clients (like Postman) can't parse it easily and tests can fail and so on.
// 	To avoid all those problems, and to send JSON values as they should be sent, this type was created.

// IMPORTANT: THIS TYPE CANNOT HANDLE A TOP LEVEL JSON ARRAY (any valid JSON document starting with '[')
//
//	It can handle arrays nested in a JSON object though.
const (
	typeAny          = "any"
	typeInt          = "int"
	typeInt64        = "int64"
	typeFloat64      = "float64"
	typeString       = "string"
	typeBool         = "bool"
	typeObject       = "object"
	typeNil          = "nil"
	typeArrayAny     = "array/any"
	typeArrayInt     = "array/int"
	typeArrayFloat64 = "array/float64"
	typeArrayString  = "array/string"
	typeArrayBool    = "array/bool"
	typeArrayObject  = "array/object"
	typeUnknown      = "unknown"
)

const topLevelArrayKey = "topLevelArrayKeyb8d8c89aea51f88b1af1144e3e0b8b74ac2a2c257d08cb80ebc99a7262e5dd8c"

type StringAnyMap map[string]any
type Typ struct {
	Valid            bool
	hasTopLevelArray bool
	StringAnyMap
	// IMPORTANT: If any map iteration and read panic happens on this type, enable the lock below
	//  and enable the Lock and Unlock calls for this type
	// sync.Mutex // To ensure two routines do not access the same object and work on it at the same time
}

var NullJsonObject Typ

// EmptyNotNullJsonObject returns a new blank Typ
// NOTE: We cannot use a var for the EmptyNotNullJsonObject value because when we copy a lot of values around and
//
//	assign the var to multiple values throughout the program in multiple goroutines, we might get the panic message
//	"concurrent map read and map write" indicating that the value is being written and read simultaneously because
//	the same variable is being used at multiple places
func EmptyNotNullJsonObject() Typ {
	return Typ{
		StringAnyMap: StringAnyMap{},
		Valid:        true,
	}
}

func init() {
	NullJsonObject = Typ{
		StringAnyMap: nil,
		Valid:        false,
	}
}

func NewJsonObject(key string, value interface{}) Typ {
	j := Typ{
		StringAnyMap: map[string]interface{}{
			key: value,
		},
		Valid: false,
	}
	return j
}

func (j *Typ) IsEmpty() bool {
	if j.Valid == false {
		// An invalid json is effectively an empty one
		return true
	}

	if len(j.StringAnyMap) == 0 || j.Valid == false {
		return true
	}
	return false
}

// ToJsonObject will convert interface{} object type to Typ using json.Marshal and json.Unmarshal
func ToJsonObject(v interface{}) (Typ, error) {
	jsonObj := Typ{
		Valid:        true,
		StringAnyMap: StringAnyMap{},
	}
	var err error

	// If the value is either a byte slice or a string, check if it is already a JSON string or not
	switch v.(type) {
	case string:
		err = jsonObj.Scan(v.(string))
		if err != nil {
			return NullJsonObject, fmt.Errorf("E#1N7FI3 - %v", err)
		}
		return jsonObj, nil
	case []byte:
		err = jsonObj.Scan(v.([]byte))
		if err != nil {
			return NullJsonObject, fmt.Errorf("E#1N7FI6 - %v", err)
		}
		return jsonObj, nil
	}

	jsonValue, err := json.Marshal(v)
	if err != nil {
		return NullJsonObject, fmt.Errorf("E#1N7FI8 - %v", err)
	}

	err = jsonObj.Scan(jsonValue)
	if err != nil {
		return NullJsonObject, fmt.Errorf("E#1N7FIB - %v", err)
	}

	return jsonObj, nil
}

func (j *Typ) IsNotEmpty() bool {
	return !j.IsEmpty()
}

func (j *Typ) SetNewTopLevelElement(key string, value interface{}) (replacedExistingKey bool) {
	if j.Valid == false {
		// We are making this object into a valid one
		j.Valid = true
		j.StringAnyMap[key] = value
		return
	}

	replacedExistingKey = false

	if _, ok := j.StringAnyMap[key]; ok {
		// Element already exists
		replacedExistingKey = true
	}

	j.StringAnyMap[key] = value
	return
}

// GetTopLevelElement will return Top-Level element identified by key. If the key does not exist, nil is returned
func (j *Typ) GetTopLevelElement(key string) any {
	if val, ok := j.StringAnyMap[key]; ok {
		return val
	}
	return nil
}

func (j *Typ) GetValueFromJsonObjectByJPath(path string) (string, any, appError.Typ) {
	vTyp := typeUnknown
	firstLetter := ""
	lastLetter := ""

	var internalObj map[string]any
	var internalArrayAny []any
	var internalArrayInt []int
	var internalArrayFloat64 []float64
	var internalArrayString []string
	var internalArrayBool []bool
	var internalArrayObject []map[string]any
	var internalInt int
	var internalAny any
	var internalFloat64 float64
	var internalString string
	var internalBool bool
	var internalNil bool
	var internalUnknown bool

	// Setting initial value, always an object
	internalObj = j.StringAnyMap

	// Useful functions
	valType := func(v any) string {
		switch v.(type) {
		case nil:
			return typeNil
		case bool:
			return typeBool
		case int:
			return typeInt
		case int64:
			return typeInt64
		case float64:
			return typeFloat64
		case string:
			return typeString
		case []int:
			return typeArrayInt
		case []float64:
			return typeArrayFloat64
		case []string:
			return typeArrayString
		case []bool:
			return typeArrayBool
		case []map[string]any:
			return typeArrayObject
		case []any:
			return typeArrayAny
		case map[string]any:
			return typeObject
		case any:
			return typeAny
		default:
			return typeUnknown
		}
	}
	setVal := func(val any) {
		vTyp = valType(val)
		switch vTyp {
		case typeObject:
			internalObj = val.(map[string]any)
		case typeArrayInt:
			internalArrayInt = val.([]int)
		case typeArrayFloat64:
			internalArrayFloat64 = val.([]float64)
		case typeArrayString:
			internalArrayString = val.([]string)
		case typeArrayBool:
			internalArrayBool = val.([]bool)
		case typeArrayObject:
			internalArrayObject = val.([]map[string]any)
		case typeArrayAny:
			internalArrayAny = val.([]any)
		case typeInt:
			internalInt = val.(int)
		case typeFloat64:
			internalFloat64 = val.(float64)
		case typeString:
			internalString = val.(string)
		case typeBool:
			internalBool = val.(bool)
		case typeAny:
			internalAny = val.(any)
		case typeNil:
			internalNil = true
		case typeUnknown:
			internalUnknown = true
		}
	}

	pathSplits := strings.Split(path, ".")

	if !j.Valid {
		return vTyp, nil, appError.NewError(appError.Error, errCode.JsonObjectInvalid, "195ARA -> Cannot process an invalid Typ")
	}

	if len(pathSplits) == 0 {
		return vTyp, nil, appError.NewError(
			appError.Error,
			errCode.JsonObjectJPathInvalid,
			"195ASQ -> Len 0 --- UNEXPECTED")
	}
	// Setting the type for first value
	setVal(map[string]any(j.StringAnyMap))

	for _, p := range pathSplits {
		if len(p) == 0 {
			return vTyp, nil, appError.NewError(
				appError.Error,
				errCode.JsonObjectJPathInvalid,
				fmt.Sprintf(
					"195AVJ -> path starts or ends with a dot (.) or there are double dots (..) in path: %v", path))
		}
		firstLetter = p[0:1]
		lastLetter = p[len(p)-1:]
		if firstLetter == "[" {
			if lastLetter != "]" {
				return vTyp, nil, appError.NewError(
					appError.Error,
					errCode.JsonObjectJPathInvalid,
					fmt.Sprintf("195AXJ -> Malformed array index notation at %v", p))
			}
			if len(p) < 3 {
				return vTyp, nil, appError.NewError(
					appError.Error,
					errCode.JsonObjectJPathInvalid,
					fmt.Sprintf("195AXV -> Blank numerical value not acceptable as array index: %v", p))
			}
			if !typs.IsPositiveNumber(p[1 : len(p)-1]) {
				return vTyp, nil, appError.NewError(
					appError.Error,
					errCode.JsonObjectJPathInvalid,
					fmt.Sprintf("195AY9 -> Non numeric index at %v", p))
			}
			// Was the previous entry an array?...
			if vTyp != typeArrayAny && vTyp != typeArrayInt && vTyp != typeArrayFloat64 &&
				vTyp != typeArrayString && vTyp != typeArrayBool && vTyp != typeArrayObject {
				// ... nope, don't think it was an array
				return vTyp, nil, appError.NewError(
					appError.Error,
					errCode.JsonObjectJPathInvalid,
					fmt.Sprintf("195B1Z -> Cannot treat entity of type %v as array", vTyp))
			}
			// ...looks like it was! Check if the length can honor numerical index
			expectedIndex, errConv := strconv.Atoi(p[1 : len(p)-1])
			if errConv != nil {
				return vTyp, nil, appError.NewError(
					appError.Error,
					errCode.JsonObjectElementNotFound,
					fmt.Sprintf("195B8W -> Cannot convert %v to integer", p[1:len(p)-1]))
			}

			// Get the value.
			switch vTyp {
			case typeArrayInt:
				if (len(internalArrayInt) - 1) < expectedIndex {
					return vTyp, nil, appError.NewError(
						appError.Error,
						errCode.JsonObjectElementNotFound,
						fmt.Sprintf("195BEU -> Cannot extract value at index %v from array of %v elements", expectedIndex, len(internalArrayInt)))
				}
				val := internalArrayInt[expectedIndex]
				setVal(val)
			case typeArrayFloat64:
				if (len(internalArrayFloat64) - 1) < expectedIndex {
					return vTyp, nil, appError.NewError(
						appError.Error,
						errCode.JsonObjectElementNotFound,
						fmt.Sprintf("195BGO -> Cannot extract value at index %v from array of %v elements", expectedIndex, len(internalArrayFloat64)))
				}
				val := internalArrayFloat64[expectedIndex]
				setVal(val)
			case typeArrayString:
				if (len(internalArrayString) - 1) < expectedIndex {
					return vTyp, nil, appError.NewError(
						appError.Error,
						errCode.JsonObjectElementNotFound,
						fmt.Sprintf("195BHI -> Cannot extract value at index %v from array of %v elements", expectedIndex, len(internalArrayString)))
				}
				val := internalArrayString[expectedIndex]
				setVal(val)
			case typeArrayBool:
				if (len(internalArrayBool) - 1) < expectedIndex {
					return vTyp, nil, appError.NewError(
						appError.Error,
						errCode.JsonObjectElementNotFound,
						fmt.Sprintf("195BI7 -> Cannot extract value at index %v from array of %v elements", expectedIndex, len(internalArrayBool)))
				}
				val := internalArrayBool[expectedIndex]
				setVal(val)
			case typeArrayObject:
				if (len(internalArrayObject) - 1) < expectedIndex {
					return vTyp, nil, appError.NewError(
						appError.Error,
						errCode.JsonObjectElementNotFound,
						fmt.Sprintf("195BIT -> Cannot extract value at index %v from array of %v elements", expectedIndex, len(internalArrayObject)))
				}
				val := internalArrayObject[expectedIndex]
				setVal(val)
			default:
				if (len(internalArrayAny) - 1) < expectedIndex {
					return vTyp, nil, appError.NewError(
						appError.Error,
						errCode.JsonObjectElementNotFound,
						fmt.Sprintf("195BJH -> Cannot extract value at index %v from array of %v elements", expectedIndex, len(internalArrayAny)))
				}
				val := internalArrayAny[expectedIndex]
				setVal(val)
			}
		} else {
			if !typs.IsAlphaNumericOrDotDashUnderscore(p) {
				return vTyp, nil, appError.NewError(
					appError.Error,
					errCode.JsonObjectJPathInvalid,
					fmt.Sprintf("195BL6 -> Unacceptable value in path: '%v'", p))
			}

			switch vTyp {
			case typeObject:
				// Get value from Object
				val, ok := internalObj[p]
				if !ok {
					return vTyp, nil, appError.NewError(
						appError.Error,
						errCode.JsonObjectElementNotFound,
						fmt.Sprintf("195BLY -> No such key in Typ: %v", p))
				}
				setVal(val)
			default:
				return vTyp, nil, appError.NewError(
					appError.Error,
					errCode.JsonObjectElementNotFound,
					fmt.Sprintf("195BNO -> Cannot get value %v from a non-object (%v) type", p, vTyp))
			}
		}
	}

	// At this point, we should have the final value
	switch vTyp {
	case typeObject:
		return vTyp, internalObj, appError.BlankError
	case typeArrayAny:
		return vTyp, internalArrayAny, appError.BlankError
	case typeArrayInt:
		return vTyp, internalArrayInt, appError.BlankError
	case typeArrayFloat64:
		return vTyp, internalArrayFloat64, appError.BlankError
	case typeArrayString:
		return vTyp, internalArrayString, appError.BlankError
	case typeArrayBool:
		return vTyp, internalArrayBool, appError.BlankError
	case typeArrayObject:
		return vTyp, internalArrayObject, appError.BlankError
	case typeString:
		return vTyp, internalString, appError.BlankError
	case typeFloat64:
		return vTyp, internalFloat64, appError.BlankError
	case typeInt:
		return vTyp, internalInt, appError.BlankError
	case typeBool:
		return vTyp, internalBool, appError.BlankError
	case typeNil:
		return vTyp, internalNil, appError.BlankError
	case typeAny:
		return vTyp, internalAny, appError.BlankError
	case typeUnknown:
		fallthrough
	default:
		return vTyp, internalUnknown, appError.BlankError
	}
}

func (j *Typ) SetValueByJPath(path string, valueToSet any) error {
	if !j.Valid {
		return fmt.Errorf("E#1RTT5F - Cannot set value in invalid jsonObject")
	}

	if j.hasTopLevelArray {
		return fmt.Errorf("E#1RTT7C - Cannot yet set value in a jsonObject with an array at top level")
	}

	jCopy := Typ{
		Valid:            j.Valid,
		hasTopLevelArray: j.hasTopLevelArray,
		StringAnyMap:     j.StringAnyMap,
	}

	newJCopy, err := SetValueAndOverrideInJsonObjectByJPath(jCopy, path, valueToSet, true)
	if err != nil {
		return fmt.Errorf("E#1RTTDB - Setting value failed. Error: %v", err)
	}

	j.StringAnyMap = newJCopy.StringAnyMap

	return nil
}

// SetValueInJsonObjectByJPath will set a value in the Typ given its JPath.
func SetValueInJsonObjectByJPath(obj Typ, path string, valueToSet any) (Typ, error) {
	return SetValueAndOverrideInJsonObjectByJPath(obj, path, valueToSet, false)
}

// SetValueAndOverrideInJsonObjectByJPath will set a value in the Typ given its JPath.
// If the JPath contains non-existing keys then original object will be overridden based override parameter.
func SetValueAndOverrideInJsonObjectByJPath(obj Typ, path string, valueToSet any, override bool) (Typ, error) {
	objToReturn := Typ{}
	vTyp := typeUnknown
	firstLetter := ""
	lastLetter := ""
	toBeCreated := false

	var internalObj map[string]any
	var internalArrayAny []any
	var internalArrayInt []int
	var internalArrayFloat64 []float64
	var internalArrayString []string
	var internalArrayBool []bool
	var internalArrayObject []map[string]any
	var internalAny any
	var internalInt int
	var internalFloat64 float64
	var internalString string
	var internalBool bool
	var internalNil bool
	var internalUnknown bool

	type jsonAction struct {
		DataType       string           // Type of value at this node e.g. object, array/int, string, int etc.
		ToBeCreated    bool             // Exists in the object = false / Needs to be created = true
		AtIndex        int              // When value is an array, this holds the index to update (-1 for creation)
		AtKey          string           // When value is non array, this is the key which is to be created/updated
		ArrayOfAny     []any            // Set when the value is an array of unknown types
		ArrayOfStrings []string         // Set when the value is an array of strings
		ArrayOfInts    []int            // Set when the value is an array of integers
		ArrayOfFloats  []float64        // Set when the value is an array of floating point values
		ArrayOfBools   []bool           // Set when the value is an array of boolean values
		ArrayOfObjects []map[string]any // Set when the value is an array of objects
		StringValue    string           // Set when the value is a string
		IntValue       int              // Set when the value is an integer
		FloatValue     float64          // Set when the value is a float
		BoolValue      bool             // Set when value is a boolean
		ObjectValue    map[string]any   // Set when the value is an object
		AnyValue       any              // Set when value is of interface{} (or `any`) type
		IsNil          bool             // If the value is nil, this is set to true
		IsUnknown      bool             // If the value is of unknown type, this is set to true
	}

	// Setting initial value, always an object
	internalObj = obj.StringAnyMap

	// Useful functions
	valType := func(v any) string {
		switch v.(type) {
		case nil:
			return typeNil
		case bool:
			return typeBool
		case int:
			return typeInt
		case int64:
			return typeInt64
		case float64:
			return typeFloat64
		case string:
			return typeString
		case map[string]any:
			return typeObject
		case []int:
			return typeArrayInt
		case []float64:
			return typeArrayFloat64
		case []string:
			return typeArrayString
		case []bool:
			return typeArrayBool
		case []map[string]any:
			return typeArrayObject
		case []any:
			return typeArrayAny
		default:
			return typeUnknown
		}
	}
	resetInternalValues := func() {
		internalArrayAny = []any{}
		internalArrayInt = []int{}
		internalArrayFloat64 = []float64{}
		internalArrayString = []string{}
		internalArrayBool = []bool{}
		internalArrayObject = []map[string]any{}
		internalInt = 0
		internalFloat64 = 0.0
		internalString = ""
		internalBool = false
		internalObj = map[string]any{}
		internalNil = false
		internalUnknown = false
	}
	setVal := func(val any) {
		resetInternalValues()
		vTyp = valType(val)
		switch vTyp {
		case typeArrayInt:
			internalArrayInt = val.([]int)
		case typeArrayFloat64:
			internalArrayFloat64 = val.([]float64)
		case typeArrayString:
			internalArrayString = val.([]string)
		case typeArrayBool:
			internalArrayBool = val.([]bool)
		case typeArrayObject:
			internalArrayObject = val.([]map[string]any)
		case typeArrayAny:
			internalArrayAny = val.([]any)
		case typeInt:
			internalInt = val.(int)
		case typeFloat64:
			internalFloat64 = val.(float64)
		case typeString:
			internalString = val.(string)
		case typeBool:
			internalBool = val.(bool)
		case typeObject:
			internalObj = val.(map[string]any)
		case typeNil:
			internalNil = true
		case typeAny:
			internalAny = val.(any)
		case typeUnknown:
			internalUnknown = true
		}
	}
	getAnyValueFromJsonAction := func(ja jsonAction) any {
		switch ja.DataType {
		case typeArrayAny:
			return reflect.ValueOf(ja.ArrayOfAny).Interface()
		case typeArrayInt:
			return reflect.ValueOf(ja.ArrayOfInts).Interface()
		case typeArrayFloat64:
			return reflect.ValueOf(ja.ArrayOfFloats).Interface()
		case typeArrayString:
			return reflect.ValueOf(ja.ArrayOfStrings).Interface()
		case typeArrayBool:
			return reflect.ValueOf(ja.ArrayOfBools).Interface()
		case typeArrayObject:
			return reflect.ValueOf(ja.ArrayOfObjects).Interface()
		case typeInt:
			return reflect.ValueOf(ja.IntValue).Interface()
		case typeFloat64:
			return reflect.ValueOf(ja.FloatValue).Interface()
		case typeString:
			return reflect.ValueOf(ja.StringValue).Interface()
		case typeBool:
			return reflect.ValueOf(ja.BoolValue).Interface()
		case typeObject:
			return reflect.ValueOf(ja.ObjectValue).Interface()
		case typeNil:
			return nil
		case typeUnknown:
			fallthrough
		default:
			return reflect.ValueOf(nil).Interface()
		}
	}

	actionPlan := make([]jsonAction, 0)

	pathSplits := strings.Split(path, ".")

	// For being used inside the loop when dealing with arrays
	expectedIndex := 0
	err := fmt.Errorf("E#1N7FIX - UnsetError")

	if !obj.Valid {
		return objToReturn, fmt.Errorf("E#1N7FJ0 - Cannot process an invalid Typ")
	}
	if len(pathSplits) == 0 {
		return objToReturn, fmt.Errorf("E#1N7FJ2 - Len 0 --- UNEXPECTED")
	}

	// Setting the type for first value
	setVal(map[string]any(obj.StringAnyMap))

	// Set the first entry in the actionPlan
	jA := jsonAction{
		DataType:       typeObject,
		ToBeCreated:    false,
		AtIndex:        0,
		AtKey:          ".",
		ArrayOfAny:     internalArrayAny,
		ArrayOfStrings: internalArrayString,
		ArrayOfInts:    internalArrayInt,
		ArrayOfFloats:  internalArrayFloat64,
		ArrayOfBools:   internalArrayBool,
		ArrayOfObjects: internalArrayObject,
		StringValue:    internalString,
		IntValue:       internalInt,
		FloatValue:     internalFloat64,
		BoolValue:      internalBool,
		ObjectValue:    obj.StringAnyMap,
		AnyValue:       internalAny,
		IsNil:          internalNil,
		IsUnknown:      internalUnknown,
	}
	actionPlan = append(actionPlan, jA)
	allocate := func(size int) {
		switch valueToSet.(type) {
		case int:
			setVal(make([]int, size))
		case float64:
			setVal(make([]float64, size))
		case string:
			setVal(make([]string, size))
		case bool:
			setVal(make([]bool, size))
		case any:
			setVal(make([]any, size))
		default:
			fmt.Println("E#1P84CQ - Found an invalid value to set.")
		}
	}
	createOverridingJsonActionPlan := func(index, arrExpectedIndex int) {
		N := len(pathSplits)
		expectedIndex := 0
		for i := index; i < N; i++ {
			atKey := pathSplits[i]
			atKeyIndex := 0
			if strings.Contains(pathSplits[i], "[") {
				atKey = ""
				atKeyIndex = expectedIndex
			}
			if i == N-1 {
				setVal(valueToSet)
			} else {
				nxtIndex := i + 1
				if nxtIndex == N-1 && strings.Contains(pathSplits[nxtIndex], "[") {
					expectedIndex, err = strconv.Atoi(pathSplits[nxtIndex][1 : len(pathSplits[nxtIndex])-1])
					if err != nil {
						fmt.Println("E#1P9WON - Found an Invalid array index in path", err)
					}
					allocate(expectedIndex)
				} else if nxtIndex != N-1 && strings.Contains(pathSplits[nxtIndex], "[") {
					expectedIndex, err = strconv.Atoi(pathSplits[nxtIndex][1 : len(pathSplits[nxtIndex])-1])
					if err != nil {
						fmt.Println("E#1P9WOT - Found an Invalid array index in path", err)
					}
					setVal(make([]map[string]any, expectedIndex))
				} else {
					setVal(map[string]any{})
				}
			}
			jA := jsonAction{
				DataType:       vTyp,
				ToBeCreated:    true,
				AtIndex:        atKeyIndex,
				AtKey:          atKey,
				ArrayOfAny:     internalArrayAny,
				ArrayOfStrings: internalArrayString,
				ArrayOfInts:    internalArrayInt,
				ArrayOfFloats:  internalArrayFloat64,
				ArrayOfBools:   internalArrayBool,
				ArrayOfObjects: internalArrayObject,
				StringValue:    internalString,
				IntValue:       internalInt,
				FloatValue:     internalFloat64,
				BoolValue:      internalBool,
				ObjectValue:    internalObj,
				AnyValue:       internalAny,
				IsNil:          internalNil,
				IsUnknown:      internalUnknown,
			}
			actionPlan = append(actionPlan, jA)
		}
	}

	// Loop to create the actionPlan
forLoop:
	for i, p := range pathSplits {
		if len(p) == 0 {
			return objToReturn, fmt.Errorf("E#1N7FJ7 - path starts or ends with a dot (.) or there are double dots (..) in path: %v", path)
		}
		firstLetter = p[0:1]
		lastLetter = p[len(p)-1:]
		if firstLetter == "[" {
			if lastLetter != "]" {
				return objToReturn, fmt.Errorf("E#1N7FJ9 - Malformed array index notation at %v", p)
			}

			if p == "[]" {
				// We want to append to the array
				toBeCreated = true
				// When trying to append to an array, the expression should be the last in the JPath
				if i < len(pathSplits)-1 {
					return objToReturn, fmt.Errorf("E#1N7FJC - Cannot append to an array whose properties are still required to be there")
				}
			} else {
				toBeCreated = false
				if !typs.IsPositiveNumber(p[1 : len(p)-1]) {
					return objToReturn, fmt.Errorf("E#1N7FJG - Non numeric index at %v", p)
				}
				expectedIndex, err = strconv.Atoi(p[1 : len(p)-1])
				if err != nil {
					return objToReturn, fmt.Errorf("E#1N7FJI - Cannot convert %v to integer", p[1:len(p)-1])
				}
			}

			// Was the previous element an array?
			if vTyp != typeArrayAny && vTyp != typeArrayInt && vTyp != typeArrayFloat64 &&
				vTyp != typeArrayString && vTyp != typeArrayBool && vTyp != typeArrayObject {
				// ... nope, don't think it was an array
				if override {
					actionPlan = actionPlan[:len(actionPlan)-1]
					createOverridingJsonActionPlan(i-1, expectedIndex)
					break forLoop
				}
				return objToReturn, fmt.Errorf("E#1N7FJK - Cannot treat entity of type %v as array", vTyp)
			}

			// Get the value.
			switch vTyp {
			case typeArrayInt:
				if (len(internalArrayInt) - 1) < expectedIndex {
					return objToReturn, fmt.Errorf("E#1N7FJQ - Cannot update value at index %v from array of %v elements", expectedIndex, len(internalArrayInt))
				}

				if toBeCreated {
					setVal(valueToSet)
				} else {
					// To be updated
					// If we are not in the last `i`, then we need to set the value of the object in the current index
					if i < len(pathSplits)-1 {
						val := internalArrayInt[expectedIndex]
						setVal(val)
					} else {
						setVal(valueToSet)
					}
				}
			case typeArrayFloat64:
				if (len(internalArrayFloat64) - 1) < expectedIndex {
					return objToReturn, fmt.Errorf("E#1N7RMR - Cannot update value at index %v from array of %v elements", expectedIndex, len(internalArrayFloat64))
				}

				if toBeCreated {
					setVal(valueToSet)
				} else {
					// To be updated
					// If we are not in the last `i`, then we need to set the value of the object in the current index
					if i < len(pathSplits)-1 {
						val := internalArrayFloat64[expectedIndex]
						setVal(val)
					} else {
						setVal(valueToSet)
					}
				}
			case typeArrayString:
				if (len(internalArrayString) - 1) < expectedIndex {
					return objToReturn, fmt.Errorf("E#1N7RMW - Cannot update value at index %v from array of %v elements", expectedIndex, len(internalArrayString))
				}

				if toBeCreated {
					setVal(valueToSet)
				} else {
					// To be updated
					// If we are not in the last `i`, then we need to set the value of the object in the current index
					if i < len(pathSplits)-1 {
						val := internalArrayString[expectedIndex]
						setVal(val)
					} else {
						setVal(valueToSet)
					}
				}
			case typeArrayBool:
				if (len(internalArrayBool) - 1) < expectedIndex {
					return objToReturn, fmt.Errorf("E#1N7RN3 - Cannot update value at index %v from array of %v elements", expectedIndex, len(internalArrayBool))
				}

				if toBeCreated {
					setVal(valueToSet)
				} else {
					// To be updated
					// If we are not in the last `i`, then we need to set the value of the object in the current index
					if i < len(pathSplits)-1 {
						val := internalArrayBool[expectedIndex]
						setVal(val)
					} else {
						setVal(valueToSet)
					}
				}
			case typeArrayObject:
				if (len(internalArrayObject) - 1) < expectedIndex {
					return objToReturn, fmt.Errorf("E#1N7RNA - Cannot update value at index %v from array of %v elements", expectedIndex, len(internalArrayObject))
				}

				if toBeCreated {
					setVal(valueToSet)
				} else {
					// To be updated
					// If we are not in the last `i`, then we need to set the value of the object in the current index
					if i < len(pathSplits)-1 {
						val := internalArrayObject[expectedIndex]
						setVal(val)
					} else {
						setVal(valueToSet)
					}
				}
			case typeArrayAny:
				if (len(internalArrayAny) - 1) < expectedIndex {
					return objToReturn, fmt.Errorf("E#1N7RNG - Cannot update value at index %v from array of %v elements", expectedIndex, len(internalArrayAny))
				}

				if toBeCreated {
					setVal(valueToSet)
				} else {
					// To be updated
					// If we are not in the last `i`, then we need to set the value of the object in the current index
					// Otherwise we need to take the value of the valueToSet and use that instead
					if i < len(pathSplits)-1 {
						val := internalArrayAny[expectedIndex]
						setVal(val)
					} else {
						setVal(valueToSet)
					}
				}
			default:
				// This piece of code is the same as above case (array/any)
				if (len(internalArrayAny) - 1) < expectedIndex {
					return objToReturn, fmt.Errorf("E#1N7RNL - Cannot update value at index %v from array of %v elements", expectedIndex, len(internalArrayAny))
				}

				if toBeCreated {
					setVal(valueToSet)
				} else {
					// To be updated
					// If we are not in the last `i`, then we need to set the value of the object in the current index
					// Otherwise we need to take the value of the valueToSet and use that instead
					if i < len(pathSplits)-1 {
						val := internalArrayAny[expectedIndex]
						setVal(val)
					} else {
						setVal(valueToSet)
					}
				}
			}
			// If we are here then we have the index to update at and the value which needs to be updated as well

			jA := jsonAction{
				DataType:       vTyp,
				ToBeCreated:    toBeCreated,
				AtIndex:        expectedIndex,
				AtKey:          "",
				ArrayOfAny:     internalArrayAny,
				ArrayOfStrings: internalArrayString,
				ArrayOfInts:    internalArrayInt,
				ArrayOfFloats:  internalArrayFloat64,
				ArrayOfBools:   internalArrayBool,
				ArrayOfObjects: internalArrayObject,
				StringValue:    internalString,
				IntValue:       internalInt,
				FloatValue:     internalFloat64,
				BoolValue:      internalBool,
				ObjectValue:    internalObj,
				AnyValue:       internalAny,
				IsNil:          internalNil,
				IsUnknown:      internalUnknown,
			}
			actionPlan = append(actionPlan, jA)
		} else {
			if !typs.IsAlphaNumeric(p) {
				return objToReturn, fmt.Errorf("E#1N7RO2 - Unacceptable non-alphanumeric value in path: '%v'", p)
			}

			switch vTyp {
			case typeObject:
				// Get value from Object
				val, ok := internalObj[p]
				if !ok {
					// No such key in the object to fetch from. We have to insert
					// Check that this is the last element in pathSplits or not
					if i < len(pathSplits)-1 {
						// We are not in the last element
						if override {
							createOverridingJsonActionPlan(i, 0)
							break forLoop
						}
						return objToReturn, fmt.Errorf("E#1N7RO8 - Cannot assign value beyond last known path element %v", strings.Join(pathSplits[0:i], "."))
					}
					// Insert
					toBeCreated = true
					setVal(valueToSet)
				} else {
					setVal(val)
				}

				// If this is the last entry and toBeCreated is false, then we need to update using input value
				if i == len(pathSplits)-1 && toBeCreated == false {
					setVal(valueToSet)
				}

				jA := jsonAction{
					DataType:       vTyp,
					ToBeCreated:    toBeCreated,
					AtIndex:        0,
					AtKey:          p,
					ArrayOfAny:     internalArrayAny,
					ArrayOfStrings: internalArrayString,
					ArrayOfInts:    internalArrayInt,
					ArrayOfFloats:  internalArrayFloat64,
					ArrayOfBools:   internalArrayBool,
					ArrayOfObjects: internalArrayObject,
					StringValue:    internalString,
					IntValue:       internalInt,
					FloatValue:     internalFloat64,
					BoolValue:      internalBool,
					ObjectValue:    internalObj,
					AnyValue:       internalAny,
					IsNil:          internalNil,
					IsUnknown:      internalUnknown,
				}
				actionPlan = append(actionPlan, jA)
			default:
				return objToReturn, fmt.Errorf("E#1N7ROE - Cannot assign value beyond last known path element %v", strings.Join(pathSplits[0:i], "."))
			}
		}
	}

	// If we are here then we should have the actionPath laid out. We need to iterate over it.
	resetInternalValues()
	vTyp = typeUnknown
	// finalJsonAction := jsonAction{}
	for i := len(actionPlan) - 1; i >= 0; i-- {
		ja := actionPlan[i]
		switch ja.DataType {
		case typeInt:
			// A scalar can't be a container of another scalar.
			// That means that the index of this item must be 0
			// NOTE: Same condition for other scalar types below
			if i != len(actionPlan)-1 {
				return objToReturn, fmt.Errorf("E#1N7ROM - int type cannot contain something else. Index: %v", i)
			}
			internalInt = ja.IntValue
			// If this value was to be created then the parent of this thing would have to create this entry
			// We don't have to do anything here
		case typeFloat64:
			if i != len(actionPlan)-1 {
				return objToReturn, fmt.Errorf("E#1N7ROR - float64 type cannot contain something else. Index: %v", i)
			}
			internalFloat64 = ja.FloatValue
		case typeString:
			if i != len(actionPlan)-1 {
				return objToReturn, fmt.Errorf("E#1N7ROW - string type cannot contain something else. Index: %v", i)
			}
			internalString = ja.StringValue
		case typeBool:
			if i != len(actionPlan)-1 {
				return objToReturn, fmt.Errorf("E#1N7RP0 - bool type cannot contain something else. Index: %v", i)
			}
			internalBool = ja.BoolValue
		case typeNil:
			if i != len(actionPlan)-1 {
				return objToReturn, fmt.Errorf("E#1N7RP5 - nil type cannot contain something else. Index: %v", i)
			}
			internalNil = ja.IsNil
		case typeObject:
			objectValue := ja.ObjectValue
			// If the previous entry existed and needed to be created, then we need to create it now.
			prevJa := jsonAction{}
			if i < len(actionPlan)-1 {
				// There is a previous entry
				// Get the value
				prevJa = actionPlan[i+1]
				// Now check the key name. It must be present
				if prevJa.AtKey == "" {
					return objToReturn, fmt.Errorf("E#1N7RPH - Could not get the AtKey value from a jsonAction where it should have been present. jsonAction: %v", prevJa)
				}

				// Check if the key was asked to be created while it already exists?
				_, ok := objectValue[prevJa.AtKey]

				if ok {
					// Key is present
					if prevJa.ToBeCreated && !override {
						// Key is present and is being asked to be created as well. That amounts to an error
						return objToReturn, fmt.Errorf("E#1N7RPM - Value is present and was asked to be created for key %v", prevJa.AtKey)
					} else {
						// We are supposed to update it
						dataTypeForSwitch := prevJa.DataType
						// Set value in current jsonAction
						switch dataTypeForSwitch {
						case typeArrayAny:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfAny
						case typeArrayInt:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfInts
						case typeArrayFloat64:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfFloats
						case typeArrayString:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfStrings
						case typeArrayBool:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfBools
						case typeArrayObject:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfObjects
						case typeInt:
							objectValue[prevJa.AtKey] = prevJa.IntValue
						case typeFloat64:
							objectValue[prevJa.AtKey] = prevJa.FloatValue
						case typeString:
							objectValue[prevJa.AtKey] = prevJa.StringValue
						case typeBool:
							objectValue[prevJa.AtKey] = prevJa.BoolValue
						case typeObject:
							objectValue[prevJa.AtKey] = prevJa.ObjectValue
						case typeAny:
							objectValue[prevJa.AtKey] = prevJa.AnyValue
						case typeNil:
							if prevJa.IsNil {
								objectValue[prevJa.AtKey] = nil
							} else {
								return objToReturn, fmt.Errorf("E#1N7RPZ - This was supposed to be nil but was something else")
							}
						case typeUnknown:
							objectValue[prevJa.AtKey] = nil
						}
					}
					// actionPlan[i].ObjectValue = objectValue
				} else {
					// Key is not present
					if !prevJa.ToBeCreated {
						// If it is not present, and we are NOT supposed to create it either, then it is an error
						return objToReturn, fmt.Errorf("E#1N7RQ7 - Value is not present and was asked to be updated for key %v", prevJa.AtKey)
					} else {
						// We are supposed to create it
						// In this case, strictness won't matter
						switch prevJa.DataType {
						case typeArrayAny:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfAny
						case typeArrayInt:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfInts
						case typeArrayFloat64:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfFloats
						case typeArrayString:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfStrings
						case typeArrayBool:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfBools
						case typeArrayObject:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfObjects
						case typeInt:
							objectValue[prevJa.AtKey] = prevJa.IntValue
						case typeFloat64:
							objectValue[prevJa.AtKey] = prevJa.FloatValue
						case typeString:
							objectValue[prevJa.AtKey] = prevJa.StringValue
						case typeBool:
							objectValue[prevJa.AtKey] = prevJa.BoolValue
						case typeObject:
							objectValue[prevJa.AtKey] = prevJa.ObjectValue
						case typeAny:
							objectValue[prevJa.AtKey] = prevJa.AnyValue
						case typeNil:
							if prevJa.IsNil {
								objectValue[prevJa.AtKey] = nil
							} else {
								return objToReturn, fmt.Errorf("E#1N7RQE - This was supposed to be nil but was something else")
							}
						case typeUnknown:
							objectValue[prevJa.AtKey] = nil
						}
						actionPlan[i].ObjectValue = objectValue
					}
				}
			} else {
				// There was no previous entry. That means that we are at the leaf and we need not do anything.
				// TODO: Delete this else condition from code later.
			}
		case typeArrayInt:
			// This one is an array of integers.
			// So if this is not the last element then...
			if i < len(actionPlan)-1 {

				// Not in strict mode, not the last item and we need to update/append an integer array with a value
				if actionPlan[i+1].DataType != typeInt {
					// The value to be updated is a non-integer data
					// That means that we have to convert this array into array of interfaces and then append
					tempArrayAny := make([]any, 0)
					for _, intVal := range ja.ArrayOfInts {
						tempArrayAny = append(tempArrayAny, reflect.ValueOf(intVal).Interface())
					}

					// Time to append or update.
					// Do we have to append or update?
					if actionPlan[i+1].ToBeCreated {
						// To be appended
						tempArrayAny = append(tempArrayAny, getAnyValueFromJsonAction(actionPlan[i+1]))
					} else {
						// To be updated
						tempArrayAny[actionPlan[i+1].AtIndex] = getAnyValueFromJsonAction(actionPlan[i+1])
					}
					// Set it into the array
					// Set the DataType of the current type to be an array of any
					actionPlan[i].DataType = typeArrayAny
					// And set the value
					actionPlan[i].ArrayOfAny = tempArrayAny
				} else {
					// Value is an integer
					if actionPlan[i+1].ToBeCreated {
						// To be appended
						ja.ArrayOfInts = append(ja.ArrayOfInts, actionPlan[i+1].IntValue)
					} else {
						// To be updated
						ja.ArrayOfInts[actionPlan[i+1].AtIndex] = actionPlan[i+1].IntValue
					}
					// Set it into the array
					actionPlan[i].ArrayOfInts = ja.ArrayOfInts
				}
			} else {
				// We are at the leaf node.
				// We have nothing to do. The next iteration will handle this.
				// TODO: Delete this else condition from code later if not required
			}
		case typeArrayFloat64:
			// This one is an array of float64s.
			// So if this is not the last element then...
			if i < len(actionPlan)-1 {
				// Not in strict mode, not the last item and we need to update/append a float array with a value
				if actionPlan[i+1].DataType != typeFloat64 {
					// The value to be updated is a non-float64 data
					// That means that we have to convert this array into array of interfaces and then append
					tempArrayAny := make([]any, 0)
					for _, floatVal := range ja.ArrayOfFloats {
						tempArrayAny = append(tempArrayAny, reflect.ValueOf(floatVal).Interface())
					}

					// Time to append or update.
					// Do we have to append or update?
					if actionPlan[i+1].ToBeCreated {
						// To be appended
						tempArrayAny = append(tempArrayAny, getAnyValueFromJsonAction(actionPlan[i+1]))
					} else {
						// To be updated
						tempArrayAny[actionPlan[i+1].AtIndex] = getAnyValueFromJsonAction(actionPlan[i+1])
					}
					// Set it into the array
					// Set the DataType of the current type to be an array of any
					actionPlan[i].DataType = typeArrayAny
					// And set the value
					actionPlan[i].ArrayOfAny = tempArrayAny
				} else {
					// Value is a float
					if actionPlan[i+1].ToBeCreated {
						// To be appended
						ja.ArrayOfFloats = append(ja.ArrayOfFloats, actionPlan[i+1].FloatValue)
					} else {
						// To be updated
						ja.ArrayOfFloats[actionPlan[i+1].AtIndex] = actionPlan[i+1].FloatValue
					}
					// Set it into the array
					actionPlan[i].ArrayOfFloats = ja.ArrayOfFloats
				}
			} else {
				// We are at the leaf node.
				// We have nothing to do. The next iteration will handle this.
				// TODO: Delete this else condition from code later if not required
			}
		case typeArrayString:
			// This one is an array of strings.
			// So if this is not the last element then...
			if i < len(actionPlan)-1 {
				// Not in strict mode, not the last item and we need to update/append a string array with a value
				if actionPlan[i+1].DataType != typeString {
					// The value to be updated is a non-string data
					// That means that we have to convert this array into array of interfaces and then append
					tempArrayAny := make([]any, 0)
					for _, stringVal := range ja.ArrayOfStrings {
						tempArrayAny = append(tempArrayAny, reflect.ValueOf(stringVal).Interface())
					}

					// Time to append or update.
					// Do we have to append or update?
					if actionPlan[i+1].ToBeCreated {
						// To be appended
						tempArrayAny = append(tempArrayAny, getAnyValueFromJsonAction(actionPlan[i+1]))
					} else {
						// To be updated
						tempArrayAny[actionPlan[i+1].AtIndex] = getAnyValueFromJsonAction(actionPlan[i+1])
					}
					// Set it into the array
					// Set the DataType of the current type to be an array of any
					actionPlan[i].DataType = typeArrayAny
					// And set the value
					actionPlan[i].ArrayOfAny = tempArrayAny
				} else {
					// Value is a string
					if actionPlan[i+1].ToBeCreated {
						// To be appended
						ja.ArrayOfStrings = append(ja.ArrayOfStrings, actionPlan[i+1].StringValue)
					} else {
						// To be updated
						ja.ArrayOfStrings[actionPlan[i+1].AtIndex] = actionPlan[i+1].StringValue
					}
					// Set it into the array
					actionPlan[i].ArrayOfStrings = ja.ArrayOfStrings
				}
			} else {
				// We are at the leaf node.
				// We have nothing to do. The next iteration will handle this.
				// TODO: Delete this else condition from code later if not required
			}
		case typeArrayBool:
			// This one is an array of bools.
			// So if this is not the last element then...
			if i < len(actionPlan)-1 {
				// Not in strict mode, not the last item and we need to update/append a bool array with a value
				if actionPlan[i+1].DataType != typeBool {
					// The value to be updated is a non-bool data
					// That means that we have to convert this array into array of interfaces and then append
					tempArrayAny := make([]any, 0)
					for _, boolVal := range ja.ArrayOfBools {
						tempArrayAny = append(tempArrayAny, reflect.ValueOf(boolVal).Interface())
					}

					// Time to append or update.
					// Do we have to append or update?
					if actionPlan[i+1].ToBeCreated {
						// To be appended
						tempArrayAny = append(tempArrayAny, getAnyValueFromJsonAction(actionPlan[i+1]))
					} else {
						// To be updated
						tempArrayAny[actionPlan[i+1].AtIndex] = getAnyValueFromJsonAction(actionPlan[i+1])
					}
					// Set it into the array
					// Set the DataType of the current type to be an array of any
					actionPlan[i].DataType = typeArrayAny
					// And set the value
					actionPlan[i].ArrayOfAny = tempArrayAny
				} else {
					// Value is a string
					if actionPlan[i+1].ToBeCreated {
						// To be appended
						ja.ArrayOfBools = append(ja.ArrayOfBools, actionPlan[i+1].BoolValue)
					} else {
						// To be updated
						ja.ArrayOfBools[actionPlan[i+1].AtIndex] = actionPlan[i+1].BoolValue
					}
					// Set it into the array
					actionPlan[i].ArrayOfBools = ja.ArrayOfBools
				}
			} else {
				// We are at the leaf node.
				// We have nothing to do. The next iteration will handle this.
				// TODO: Delete this else condition from code later if not required
			}
		case typeArrayObject:
			// This one is an array of objects.
			// So if this is not the last element then...
			if i < len(actionPlan)-1 {
				// Not in strict mode, not the last item and we need to update/append an object array with a value
				if actionPlan[i+1].DataType != typeObject {
					// The value to be updated is a non-object data
					// That means that we have to convert this array into array of interfaces and then append
					tempArrayAny := make([]any, 0)
					for _, objVal := range ja.ArrayOfObjects {
						tempArrayAny = append(tempArrayAny, reflect.ValueOf(objVal).Interface())
					}

					// Time to append or update.
					// Do we have to append or update?
					if actionPlan[i+1].ToBeCreated {
						// To be appended
						tempArrayAny = append(tempArrayAny, getAnyValueFromJsonAction(actionPlan[i+1]))
					} else {
						// To be updated
						tempArrayAny[actionPlan[i+1].AtIndex] = getAnyValueFromJsonAction(actionPlan[i+1])
					}
					// Set it into the array
					// Set the DataType of the current type to be an array of any
					actionPlan[i].DataType = typeArrayAny
					// And set the value
					actionPlan[i].ArrayOfAny = tempArrayAny
				} else {
					// Value is a string
					if actionPlan[i+1].ToBeCreated {
						// To be appended
						ja.ArrayOfObjects = append(ja.ArrayOfObjects, actionPlan[i+1].ObjectValue)
					} else {
						// To be updated
						ja.ArrayOfObjects[actionPlan[i+1].AtIndex] = actionPlan[i+1].ObjectValue
					}
					// Set it into the array
					actionPlan[i].ArrayOfObjects = ja.ArrayOfObjects
				}
			} else {
				// We are at the leaf node.
				// We have nothing to do. The next iteration will handle this.
				// TODO: Delete this else condition from code later if not required
			}
		case typeArrayAny:
			// This one is an array of anys.
			// So if this is not the last element then...
			if i < len(actionPlan)-1 {
				// Not in strict mode, not the last item and we need to update/append any value array with a value
				// Do we have to append or update?
				if actionPlan[i+1].ToBeCreated {
					// To be appended
					ja.ArrayOfAny = append(ja.ArrayOfAny, getAnyValueFromJsonAction(actionPlan[i+1]))
				} else {
					// To be updated
					ja.ArrayOfAny[actionPlan[i+1].AtIndex] = getAnyValueFromJsonAction(actionPlan[i+1])
				}
				// Set it into the array
				// Set the DataType of the current type to be an array of any
				actionPlan[i].DataType = typeArrayAny
				// And set the value
				actionPlan[i].ArrayOfAny = ja.ArrayOfAny
			} else {
				// We are at the leaf node.
				// We have nothing to do. The next iteration will handle this.
				// TODO: Delete this else condition from code later if not required
			}
		case typeAny:
			// NOTE: This should not have happened! We should not have a scalar typeAny
			return objToReturn, fmt.Errorf("E#1N7RQW - Don't know what to do")
		case typeUnknown:
			fallthrough
		default:
			if i != 0 {
				return objToReturn, fmt.Errorf("E#1N7RR1 - any type cannot contain something else. Index: %v", i)
			}
			internalInt = ja.IntValue
		}
	}

	objToReturn = obj
	objToReturn.StringAnyMap = actionPlan[0].ObjectValue

	return objToReturn, nil
}

// ValidateJPathSyntax validates JPath syntax (not against any specific Typ)
func ValidateJPathSyntax(path string) error {
	var firstLetter string
	var lastLetter string
	pathSplits := strings.Split(path, ".")
	if len(pathSplits) == 0 {
		return fmt.Errorf("E#1N7RS3 - Len 0 --- UNEXPECTED")
	}
	// fmt.Printf("pathSplits[0]:'%v'\n", pathSplits[0])
	// fmt.Printf("len:'%v'\n", len(pathSplits))
	for _, p := range pathSplits {
		if len(p) == 0 {
			return fmt.Errorf("E#1N7RSA - possible double dots (..) in path: %v", path)
		}
		firstLetter = p[0:1]
		lastLetter = p[len(p)-1:]
		if firstLetter == "[" {
			if lastLetter != "]" {
				return fmt.Errorf("E#1N7RSI - Malformed array index notation at %v", p)
			}
			if len(p) < 3 {
				return fmt.Errorf("E#1N7RSS - Blank numerical value not acceptable as array index: %v", p)
			}
			if !typs.IsPositiveNumber(p[1 : len(p)-1]) {
				return fmt.Errorf("E#1N7RSX - Non numeric index at %v", p)
			}
		} else {
			if !typs.IsAlphaNumeric(p) {
				return fmt.Errorf("E#1N7RT2 - Unacceptable non-alphanumeric value in path: '%v'", p)
			}
		}
	}
	return nil
}

// MARKER: Stringer interface implementation

func (j *Typ) String() string {
	if !j.Valid {
		return ""
	}

	bytes, err := json.Marshal(j.StringAnyMap)
	if err != nil {
		return ""
	}

	return string(bytes)
}

func (j *Typ) StringOrBlankObject() string {
	retVal := j.String()
	if retVal == "" {
		return "{}"
	}
	return retVal
}

func (j *Typ) StringOrNil() *string {
	s := j.String()
	if j.Valid {
		return &s
	}
	return nil
}

// PrettyString will give the formatted string for this Typ
func (j *Typ) PrettyString() string {
	if !j.Valid {
		return ""
	}

	bytes, err := json.MarshalIndent(j.StringAnyMap, "", "    ")
	if err != nil {
		return ""
	}

	return string(bytes)
}

func (j *Typ) HasTopLevelArray() bool {
	if j.Valid && len(j.StringAnyMap) == 1 && j.hasTopLevelArray {
		if _, ok := j.StringAnyMap[topLevelArrayKey]; ok {
			return true
		}
	}

	return false
}

func (j *Typ) AsByteSlice() []byte {
	return []byte(j.String())
}

// MARKER: DB Interface implementations

// Value implements the driver.Valuer interface. This method returns the JSON-encoded representation of the struct.
func (j *Typ) Value() (driver.Value, error) {
	if j.Valid == false {
		return nil, nil
	}

	if len(j.StringAnyMap) == 0 {
		// Valid empty JSON
		return []byte("{}"), nil
	}

	return j.MarshalJSON()
}

//func (Typ) ConvertValue(v any) (driver.Value, error) {
//	rv := reflect.ValueOf(v)
//	rv := reflect.TypeOf(v)
//	switch rv.Kind() {
//	case reflect.Pointer:
//		val:=rv.Elem()
//		if
//	case reflect.Struct:
//
//	default:
//	}
//}

// Scan implements the sql.Scanner interface. This method decodes a JSON-encoded value into the struct fields.
func (j *Typ) Scan(value interface{}) error {
	var arrAnys []any = make([]any, 0)
	switch value.(type) {
	case nil:
		j.Valid = false
		return nil
	case string:
		// Convert to byte slice and try
		err := json.Unmarshal([]byte(value.(string)), &j.StringAnyMap)
		if err != nil {
			err2 := json.Unmarshal(value.([]byte), &arrAnys)
			if err2 != nil {
				return errors.New(fmt.Sprintf("E#1N7RT9 - Unmarshalling failed: %v", err))
			}
			j.hasTopLevelArray = true
			j.StringAnyMap = StringAnyMap{
				topLevelArrayKey: arrAnys,
			}
		}
		j.Valid = true
		return nil
	case []byte:
		err := json.Unmarshal(value.([]byte), &j.StringAnyMap)
		if err != nil {
			err2 := json.Unmarshal(value.([]byte), &arrAnys)
			if err2 != nil {
				return errors.New(fmt.Sprintf("E#1N7RTE - Unmarshalling failed: %v", err))
			}
			j.hasTopLevelArray = true
			j.StringAnyMap = StringAnyMap{
				topLevelArrayKey: arrAnys,
			}
		}
		j.Valid = true
		return nil
	default:
		// Attempt to convert
		b, ok := value.([]byte)
		if !ok {
			return errors.New("E#1N7RTK - Type assertion to []byte failed")
		}

		// return json.Unmarshal(b, &j.StringAnyMap)
		err := json.Unmarshal(b, &j.StringAnyMap)
		if err != nil {
			return errors.New(fmt.Sprintf("E#1N7RTQ - Unmarshalling failed after assertion passed: %v", err))
		}
		j.Valid = true
		return nil
	}
}

// MARKER: Custom implementation of JSON Encoder for this type

// MarshalJSON implements json.Marshaler interface
// IMPORTANT: PLEASE DO NOT CONVERT THE RECEIVER TO POINTER TYPE (DESPITE WARNINGS)
func (j Typ) MarshalJSON() ([]byte, error) {
	if !j.Valid {
		return []byte("null"), nil
	}
	if j.HasTopLevelArray() {
		return json.Marshal(j.StringAnyMap[topLevelArrayKey])
	}
	return json.Marshal(j.StringAnyMap)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Typ) UnmarshalJSON(dataToUnmarshal []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(dataToUnmarshal, &v); err != nil {
		return err
	}
	switch v.(type) {
	case StringAnyMap:
		err = json.Unmarshal(dataToUnmarshal, &j.StringAnyMap)
	case map[string]interface{}:
		err = json.Unmarshal(dataToUnmarshal, &j.StringAnyMap)
	// TODO: Try to implement the case for []interface{} here
	//  (case []interface{}:)
	case nil:
		j.Valid = false
		j.StringAnyMap = nil
		return nil
	default:
		err = fmt.Errorf("E#1N7RTW - Cannot convert object of type %v to Typ", reflect.TypeOf(v).Name())
	}

	j.Valid = true
	if err != nil {
		j.Valid = false
	}

	return err
}

func (j *Typ) FindKeyByValue(value interface{}) string {
	for key, val := range j.StringAnyMap {
		// fmt.Println("key...", key, "val...", val)
		switch val := val.(type) {
		case map[string]interface{}:
			if nestedKey := j.FindKeyByValue(value); nestedKey != "" {
				return nestedKey
			}
		default:
			if val == value {
				return key
			}
		}
	}
	return ""
}

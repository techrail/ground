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
	TypeAny          = "any"
	TypeInt          = "int"
	TypeFloat64      = "float64"
	TypeString       = "string"
	TypeBool         = "bool"
	TypeObject       = "object"
	TypeNil          = "nil"
	TypeArrayAny     = "array/any"
	TypeArrayInt     = "array/int"
	TypeArrayFloat64 = "array/float64"
	TypeArrayString  = "array/string"
	TypeArrayBool    = "array/bool"
	TypeArrayObject  = "array/object"
	TypeUnknown      = "unknown"
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

// FillFromJsonObject will take a jsonObject Typ and fill the pointerToStruct with the value extracted
// from the jsonObject
func FillFromJsonObject(j Typ, pointerToStruct any) error {
	if reflect.TypeOf(pointerToStruct).Kind() != reflect.Ptr {
		return fmt.Errorf("E#1UNYJ3 - operation allowed only on a pointer to struct")
	}

	if reflect.ValueOf(pointerToStruct).IsNil() {
		return fmt.Errorf("E#1UO0A8 - operation not allowed on nil pointers")
	}

	if reflect.TypeOf(reflect.ValueOf(pointerToStruct).Elem()).Kind() != reflect.Struct {
		return fmt.Errorf("E#1UO9T9 - Value must be pointer to a struct")
	}

	jsonBytes, err := json.Marshal(j)
	if err != nil {
		return fmt.Errorf("E#1UO9VA - Marshalling error: %v", err)
	}

	err = json.Unmarshal(jsonBytes, pointerToStruct)
	if err != nil {
		return fmt.Errorf("E#1UO9XC - Could not unmarshal. Error: %v", err)
	}

	return nil
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

// GetValueByJPath returns 3 values, in that order: dataType, value and an appError
// If there is any error when trying to get the value, the appError value contains the error and in that case
// the other two values should be ignored. In other cases the dataType indicates data type detected
// and the value can then be safely casted using value.(correspondingGoDataType) expression
// where `correspondingGoDataType` is the data type corresponding to dataType.
func (j *Typ) GetValueByJPath(path string) (string, any, appError.Typ) {
	vTyp := TypeUnknown
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
			return TypeNil
		case bool:
			return TypeBool
		case int:
			return TypeInt
		case float64:
			return TypeFloat64
		case string:
			return TypeString
		case []int:
			return TypeArrayInt
		case []float64:
			return TypeArrayFloat64
		case []string:
			return TypeArrayString
		case []bool:
			return TypeArrayBool
		case []map[string]any:
			return TypeArrayObject
		case []any:
			return TypeArrayAny
		case map[string]any:
			return TypeObject
		case any:
			return TypeAny
		default:
			return TypeUnknown
		}
	}
	setVal := func(val any) {
		vTyp = valType(val)
		switch vTyp {
		case TypeObject:
			internalObj = val.(map[string]any)
		case TypeArrayInt:
			internalArrayInt = val.([]int)
		case TypeArrayFloat64:
			internalArrayFloat64 = val.([]float64)
		case TypeArrayString:
			internalArrayString = val.([]string)
		case TypeArrayBool:
			internalArrayBool = val.([]bool)
		case TypeArrayObject:
			internalArrayObject = val.([]map[string]any)
		case TypeArrayAny:
			internalArrayAny = val.([]any)
		case TypeInt:
			internalInt = val.(int)
		case TypeFloat64:
			internalFloat64 = val.(float64)
		case TypeString:
			internalString = val.(string)
		case TypeBool:
			internalBool = val.(bool)
		case TypeAny:
			internalAny = val.(any)
		case TypeNil:
			internalNil = true
		case TypeUnknown:
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
			if vTyp != TypeArrayAny && vTyp != TypeArrayInt && vTyp != TypeArrayFloat64 &&
				vTyp != TypeArrayString && vTyp != TypeArrayBool && vTyp != TypeArrayObject {
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
			case TypeArrayInt:
				if (len(internalArrayInt) - 1) < expectedIndex {
					return vTyp, nil, appError.NewError(
						appError.Error,
						errCode.JsonObjectElementNotFound,
						fmt.Sprintf("195BEU -> Cannot extract value at index %v from array of %v elements", expectedIndex, len(internalArrayInt)))
				}
				val := internalArrayInt[expectedIndex]
				setVal(val)
			case TypeArrayFloat64:
				if (len(internalArrayFloat64) - 1) < expectedIndex {
					return vTyp, nil, appError.NewError(
						appError.Error,
						errCode.JsonObjectElementNotFound,
						fmt.Sprintf("195BGO -> Cannot extract value at index %v from array of %v elements", expectedIndex, len(internalArrayFloat64)))
				}
				val := internalArrayFloat64[expectedIndex]
				setVal(val)
			case TypeArrayString:
				if (len(internalArrayString) - 1) < expectedIndex {
					return vTyp, nil, appError.NewError(
						appError.Error,
						errCode.JsonObjectElementNotFound,
						fmt.Sprintf("195BHI -> Cannot extract value at index %v from array of %v elements", expectedIndex, len(internalArrayString)))
				}
				val := internalArrayString[expectedIndex]
				setVal(val)
			case TypeArrayBool:
				if (len(internalArrayBool) - 1) < expectedIndex {
					return vTyp, nil, appError.NewError(
						appError.Error,
						errCode.JsonObjectElementNotFound,
						fmt.Sprintf("195BI7 -> Cannot extract value at index %v from array of %v elements", expectedIndex, len(internalArrayBool)))
				}
				val := internalArrayBool[expectedIndex]
				setVal(val)
			case TypeArrayObject:
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
			case TypeObject:
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
	case TypeObject:
		return vTyp, internalObj, appError.BlankError
	case TypeArrayAny:
		return vTyp, internalArrayAny, appError.BlankError
	case TypeArrayInt:
		return vTyp, internalArrayInt, appError.BlankError
	case TypeArrayFloat64:
		return vTyp, internalArrayFloat64, appError.BlankError
	case TypeArrayString:
		return vTyp, internalArrayString, appError.BlankError
	case TypeArrayBool:
		return vTyp, internalArrayBool, appError.BlankError
	case TypeArrayObject:
		return vTyp, internalArrayObject, appError.BlankError
	case TypeString:
		return vTyp, internalString, appError.BlankError
	case TypeFloat64:
		return vTyp, internalFloat64, appError.BlankError
	case TypeInt:
		return vTyp, internalInt, appError.BlankError
	case TypeBool:
		return vTyp, internalBool, appError.BlankError
	case TypeNil:
		return vTyp, internalNil, appError.BlankError
	case TypeAny:
		return vTyp, internalAny, appError.BlankError
	case TypeUnknown:
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
	vTyp := TypeUnknown
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
			return TypeNil
		case bool:
			return TypeBool
		case int:
			return TypeInt
		case float64:
			return TypeFloat64
		case string:
			return TypeString
		case map[string]any:
			return TypeObject
		case []int:
			return TypeArrayInt
		case []float64:
			return TypeArrayFloat64
		case []string:
			return TypeArrayString
		case []bool:
			return TypeArrayBool
		case []map[string]any:
			return TypeArrayObject
		case []any:
			return TypeArrayAny
		default:
			return TypeUnknown
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
		case TypeArrayInt:
			internalArrayInt = val.([]int)
		case TypeArrayFloat64:
			internalArrayFloat64 = val.([]float64)
		case TypeArrayString:
			internalArrayString = val.([]string)
		case TypeArrayBool:
			internalArrayBool = val.([]bool)
		case TypeArrayObject:
			internalArrayObject = val.([]map[string]any)
		case TypeArrayAny:
			internalArrayAny = val.([]any)
		case TypeInt:
			internalInt = val.(int)
		case TypeFloat64:
			internalFloat64 = val.(float64)
		case TypeString:
			internalString = val.(string)
		case TypeBool:
			internalBool = val.(bool)
		case TypeObject:
			internalObj = val.(map[string]any)
		case TypeNil:
			internalNil = true
		case TypeAny:
			internalAny = val.(any)
		case TypeUnknown:
			internalUnknown = true
		}
	}
	getAnyValueFromJsonAction := func(ja jsonAction) any {
		switch ja.DataType {
		case TypeArrayAny:
			return reflect.ValueOf(ja.ArrayOfAny).Interface()
		case TypeArrayInt:
			return reflect.ValueOf(ja.ArrayOfInts).Interface()
		case TypeArrayFloat64:
			return reflect.ValueOf(ja.ArrayOfFloats).Interface()
		case TypeArrayString:
			return reflect.ValueOf(ja.ArrayOfStrings).Interface()
		case TypeArrayBool:
			return reflect.ValueOf(ja.ArrayOfBools).Interface()
		case TypeArrayObject:
			return reflect.ValueOf(ja.ArrayOfObjects).Interface()
		case TypeInt:
			return reflect.ValueOf(ja.IntValue).Interface()
		case TypeFloat64:
			return reflect.ValueOf(ja.FloatValue).Interface()
		case TypeString:
			return reflect.ValueOf(ja.StringValue).Interface()
		case TypeBool:
			return reflect.ValueOf(ja.BoolValue).Interface()
		case TypeObject:
			return reflect.ValueOf(ja.ObjectValue).Interface()
		case TypeNil:
			return nil
		case TypeUnknown:
			fallthrough
		default:
			return reflect.ValueOf(nil).Interface()
		}
	}

	val := reflect.ValueOf(valueToSet)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {
		jo, err := ToJsonObject(valueToSet)
		if err != nil {
			return objToReturn, fmt.Errorf("E#1SCFLZ - Cannot convert value to jsonObject. Error: %v", err)
		}

		x := map[string]any(jo.StringAnyMap)

		// Call with the string any map as the value to be set
		return SetValueAndOverrideInJsonObjectByJPath(obj, path, x, override)
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
		DataType:       TypeObject,
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
			if vTyp != TypeArrayAny && vTyp != TypeArrayInt && vTyp != TypeArrayFloat64 &&
				vTyp != TypeArrayString && vTyp != TypeArrayBool && vTyp != TypeArrayObject {
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
			case TypeArrayInt:
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
			case TypeArrayFloat64:
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
			case TypeArrayString:
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
			case TypeArrayBool:
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
			case TypeArrayObject:
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
			case TypeArrayAny:
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
			case TypeObject:
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
	vTyp = TypeUnknown
	// finalJsonAction := jsonAction{}
	for i := len(actionPlan) - 1; i >= 0; i-- {
		ja := actionPlan[i]
		switch ja.DataType {
		case TypeInt:
			// A scalar can't be a container of another scalar.
			// That means that the index of this item must be 0
			// NOTE: Same condition for other scalar types below
			if i != len(actionPlan)-1 {
				return objToReturn, fmt.Errorf("E#1N7ROM - int type cannot contain something else. Index: %v", i)
			}
			internalInt = ja.IntValue
			// If this value was to be created then the parent of this thing would have to create this entry
			// We don't have to do anything here
		case TypeFloat64:
			if i != len(actionPlan)-1 {
				return objToReturn, fmt.Errorf("E#1N7ROR - float64 type cannot contain something else. Index: %v", i)
			}
			internalFloat64 = ja.FloatValue
		case TypeString:
			if i != len(actionPlan)-1 {
				return objToReturn, fmt.Errorf("E#1N7ROW - string type cannot contain something else. Index: %v", i)
			}
			internalString = ja.StringValue
		case TypeBool:
			if i != len(actionPlan)-1 {
				return objToReturn, fmt.Errorf("E#1N7RP0 - bool type cannot contain something else. Index: %v", i)
			}
			internalBool = ja.BoolValue
		case TypeNil:
			if i != len(actionPlan)-1 {
				return objToReturn, fmt.Errorf("E#1N7RP5 - nil type cannot contain something else. Index: %v", i)
			}
			internalNil = ja.IsNil
		case TypeObject:
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
						case TypeArrayAny:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfAny
						case TypeArrayInt:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfInts
						case TypeArrayFloat64:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfFloats
						case TypeArrayString:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfStrings
						case TypeArrayBool:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfBools
						case TypeArrayObject:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfObjects
						case TypeInt:
							objectValue[prevJa.AtKey] = prevJa.IntValue
						case TypeFloat64:
							objectValue[prevJa.AtKey] = prevJa.FloatValue
						case TypeString:
							objectValue[prevJa.AtKey] = prevJa.StringValue
						case TypeBool:
							objectValue[prevJa.AtKey] = prevJa.BoolValue
						case TypeObject:
							objectValue[prevJa.AtKey] = prevJa.ObjectValue
						case TypeAny:
							objectValue[prevJa.AtKey] = prevJa.AnyValue
						case TypeNil:
							if prevJa.IsNil {
								objectValue[prevJa.AtKey] = nil
							} else {
								return objToReturn, fmt.Errorf("E#1N7RPZ - This was supposed to be nil but was something else")
							}
						case TypeUnknown:
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
						case TypeArrayAny:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfAny
						case TypeArrayInt:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfInts
						case TypeArrayFloat64:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfFloats
						case TypeArrayString:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfStrings
						case TypeArrayBool:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfBools
						case TypeArrayObject:
							objectValue[prevJa.AtKey] = prevJa.ArrayOfObjects
						case TypeInt:
							objectValue[prevJa.AtKey] = prevJa.IntValue
						case TypeFloat64:
							objectValue[prevJa.AtKey] = prevJa.FloatValue
						case TypeString:
							objectValue[prevJa.AtKey] = prevJa.StringValue
						case TypeBool:
							objectValue[prevJa.AtKey] = prevJa.BoolValue
						case TypeObject:
							objectValue[prevJa.AtKey] = prevJa.ObjectValue
						case TypeAny:
							objectValue[prevJa.AtKey] = prevJa.AnyValue
						case TypeNil:
							if prevJa.IsNil {
								objectValue[prevJa.AtKey] = nil
							} else {
								return objToReturn, fmt.Errorf("E#1N7RQE - This was supposed to be nil but was something else")
							}
						case TypeUnknown:
							objectValue[prevJa.AtKey] = nil
						}
						actionPlan[i].ObjectValue = objectValue
					}
				}
			} else {
				// There was no previous entry. That means that we are at the leaf and we need not do anything.
				// TODO: Delete this else condition from code later.
			}
		case TypeArrayInt:
			// This one is an array of integers.
			// So if this is not the last element then...
			if i < len(actionPlan)-1 {

				// Not in strict mode, not the last item and we need to update/append an integer array with a value
				if actionPlan[i+1].DataType != TypeInt {
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
					actionPlan[i].DataType = TypeArrayAny
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
		case TypeArrayFloat64:
			// This one is an array of float64s.
			// So if this is not the last element then...
			if i < len(actionPlan)-1 {
				// Not in strict mode, not the last item and we need to update/append a float array with a value
				if actionPlan[i+1].DataType != TypeFloat64 {
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
					actionPlan[i].DataType = TypeArrayAny
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
		case TypeArrayString:
			// This one is an array of strings.
			// So if this is not the last element then...
			if i < len(actionPlan)-1 {
				// Not in strict mode, not the last item and we need to update/append a string array with a value
				if actionPlan[i+1].DataType != TypeString {
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
					actionPlan[i].DataType = TypeArrayAny
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
		case TypeArrayBool:
			// This one is an array of bools.
			// So if this is not the last element then...
			if i < len(actionPlan)-1 {
				// Not in strict mode, not the last item and we need to update/append a bool array with a value
				if actionPlan[i+1].DataType != TypeBool {
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
					actionPlan[i].DataType = TypeArrayAny
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
		case TypeArrayObject:
			// This one is an array of objects.
			// So if this is not the last element then...
			if i < len(actionPlan)-1 {
				// Not in strict mode, not the last item and we need to update/append an object array with a value
				if actionPlan[i+1].DataType != TypeObject {
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
					actionPlan[i].DataType = TypeArrayAny
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
		case TypeArrayAny:
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
				actionPlan[i].DataType = TypeArrayAny
				// And set the value
				actionPlan[i].ArrayOfAny = ja.ArrayOfAny
			} else {
				// We are at the leaf node.
				// We have nothing to do. The next iteration will handle this.
				// TODO: Delete this else condition from code later if not required
			}
		case TypeAny:
			// NOTE: This should not have happened! We should not have a scalar TypeAny
			return objToReturn, fmt.Errorf("E#1N7RQW - Don't know what to do")
		case TypeUnknown:
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

// StringOrBlankObject will return "{}" if the Typ is not valid, otherwise it will return the JSON string representation
// NOTE: This method is used by the generator and is not supposed to be removed
func (j *Typ) StringOrBlankObject() string {
	retVal := j.String()
	if retVal == "" {
		return "{}"
	}
	return retVal
}

// StringOrNil will return nil if the Typ is not valid, otherwise it will return the JSON string representation
// NOTE: This method is used by the generator and is not supposed to be removed
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

package jsonObject

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// The reason we have this type here is that we have to deal with JSON data in the database with PostgreSQL
// JSON and JSONB types. While it is possible to use the `string` data type with some validations, it can cause problems
// when encoding and decoding at multiple places. In addition, we need JSON values to be saved to the DB as well.
// This type helps in such cases.

const topLevelArrayKey = "topLevelArrayKey5a4ff523b0e7e5f0813a49e82519b11c6715e08b301458e0f00ea7abe5f3b184d78"

type StringAnyMap map[string]any
type Typ struct {
	Valid            bool
	hasTopLevelArray bool
	StringAnyMap
}

var NullJsonObject Typ

// EmptyNotNullJsonObject returns a new blank Typ
// We can't use a value here because we need a new copy every time. Under high concurrency and some conditions,
// the values cause reference problems where updating one copy also updates another and we can also get concurrent
// map access panic case. Returning a new copy every single time never causes those conflicts.
func EmptyNotNullJsonObject() Typ {
	return Typ{
		StringAnyMap: StringAnyMap{},
		Valid:        true,
	}
}

func init() {
	NullJsonObject = Typ{
		StringAnyMap:     nil,
		Valid:            false,
		hasTopLevelArray: false,
	}
}

func NewJsonObject(key string, value any) Typ {
	j := Typ{
		StringAnyMap: map[string]any{
			key: value,
		},
		Valid:            false,
		hasTopLevelArray: false,
	}
	return j
}

func (j *Typ) IsEmpty() bool {
	if j.Valid == false || len(j.StringAnyMap) == 0 {
		return true
	}
	return false
}

func (j *Typ) IsNotEmpty() bool {
	return !j.IsEmpty()
}

// ToJsonObject will try to convert any object type to the JsonObject type using json.Marshal and json.Unmarshal
func ToJsonObject(v any) (Typ, error) {
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
			return NullJsonObject, fmt.Errorf("E#1MQFQW - %v", err)
		}
		return jsonObj, nil
	case []byte:
		err = jsonObj.Scan(v.([]byte))
		if err != nil {
			return NullJsonObject, fmt.Errorf("E#1MQFQY - %v", err)
		}
		return jsonObj, nil
	}

	jsonValue, err := json.Marshal(v)
	if err != nil {
		return NullJsonObject, fmt.Errorf("E#1MQFR1 - %v", err)
	}

	err = jsonObj.Scan(jsonValue)
	if err != nil {
		return NullJsonObject, fmt.Errorf("E#1MQFR6 - %v", err)
	}

	return jsonObj, nil
}

func (j *Typ) SetNewTopLevelElement(key string, value any) (replacedExistingKey bool) {
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

// String implements the fmt.Stringer interface
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

// PrettyString will give the formatted string for this Typ. Indentation is set to 4 spaces.
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

// HasTopLevelArray tells if the top level element is an array or not
func (j *Typ) HasTopLevelArray() bool {
	if j.Valid && len(j.StringAnyMap) == 1 && j.hasTopLevelArray {
		if _, ok := j.StringAnyMap[topLevelArrayKey]; ok {
			return true
		}
	}

	return false
}

// AsByteSlice returns the []byte representation of the JsonObject as string
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
func (j *Typ) Scan(value any) error {
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
				return errors.New(fmt.Sprintf("E#1MQFRF - Unmarshalling failed: %v", err))
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
				return errors.New(fmt.Sprintf("E#1MQFRJ - Unmarshalling failed: %v", err))
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
			return errors.New("E#1MQFRO - Type assertion to []byte failed")
		}

		// return json.Unmarshal(b, &j.StringAnyMap)
		err := json.Unmarshal(b, &j.StringAnyMap)
		if err != nil {
			return errors.New(fmt.Sprintf("E#1MQFRR - Unmarshalling failed after assertion passed: %v", err))
		}
		j.Valid = true
		return nil
	}
}

// MARKER: Custom implementation of JSON Encoder for this type

// MarshalJSON implements json.Marshaler interface
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
	var v any
	if err = json.Unmarshal(dataToUnmarshal, &v); err != nil {
		return err
	}
	switch v.(type) {
	case StringAnyMap:
		err = json.Unmarshal(dataToUnmarshal, &j.StringAnyMap)
	case map[string]any:
		err = json.Unmarshal(dataToUnmarshal, &j.StringAnyMap)
	case []any:
		j.StringAnyMap = StringAnyMap{
			topLevelArrayKey: v.([]any),
		}
		j.hasTopLevelArray = true
	case nil:
		j.Valid = false
		j.StringAnyMap = nil
		return nil
	default:
		err = fmt.Errorf("E#1MQFRZ - Cannot convert object of type %v to Typ", reflect.TypeOf(v).Name())
	}

	j.Valid = true
	if err != nil {
		j.Valid = false
	}

	return err
}

package jsonObject

import (
	"fmt"
	"log"
	"runtime"
	"testing"

	"github.com/techrail/ground/typs/appError"
)

// arrObjs.[0].contents.[2].[3].[1].[2]
var jsonString = `
{
  "str": "stringValue",
  "bool1": true,
  "bool2": false,
  "int": 9741,
  "float": 12.34,
  "arrInts": [1, 2, 4, 8, 16],
  "arrFloats": [1.2, 2.4, 4.8, 8.16, 16.32],
  "arrStrings": ["vaibhav", "kaushal", "archana", "anahat", "bharti"],
  "arrBools": [true, false, true, true, false],
  "arrObjs": [
    {
      "objName": "firstObject",
      "kind": "test",
      "contents": [
        [101, 102, 103, 104, 105],
        ["this", "is", "a", "string", "array!!"],
        [
          90.8, 89.7, 78.6,
          [
            [11.1, 22.2, 33.3],
            [44.4, 55.5, 66.6],
            [77.7, 88.8, 99.9]
          ]
        ],
        "vaibhav",
        "kaushal"
      ]
    },
    {
      "objName": "secondObject",
      "kind": "test2",
      "intVal": 123,
      "fltVal": 345.6,
      "moreArrays": [
        {
          "k": "v"
        }, {
          "k2": 9035
        }
      ]
    }
  ],
  "obj": {
    "key": "value",
    "key2": "value2",
    "intVal": 90351,
    "arrInts": [11, 22, 33, 44, 55],
    "arrFloats": [11.1, 22.2, 33.3, 44.5, 55.5],
    "arrStrings": ["vaibhav", "anuj", "dilip", "ajay", "uday", "ravi"],
    "arr2": [
      {
        "k": "1v",
        "k2": "1v2",
        "list": [9, 8, 7, 6, 5, 44, 33]
      },
      {
        "ki": 2,
        "ki2": 20,
        "nested": [
          [1, 2, 3, 4, 5],
          [6, 7, 8, 9, 10],
          [11, 12, 13, 14, 15],
          [
            [101, 102, 103, 104],
            [201, 202, 203, 204, "test"],
            ["this", "is", "string", "array"],
            [12.3, 23.4, 34.5, 45.6, 56.7]
          ]
        ]
      }
    ],
    "nestedObj": {
      "1this": "2object",
      "3has" : "4nesting",
      "child0": {
        "lvl": 1,
        "child1": {
          "lvl": 2,
          "child2": {
            "lvl": 3,
            "finalChild": {
              "k5": "v5",
              "k6": 6,
              "k7": 10.01
            }
          }
        }
      }
    }
  }
}
`

// getTestsSuccessful indicates if the TestGetValueFromJsonObjectByJPath encountered no errors
//
//	The set value tests can't run if get value tests are failing
var getTestsSuccessful bool

func init() {
	getTestsSuccessful = false
}

func TestGetValueFromJsonObjectByJPath(t *testing.T) {
	var errTy appError.Typ
	jo, err := ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, err := jo.GetValueFromJsonObjectByJPath("nosuchkey")
	if err == nil {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#145EKW - ==> [[ERROR]] <== Should have errored!")
	} else {
		fmt.Printf("~~~~~~~\nL#145ELG - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("str")
	if errTy.IsNotBlank() || val != "stringValue" {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#145EZ9 - ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#145EWS - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("bool1")
	if errTy.IsNotBlank() || val != true {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#145F21 -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#145F27 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("bool2")
	if errTy.IsNotBlank() || val != false {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#145F5T - ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#145F5Y - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("int")
	if errTy.IsNotBlank() || (val != int(9741) && val != float64(9741)) {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#145F74 - ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#145F79 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("float")
	if errTy.IsNotBlank() || val != float64(12.34) {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#145FAP - ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#145FAT - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrInts")
	if errTy.IsNotBlank() || (typ != "array/int" && typ != "array/float64" && typ != "array/any") {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#145FDJ - ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#145FDN - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrFloats")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#145FE2 - ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#145FE6 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrStrings")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#145FK0 - ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#145FK4 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrBools")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146PMX - ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146PMZ - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146PNL - ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146PNN - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146PST - ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146PSW - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146PST -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146PSW - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	// /////////////////////////////////////////////////////
	// First Level Keys done. Nested types to be tested now.
	// /////////////////////////////////////////////////////

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[]")
	if errTy.IsBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146PST -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146PSW - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[0]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146Q14 -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146Q17 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[1]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146Q57 -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146Q59 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[abcd]")
	if errTy.IsBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146Q8T -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146Q8V - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[-1]")
	if errTy.IsBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146QA0 -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146QA2 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[2]")
	if errTy.IsBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146QA6 -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146QA8 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[1].")
	if errTy.IsBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146QNN - ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146QNP - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[1].objName")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146QNW - ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146QNY - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[1].intVal")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146QP1 -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146QP3 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[1].fltVal")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146QPN -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146QPQ - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[1].fltVal2")
	if errTy.IsBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146QQ5 -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146QQ7 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[1].moreArrays")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146QR2 -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146QR4 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[1].moreArrays.[0]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146QS4 -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146QS7 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[1].moreArrays.[0].k")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146QSO -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146QSQ - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[1].moreArrays.[1].k")
	if errTy.IsBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146QTH -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146QTJ - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[1].moreArrays.[1].k2")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146QU0 -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146QU2 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[0].contents")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146QV5 -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146QV7 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[0].contents.[0]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146QWA -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146QWC - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[0].contents.[0].[0]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146QWT -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146QWV - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[0].contents.[0].[4]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146QX5 -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146QX7 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[0].contents.[0].[5]")
	if errTy.IsBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146QXT -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146QXV - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[0].contents.[1].[4]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146QY5 -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146QY6 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[0].contents.[2].[0]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146R1F -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146R1H - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[0].contents.[2].[3]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146R4Y -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146R50 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[0].contents.[2].[3].[1]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146R5W -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146R5Y - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("arrObjs.[0].contents.[2].[3].[1].[1]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146R5W -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146R5Y - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	// Now we need to test from within the top level object

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.k")
	if errTy.IsBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146R87 -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146R89 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.key")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146R8T -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146R8V - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.key2")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146R9H -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146R9J - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.intVal")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146RA2 -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146RA5 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.arrInts")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146RAX -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146RAZ - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.arrInts.[0]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146RDK -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146RDM - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.arrInts.[4]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146RDB -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146RDD - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.arrInts.[5]")
	if errTy.IsBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146RG7 -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146RG9 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.arrFloats.[3]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146RGX -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146RGZ - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.arrStrings.[3]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146RJ1 -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146RJ3 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.arr2")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146SSH -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146SSK - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.arr2.[1]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146SU8 -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146SUA - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.arr2.[1].ki")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146SUN -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146SUP - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.arr2.[1].nested")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146SVM -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146SVP - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.arr2.[1].nested.[2]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146SYX -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146SYZ - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.arr2.[1].nested.[3]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146SW9 -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146SWC - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.arr2.[1].nested.[3].[1]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146SX7 -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146SX9 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.arr2.[1].nested.[3].[1].[4]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146SXP -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146SXR - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.arr2.[1].nested.[3].[1].[4]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146SZW -  ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146SZY - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.arr2.[1].nested.[3].[3].[4]")
	if errTy.IsNotBlank() {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#146T04 - ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#146T06 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	jo, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] %s:%d %v", filename, line, err)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	typ, val, errTy = jo.GetValueFromJsonObjectByJPath("obj.nestedObj.child0.child1.child2.finalChild.k7")
	if errTy.IsNotBlank() || (typ != "float64") || val.(float64) != 10.01 {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#148MG5 - ==> [[ERROR]] <== \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	} else {
		fmt.Printf("~~~~~~~\nL#148MG7 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", err, typ, val)
	}

	getTestsSuccessful = true
}

func TestSetValueInJsonObjectByJPath(t *testing.T) {
	if !getTestsSuccessful {
		// IMPORTANT: enable it later
		// t.Errorf("~~~~~~~\nE#147KEZ - Cannot proceed if Get Tests failed\n")
	}

	// ==========================================
	// IMPORTANT: UPDATES ONLY
	// ==========================================
	// NOTE: TOP LEVEL ONLY
	// ==========================================
	joOrig, err := ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr := SetValueInJsonObjectByJPath(joOrig, "str", "newStringValue")
	if setErr != nil {
		getTestsSuccessful = false
		t.Errorf("~~~~~~~\nE#147CJB - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("str")
		if getErr.IsNotBlank() || val != "newStringValue" {
			t.Errorf("~~~~~~~\nE#147CMI - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147COR - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "bool1", "someStringValue Because type strictness is not there")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147CXR - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("bool1")
		if getErr.IsNotBlank() || val != "someStringValue Because type strictness is not there" {
			t.Errorf("~~~~~~~\nE#147CXT - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147CXV - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "bool1", false)
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147D19 - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("bool1")
		if getErr.IsNotBlank() || val != false {
			t.Errorf("~~~~~~~\nE#147D1L - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147D1N - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "int", "someString_because_no_strictness_on_data_type")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147D46 - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("int")
		if getErr.IsNotBlank() || val != "someString_because_no_strictness_on_data_type" {
			t.Errorf("~~~~~~~\nE#147D4G - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147D4K - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "int", 987654)
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147D5G - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("int")
		if getErr.IsNotBlank() || val != 987654 {
			t.Errorf("~~~~~~~\nE#147D5O - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147D5Q - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "float", 9876.54)
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147D7J - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("float")
		if getErr.IsNotBlank() || val != 9876.54 {
			t.Errorf("~~~~~~~\nE#147D7L - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147D7N - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "float", "someString_no_strictness_on_data_type")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147D8K - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("float")
		if getErr.IsNotBlank() || val != "someString_no_strictness_on_data_type" {
			t.Errorf("~~~~~~~\nE#147D8N - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147D8P - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrInts", "someString_no_strictness_on_data_type")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147DBY - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrInts")
		if getErr.IsNotBlank() || val != "someString_no_strictness_on_data_type" {
			t.Errorf("~~~~~~~\nE#147DCF - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147DCH - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrInts", []int{1, 2, 3, 4, 5, 6, 7, 8})
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147DGC - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrInts")
		if getErr.IsNotBlank() || (typ != "array/int" && typ != "array/float64" && typ != "array/any") {
			t.Errorf("~~~~~~~\nE#147DGJ - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147DGL - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrFloats", "someString_no_strictness_on_data_type")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147DMV - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrFloats")
		if getErr.IsNotBlank() || val != "someString_no_strictness_on_data_type" {
			t.Errorf("~~~~~~~\nE#147DMT - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147DMQ - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrFloats", []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8})
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147DMM - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrFloats")
		if getErr.IsNotBlank() || (typ != "array/float64" && typ != "array/any") {
			t.Errorf("~~~~~~~\nE#147DMJ - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147DMH - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrStrings", []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8})
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147E0Z - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrStrings")
		if getErr.IsNotBlank() || (typ != "array/float64" && typ != "array/any") {
			t.Errorf("~~~~~~~\nE#147E12 - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147E14 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrStrings", []string{"three", "element", "string"})
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147E0Z - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrStrings")
		if getErr.IsNotBlank() || typ != "array/string" {
			t.Errorf("~~~~~~~\nE#147E12 - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147E14 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrBools", "someString_no_strictness_on_data_type")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147EM1 - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrBools")
		if getErr.IsNotBlank() || typ != "string" {
			t.Errorf("~~~~~~~\nE#147EM4 - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147EM7 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrBools", []string{"strict", "mode", "off"})
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147EIP - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrBools")
		if getErr.IsNotBlank() || typ != "array/string" {
			t.Errorf("~~~~~~~\nE#147EIU - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147EIW - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrBools", []bool{true, true, false})
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147EKP - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrBools")
		if getErr.IsNotBlank() || typ != "array/bool" {
			t.Errorf("~~~~~~~\nE#147EKS - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147EKU - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrObjs", []bool{true, true, false})
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147EX1 - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrObjs")
		if getErr.IsNotBlank() || typ != "array/bool" {
			t.Errorf("~~~~~~~\nE#147EX4 - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147EX4 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrObjs", "someString_no_strictness_on_data_type")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147EYH - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrObjs")
		if getErr.IsNotBlank() || typ != "string" {
			t.Errorf("~~~~~~~\nE#147EYL - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147EYN - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrObjs", []map[string]any{{"key": "value"}, {"k2": 1234}})
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147F06 - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrObjs")
		if getErr.IsNotBlank() || typ != "array/object" {
			t.Errorf("~~~~~~~\nE#147F0F - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147F0I - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj", map[string]any{"key": "value"})
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147K6I - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj")
		if getErr.IsNotBlank() || typ != "object" {
			t.Errorf("~~~~~~~\nE#147K6S - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147K6U - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj", "stringValue_no_strictness_on_data_type")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147K7X - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj")
		if getErr.IsNotBlank() || typ != "string" {
			t.Errorf("~~~~~~~\nE#147K7Z - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147K81 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	// ==========================================
	// NOTE: NESTED VALUES
	// ==========================================

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrObjs.[0].objName", 1234)
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147LVW - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrObjs.[0].objName")
		if getErr.IsNotBlank() || (typ != "float64" && typ != "int") {
			t.Errorf("~~~~~~~\nE#147LW3 - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147LW5 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrObjs.[0].contents.[1]", 1234)
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147M5K - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrObjs.[0].contents.[1]")
		if getErr.IsNotBlank() || (typ != "float64" && typ != "int") {
			t.Errorf("~~~~~~~\nE#147M5N - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147M5P - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrObjs.[0].contents.[1].[2]", 2345)
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147MO9 - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrObjs.[0].contents.[1].[2]")
		if getErr.IsNotBlank() || (typ != "float64" && typ != "int") {
			t.Errorf("~~~~~~~\nE#147MOB - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147MOD - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrObjs.[0].contents.[2].[3].[1]", 3456)
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147MYM - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrObjs.[0].contents.[2].[3].[1]")
		if getErr.IsNotBlank() || (typ != "float64" && typ != "int") {
			t.Errorf("~~~~~~~\nE#147MYV - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147MYX - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrObjs.[0].contents.[2].[3].[1].[1]", 3456)
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147N3H - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrObjs.[0].contents.[2].[3].[1].[1]")
		if getErr.IsNotBlank() || (typ != "float64" && typ != "int") {
			t.Errorf("~~~~~~~\nE#147N3J - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147N3L - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrObjs.[0].contents.[2].[3].[1].[1]", "someString_no_strictness_on_data_type")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147N5Q - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrObjs.[0].contents.[2].[3].[1].[1]")
		if getErr.IsNotBlank() || (typ != "string") {
			t.Errorf("~~~~~~~\nE#147N60 - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#14D54E - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrObjs.[0].contents.[2].[3].[1].[1]", nil)
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#1481KT - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrObjs.[0].contents.[2].[3].[1].[1]")
		if getErr.IsNotBlank() || (typ != "nil") {
			t.Errorf("~~~~~~~\nE#1481KX - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#1481L0 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrObjs.[1].moreArrays.[1].k2", "vaibhav_no_strictness_on_data_type")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147NCG - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrObjs.[1].moreArrays.[1].k2")
		if getErr.IsNotBlank() || (typ != "string") {
			t.Errorf("~~~~~~~\nE#147NCI - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147NCK - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrObjs.[1].moreArrays.[1].k2", nil)
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#148190 - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrObjs.[1].moreArrays.[1].k2")
		if getErr.IsNotBlank() || (typ != "nil") {
			t.Errorf("~~~~~~~\nE#148196 - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#148198 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.[1].moreArrays.[1].k2", "something")
	if setErr == nil {
		t.Errorf("~~~~~~~\nE#14890K - There should have been an error but there was none")
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.arrInts", "something")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#1489AP - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.arrInts")
		if getErr.IsNotBlank() || (typ != "string") {
			t.Errorf("~~~~~~~\nE#1489AZ - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#1489B2 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.arrInts.[2]", "someString_no_strictness_on_data_type")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#1489HX - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.arrInts.[2]")
		if getErr.IsNotBlank() || (typ != "string") {
			t.Errorf("~~~~~~~\nE#1489I6 - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#1489I8 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.arrInts.[2]", 999)
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#1489MO - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.arrInts.[2]")
		if getErr.IsNotBlank() || (typ != "int" && typ != "float64") {
			t.Errorf("~~~~~~~\nE#1489N9 - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#1489NB - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.arr2.[2]", 999)
	if setErr == nil {
		t.Errorf("~~~~~~~\nE#1489Q4 - Should have errored because the path does not exist")
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.arr2.[1]", 999)
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#1489QU - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.arr2.[1]")
		if getErr.IsNotBlank() || (typ != "int" && typ != "float64") {
			t.Errorf("~~~~~~~\nE#1489QX - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#1489QZ - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.arr2.[1]", 999)
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#148A09 - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.arr2.[1]")
		if getErr.IsNotBlank() || (typ != "int" && typ != "float64") || val.(int) != 999 {
			t.Errorf("~~~~~~~\nE#148A0D - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#148A0F - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.arr2.[0].list", 999)
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#148A0I - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.arr2.[0].list")
		if getErr.IsNotBlank() || (typ != "int" && typ != "float64") || val.(int) != 999 {
			t.Errorf("~~~~~~~\nE#148A0L - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#148A0O - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.arr2.[0].list.[2]", 999)
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#148A1M - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.arr2.[0].list.[2]")
		if getErr.IsNotBlank() || (typ != "int" && typ != "float64") || val.(int) != 999 {
			t.Errorf("~~~~~~~\nE#148A1O - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#148A1Q - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.arr2.[0].list.[2]", "someString_no_strictness_on_data_type")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#148A3N - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.arr2.[0].list.[2]")
		if getErr.IsNotBlank() || (typ != "string") || val.(string) != "someString_no_strictness_on_data_type" {
			t.Errorf("~~~~~~~\nE#148A3P - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#148A3R - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.arr2.[1].nested.[2]", "someString_no_strictness_on_data_type")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#148A8W - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.arr2.[1].nested.[2]")
		if getErr.IsNotBlank() || (typ != "string") || val.(string) != "someString_no_strictness_on_data_type" {
			t.Errorf("~~~~~~~\nE#148A8Z - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#148A91 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.arr2.[1].nested.[2].[3]", "someString_no_strictness_on_data_type")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#148LYP - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.arr2.[1].nested.[2].[3]")
		if getErr.IsNotBlank() || (typ != "string") || val.(string) != "someString_no_strictness_on_data_type" {
			t.Errorf("~~~~~~~\nE#148LYT - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#148LYV - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.arr2.[1].nested.[2].[3]", 555)
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#148M18 - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.arr2.[1].nested.[2].[3]")
		if getErr.IsNotBlank() || (typ != "int" && typ != "float64") || val.(int) != 555 {
			t.Errorf("~~~~~~~\nE#148M1A - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#148M1C - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.arr2.[1].nested.[3].[1].[4]", 666)
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#148M36 - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.arr2.[1].nested.[3].[1].[4]")
		if getErr.IsNotBlank() || (typ != "int" && typ != "float64") || val.(int) != 666 {
			t.Errorf("~~~~~~~\nE#148M39 - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#148M3B - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.arr2.[1].nested.[3].[1].[4]", map[string]any{"k8": "v8", "k9": "v9"})
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#148M61 - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.arr2.[1].nested.[3].[1].[4]")
		if getErr.IsNotBlank() || (typ != "object") || val.(map[string]any)["k8"] != "v8" {
			t.Errorf("~~~~~~~\nE#148M64 - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#148M66 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.nestedObj.child0.child1.child2.finalChild.k7", "someString_no_strictness_on_data_type")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#148MJ8 - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.nestedObj.child0.child1.child2.finalChild.k7")
		if getErr.IsNotBlank() || (typ != "string") || val.(string) != "someString_no_strictness_on_data_type" {
			t.Errorf("~~~~~~~\nE#148MJC - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#148MJE - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	// IMPORTANT: APPEND TO ARRAY
	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrObjs.[0].contents.[2].[3].[1].[]", "appendString_no_strictness_on_data_type")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#147N5Q - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrObjs.[0].contents.[2].[3].[1].[3]")
		if getErr.IsNotBlank() || (typ != "string") || val != "appendString_no_strictness_on_data_type" {
			t.Errorf("~~~~~~~\nE#147N60 - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#147N63 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.arr2.[1].nested.[]", "appendStringToArrayOfArrays_no_strictness_on_data_type")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#148MXM - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.arr2.[1].nested.[4]")
		if getErr.IsNotBlank() || (typ != "string") || val != "appendStringToArrayOfArrays_no_strictness_on_data_type" {
			t.Errorf("~~~~~~~\nE#148MXP - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#148MXR - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.arr2.[1].nested.[3].[]", "appendStringToArrayOfArrays_no_strictness_on_data_type")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#148N0N - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.arr2.[1].nested.[3].[4]")
		if getErr.IsNotBlank() || (typ != "string") || val != "appendStringToArrayOfArrays_no_strictness_on_data_type" {
			t.Errorf("~~~~~~~\nE#148N0P - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#148N0R - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.arr2.[1].nested.[3].[2].[]", 12.21)
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#148N2I - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.arr2.[1].nested.[3].[2].[4]")
		if getErr.IsNotBlank() || (typ != "float64") || val != 12.21 {
			t.Errorf("~~~~~~~\nE#148N2K - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#148N2N - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.arr2.[1].nested.[3].[2].[]", "with 5 elements ")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#148N6Z - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.arr2.[1].nested.[3].[2].[4]")
		if getErr.IsNotBlank() || (typ != "string") || val != "with 5 elements " {
			t.Errorf("~~~~~~~\nE#148N72 - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#148N74 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	// IMPORTANT: SET NEW ELEMENT IN OBJECT
	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "newTopLevelKey", "newTopLevelValue")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#14BK3Q - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("newTopLevelKey")
		if getErr.IsNotBlank() || (typ != "string") || val != "newTopLevelValue" {
			t.Errorf("~~~~~~~\nE#14BK4B - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#14BK4E - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrObjs.[0].newKey", "newValue")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#14BK8C - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrObjs.[0].newKey")
		if getErr.IsNotBlank() || (typ != "string") || val != "newValue" {
			t.Errorf("~~~~~~~\nE#14BK8E - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#14BK8H - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "arrObjs.[1].moreArrays.[1].newName", "Vaibhav")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#14BKB4 - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("arrObjs.[1].moreArrays.[1].newName")
		if getErr.IsNotBlank() || (typ != "string") || val != "Vaibhav" {
			t.Errorf("~~~~~~~\nE#14BKB7 - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#14BKBA - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.newEmail", "hi@vaibhavkaushal.com")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#14BKEQ - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.newEmail")
		if getErr.IsNotBlank() || (typ != "string") || val != "hi@vaibhavkaushal.com" {
			t.Errorf("~~~~~~~\nE#14BKES - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#14BKEU - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.newFloat", 909.808)
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#14BKGY - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.newFloat")
		if getErr.IsNotBlank() || (typ != "float64") || val != 909.808 {
			t.Errorf("~~~~~~~\nE#14BKH0 - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#14BKH2 - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.arr2.[1].newNumericString", "909.808")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#14BKOJ - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.arr2.[1].newNumericString")
		if getErr.IsNotBlank() || (typ != "string") || val != "909.808" {
			t.Errorf("~~~~~~~\nE#14BKOL - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#14BKON - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}

	joOrig, err = ToJsonObject(jsonString)
	if err != nil {
		_, filename, line, _ := runtime.Caller(0)
		t.Errorf("~~~~~~~\nToJsonObject Conversion failed in file %v at line %v | Error: %v", filename, line, err)
	}
	jo, setErr = SetValueInJsonObjectByJPath(joOrig, "obj.nestedObj.child0.child1.child2.finalChild.creator", "Vaibhav Kaushal")
	if setErr != nil {
		t.Errorf("~~~~~~~\nE#14BKRQ - Failed. Error: %v", setErr)
	} else {
		// Check what we got there
		typ, val, getErr := jo.GetValueFromJsonObjectByJPath("obj.nestedObj.child0.child1.child2.finalChild.creator")
		if getErr.IsNotBlank() || (typ != "string") || val != "Vaibhav Kaushal" {
			t.Errorf("~~~~~~~\nE#14BKRS - ==> [[ERROR]] <== Could not get UPDATED VALUE\nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		} else {
			fmt.Printf("~~~~~~~\nL#14BKRU - ==> AS EXPECTED!! <==  \nErr: %v \nTyp: %v \nValue : %v\n", getErr, typ, val)
		}
	}
}

func TestSetValueAndOverrideInJsonObjectByJPath(t *testing.T) {
	json, err := ToJsonObject(jsonString)
	if err != nil {
		fmt.Println("Failed to load JSON object")
	}
	fmt.Printf("Before value:\n %s\n", json.String())
	json, err = SetValueAndOverrideInJsonObjectByJPath(json, "obj.nestedObj.child0.key1.key2.key3", "Hurrah!", true)
	if err != nil {
		fmt.Println("Failed to set value", err)
	}
	fmt.Printf("After value:\n %s\n", json.String())
}

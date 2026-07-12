package syntax

import (
	"reflect"
)

/*       
GreaterOrEqual (>=), when input to an [OrderingRuleAssertion] function,
results in a greater or equal (GE) comparison.
*/
const GreaterOrEqual byte = 0x0
  
/*              
LessOrEqual (<=), when input to an [OrderingRuleAssertion] function,
results in a less or equal (LE) comparison.
*/                      
const LessOrEqual byte = 0x1 

/*
assertFirstStructField is a private function used for
firstComponent EQUALITY matching, in which the first
struct (ASN.1 SEQUENCE) field is matched.
*/
func assertFirstStructField(x any) (first any) {
        if isStruct(x) {
                if typ := reflect.TypeOf(x); typ.NumField() > 0 {
                        first = reflect.ValueOf(x).Field(0).Interface()
                }
        }

        return
}

/*
isStruct is a private function which returns a Boolean
value indicative of whether kind reflection revealed
the presence of a struct type.
*/
func isStruct(x any) (is bool) {
        if x != nil {
                is = reflect.TypeOf(x).Kind() == reflect.Struct
        }

        return
}

func assertString(x any, min int, name string) (str string, err error) {
        switch tv := x.(type) {
        case []byte:
                str, err = assertString(string(tv), min, name) 
        case string:
                if len(tv) < min && min != 0 {
                        err = errorBadLength(name, 0)
                        break
                }  
                str = tv
        default:   
                err = errorBadType(name)
        }

        return
}

func isNegativeInteger(x any) (is bool) {
        switch tv := x.(type) {
        case int:
                is = tv < 0
        case int8:
                is = tv < 0
        case int16:
                is = tv < 0
        case int32:
                is = tv < 0
        case int64:
                is = tv < 0
        }

        return
}

func castInt64(x any) (i int64, err error) {
        switch tv := x.(type) {
        case int:
                i = int64(tv)
        case int8:
                i = int64(tv)
        case int16:
                i = int64(tv)
        case int32:
                i = int64(tv)
        case int64: 
                i = tv
        default:
                err = errorBadType("castInt64")
        }

        return
}

func castUint64(x any) (i uint64, err error) {
        switch tv := x.(type) {
        case uint:
                i = uint64(tv)
        case uint8:
                i = uint64(tv)
        case uint16:
                i = uint64(tv)
        case uint32:
                i = uint64(tv)
        case uint64:
                i = tv
        default:
                err = errorBadType("castUint64")
        }

        return
}


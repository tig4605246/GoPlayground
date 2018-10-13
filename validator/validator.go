// Package util provides a variety of handy functions while developing backends
package validator

import (
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"strconv"
)

// InitValidator is used to initialize the validator.
// Call this funtion before using any validator.
func InitValidator() {
	govalidator.TagMap["required"] = govalidator.Validator(required)
	return
}

// CheckParam can help to verify if the input is valid or not.
//
// Input:
//  s      : input map for validation.
//  target : A map indicates that what validation should be applied to the input.
// Output:
//  bool   : Return true if valid, otherwise a false is returned.
//  error  : the output of the detail of the error.
func CheckParam(param map[string]interface{}, target map[string][]string) (bool, map[string]string) {
	var res bool = true
	var detail string
	var errOutput map[string]string
	errOutput = make(map[string]string)
	for key, value := range target {
		if input, ok := param[key]; ok {
			res, detail = checkParamWithArray(input, value)
			if !res {
				errOutput[key] = detail
			}
		}
	}
	if len(errOutput) > 0 {
		return false, errOutput
	}
	return true, errOutput
}

// checkParamWithArray helps CheckParam to verify the inputs recursively.
// It should not be called by other funcions.
// Currently supported tags:  DNS,Alphanumeric,Numeric,Alpha,Email.
func checkParamWithArray(param interface{}, target []string) (bool, string) {
	var res bool = true
	var detail string
	var subString string
	//The last one
	if input, ok := param.(string); ok {
		for i := 0; i < len(target); i++ {
			switch test := target[i]; test {
			case "DNS":
				res = govalidator.IsDNSName(input)
				break
			case "Alphanumeric":
				res = govalidator.IsAlphanumeric(input)
				break
			case "Numeric":
				res = govalidator.IsNumeric(input)
				break
			case "Alpha":
				res = govalidator.IsAlpha(input)
				break
			case "Email":
				res = govalidator.IsEmail(input)
				break
			case "Required":
				res = required(input)
				break
			case "MongoID":
				res = govalidator.IsMongoID(input)
				break
			case "Bool":
				res = isBool(input)
				break
			default:
				subString = fmt.Sprintf("%s %s\n", "No such validator: ", target[i])
				res = false
				break
			}
			if !res {
				subString = "[Check " + input + " with " + target[i] + " failed]"
				detail += subString
			}
		}
	} else if input, ok := param.([]string); ok { //Recursively disemble the []string
		if len(input) <= 0 {
			return true, detail
		}
		pop := input[0]
		res, subString = checkParamWithArray(pop, target)
		detail += subString
		input = append(input[:0], input[1:]...)
		res, subString = checkParamWithArray(input, target)
		detail += subString
	}
	if len(detail) == 0 {
		return true, detail
	} else {
		return false, detail
	}
}

// IsInt can help to verify if the input is Int and the value is within the range or not.
//
// If the bound array is not 2, it will skip the ranging test
func IsInt(param string, bound []int) (int, error) {
	//Check if the content type is match with the target
	var res bool = false
	res = govalidator.IsInt(param)
	if !res {
		return -1, errors.New("Input is not a Integer")
	}
	num, err := strconv.Atoi(param)
	if err != nil {
		return -1, err
	}
	if len(bound) == 2 {
		res = govalidator.InRangeInt(num, bound[0], bound[1])
		if !res {
			return num, errors.New("Input is out of range")
		}
	}
	return num, nil
}

// IsFloat64 can help to verify if the input is Float64 and the value is within the range or not.
//
// If the bound array is not 2, it will skip the ranging test
func IsFloat64(param string, bound []float64) (float64, error) {
	//Check if the content type is match with the target
	var res bool = false
	res = govalidator.IsFloat(param)
	if !res {
		return -1, errors.New("Input is not a Float")
	}
	num, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return -1, err
	}
	if len(bound) == 2 {
		res = govalidator.InRange(num, bound[0], bound[1])
		if !res {
			return num, errors.New("Input is out of range")
		}
	}
	return num, nil
}

//IsBool transform strings of true and false into bool.
func IsBool(str string) (bool, error) {
	return strconv.ParseBool(str)
}

// CheckStruct can help to verify the content of the struct with tags.
// It accepts recursive structs.
//
// Remember to Add valid:"[tags]" to the struct.
func CheckStruct(s interface{}) (bool, error) {
	result, err := govalidator.ValidateStruct(s)
	return result, err
	//BUG(Kevin Xu): haha
}

//IsStrInList can help to see if input is one of the target or not.
//a non-empty error will be returned if input does match any of the target.
func CheckStrInList(input string, target ...string) bool {
	for _, paramName := range target {
		if input == paramName {
			return true
		}
	}
	return false

}

//isMail can help to verify if the input is Email format or not
func isMail(param string) (bool, error) {
	res := govalidator.IsEmail(param)
	return res, errors.New("Nothing")
}

//Check if str is empty or not
func required(str string) bool {
	if len(str) > 0 {
		return true
	}
	return false
}

//Check bool
func isBool(str string) bool {
	_, err := strconv.ParseBool(str)
	if err != nil {
		return false
	}
	return true
}

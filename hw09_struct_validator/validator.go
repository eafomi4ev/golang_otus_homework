package hw09_struct_validator //nolint:golint,stylecheck
import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}

	accumulator := "Validation failed.\n"
	for _, errItem := range v {
		accumulator += fmt.Sprintf("field: %s; %v\n", errItem.Field, errItem.Err)
	}

	return accumulator
}

func Validate(v interface{}) error {
	rt := reflect.TypeOf(v)  //  reflect type
	rv := reflect.ValueOf(v) //  reflect value

	if rt.Kind() != reflect.Struct {
		fmt.Println(rv)
		return fmt.Errorf("passed argument is not a structure")
	}

	errorsAccumulator := ValidationErrors{}

	for i := 0; i < rt.NumField(); i++ {
		fieldT := rt.Field(i)
		fieldV := rv.Field(i)

		tagValue := fieldT.Tag.Get("validate")
		if len(tagValue) == 0 {
			continue
		}
		errorsAccumulator = append(errorsAccumulator, validateVal(tagValue, fieldV, fieldT.Name)...)
	}

	if len(errorsAccumulator) != 0 {
		return errorsAccumulator
	}

	return nil
}

type IValidator interface {
	Validate(pattern string, value interface{}) bool
}
type SValidator struct {
}

func validateVal(pattern string, fieldValue reflect.Value, fieldName string) ValidationErrors {
	patternParts := strings.Split(pattern, "|")
	errorsAccumulator := ValidationErrors{}

	for _, patternPart := range patternParts {
		rule := strings.Split(patternPart, ":") // e.g. patternPart == max:25 || in:tomato,omelet
		ruleName := rule[0]

		switch fieldValue.Kind() {
		case reflect.Slice:
			for i := 0; i < fieldValue.Len(); i++ {
				errorsAccumulator = append(errorsAccumulator, validateVal(pattern, fieldValue.Index(i), fieldName)...)
			}
		case reflect.Int:
			validator := IntValidator{}
			value := int(fieldValue.Int())

			switch ruleName {
			case "min":
				ruleValue, _ := strconv.Atoi(rule[1])
				if ok := validator.Min(ruleValue, value); !ok {
					errorsAccumulator = append(errorsAccumulator, ValidationError{
						Field: fieldName,
						Err:   fmt.Errorf("%w. Rule value: %d. Struct value: %d", ErrValidationMin, ruleValue, value),
					})
				}
			case "max":
				ruleValue, _ := strconv.Atoi(rule[1])
				if ok := validator.Max(ruleValue, value); !ok {
					errorsAccumulator = append(errorsAccumulator, ValidationError{
						Field: fieldName,
						Err:   fmt.Errorf("%w. Rule value: %d. Struct value: %d", ErrValidationMax, ruleValue, value),
					})
				}
			case "in":
				var ruleValues []int
				ruleV := fmt.Sprintf("[%s]", rule[1])
				if err := json.Unmarshal([]byte(ruleV), &ruleValues); err != nil {
					// todo: убрать панику
					panic(err)
				}

				if ok := validator.In(ruleValues, value); !ok {
					errorsAccumulator = append(errorsAccumulator, ValidationError{
						Field: fieldName,
						Err:   fmt.Errorf("%w. Rule values: %d. Struct value: %d", ErrValidationIn, ruleValues, value),
					})
				}
			}
		case reflect.String:
			validator := StringValidator{}
			value := fieldValue.String()

			switch ruleName {
			case "len":
				ruleValue, _ := strconv.Atoi(rule[1])
				if ok := validator.Len(ruleValue, value); !ok {
					errorsAccumulator = append(errorsAccumulator, ValidationError{
						Field: fieldName,
						Err:   fmt.Errorf("%w. Rule value: %d. Struct value: %s", ErrValidationLen, ruleValue, value),
					})
				}
			case "in":
				ruleValues := strings.Split(rule[1], ",")

				if ok := validator.In(ruleValues, value); !ok {
					errorsAccumulator = append(errorsAccumulator, ValidationError{
						Field: fieldName,
						Err:   fmt.Errorf("%w. Rule values: %s. Struct value: %s", ErrValidationIn, rule[1], value),
					})
				}
			case "regexp":
				ok, err := validator.Regexp(rule[1], value)
				if err != nil {
					errorsAccumulator = append(errorsAccumulator, ValidationError{
						Field: fieldName,
						Err:   err,
					})
				} else if !ok {
					errorsAccumulator = append(errorsAccumulator, ValidationError{
						Field: fieldName,
						Err:   fmt.Errorf("%w. Rule value: %s. Struct value: %s", ErrValidationRegexp, rule[1], value),
					})
				}
			}
		}
	}

	return errorsAccumulator
}

type IValidationInt interface {
	Max(max int, value int) bool
	Min(min int, value int) bool
	In(values []int, value int) bool
}

type IValidationString interface {
	Len(len int, value string) bool
	Regexp(pattern string, value string) (bool, error)
	In(values []string, value string) bool
}

// -------------------
// fmt.Println(reflect.TypeOf(fieldT))
// fmt.Println(rv.Type())
// fmt.Println("NAME:", fieldT.Name)
// fmt.Println(fieldT.Type)
// fmt.Println(fieldV.Kind())
// fmt.Println(fieldT.Name)
// fmt.Println("---------")
//
// if fieldT.Type.String() == "[]string" {
//	fmt.Println("max, slice")
// }
// if rt.Field(i).Name == "Phones" && rv.Field(i).Kind() == reflect {
//	fmt.Println("Yeeees")
// }
// -------------------

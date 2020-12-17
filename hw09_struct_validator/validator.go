package hw09_struct_validator //nolint:golint,stylecheck
import (
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
	rv := reflect.ValueOf(v) // reflect value
	rt := reflect.TypeOf(v)  //  reflect type

	errorsAccumulator := ValidationErrors{}

	if rt.Kind() != reflect.Struct {
		return fmt.Errorf("passed argument is not a structure")
	}

	var validator Validator

	for i := 0; i < rt.NumField(); i++ {
		fieldT := rt.Field(i)
		fieldV := rv.Field(i)

		// -------------------
		//fmt.Println(reflect.TypeOf(fieldT))
		//fmt.Println(rv.Type())
		//fmt.Println("NAME:", fieldT.Name)
		//fmt.Println(fieldT.Type)
		//fmt.Println(fieldV.Kind())
		//fmt.Println(fieldT.Name)
		//fmt.Println("---------")
		//
		//if fieldT.Type.String() == "[]string" {
		//	fmt.Println("max, slice")
		//}
		//if rt.Field(i).Name == "Phones" && rv.Field(i).Kind() == reflect {
		//	fmt.Println("Yeeees")
		//}
		// -------------------

		tagValue := fieldT.Tag.Get("validate")

		if len(tagValue) == 0 {
			continue
		}

		rules := strings.Split(tagValue, "|")

		for _, rule := range rules {
			tmp := strings.Split(rule, ":")
			ruleName := tmp[0]

			switch fieldV.Kind() {
			case reflect.Slice:
				fmt.Println("Обработка слайса стрингов")
			case reflect.Int:
				validator = IntValidator{}.Min
			case reflect.String:
				fmt.Println("Обработка стринга")
			}

			var ok bool
			switch ruleName {
			case "max":
				limit, _ := strconv.Atoi(tmp[1])
				val := rv.Field(i).Int()
				ok = validator.(IntValidation).Max(limit, int(val))
			case "min":
				limit, _ := strconv.Atoi(tmp[1])
				val := rv.Field(i).Int()
				ok = validator.(IntValidation).Min(limit, int(val))
			}

			if !ok {
				errorsAccumulator = append(errorsAccumulator, ValidationError{
					Field: fieldT.Name,
					Err:   ErrMaxValidation,
				})
			}
		}
	}

	if len(errorsAccumulator) != 0 {
		return errorsAccumulator
	}

	return nil
}

func getValidatorForType(rule string) {
	switch fieldV.Kind() {
	case reflect.Slice:
		fmt.Println("Обработка слайса стрингов")
	case reflect.Int:
		validator = IntValidator{}.Min
	case reflect.String:
		fmt.Println("Обработка стринга")
	}
}

type Validator interface {
	Validate(ruleName string, ruleValue interface{}, value interface{}) bool
}

type IntValidation interface {
	Max(max int, value int) bool
	Min(min int, value int) bool
	In(values []int, value int) bool
}

type StringValidation interface {
	Len(len int, value string) bool
	In(values []string, value string) bool
}

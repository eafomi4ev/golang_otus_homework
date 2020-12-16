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

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		fmt.Println(reflect.TypeOf(field))

		tagValue := field.Tag.Get("validate")

		if len(tagValue) == 0 {
			continue
		}

		rules := strings.Split(tagValue, "|")
		validationErrors := make([]error, 0)

		for _, rule := range rules {
			tmp := strings.Split(rule, ":")
			ruleName := tmp[0]

			switch ruleName {
			case "max":

				limit, _ := strconv.Atoi(tmp[1])
				val := rv.Field(i).Int()
				err := ValidateMax(int64(limit), val)
				if err != nil {
					errorsAccumulator = append(errorsAccumulator, ValidationError{
						Field: field.Name,
						Err:   err,
					})
				}
			case "min":
				limit, _ := strconv.Atoi(tmp[1])
				val := rv.Field(i).Int()
				err := ValidateMin(int64(limit), val)
				if err != nil {
					errorsAccumulator = append(errorsAccumulator, ValidationError{
						Field: field.Name,
						Err:   err,
					})
				}
			case "len":
				limit, _ := strconv.Atoi(tmp[1])
				val := rv.Field(i).String()
				validationErrors = ValidateLen(int64(limit), val)
			}

			if validationErrors != nil {
				for _, errItem := range validationErrors {
					errorsAccumulator = append(errorsAccumulator, ValidationError{
						Field: field.Name,
						Err:   errItem,
					})
				}
			}
		}
	}

	if len(errorsAccumulator) != 0 {
		return errorsAccumulator
	}

	return nil
}

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

	for i := 0; i < rt.NumField(); i++ {
		fieldT := rt.Field(i)
		//fieldV := rv.Field(i)

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
		validationErrors := make([]error, 0)

		for _, rule := range rules {
			tmp := strings.Split(rule, ":")
			ruleName := tmp[0]

			switch ruleName {
			case "max":
				switch fieldT.Type.String() {
				case "[]string":
					fmt.Println("Обработка слайса стрингов")
				case "[]int":
					fmt.Println("Обработка слайса интов")
				}
				limit, _ := strconv.Atoi(tmp[1])
				val := rv.Field(i).Int()
				err := ValidateMax(int64(limit), val)
				if err != nil {
					errorsAccumulator = append(errorsAccumulator, ValidationError{
						Field: fieldT.Name,
						Err:   err,
					})
				}
			case "min":
				limit, _ := strconv.Atoi(tmp[1])
				val := rv.Field(i).Int()
				err := ValidateMin(int64(limit), val)
				if err != nil {
					errorsAccumulator = append(errorsAccumulator, ValidationError{
						Field: fieldT.Name,
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
						Field: fieldT.Name,
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

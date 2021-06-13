package hw09structvalidator

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Value interface{}
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("поле %v: \"%v\" %v", v.Field, v.Value, v.Err)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	buffer := bytes.Buffer{}
	for _, val := range v {
		buffer.WriteString(val.Error())
		buffer.WriteString("\n")
	}
	return fmt.Sprint(buffer.String())
}

func (v ValidationErrors) wrongFields() string {
	buffer := bytes.Buffer{}
	for _, val := range v {
		switch v := val.Value.(type) {
		case string:
			buffer.WriteString(v)
		case int64:
			buffer.WriteString(fmt.Sprint(v))
		default:
			panic("you should add more types")
		}
		buffer.WriteString("\n")
	}
	return fmt.Sprint(buffer.String())
}

var (
	ErrItIsNotStruc = errors.New("object on input is not a struc")
	ErrInvalidCheck = errors.New("invalid check in validate tag")

	ErrValidationWrongIn           = errors.New("has wrong value of field")
	ErrValidationSliseIsEmpty      = errors.New("slise must be checked, but it is empty")
	ErrValidationStringWrongLen    = errors.New("has wrong length of field")
	ErrValidationStringWrongRegexp = errors.New("has wrong regexp of field")
	ErrValidationIntWrongMin       = errors.New("has wrong value of field, needs min=")
	ErrValidationIntWrongMax       = errors.New("has wrong value of field, needs max=")
)

const (
	strucTag            = "validate"
	logicalANDSeparator = "|"
	logicalORSeparator  = ","
	lenValidator        = "len:"
	inValidator         = "in:"
	regexpValidator     = "regexp:"
	minValidator        = "min:"
	maxValidator        = "max:"
)

// Validate checks the structure by special strucTag on fields (for ex. "validate:").
// It retruns either a program error or ValidationErrors, which is a slice of structures
// with field name and error of its validation.
// It can work with:
// - `int`, `[]int`;
// - `string`, `[]string`.
// With another types it panics.
// Available validators:
// - For int:
//     * `min:10` - integer must be above 10;
//     * `max:20` - integer must be below 20;
//     * `in:256,1024` - integer can be 256 or 1024;
// - For string:
//     * `len:32` - length of string must be 32;
//     * `regexp:\\d+` - by regular expression, string must consist of numbers
//     * `in:foo,bar` - integer can be "foo" or "bar"
// - For slices, it validates every item of slice
//
// Combinations of validators with logical AND is available by `|` sign, examle:
// * `min:0|max:10` - int must be between 0 and 10;
// * `regexp:\\d+|len:20` - string must consist of numbers and have length of 20.
func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return ErrItIsNotStruc
	}
	rto := reflect.TypeOf(v)
	var errorsFromFields ValidationErrors
	for i := 0; i < rto.NumField(); i++ {
		field := rto.Field(i)
		if validate, ok := field.Tag.Lookup(strucTag); ok && validate != "" {
			errorsFromField, err := validateFieldOrSlice(field.Name, rv.Field(i), field.Type, validate)
			if err != nil {
				return fmt.Errorf("in validation of field %v has got an error: %w", field.Name, err)
			}
			errorsFromFields = append(errorsFromFields, errorsFromField...)
		}
	}
	if len(errorsFromFields) > 0 {
		return errorsFromFields
	}
	return nil
}

func validateFieldOrSlice(name string, val reflect.Value, tp reflect.Type, checks string) (ValidationErrors, error) {
	var validationErrors ValidationErrors
	var err error
	switch tp.Kind() {
	case reflect.String, reflect.Int:
		err = validateField(name, val, checks, &validationErrors)
	case reflect.Slice:
		if val.Len() == 0 {
			validationErrors = append(validationErrors, ValidationError{
				Field: name,
				Value: fmt.Sprint(val.Interface()),
				Err:   ErrValidationSliseIsEmpty,
			})
		}
		for i := 0; i < val.Len(); i++ {
			switch val.Index(0).Kind() {
			case reflect.String, reflect.Int:
				err = validateField(name, val.Index(i), checks, &validationErrors)
			default:
				panic("you should add more types")
			}
			if err != nil {
				return nil, err
			}
		}
	default:
		panic("you should add more types")
	}
	if err != nil {
		return nil, err
	}
	return validationErrors, nil
}

func validateField(name string, rV reflect.Value, aLLchecks string, validationErrors *ValidationErrors) error {
	var err error
	for _, someTypeOfCheck := range strings.Split(aLLchecks, logicalANDSeparator) {
		switch {
		case rV.Kind() == reflect.String && strings.HasPrefix(someTypeOfCheck, lenValidator):
			seqOfLenChecks := strings.TrimPrefix(someTypeOfCheck, lenValidator)
			err = validateStringForLen(name, rV, seqOfLenChecks, validationErrors)
		case (rV.Kind() == reflect.String || rV.Kind() == reflect.Int) && strings.HasPrefix(someTypeOfCheck, inValidator):
			seqOfInChecks := strings.TrimPrefix(someTypeOfCheck, inValidator)
			validateForIn(name, rV, seqOfInChecks, validationErrors)
		case rV.Kind() == reflect.String && strings.HasPrefix(someTypeOfCheck, regexpValidator):
			regexpChecks := strings.TrimPrefix(someTypeOfCheck, regexpValidator)
			validateStringForRegexp(name, rV, regexpChecks, validationErrors)
		case rV.Kind() == reflect.Int && (strings.HasPrefix(someTypeOfCheck, minValidator) || strings.HasPrefix(someTypeOfCheck, maxValidator)):
			err = validateIntForMinMax(name, rV, someTypeOfCheck, validationErrors)
		default:
			return ErrInvalidCheck
		}
	}
	return err
}

func validateStringForLen(name string, rV reflect.Value, seqOfLenChecks string, validationErrors *ValidationErrors) error {
	bufOfChecks := make(ValidationErrors, 0, 1)
	for _, oneOfLenCheck := range strings.Split(seqOfLenChecks, logicalORSeparator) {
		needLen, err := strconv.Atoi(oneOfLenCheck)
		if err != nil {
			return fmt.Errorf("in conversion of string %v to int has got an error: %w", oneOfLenCheck, err)
		}
		if needLen == len(rV.String()) {
			return nil
		}
		err = fmt.Errorf("%w, needs len=%v", ErrValidationStringWrongLen, needLen)
		bufOfChecks = append(bufOfChecks, ValidationError{
			Field: name,
			Value: rV.String(),
			Err:   err,
		})
	}
	*validationErrors = append(*validationErrors, bufOfChecks...)
	return nil
}

func validateForIn(name string, rV reflect.Value, seqOfInChecks string, validationErrors *ValidationErrors) {
	bufOfChecks := make(ValidationErrors, 0, 1)
	for _, oneOfInCheck := range strings.Split(seqOfInChecks, logicalORSeparator) {
		if oneOfInCheck == fmt.Sprint(rV.Interface()) {
			return
		}
		err := fmt.Errorf("%w, needs value=%v", ErrValidationWrongIn, oneOfInCheck)
		bufOfChecks = append(bufOfChecks, ValidationError{
			Field: name,
			Value: fmt.Sprint(rV.Interface()),
			Err:   err,
		})
	}
	*validationErrors = append(*validationErrors, bufOfChecks...)
}

func validateStringForRegexp(name string, rV reflect.Value, regexpChecks string, validationErrors *ValidationErrors) {
	bufOfChecks := make(ValidationErrors, 0, 1)
	regexp := regexp.MustCompile(regexpChecks)
	if regexp.MatchString(rV.String()) {
		return
	}
	err := fmt.Errorf("%w, needs regexp=%v", ErrValidationStringWrongRegexp, regexpChecks)
	bufOfChecks = append(bufOfChecks, ValidationError{
		Field: name,
		Value: rV.String(),
		Err:   err,
	})
	*validationErrors = append(*validationErrors, bufOfChecks...)
}

func validateIntForMinMax(name string, rV reflect.Value, someTypeOfCheck string, validationErrors *ValidationErrors) error {
	var needMin int
	var needMax int
	var checkIsMin bool
	var checkIsMax bool
	var maxOrMinErr error
	var err error
	switch {
	case strings.HasPrefix(someTypeOfCheck, minValidator):
		needMin, err = strconv.Atoi(strings.TrimPrefix(someTypeOfCheck, minValidator))
		checkIsMin = true
		maxOrMinErr = fmt.Errorf("%w%v", ErrValidationIntWrongMin, needMin)
	case strings.HasPrefix(someTypeOfCheck, maxValidator):
		needMax, err = strconv.Atoi(strings.TrimPrefix(someTypeOfCheck, maxValidator))
		checkIsMax = true
		maxOrMinErr = fmt.Errorf("%w%v", ErrValidationIntWrongMax, needMax)
	}
	if err != nil {
		return fmt.Errorf("in conversion of string %v to int has got an error: %w", someTypeOfCheck, err)
	}
	if !(checkIsMin && int64(needMin) <= rV.Int() || checkIsMax && int64(needMax) >= rV.Int()) {
		*validationErrors = append(*validationErrors, ValidationError{
			Field: name,
			Value: rV.Int(),
			Err:   maxOrMinErr,
		})
	}
	return nil
}

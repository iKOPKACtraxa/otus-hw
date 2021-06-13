package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}
	App struct {
		Version string `validate:"len:5"`
	}
	AppNoValCheck struct {
		Version string `validate:"errCheck:5"`
	}
	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}
	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
	TestStringComplex struct {
		StringSlice []string `validate:"len:5,10|in:abcd5,abcdefgi10"`
	}
	TestIntComplex struct {
		IntSlice []int `validate:"min:10|max:100"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name                 string
		in                   interface{}
		expectedErrs         []error
		expectedAmountOfErrs int
		thisIsWrong          []string
		thisIsRight          []string
	}{
		{
			name:         "not struc at input, must err",
			in:           "string, not struc",
			expectedErrs: []error{ErrItIsNotStruc},
		},
		{
			name:         "String with Len, is right",
			in:           App{Version: "abcde"},
			expectedErrs: nil,
			thisIsWrong:  []string{},
			thisIsRight:  []string{"abcde"},
		},
		{
			name:                 "String with Len, must err",
			in:                   App{Version: "abcdefg1234567"},
			expectedErrs:         []error{ErrValidationStringWrongLen},
			expectedAmountOfErrs: 1,
			thisIsWrong:          []string{"abcdefg1234567"},
			thisIsRight:          []string{},
		},
		{
			name:         "Validate chick is not valid, must err",
			in:           AppNoValCheck{Version: "foobar"},
			expectedErrs: []error{ErrInvalidCheck},
		},
		{
			name:                 "String (in slice) with complex checks, must err",
			in:                   TestStringComplex{StringSlice: []string{"abcd5", "abcdefgi10", "ab3"}},
			expectedErrs:         []error{ErrValidationStringWrongLen, ErrValidationWrongIn},
			expectedAmountOfErrs: 4,
			thisIsWrong:          []string{"ab3"},
			thisIsRight:          []string{"abcd5", "abcdefgi10"},
		},
		{
			name:                 "Test for regexp, empty slice and UserRole type, must err",
			in:                   User{ID: "incoorectID", Name: "", Age: 0, Email: "ivanovAmail.ru", Role: "auditor", Phones: []string{}, meta: []byte{}},
			expectedErrs:         []error{ErrValidationIntWrongMin, ErrValidationStringWrongRegexp, ErrValidationWrongIn, ErrValidationSliseIsEmpty, ErrValidationStringWrongLen},
			expectedAmountOfErrs: 6,
			thisIsWrong:          []string{"ivanovAmail.ru", "auditor", "0", "[]"},
			thisIsRight:          []string{"123456789012345678901234567890123-36"},
		},
		{
			name:         "Test for regexp, is right",
			in:           User{ID: "123456789012345678901234567890123-36", Name: "", Age: 20, Email: "ivanov@mail.ru", Role: "admin", Phones: []string{"12345678-11"}, meta: []byte{}},
			expectedErrs: nil,
		},
		{
			name:         "Int with In check, is right",
			in:           Response{Code: 500, Body: ""},
			expectedErrs: nil,
			thisIsWrong:  []string{},
			thisIsRight:  []string{"500"},
		},
		{
			name: "Int with In check, must err",
			in: Response{
				Code: 666,
				Body: "",
			},
			expectedErrs:         []error{ErrValidationWrongIn},
			expectedAmountOfErrs: 3,
			thisIsWrong:          []string{"666"},
			thisIsRight:          []string{},
		},
		{
			name: "Int (in slice) with complex checks, must err",
			in: TestIntComplex{
				IntSlice: []int{5, 1, 20, 10, 56, 133, 100, 29},
			},
			expectedErrs:         []error{ErrValidationIntWrongMin, ErrValidationIntWrongMax},
			expectedAmountOfErrs: 3,
			thisIsWrong:          []string{"5", "1", "133"},
			thisIsRight:          []string{"20", "10", "56", "100", "29"},
		},
		{
			name: "Token should be checked",
			in: Token{
				Header:    []byte{'f', 'd'},
				Payload:   []byte{'h', 'e', 'e'},
				Signature: []byte{'d', 'h'},
			},
			expectedErrs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("subtest: %v", tt.name), func(t *testing.T) {
			tt := tt
			t.Parallel()
			errFromValidate := Validate(tt.in)
			if tt.expectedErrs == nil {
				require.NoErrorf(t, errFromValidate, "need: no error, got: ", errFromValidate)
			}
			var pValidationErrors ValidationErrors
			isValidationErrors := errors.As(errFromValidate, &pValidationErrors)
			for _, oneExpErr := range tt.expectedErrs {
				switch {
				case errors.Is(oneExpErr, ErrItIsNotStruc):
					require.ErrorIs(t, errFromValidate, ErrItIsNotStruc, "need: %v, got: %v", ErrItIsNotStruc, errFromValidate)
				case errors.Is(oneExpErr, ErrInvalidCheck):
					require.ErrorIs(t, errFromValidate, ErrInvalidCheck, "need: %v, got: %v", ErrInvalidCheck, errFromValidate)
				case isValidationErrors:
					pExpectedErr := oneExpErr
					require.ErrorAsf(t, errFromValidate, &pExpectedErr, "need: %v, got: %v", pExpectedErr, errFromValidate)
					require.Equal(t, tt.expectedAmountOfErrs, len(pValidationErrors), "expected %v errors, got %v", tt.expectedAmountOfErrs, len(pValidationErrors))
				}
			}
			for _, v := range tt.thisIsWrong {
				require.Truef(t, strings.Contains(pValidationErrors.wrongFields(), v), "this value is missed in errors: %v", v)
			}
			for _, v := range tt.thisIsRight {
				require.Falsef(t, strings.Contains(pValidationErrors.wrongFields(), v), "this value must not to be in errors: %v", v)
			}
		})
	}
}

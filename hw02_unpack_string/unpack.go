package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inStr string) (string, error) {
	var symbolToWrite string
	var stringToBuild strings.Builder
	var flagOfEscapeExist bool
	const escapeCharacter string = "\\"
	for i, runeSymbol := range inStr {
		symbol := string(runeSymbol)
		intSymbol, err := strconv.Atoi(symbol)
		gotNumber := err == nil
		if gotNumber && symbolToWrite != "" { // обычный случай, для вывода, например, a4 -> aaaa
			if intSymbol > 0 { // условие нужно для того, чтобы не было вывода при 0, например, a0 -> ""
				stringToBuild.WriteString(strings.Repeat(symbolToWrite, intSymbol))
			}
			symbolToWrite = ""
		} else if gotNumber && symbolToWrite == "" { // случаи "некорректных строк" или экранированных символов
			if flagOfEscapeExist { // случай экранированных символов
				symbolToWrite = symbol
				flagOfEscapeExist = false
			} else { // случай "некорректных строк"
				return "", ErrInvalidString
			}
		}
		if !gotNumber { // обычный случай, когда текущий символ не число, например ab -> ab
			stringToBuild.WriteString(symbolToWrite) // нужно вывести сохраненный в прошлой итерации символ (a)
			symbolToWrite = symbol                   // а текущий символ сохранить для следующей итерации (b)
			// логика экранирования "\"
			if symbol == escapeCharacter {
				if !flagOfEscapeExist {
					flagOfEscapeExist = true
					symbolToWrite = ""
				} else {
					flagOfEscapeExist = false
				}
			}
		}
		// если это была последняя итерация, то выводим сразу символ к выводу
		if i == len(inStr)-1 {
			stringToBuild.WriteString(symbolToWrite)
		}
		// fmt.Println(stringToBuild.String()) // uncomment to see stringToBuild at every step
	}
	return stringToBuild.String(), nil
}

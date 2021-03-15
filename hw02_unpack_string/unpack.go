package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inStr string) (string, error) {
	var symbolToWrite string          // это символ, который нужно вывести
	var stringToBuild strings.Builder // это строка к выводу
	var flag bool                     // флаг экранирования
	for i, runeSymbol := range inStr {
		symbol := string(runeSymbol)           // это символ, но пока не понятно, число это или буква
		intSymbol, err := strconv.Atoi(symbol) // пытаемся получить из символа число
		gotNumber := err == nil                // обозначает что удалось достать число
		if gotNumber && symbolToWrite != "" {  // если получилось число и из предыдущей итерации есть символ для вывода...
			if intSymbol > 0 { // ...и если нужно вывести >0 раз, то выводим intSymbol раз, иначе нет вывода (например для a0)
				stringToBuild.WriteString(strings.Repeat(symbolToWrite, intSymbol))
			}
			// так как вывели символы "через распаковку", то переменную нужно очистить
			// иначе она выведется на следующей итерации еще раз
			symbolToWrite = ""
		} else if gotNumber && symbolToWrite == "" { // если получилось число и если нет символа для вывода,
			// сначала проверим был ли флаг, если был,
			// значит это заэкранированное число и нужно передать symbol в symbolToWrite, и снять флаг
			if flag {
				symbolToWrite = symbol
				flag = false
			} else { // а если и флага не было, то это "некорректная строка"
				return "", ErrInvalidString
			}
		}
		if !gotNumber { // если текущий символ не число,
			// значит из предыдущей итерации нужно вывести символ к выводу
			stringToBuild.WriteString(symbolToWrite)
			// и так как текущий символ не число,
			// то мы сохраняем его к выводу на следующей итерации
			symbolToWrite = symbol
			// логика экранирования:
			if symbol == "\\" { // если текущий символ "\"...
				if !flag { // ...и флага не было, то поднимаем флаг и symbolToWrite=""
					flag = true
					symbolToWrite = ""
				} else { // ...и флаг был то снимаем флаг
					flag = false
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

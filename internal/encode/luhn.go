package encode

import (
	"log"
	"strconv"
)

// TOREDO corrently uses 	"github.com/theplant/luhn"
// LuhnValid return true if input is Luhn valid
func LuhnValid(inter interface{}) bool {
	switch t := inter.(type) {
	case string:
		return LuhnValidString(inter.(string))
	case int:
		return LuhnValidInt(inter.(int))
	default:
		log.Println("unknown type: ", t)
	}
	return false
}

// LuhnValidInt return true if valid on int
func LuhnValidInt(number int) bool {
	return number%10 == LuhnReturnNumber(number/10)
}

// LuhnCheckNumber returns control number for input int
func LuhnReturnNumber(number int) int {
	sum := 0
	for i := 0; number > 0; i++ {
		sum += (2*(i%2)*(number%10))/10 + (2*(i%2)*(number%10))%10
		number = number / 10
	}
	return sum % 10
}

// LuhnValidString return true if valid  on string
func LuhnValidString(s string) bool {
	in, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return LuhnValidInt(in)
}

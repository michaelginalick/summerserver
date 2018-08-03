package nums

import (
	"strconv"
)

// GetIndividualDays :  gets days from range 9-11; [9,10,11]
func GetIndividualDays(firstInt, lastInt int) []string {
	var s []string

	i := firstInt

	for i <= lastInt {
		j := strconv.Itoa(i)
		s = append(s, j)
		i++
	}
	return s
}

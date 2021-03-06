package calendar

// Value : a struct that is the value in the calander map
type Value struct {
	Name string
	Days int
}

// GetMonth : given a string it will return the value in the map
func GetMonth(s string) Value {

	calendarMap := map[string]Value{
		"January":   {"January", 31},
		"February":  {"February", 28},
		"March":     {"March", 31},
		"April":     {"April", 30},
		"May":       {"May", 31},
		"June":      {"June", 30},
		"July":      {"July", 31},
		"August":    {"August", 31},
		"September": {"September", 30},
		"October":   {"October", 31},
		"November":  {"November", 30},
		"December":  {"December", 31},
		"Aug":       {"August", 31},
		"Sept":      {"September", 30},
		"Feb":       {"February", 28},
		"Apr":       {"April", 30},
		"Jan.":			 {"January", 31},
		"JAN.":			 {"January", 31},
		"Feb.":			 {"February", 28},
		"FEB.":			 {"February", 28},
	}

	return calendarMap[s]

}

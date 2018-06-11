package calendar


// GetMonth : given a string it will return the value in the map
func GetMonth(s string) string {
	
	calendarMap := map[string]string {
		"January":"January",
		"February":"February",
		"March":"March",
		"April":"April",
		"May":"May",
		"June":"June",
		"July":"July",
		"August":"August",
		"September":"September",
		"October":"October",
		"November":"November",
		"December":"December",
	}

	return calendarMap[s]


}




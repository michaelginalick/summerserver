package event

// Event : an event struct
type Event struct {
	Name   string
	Link   string
	Month  string
	Days   []string
	IndividualDays []int
	FestivalLength int
	Next   *Event
}

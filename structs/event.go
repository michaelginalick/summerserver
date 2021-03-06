package event

import (
	"fmt"
)

// Event : an event struct
type Event struct {
	ID             int
	Name           string
	Link           string
	Month          string
	Days           []string
	Year           string
	IndividualDays []string
	Day            int
	Location			 string
	Next           *Event
}

// List : return a new event list
func List() *Event {
	return &Event{}
}

// AddBeginning :  append new event to beginning of list
func AddBeginning(newEvent, eventList *Event) *Event {
	newEvent.Next = eventList
	return newEvent
}

// PrintList : print list one event at a time
func PrintList(eventList *Event) {
	for i := eventList; i != nil; i = i.Next {
		fmt.Println(i.Name, i.Link, i.Month, i.Days, i.Year, i.IndividualDays)
	}
}

// PrintListByMonth :  print all event for a given month
func PrintListByMonth(e *Event, s string) {
	for i := e; i != nil; i = i.Next {
		if i.Month == s {
			fmt.Println(i.Name, i.Link, i.Month, i.Days, i.Year, i.IndividualDays)
		}
	}
}

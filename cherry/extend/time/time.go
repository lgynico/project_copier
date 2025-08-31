package cherryTime

import "time"

type CherryTime struct {
	time.Time
}

func NewTime(tt time.Time, setGlobal bool) CherryTime {
	ct := CherryTime{}

	if setGlobal {
		ct.Time = tt.In(offsetLocation).Add(offsetTime)
	} else {
		ct.Time = tt
	}

	return ct
}

func Now() CherryTime {
	return NewTime(time.Now(), true)
}

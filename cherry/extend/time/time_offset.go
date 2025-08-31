package cherryTime

import "time"

var (
	offsetTime     time.Duration
	offsetLocation *time.Location
)

func init() {
	SetOffsetLocation("Local")
}

func AddOffsetTime(t time.Duration) {
	offsetTime = t
}

func SubOffsetTime(t time.Duration) {
	offsetTime = -t
}

func SetOffsetLocation(name string) (err error) {
	offsetLocation, err = time.LoadLocation(name)
	return
}

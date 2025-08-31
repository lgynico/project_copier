package cherryTime

func (p CherryTime) ToMillisecond() int64 {
	return p.Time.UnixMilli()
}

package utils

import "time"

func TimeToEpochUTC(timeStamp int64) time.Time {
	return time.Unix(timeStamp, 0).UTC()
}

package types

import (
	"math"
	"net/http"
	"time"
)

type Backoff func(min, max time.Duration, retryNum int, resp *http.Response) time.Duration

func DefaultBackoff(min, max time.Duration, retryNum int, resp *http.Response) time.Duration {

	multiply := math.Pow(2, float64(retryNum)) * float64(min)
	wait := time.Duration(multiply)
	if float64(wait) != multiply || wait > max {
		wait = max
	}
	return wait
}

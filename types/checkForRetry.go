package types

import "net/http"

/*

 The default retry policy which checks
 for a range of status codes in the response and retry on 500 range responses.
*/

type CheckForRetry func(resp *http.Response, err error) (bool, error)

func DefaultRetryPolicy(resp *http.Response, err error) (bool, error) {
	if err != nil {
		return true, err
	}
	if resp.StatusCode == 0 || resp.StatusCode >= 500 {
		return true, nil
	}
	return false, nil
}

/*

   A backoff strategy specifies how long to wait between retries.

*/

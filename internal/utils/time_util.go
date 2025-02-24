package utils

import "time"

func DataNow() (t time.Time, err error) {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return time.Time{}, err
	}

	return time.Now().In(loc), err
}

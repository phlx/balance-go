package time

import (
	"log"
	"time"

	"github.com/pkg/errors"

	"balance/internal/times"
)

func Location() *time.Location {
	location, err := time.LoadLocation(times.TimeZone)
	if err != nil {
		err = errors.Wrapf(err, "unable to load location with timezone %s", times.TimeZone)
		log.Fatal(err)
	}
	return location
}

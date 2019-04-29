package customtime

import (
	"sync"
	"time"

	"bitbucket.org/kudoindonesia/koolkit/koollog"
	"go.uber.org/zap"
)

const (
	defaultTimezone = "Asia/Jakarta"
)

var (
	location   *time.Location
	locationMu sync.RWMutex
)

func init() {
	var err error

	location, err = time.LoadLocation(defaultTimezone)
	if err != nil {
		koollog.Fatal("fail to load timezone location", zap.Error(err))
	}
}

// SetLocation will sets the location
func SetLocation(tz string) error {
	// set Lock, so no one will read/write location until this function is done
	locationMu.Lock()
	defer locationMu.Unlock()

	var err error
	location, err = time.LoadLocation(tz)

	return err
}

// Location returns pointer of time.Location, the default location is "Asia/Jakarta"
func Location() *time.Location {
	// lock for read, allow other go routine to read location, but block other go routine write location
	locationMu.RLock()
	defer locationMu.RUnlock()

	return location
}

// Now return current time using current location
func Now() time.Time {
	// lock for read, allow other go routine to read location, but block other go routine write location
	locationMu.RLock()
	defer locationMu.RUnlock()

	return time.Now().In(location)
}

// Parse return current time using current location
func Parse(layout, value string) (time.Time, error) {
	// lock for read, allow other go routine to read location, but block other go routine write location
	locationMu.RLock()
	defer locationMu.RUnlock()

	return time.ParseInLocation(layout, value, location)
}

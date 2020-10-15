package namecheck

import "fmt"

type ErrUnknownAvailability struct {
	Username string
	Platform string
	Cause    error
}

func (err *ErrUnknownAvailability) Error() string {
	const tmpl = "namecheck: availablity of %q on %s could not be determined"
	return fmt.Sprintf(tmpl, err.Username, err.Platform)
}

func (err *ErrUnknownAvailability) Unwrap() error {
	return err.Cause
}

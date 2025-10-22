package util

import "errors"

// JoinErrors join given err with others
func JoinErrors(err error, errs ...error) error {
	joined := append([]error{err}, errs...)
	return errors.Join(joined...)
}

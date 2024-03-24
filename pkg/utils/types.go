package utils

import "errors"

func ErrorToString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func StringToError(err string) error {
	if err == "" {
		return nil
	}
	return errors.New(err)
}

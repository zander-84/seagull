package validator

import (
	"errors"
	"net"
	"reflect"
	"unicode/utf8"
)

func Require(data any, msg string) error {
	if data == nil {
		return errors.New(msg)
	}
	switch v := data.(type) {
	case bool:
		if v == false {
			return errors.New(msg)
		}
		return nil
	case *bool:
		if v == nil || *v == false {
			return errors.New(msg)
		}
		return nil
	case int, int8, int16, int32, int64:
		if v == 0 {
			return errors.New(msg)
		}
		return nil
	case *int, *int8, *int16, *int32, *int64:
		if v == nil || v == 0 {
			return errors.New(msg)
		}
		return nil
	case uint, uint8, uint16, uint32, uint64:
		if v == 0 {
			return errors.New(msg)
		}
		return nil
	case *uint, *uint8, *uint16, *uint32, *uint64:
		if v == nil || v == 0 {
			return errors.New(msg)
		}
		return nil
	case float32, float64:
		if v == 0 {
			return errors.New(msg)
		}
		return nil
	case *float32, *float64:
		if v == nil || v == 0 {
			return errors.New(msg)
		}
		return nil
	case string:
		if len(v) == 0 {
			return errors.New(msg)
		}
		return nil
	case *string:
		if v == nil || len(*v) == 0 {
			return errors.New(msg)
		}
		return nil
	default:
		kind := reflect.TypeOf(data).Kind()
		if kind == reflect.Ptr {
			dataKind := reflect.ValueOf(data).Elem().Kind()
			if dataKind == reflect.Invalid || dataKind == reflect.Struct {
				if reflect.ValueOf(data).IsNil() {
					return errors.New(msg)
				}
			}
		} else if kind == reflect.Slice || kind == reflect.Map {
			if reflect.ValueOf(data).Len() == 0 {
				return errors.New(msg)
			}
		}
		return nil
	}
}

func IsIpv4(v string, msg string) error {
	ip := net.ParseIP(v)
	if ip != nil && ip.To4() != nil {
		return nil
	}
	return errors.New(msg)
}

func HasIntMax(v int, max int, msg string) error {
	if v >= max {
		return errors.New(msg)
	}
	return nil
}

func HasIntMin(v int, min int, msg string) error {
	if v >= min {
		return nil
	}
	return errors.New(msg)
}

func HasIntRange(v int, min int, max int, msg string) error {
	if err := HasIntMin(v, min, msg); err != nil {
		return err
	} else if err = HasIntMax(v, max, msg); err != nil {
		return err
	}
	return nil
}

func HasStrRange(v string, min int, max int, msg string) error {
	if err := HasStrMin(v, min, msg); err != nil {
		return err
	} else if err = HasStrMax(v, max, msg); err != nil {
		return err
	}
	return nil
}

func HasStrMin(v string, min int, msg string) error {
	if utf8.RuneCountInString(v) >= min {
		return nil
	}
	return errors.New(msg)
}

func HasStrMax(v string, max int, msg string) error {
	if utf8.RuneCountInString(v) >= max {
		return errors.New(msg)
	}
	return nil
}

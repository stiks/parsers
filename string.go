package parsers

/*
   Come from here:
   https://github.com/gomodule/redigo/blob/master/redis/reply.go#L377
*/

import (
	"errors"
	"fmt"
	"strconv"
)

// Error represents an error returned in a command data.
type Error string

func (err Error) Error() string { return string(err.Error()) }

// ErrNil indicates that a data value is nil.
var ErrNil = errors.New("nil returned")

// Int is a helper that converts a command data to an integer. If err is not
// equal to nil, then Int returns 0, err. Otherwise, Int converts the
// data to an int as follows:
//
//  Reply type    Result
//  integer       int(data), nil
//  bulk string   parsed data, nil
//  nil           0, ErrNil
//  other         0, error
func Int(data interface{}, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	switch data := data.(type) {
	case int64:
		x := int(data)
		if int64(x) != data {
			return 0, strconv.ErrRange
		}
		return x, nil
	case []byte:
		n, err := strconv.ParseInt(string(data), 10, 0)
		return int(n), err
	case nil:
		return 0, ErrNil
	case Error:
		return 0, data
	}
	return 0, fmt.Errorf("unexpected type for Int, got type %T", data)
}

// Int64 is a helper that converts a command data to 64 bit integer. If err is
// not equal to nil, then Int returns 0, err. Otherwise, Int64 converts the
// data to an int64 as follows:
//
//  Reply type    Result
//  integer       data, nil
//  bulk string   parsed data, nil
//  nil           0, ErrNil
//  other         0, error
func Int64(data interface{}, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	switch data := data.(type) {
	case float32:
		return int64(data), err
	case float64:
		return int64(data), err
	case int32:
		return int64(data), nil
	case int64:
		return data, nil
	case []byte:
		n, err := strconv.ParseInt(string(data), 10, 64)
		return n, err
	case nil:
		return 0, ErrNil
	case Error:
		return 0, data
	}
	return 0, fmt.Errorf("unexpected type for int64, got type %T", data)
}

// Float64 is a helper that converts a command data to 64 bit float. If err is
// not equal to nil, then Float64 returns 0, err. Otherwise, Float64 converts
// the data to an int as follows:
//
//  Reply type    Result
//  bulk string   parsed data, nil
//  nil           0, ErrNil
//  other         0, error
func Float64(data interface{}, err error) (float64, error) {
	if err != nil {
		return 0, err
	}
	switch data := data.(type) {
	case float32:
		return float64(data), err
	case float64:
		return float64(data), err
	case int32:
		return float64(data), nil
	case int64:
		return float64(data), nil
	case []byte:
		n, err := strconv.ParseFloat(string(data), 64)
		return n, err
	case string:
		n, err := strconv.ParseFloat(data, 64)
		return n, err
	case nil:
		return 0, ErrNil
	case Error:
		return 0, data
	}
	return 0, fmt.Errorf("unexpected type for float64, got type %T", data)
}

// String is a helper that converts a command data to a string. If err is not
// equal to nil, then String returns "", err. Otherwise String converts the
// data to a string as follows:
//
//  Reply type      Result
//  bulk string     string(data), nil
//  simple string   data, nil
//  nil             "",  ErrNil
//  other           "",  error
func String(data interface{}, err error) (string, error) {
	if err != nil {
		return "", err
	}
	switch data := data.(type) {
	case []byte:
		return string(data), nil
	case string:
		return data, nil
	case nil:
		return "", ErrNil
	case Error:
		return "", data
	}
	return "", fmt.Errorf("unexpected type for String, got type %T", data)
}

// Bool is a helper that converts a command data to a boolean. If err is not
// equal to nil, then Bool returns false, err. Otherwise Bool converts the
// data to boolean as follows:
//
//  Reply type      Result
//  integer         value != 0, nil
//  bulk string     strconv.ParseBool(data)
//  nil             false, ErrNil
//  other           false, error
func Bool(data interface{}, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	switch data := data.(type) {
	case int64:
		return data != 0, nil
	case []byte:
		return strconv.ParseBool(string(data))
	case nil:
		return false, ErrNil
	case Error:
		return false, data
	}
	return false, fmt.Errorf("unexpected type for Bool, got type %T", data)
}

// Strings is a helper that converts an array command data to a []string. If
// err is not equal to nil, then Strings returns nil, err. Nil array items are
// converted to "" in the output slice. Strings returns an error if an array
// item is not a bulk string or nil.
func Strings(data interface{}, err error) ([]string, error) {
	var result []string
	err = sliceHelper(data, err, "Strings", func(n int) { result = make([]string, n) }, func(i int, v interface{}) error {
		switch v := v.(type) {
		case string:
			result[i] = v
			return nil
		case []byte:
			result[i] = string(v)
			return nil
		default:
			return fmt.Errorf("unexpected element type for Strings, got type %T", v)
		}
	})
	return result, err
}

func sliceHelper(data interface{}, err error, name string, makeSlice func(int), assign func(int, interface{}) error) error {
	if err != nil {
		return err
	}
	switch data := data.(type) {
	case []interface{}:
		makeSlice(len(data))
		for i := range data {
			if data[i] == nil {
				continue
			}
			if err := assign(i, data[i]); err != nil {
				return err
			}
		}
		return nil
	case nil:
		return ErrNil
	case Error:
		return data
	}
	return fmt.Errorf("unexpected type for %s, got type %T", name, data)
}

// Values is a helper that converts an array command reply to a []interface{}.
// If err is not equal to nil, then Values returns nil, err. Otherwise, Values
// converts the reply as follows:
//
//  Reply type      Result
//  array           reply, nil
//  nil             nil, ErrNil
//  other           nil, error
func Values(reply interface{}, err error) ([]interface{}, error) {
	if err != nil {
		return nil, err
	}
	switch reply := reply.(type) {
	case []interface{}:
		return reply, nil
	case nil:
		return nil, ErrNil
	case Error:
		return nil, reply
	}
	return nil, fmt.Errorf("redigo: unexpected type for Values, got type %T", reply)
}

func errNegativeInt(v int64) error {
	return fmt.Errorf("redigo: unexpected negative value %v for Uint64", v)
}

// Uint64 is a helper that converts a command reply to 64 bit unsigned integer.
// If err is not equal to nil, then Uint64 returns 0, err. Otherwise, Uint64 converts the
// reply to an uint64 as follows:
//
//  Reply type    Result
//  +integer      reply, nil
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
func Uint64(reply interface{}, err error) (uint64, error) {
	if err != nil {
		return 0, err
	}
	switch reply := reply.(type) {
	case int64:
		if reply < 0 {
			return 0, errNegativeInt(reply)
		}
		return uint64(reply), nil
	case []byte:
		n, err := strconv.ParseUint(string(reply), 10, 64)
		return n, err
	case string:
		n, err := strconv.ParseUint(reply, 10, 64)
		return n, err
	case nil:
		return 0, ErrNil
	case Error:
		return 0, reply
	}
	return 0, fmt.Errorf("redigo: unexpected type for Uint64, got type %T", reply)
}

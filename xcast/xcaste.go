package xcast

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"unicode"

	"github.com/spf13/cast"
)

// DeepCopyInto deep copy value into ptr. ptr must be a pointer.
func DeepCopyIntoE[T any](ptr *T, value any) error {
	if ptr == nil {
		return errors.New("ptr is nil")
	}
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(jsonValue, ptr); err != nil {
		return err
	}
	return nil
}

func DeepCopyE[T any](value any) (T, error) {
	var ptr T
	err := DeepCopyIntoE(&ptr, value)
	return ptr, err
}

func ToAnyE[T any](value any) (T, error) {
	var ptr T
	err := DeepCopyIntoE(&ptr, value)
	return ptr, err
}

func CamelToSnakeE(str string) (string, error) {
	reg, err := regexp.Compile("([a-z0-9])([A-Z])")
	if err != nil {
		return "", err
	}
	snake := reg.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake), nil
}

func SnakeToCamelE(str string) (string, error) {
	words := strings.Split(str, "_")
	for i := range words {
		if len(words[i]) > 0 {
			words[i] = string(unicode.ToUpper(rune(words[i][0]))) + words[i][1:]
		}
	}
	return strings.Join(words, ""), nil
}

func ToBoolE(i interface{}) (bool, error) {
	return cast.ToBoolE(i)
}

// ToFloat64 casts an interface to a float64 type.
func ToFloat64E(i interface{}) (float64, error) {
	return cast.ToFloat64E(i)
}

// ToFloat32 casts an interface to a float32 type.
func ToFloat32E(i interface{}) (float32, error) {
	return cast.ToFloat32E(i)
}

// ToInt64 casts an interface to an int64 type.
func ToInt64E(i interface{}) (int64, error) {
	return cast.ToInt64E(i)
}

// ToInt32 casts an interface to an int32 type.
func ToInt32E(i interface{}) (int32, error) {
	return cast.ToInt32E(i)
}

// ToInt16 casts an interface to an int16 type.
func ToInt16E(i interface{}) (int16, error) {
	return cast.ToInt16E(i)
}

// ToInt8 casts an interface to an int8 type.
func ToInt8E(i interface{}) (int8, error) {
	return cast.ToInt8E(i)
}

// ToInt casts an interface to an int type.
func ToIntE(i interface{}) (int, error) {
	return cast.ToIntE(i)
}

// ToUint64 casts an interface to a uint64 type.
func ToUint64E(i interface{}) (uint64, error) {
	return cast.ToUint64E(i)
}

// ToUint32 casts an interface to a uint32 type.
func ToUint32E(i interface{}) (uint32, error) {
	return cast.ToUint32E(i)
}

// ToUint16 casts an interface to a uint16 type.
func ToUint16E(i interface{}) (uint16, error) {
	return cast.ToUint16E(i)
}

// ToUint8 casts an interface to a uint8 type.
func ToUint8E(i interface{}) (uint8, error) {
	return cast.ToUint8E(i)
}

// ToUint casts an interface to a uint type.
func ToUintE(i interface{}) (uint, error) {
	return cast.ToUintE(i)
}

// ToString casts an interface to a string type.
func ToStringE(i interface{}) (string, error) {
	return cast.ToStringE(i)
}

package mtg

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"fmt"
	"reflect"
)

func Encode(values ...interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}

	for _, v := range values {
		b, err := encode(v)
		if err != nil {
			return nil, err
		}

		n := len(b)
		if n > 255 {
			return nil, fmt.Errorf("mtg: no enough bytes to encode %v", v)
		}

		if err := buf.WriteByte(byte(len(b))); err != nil {
			return nil, err
		}

		if _, err := buf.Write(b); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func encode(value interface{}) ([]byte, error) {
	if u, ok := value.(encoding.BinaryMarshaler); ok {
		return u.MarshalBinary()
	}

	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		b := make([]byte, binary.MaxVarintLen64)
		n := binary.PutVarint(b, v.Int())
		return b[:n], nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		b := make([]byte, binary.MaxVarintLen64)
		n := binary.PutUvarint(b, v.Uint())
		return b[:n], nil
	case reflect.String:
		return []byte(v.String()), nil
	}

	return nil, fmt.Errorf("mtg: cannot encode %T %v", value, value)
}

func EncodeStruct(s interface{}) ([]byte, error) {
	if s == nil {
		return nil, nil
	}

	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("function only accepts structs; got %s", val.Kind())
	}

	n := val.NumField()
	values := make([]interface{}, 0, n)

	for i := 0; i < n; i++ {
		valueField := val.Field(i)
		if valueField.CanInterface() {
			values = append(values, valueField.Interface())
		}
	}

	return Encode(values...)
}

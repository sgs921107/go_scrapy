package gcommon

import (
	"bytes"
	"errors"
	"io"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// ErrorStructToMapSA item转map失败
var ErrorStructToMapSA = errors.New("ValueError: item must be a struct")

// MapSA 别名
type MapSA = map[string]interface{}

// StructToMapSA item to map
func StructToMapSA(item interface{}) (MapSA, error) {
	t := reflect.TypeOf(item)
	if t.Kind() != reflect.Struct {
		return nil, ErrorStructToMapSA
	}
	data := make(MapSA)
	d := reflect.ValueOf(item)
	for i := 0; i < d.NumField(); i++ {
		data[t.Field(i).Name] = d.Field(i).Interface()
	}
	return data, nil
}

/*
MapToBytes convert type map to []byte
*/
func MapToBytes(data map[string]string) []byte {
	reader := MapToReader(data)
	return ReaderToBytes(reader)
}

/*
MapToReader convert type map to io.reader
*/
func MapToReader(data map[string]string) io.Reader {
	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}
	return strings.NewReader(form.Encode())
}

/*
ReaderToBytes convert type io.reader to []byte
*/
func ReaderToBytes(reader io.Reader) []byte {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(reader)
	return buffer.Bytes()
}

/*
ReaderToString convert type io.reader to string
*/
func ReaderToString(reqBody io.Reader) string {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(reqBody)
	return buffer.String()
}

/*
DurationToIntSecond 将time.Duration类型转换为值为秒数的int类型
*/
func DurationToIntSecond(duration time.Duration) int {
	return int(duration) / 1e9
}

// MapSAToSS map[string]interface{} to map[string]string
func MapSAToSS(data map[string]interface{}) (map[string]string, bool) {
	var newData = make(map[string]string, len(data))
	for k, v := range data {
		val, ok := v.(string)
		if !ok {
			return nil, false
		}
		newData[k] = val
	}
	return newData, true
}

// MapSSToSA map[string]string to map[string]interface{}
func MapSSToSA(data map[string]string) map[string]interface{} {
	var newData = make(map[string]interface{}, len(data))
	for k, v := range data {
		newData[k] = v
	}
	return newData
}

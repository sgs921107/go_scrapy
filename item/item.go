/*************************************************************************
	> File Name: item.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月08日 星期二 14时55分41秒
 ************************************************************************/

/*
一个实现解析时item结构的列子
这里没有任何实际意义
*/

package item

import (
	"encoding/json"
	"errors"
	"reflect"
)

// ErrorItemToMap item转map失败
var ErrorItemToMap = errors.New("ItemToMapError: item must be a struct")

// Map 定义ItemMap类型
type Map map[string]interface{}

// ToMap item to map
func ToMap(item interface{}) (Map, error) {
	t := reflect.TypeOf(item)
	if t.Kind() != reflect.Struct {
		return nil, ErrorItemToMap
	}
	data := make(Map)
	d := reflect.ValueOf(item)
	for i := 0; i < d.NumField(); i++ {
		data[t.Field(i).Name] = d.Field(i)
	}
	return data, nil
}

// ToJSON item to json
func ToJSON(item interface{}) ([]byte, error) {
	return json.Marshal(&item)
}

type itemDemo struct {
	Name string
	Age  int
}

// ToJson 转json
func (i itemDemo) ToJSON() ([]byte, error) {
	return ToJSON(i)
}

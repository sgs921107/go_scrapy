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

var ItemToMapError = errors.New("ItemToMapError: item must be a struct")

type ItemMap map[string]interface{}

func ItemToMap(item interface{}) (ItemMap, error) {
	t := reflect.TypeOf(item)
	if t.Kind() != reflect.Struct {
		return nil, ItemToMapError
	}
	data := make(ItemMap)
	d := reflect.ValueOf(item)
	for i := 0; i < d.NumField(); i++ {
		data[t.Field(i).Name] = d.Field(i)
	}
	return data, nil
}

func ItemToJson(item interface{}) ([]byte, error) {
	return json.Marshal(&item)
}

type itemDemo struct {
	Name string
	Age  int
}

func (i itemDemo) ToJson() ([]byte, error) {
	return ItemToJson(i)
}

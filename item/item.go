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
	"github.com/sgs921107/gcommon"
)

type itemDemo struct {
	Name string
	Age  int
}

// ToJson 转json
func (i *itemDemo) ToJSON() ([]byte, error) {
	return json.Marshal(i)
}

// ToMapSA to msp[string]interface{}
func (i *itemDemo) ToMapSA() (gcommon.MapSA, error) {
	return gcommon.StructToMapSA(*i)
}

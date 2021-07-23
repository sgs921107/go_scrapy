/*
 * @Author: xiangcai
 * @Date: 2021-07-23 17:16:10
 * @LastEditors: xiangcai
 * @LastEditTime: 2021-07-23 17:43:49
 * @Description: file content
 */

package gspider

type Extension interface{
	Run(*BaseSpider)
}

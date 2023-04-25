package struct2Map

import (
	"fmt"
	"github.com/wxnacy/wgo/arrays"
	"reflect"
	"strings"
)

var emptySlice []struct{}

const (
	TagOfIgnore    = "-"
	TagOfOmitempty = "omitempty"
)

func StructToMap(iface interface{}, tag string) (map[string]interface{}, error) {
	val := reflect.ValueOf(iface)
	//(1) iface合法性校验
	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			return nil, fmt.Errorf("nil pointer")
		}
		val = val.Elem()
	case reflect.Struct:
	default:
		return nil, fmt.Errorf("is not struct")
	}
	//(2) 遍历iface各个字段
	var res = make(map[string]interface{}, val.NumField())
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		//(3) 获取对应字段的tag值
		fieldType := typ.Field(i)
		tagVal, ok := fieldType.Tag.Lookup(tag)
		if !ok {
			continue
		}
		tagValList := strings.Split(tagVal, ",")
		var hasIgnoreTag = arrays.ContainsString(tagValList, TagOfIgnore) != -1
		var hasOmitemptyTag = arrays.ContainsString(tagValList, TagOfOmitempty) != -1
		if hasIgnoreTag {
			continue
		}
		fieldVal := val.Field(i)
		// 字段合法性校验
		if fieldVal.Kind() == reflect.Ptr && fieldVal.IsNil() {
			//需要判断其是否存在omitempty
			if !hasOmitemptyTag {
				res[tagVal] = nil
				continue
			}
		}
		if fieldVal.Kind() == reflect.Ptr {
			fieldVal = fieldVal.Elem()
		}
		//(4) 判断对应字段的类型
		switch fieldVal.Kind() {
		case reflect.Struct:
			//如果是struct 递归调用
			m, err := StructToMap(fieldVal.Interface(), tag)
			if err != nil {
				return nil, err
			}
			res[tagValList[0]] = m
		case reflect.Slice, reflect.Array:
			//如果切片长度为0 那么返回一个空的切片
			if fieldVal.Len() == 0 {
				res[tagValList[0]] = emptySlice
			}
			//如果切片元素非struct的话
			if fieldType.Type.Elem().Kind() != reflect.Struct {
				res[tagValList[0]] = fieldVal.Interface()
			} else {
				list := make([]interface{}, 0, fieldVal.Len())
				for i := 0; i < fieldVal.Len(); i++ {
					entry := fieldVal.Index(i)
					if entry.Kind() == reflect.Ptr && entry.IsNil() {
						list = append(list, nil)
						continue
					}
					if entry.Kind() == reflect.Ptr {
						entry = entry.Elem()
					}
					m, err := StructToMap(entry.Interface(), tag)
					if err != nil {
						return nil, err
					}
					list = append(list, m)
				}
				res[tagValList[0]] = list
			}
		default:
			if !fieldVal.IsValid() {
				if !hasOmitemptyTag {
					res[tagValList[0]] = nil
				}
			} else if !fieldVal.IsZero() || !hasOmitemptyTag {
				res[tagValList[0]] = fieldVal.Interface()
			}
		}
	}
	return res, nil
}

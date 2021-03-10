/*
 @Author : wangdy
 @Time   : 2021/2/5 4:39 下午
 @Desc   :
*/

package engine_func

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

func Len(i interface{}) int {
	v := reflect.ValueOf(i)
	//fmt.Println("name:", v.Type())
	switch v.Kind() {
	case reflect.Array:
	case reflect.Slice:
	case reflect.Map:
	default:
		return -1
	}
	fmt.Println("len:", v.Len())
	return v.Len()
}

func InArr(meta interface{}, e interface{}) bool {
	items := reflect.ValueOf(meta)
	if items.Len() == 0 || items.Len() > 1000 {
		return false
	}
	switch items.Kind() {
	case reflect.Slice:
		for i := 0; i < items.Len(); i++ {
			item := items.Index(i)
			//fmt.Println("item:", item, item.Type())
			if item.Interface() == e {
				return true
			}
		}
	default:
		return false
	}
	return false
}

func InMap(meta interface{}, k, v interface{}) bool {
	now := time.Now()
	defer func() {
		fmt.Println("timeSub:", time.Now().Sub(now))
	}()
	items := reflect.ValueOf(meta)
	if items.Len() == 0 || items.Len() > 1000 {
		return false
	}
	switch items.Kind() {
	case reflect.Map:
		iter := items.MapRange()
		for iter.Next() {
			//fmt.Println("value:", iter.Value(), v)
			if iter.Key().Interface() == k && iter.Value().Interface() == v {
				return true
			}
		}
	}
	return false
}

func Int64(i interface{}) (int64, error) {
	fun := "Int64 "
	iv := reflect.ValueOf(i)
	if !iv.IsValid() {
		return 0, errors.New("无效的输入")
	}
	fmt.Println(fun, iv.Type(), iv.Interface())
	k := iv.Kind()
	switch k {
	case reflect.Uint:
		return int64(iv.Uint()), nil
	case reflect.Uint8:
		return int64(iv.Uint()), nil
	case reflect.Uint16:
		return int64(iv.Uint()), nil
	case reflect.Uint32:
		return int64(iv.Uint()), nil
	case reflect.Uint64:
		return int64(iv.Uint()), nil
	case reflect.Uintptr:
		return int64(iv.Uint()), nil
	case reflect.Int:
		return iv.Int(), nil
	case reflect.Int8:
		return iv.Int(), nil
	case reflect.Int16:
		return iv.Int(), nil
	case reflect.Int32:
		return iv.Int(), nil
	case reflect.Int64:
		return iv.Int(), nil
	default:
		return 0, errors.New("未知的类型，不是数字")
	}
}

func ReTry(num int, sleep int, fn interface{}, params ...interface{}) int64 {
	fun := "ReTry "
	fnt := reflect.TypeOf(fn)
	fmt.Println(fnt.Kind())
	fnv := reflect.ValueOf(fn)
	if fnt.Kind() != reflect.Func {
		fmt.Println(errors.New("fun不是函数"))
		return 0
	}
	paramsValues := make([]reflect.Value, len(params))
	for i, item := range params {
		paramsValues[i] = reflect.ValueOf(item)
	}
	fmt.Println(fun, "len:", len(paramsValues))
	var counter int64
	for i := 0; i < num; i++ {
		for _, vv := range fnv.Call(paramsValues) {
			counter++
			fmt.Println(vv)
			//fmt.Println(fun, vv.Kind(), vv.Interface())
		}
		if sleep > 0 {
			time.Sleep(time.Duration(sleep) * time.Millisecond)
		}
	}
	return counter
}

/*
 @Author : wangdy
 @Time   : 2021/2/2 9:51 下午
 @Desc   :
*/
package main

import (
	"fmt"
	"github.com/wangdyqxx/common/engine_func"
	"github.com/wangdyqxx/common/example/rule"
	"strings"
)

func main() {
	//arr := int64(1)
	//err := base.ReTry(4, 1, base.Int64, arr)
	//fmt.Println(err)
	//return

	strs := rule.LoadFile("./rule.txt")
	rules := strings.Join(strs, "\r\n")
	//fmt.Println("rules:\r", rules)
	ser, err := rule.EngineInit(rules, nil)
	if err != nil {
		fmt.Println(fmt.Sprintf("EngineInit err:%+v", err))
		return
	}
	//基于需要注入接口或数据,data这里最好仅注入与本次请求相关的结构体或数据，便于状态管理
	arr := []int64{1, 2, 3}
	req := &Request{Rid: 123}
	res := &Response{}
	ruleData["req"] = req
	ruleData["res"] = res
	ruleData["room"] = &Room{AccountIds: arr}
	ruleData["arr"] = arr
	ruleData["mm"] = map[int64]int{1: 11, 2: 12, 3: 13, 4: 14}
	err = ser.Service(ruleData, ruleNames...)
	if err != nil {
		fmt.Println(fmt.Sprintf("service err:%+v", err))
		return
	}
	println("res = ", res.At, res.Num)
}

var (
	ruleNames = []string{"1"}
	ruleData  = map[string]interface{}{}
)

func init() {
	rule.RegBaseApi("joinEvent", engine_func.JoinEvent)
}

//request
type Request struct {
	Rid int64
	//other params
}

//resp
type Response struct {
	At  int64
	Num int64
	//other params
}

//特定的场景服务
type Room struct {
	AccountIds []int64
}

func (r *Room) GetAttention( /*params*/ ) int64 {
	// logic
	return 100
}

func (r *Room) GetNum( /*params*/ ) int64 {
	//logic
	return 111
}

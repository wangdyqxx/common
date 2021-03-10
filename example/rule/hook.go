/*
 @Author : wangdy
 @Time   : 2021/2/2 8:18 下午
 @Desc   :
*/

package rule

import (
	"errors"
	"fmt"
	"github.com/bilibili/gengine/engine"
	"github.com/wangdyqxx/common/engine_func"
)

//业务接口
type MyService struct {
	//gengine pool
	Pool *engine.GenginePool
	//other params
}

//初始化业务服务
//apiOuter这里最好仅注入一些无状态函数，方便应用中的状态管理
func NewMyService(poolMinLen, poolMaxLen int64, em int, rulesStr string, apiOuter map[string]interface{}) (*MyService, error) {
	pool, err := engine.NewGenginePool(poolMinLen, poolMaxLen, em, rulesStr, apiOuter)
	if err != nil {
		fmt.Println(fmt.Sprintf("初始化gengine失败，err:%+v", err))
		return nil, err
	}
	myService := &MyService{Pool: pool}
	return myService, nil
}

var (
	apis = map[string]interface{}{
		"println": fmt.Println,
		"len":     engine_func.Len,
		"inArr":   engine_func.InArr,
		"inMap":   engine_func.InMap,
		"int64":   engine_func.Int64,
		"reTry":   engine_func.ReTry,
	}

	baseApis = map[string]interface{}{
		"println": fmt.Println,
		"len":     engine_func.Len,
		"inArr":   engine_func.InArr,
		"inMap":   engine_func.InMap,
		"int64":   engine_func.Int64,
		"reTry":   engine_func.ReTry,
	}
)

func RegBaseApi(str string, fun interface{}) bool {
	_, ok := apis[str]
	if ok {
		return false
	}
	apis[str] = fun
	return true
}

//初始化，注入api，确保注入的API属于并发安全
func EngineInit(rules string, regApis map[string]interface{}) (ser *MyService, err error) {
	for apiName, function := range regApis {
		if _, ok := baseApis[apiName]; ok {
			return nil, errors.New(fmt.Sprint("已存在的基础api：", apiName))
		}
		_, ok := apis[apiName]
		if ok {
			return nil, errors.New("已存在的api")
		}
		apis[apiName] = function
	}
	return NewMyService(10, 20, 1, rules, apis)
}

//service
func (ms *MyService) Service(data map[string]interface{}, ruleNames ...string) error {
	if len(ruleNames) == 0 {
		return errors.New("缺少运行规则名称")
	}
	err, returnResultMap := ms.Pool.ExecuteSelectedRules(data, ruleNames)
	if err != nil {
		fmt.Println(fmt.Sprintf("pool execute rules error: %+v", err))
		return err
	}
	fmt.Println(fmt.Sprintf("returnResultMap: %+v", returnResultMap))
	return nil
}

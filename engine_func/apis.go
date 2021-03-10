/*
 @Author : wangdy
 @Time   : 2021/2/4 4:46 下午
 @Desc   :
*/
package engine_func

import (
	"fmt"
)

func JoinEvent(uid int64) bool {
	fmt.Println("joinEvent uid:", uid)
	if uid%2 == 0 {
		return true
	}
	return false
}


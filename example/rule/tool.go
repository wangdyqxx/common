/*
 @Author : wangdy
 @Time   : 2021/2/4 4:57 下午
 @Desc   :
*/
package rule

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var step = 10

//加载文件，注意是否全部读取 fileDetails key为行数，value
func LoadFile(path string) []string {
	fun := "LoadFile ->"
	var fileDetails []string
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(fmt.Sprintf("%s err::%v", fun, err))
		return fileDetails
	}
	defer f.Close()
	br := bufio.NewReader(f)
	for {
		data, _, err := br.ReadLine()
		if len(fileDetails)%step == 0 {
			//fmt.Println(fmt.Sprintf("%s data: %s , idx: %v ", fun, data, len(fileDetails)+1))
		}
		if err == io.EOF {
			//fmt.Println(fmt.Sprintf("loadClose:%v",path))
			break
		} else {
			if len(data) == 0 {
				continue
			}
			fileDetails = append(fileDetails, string(data))
		}
	}
	fmt.Println(fmt.Sprintf("%s path:%v, fileDetails:%+v", fun, path, len(fileDetails)))
	return fileDetails
}
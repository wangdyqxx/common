package util

import (
	"crypto/sha1"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/wangdyqxx/common/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

var token = "wangdy"
var appId = "wxc8a78f278759a02b"
var appSecret = "94c0eb31e2a740ad643c59bb2b8abcc5"
var url = "https://api.weixin.qq.com/cgi-bin/token"
var Json = jsoniter.ConfigCompatibleWithStandardLibrary

//校验签名
func MakeSig(timestamp, nonce string) string {
	//1. 将 plat_token、timestamp、nonce三个参数进行字典序排序
	sl := []string{token, timestamp, nonce}
	sort.Strings(sl)
	//2. 将三个参数字符串拼接成一个字符串进行sha1加密
	s := sha1.New()
	_, err := io.WriteString(s, strings.Join(sl, ""))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", s.Sum(nil))
}

var requestLine = strings.Join([]string{
	url, "?grant_type=client_credential&appid=",
	appId, "&secret=", appSecret}, "")

//获取微信accessToken
func GetAccessToken() (string, error) {
	resp, err := http.Get(requestLine)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Info("发送get请求获取 atoken 错误", err)
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Info("发送get请求获取 atoken 读取返回body错误", err)
		return "", err
	}
	s := string(body)
	log.Info("body:", s)
	return s, nil
}

//{"access_token":"20_fuSdVSLREwK_Q8PJ6Bt0tYhqdeL-QnSoDU3vjOgey_MbbKm1H8tDGLCCMYqLo6VmdFXX1WclJe3laNrYM_crCa9NFPM4VLshTFkDST3w2IxQTpNZQWU_1BJB
//rRBriZtGq8ZE0_bS7VH7bDOERRAcAJAKKT","expires_in":7200}

//判断路径是否存在,返回true为存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//解析二维字符串
func ParseXmlStr2(str string) [][]string {
	str1 := strings.Split(str, "\n")
	var excelData [][]string
	for _, v := range str1 {
		//todo 去除\r、空格，以+为分隔符
		v = strings.Replace(v, "\r", "", -1)
		v = strings.Replace(v, " ", "", -1)
		str2 := strings.Split(v, "+")
		excelData = append(excelData, str2)
	}
	return excelData
}

//解析一维字符串
func ParseXmlStr1(str string) []string {
	str = strings.Replace(str, " ", "", -1)
	//todo 以十、＋、₊、+、➕为分隔符
	rowData := strings.Split(str, "+")
	return rowData
}

func Int64ToStr(num int64) string {
	return strconv.FormatInt(num, 10)
}

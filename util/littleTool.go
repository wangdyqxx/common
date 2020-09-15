package util

import (
	"context"
	"fmt"
	"github.com/silenceper/wechat/material"
	wechatCtx "github.com/silenceper/wechat/context"
	"github.com/wangdyqxx/common/log"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"os/exec"
)

var UserTool map[string]string
var ToolDesc = map[string]string{
	"TinyProto": "无损压缩png图片",
}
var Tools = map[string]func(str ...string) bool{
	"TinyProto": TinyProto,
}

func InTool(str string) string {
	_, ok := ToolDesc[str]
	if ok {
		return str
	}
	return ""
}

func UseTool(userName string) string {
	v, ok := UserTool[userName]
	if ok {
		return v
	}
	return ""
}

func init() {
	UserTool = make(map[string]string, 8)
}

//func ToolJsonTrans(str string) string {
//	buf := bytes.NewBufferString(str)
//	de := json.NewDecoder(buf)
//	de.Decode()
//}

const (
	tinyProto = "pngquant"
)

func TinyProto(strs ...string) bool {
	if len(strs) < 2 {
		return false
	}
	log.Info("strs:", strs)
	oldPath, extName := strs[0], strs[1]
	ctx := context.TODO()
	log.Infoln(ctx, "TinyProto->")
	params := []string{oldPath, "--ext", extName}
	cmd := exec.Command(tinyProto, params...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Errorf(ctx, "err1:", err)
		return false
	}
	defer stdout.Close()
	if err = cmd.Start(); err != nil {
		log.Errorf(ctx, "err3:", err)
		return false
	}
	content, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Errorf(ctx, "err:", err)
		return false
	}
	log.Infoln(ctx, "content:", string(content))
	return true
}

func ConvertPng(old, new string) bool {
	f1, err := os.Open(old)
	if err != nil {
		fmt.Println("err:", err)
		return false
	}
	defer f1.Close()
	im1, err := jpeg.Decode(f1)
	if err != nil {
		fmt.Println("err:", err)
		return false
	}
	f2, err := os.Create(new)
	if err != nil {
		fmt.Println("err:", err)
		return false
	}
	defer func() {
		err := f2.Close()
		fmt.Println("err2:",err)
	}()
	err = png.Encode(f2, im1)
	if err != nil {
		fmt.Println("err:", err)
		return false
	}
	fmt.Println("path:",f2.Name())
	return true
}

func NewMaterial(wxCtx *wechatCtx.Context, path string) {
	//上传素材
	newMaterial := material.NewMaterial(wxCtx)
	mediaID, url, err := newMaterial.AddMaterial(material.MediaTypeImage, path)
	log.Infof(context.TODO(), "newMaterial id:%v, url:%v, err:%v", mediaID, url, err)
	return
}
package cmd

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/wangdyqxx/util/slog"
)

var (
	command = "/bin/bash"
	//params = []string{"-c", "sh ./test.sh"}
)

func ExecCommand(ctx context.Context, params ...string) (bool, error) {
	fun := "ExecCommand ->"
	args := []string{"-c"}
	args = append(args, params...)
	cmd := exec.Command(command, args...)
	//显示运行的命令
	slog.Infof(ctx, "%s 执行命令: %s\n", fun, strings.Join(cmd.Args, " "))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		slog.Errorf(ctx, "%s err:%v, %v", fun, os.Stderr, err)
		return false, err
	}
	err = cmd.Start() // Start开始执行c包含的命令，但并不会等待该命令完成即返回。Wait方法会返回命令的返回状态码并在命令返回后释放相关的资源。
	if err != nil {
		slog.Errorf(ctx, "%s Start err:%v", fun, err)
		return false, err
	}
	reader := bufio.NewReader(stdout)
	//循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		slog.Infof(ctx, "%s line:%v", fun, line)
	}
	err = cmd.Wait()
	if err != nil {
		slog.Errorf(ctx, "%s Wait err:%v", fun, err)
		return false, err
	}
	return true, nil
}

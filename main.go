package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"github.com/ci-plugins/golang-plugin-sdk/api"
	"github.com/ci-plugins/golang-plugin-sdk/log"
)

type greetingParam struct {
	UserName string `json:"userName"`
	Greeting string `json:"greeting"`
}

func (a *greetingParam) String() string {
	return fmt.Sprintf("userName: %v, greeting: %v", a.UserName, a.Greeting)
}

func main() {
	runtime.GOMAXPROCS(4)
	log.Info("atom-demo-glang starts")
	defer func() {
		if err := recover(); err != nil {
			log.Error("panic: ", err)
			api.FinishBuild(api.StatusError, "panic occurs")
		}
	}()

	helloBuild()
}

func helloBuild() {
	// 获取单个输入参数
	userName := api.GetInputParam("userName")
	log.Info("userName: ", userName)

	// 打屏
	log.Info("\nBuildInfo:")
	log.Info("Project Name:     ", api.GetProjectDisplayName())
	log.Info("Pipeline Id:      ", api.GetPipelineId())
	log.Info("Pipeline Name:    ", api.GetPipelineName())
	log.Info("Pipeline Version: ", api.GetPipelineVersion())
	log.Info("Build Id:         ", api.GetPipelineBuildId())
	log.Info("Build Num:        ", api.GetPipelineBuildNumber())
	log.Info("Start Type:       ", api.GetPipelineStartType())
	log.Info("Start UserId:     ", api.GetPipelineStartUserId())
	log.Info("Start UserName:   ", api.GetPipelineStartUserName())
	log.Info("Start Time:       ", api.GetPipelineStartTimeMills())
	log.Info("Workspace:        ", api.GetWorkspace())

	// 输入参数解析到对象
	paramData := new(greetingParam)
	api.LoadInputParam(paramData)
	log.Info(fmt.Sprintf("\n%v，%v\n", paramData.Greeting, paramData.UserName))

	// 业务逻辑
	log.Info("start build")
	build()
	time.Sleep(2 * time.Second)

	// 输出
	// 字符串输出
	strData := api.NewStringData("test")
	api.AddOutputData("strData_01", strData)

	// 文件归档输出
	artifactData := api.NewArtifactData()
	artifactData.AddArtifact("result.dat")
	api.AddOutputData("artifactData_02", artifactData)

	// 报告输出
	reportData := api.NewReportData("label_01", api.GetWorkspace()+"/report", "report.htm")
	api.AddOutputData("report_01", reportData)

	api.WriteOutput()
	log.Info("build done")
}

func build() {
	log.Info("write result.dat")
	ioutil.WriteFile(api.GetWorkspace()+"/result.dat", []byte("content"), 0644)
	log.Info("write report.htm")
	os.Mkdir(api.GetWorkspace()+"/report", os.ModePerm)
	ioutil.WriteFile(api.GetWorkspace()+"/report/report.htm", []byte("<html><head><title>Report</title></head><body><H1>This is a Report</H1></body></html>"), 0644)
}

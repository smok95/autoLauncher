package main

import (
	"flag"
	"os"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	autolauncher "github.com/smok95/autoLauncher"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Auto Launcher")
	systray.SetTooltip("Auto Launcher")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	// Sets the icon of a menu item. Only available on Mac and Windows.
	mQuit.SetIcon(icon.Data)

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	// run autoLauncher
	opts := parseOptions()
	autolauncher.Run(opts)
}

func onExit() {
	// clean up here

}

// 실행 인자로부터 Options 구조체 값 파싱
func parseOptions() autolauncher.Options {
	opts := autolauncher.Options{}

	// 실행 인자 파싱
	processPath := flag.String("path", "", "프로세스 경로(프로세스명 포함된 전체경로)")
	checkStartTime := flag.String("start", "", "체크 시작 시간(형식: 15:04)")
	checkEndTime := flag.String("end", "", "체크 종료 시간(형식: 15:04)")
	checkInterval := flag.Int("interval", 0, "체크 주기(단위: 초)")
	flag.Parse()

	// 필수 입력값 확인
	if *processPath == "" || *checkStartTime == "" || *checkEndTime == "" || *checkInterval == 0 {
		flag.Usage()
		os.Exit(1)
	}

	opts.ProcessPath = *processPath
	opts.CheckStartTime = *checkStartTime
	opts.CheckEndTime = *checkEndTime
	opts.CheckInterval = *checkInterval

	return opts
}

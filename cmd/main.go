package main

import (
	"flag"
	"log"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	autolauncher "github.com/smok95/autoLauncher"
	"gopkg.in/ini.v1"
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
	if opts == nil {
		// config.ini에서 설정값 읽어오기
		opts = ReadConfig()
	}

	if opts == nil {
		log.Fatalf("설정값이 없습니다.")
		systray.Quit()
	}

	autolauncher.Run(*opts)
}

func onExit() {
	// clean up here

}

func ReadConfig() *autolauncher.Options {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("INI 파일 로드 실패: %v", err)
		return nil
	}

	opts := autolauncher.Options{}
	err = cfg.Section("Options").MapTo(&opts)
	if err != nil {
		log.Fatalf("INI 값 매핑 실패: %v", err)
		return nil
	}

	return &opts
}

// 실행 인자로부터 Options 구조체 값 파싱
func parseOptions() *autolauncher.Options {

	// 실행 인자 파싱
	processPath := flag.String("path", "", "프로세스 경로(프로세스명 포함된 전체경로)")
	checkStartTime := flag.String("start", "", "체크 시작 시간(형식: 15:04)")
	checkEndTime := flag.String("end", "", "체크 종료 시간(형식: 15:04)")
	checkInterval := flag.Int("interval", 0, "체크 주기(단위: 초)")
	flag.Parse()

	// 필수 입력값 확인
	if *processPath == "" || *checkStartTime == "" || *checkEndTime == "" ||
		*checkInterval == 0 {
		flag.Usage()
		return nil
	}

	return &autolauncher.Options{
		ProcessPath:    *processPath,
		CheckStartTime: *checkStartTime,
		CheckEndTime:   *checkEndTime,
		CheckInterval:  *checkInterval,
	}
}

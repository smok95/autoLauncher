package autolauncher

import (
	"bytes"
	"log"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/shirou/gopsutil/process"
)

type Options struct {
	ProcessPath    string // 프로세스 경로(프로세스명 포함된 전체경로)
	CheckStartTime string // 체크 시작 시간(형식: 15:04)
	CheckEndTime   string // 체크 종료 시간(형식: 15:04)
	CheckInterval  int    // 체크 주기(단위:초)
}

func Run(opts Options) {
	procName := filepath.Base(opts.ProcessPath)

	// 체크 시작 시간 파싱
	startTime, err := time.Parse("15:04", opts.CheckStartTime)
	if err != nil {
		log.Fatalf("체크 시작 시간 파싱 실패: %v", err)
	}

	// 체크 종료 시간 파싱
	endTime, err := time.Parse("15:04", opts.CheckEndTime)
	if err != nil {
		log.Fatalf("체크 종료 시간 파싱 실패: %v", err)
	}

	for {
		now := time.Now().Local()

		sTime := time.Date(now.Year(), now.Month(), now.Day(), startTime.Hour(), startTime.Minute(), 0, 0, time.Local)
		eTime := time.Date(now.Year(), now.Month(), now.Day(), endTime.Hour(), endTime.Minute(), 0, 0, time.Local)

		// 입력한 시간 범위에 체크
		if now.After(sTime) && now.Before(eTime) {
			isRunning, err := isProcessRunning(procName)
			if err != nil {
				log.Fatalf("프로세스 상태 확인 실패: %v", err)
			}

			if !isRunning {
				if err := startProcess(opts.ProcessPath); err != nil {
					log.Fatalf("프로세스 실행 실패: %v", err)
				}
			}
		}

		// 일정 시간 대기
		time.Sleep(time.Duration(opts.CheckInterval) * time.Second)
	}
}

// 프로세스 실행 여부 확인
func isProcessRunning(exeName string) (bool, error) {
	processes, err := process.Processes()
	if err != nil {
		return false, err
	}

	for _, proc := range processes {
		processName, err := proc.Name()
		if err != nil {
			log.Printf("프로세스 이름 조회 실패: %v", err)
			continue
		}

		if processName == exeName {
			return true, nil
		}
	}

	return false, nil
}

// 실행 중인 프로세스 목록에서 해당 프로세스 이름 또는 전체 경로를 검색
func containsProcess(processList []byte, processName string) bool {
	return bytes.Contains(processList, []byte(processName))
}

// 프로세스 실행
func startProcess(processPath string) error {
	cmd := exec.Command(processPath)
	err := cmd.Start()
	if err != nil {
		return err
	}

	log.Printf("프로세스 %s 시작됨.", processPath)
	return nil
}

package monitor

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 기능이 없어도 init 을 넣고, false 와 err 를 호출하도록

type C_monitor struct {
	C_monitor__db
	C_monitor__log

	s_monitor__url    string
	s_monitor__name   string
	s_monitor__status string

	n_monitor__rate int

	arrs_monitor__urls       []string
	arrs_monitor__status_grp []string
}

func (t *C_monitor) Init_check() (bool, error) {
	var err error
	if err != nil {
		fmt.Print(err)
		return false, err
	}
	return true, nil
}

// URL HTTP 상태 체크 실행
func (t *C_monitor) Monitor__checkUrl(_n_monitor__rate int) error {

	_, err := t.Init_check()
	if err != nil {
		return err
	}

	// DB에서 모니터링 대상 URL 호출
	target, err := t.Get__urls()
	if err != nil {
		return err
	}

	// 반복 시간 설정
	ticker := time.NewTicker(time.Second * time.Duration(_n_monitor__rate))

	for range ticker.C {

		// 모니터링 대상 URL string 으로 변경 후 http 상태 조회
		for _, url := range target {
			resp, err := http.Get(url)
			if err != nil || resp.StatusCode >= 400 {
				// http status 오류의 경우 DB status 값을 0(false)로 변경
				log.Println("URL :", url, ", STATUS : ERR ")
				err = t.Change_status__false(url)
				if err != nil {
					return err
				}
			} else {
				// http status 정상의 경우 DB status 값을 0(false)로 변경
				log.Println("URL :", url, ", STATUS :", resp.Status)
				err = t.Change_status__true(url)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// DB URL status 값 체크 및 알림 발송
func (t *C_monitor) Monitor__checkStatus(_n_monitor__rate int) error {

	_, err := t.Init_check()
	if err != nil {
		return err
	}

	// 반복시간 설정
	ticker := time.NewTicker(time.Second * time.Duration(_n_monitor__rate))

	for range ticker.C {

		// DB url, status 값 호출 및 string 변환
		url, status, err := t.Get__status()
		if err != nil {
			return err
		}

		for i, _status := range status {
			// status 값 체크하여 false의 경우 알림 발송
			if _status == "true" {
				fmt.Print()
			} else {
				// 에러 메시지 조합
				message := "SERVER Error Alert!\n" + "URL :" + url[i] + "\n" + "SERVER STATUS :" + _status
				err = t.Monitor__sendAlert(message)
				if err != nil {
					return err
				}
				log.Println(message)
			}
		}
	}
	return nil
}

// 메일 및 SMS 발송
func (t *C_monitor) Monitor__sendAlert(_s_message string) error {
	var err error

	_, err = t.Init_check()
	if err != nil {
		return err
	}

	c_sendmail := C_Sendmail{}

	// DB 연락처, 메일 데이터 쿼리하여 변수 저장
	mail, number, err := t.Get__contact_info()
	if err != nil {
		return err
	}

	// 메일 발송 함수 실행
	err = c_sendmail.Send_mail(_s_message, mail)
	if err != nil {
		return err
	}

	// 연락처 string 변환 후 SMS 발송
	for _, _number := range number {
		err = Send_sns(_s_message, _number)
		if err != nil {
			return err
		}
	}
	return nil
}
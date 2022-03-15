[필요 기능 정의]

1. 알림 발송 기능
2. 상태 체크 기능
3. DB 연동(target 정보)

 
[프로세스]

1. 타겟의 http 상태를 체크한다 (1분에 1회)
  1.1 http status code 200 ~ 299 (정상)
  1.2 그외 코드 및 접속 오류 (장애)

2. 타겟의 http 상태가
  2.1 정상의 경우, status = green
  2.2 장애의 경우, status = red

3.메인 모니터링
  3.1 호스트 별 상태를 체크한다. (100초 1회)
  3.2 상태가 red 인 호스트가 발견되면 aws sns 알림을 발송한다.
  3.3 발송 메세지의 red 호스트의 정보를 message 내용에 추가한다.


[함수 정의]
1. Send_sns
  1.1 Send
  1.2 Init
  1.3 New_C_sns

2. Run_CheckUrl
  2.1 checkUrl
  2.2 New_C_checkUrl

3. Main_CheckTarget (작성중)

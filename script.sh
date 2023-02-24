#!/bin/sh

userId=

password=

captiveReturnCode=`curl -s -I -m 10 -o /dev/null -s -w %{http_code} https://www.baidu.com/`
if [ "${captiveReturnCode}" = "200" ]; then
  echo 'online'
  exit 0
fi

loginURL='http://10.0.1.52:801/eportal/?c=ACSetting&a=Login&protocol=http:&hostname=10.0.1.52&iTermType=1&mac=00-00-00-00-00-00&ip=10.18.241.254&enAdvert=0&queryACIP=0&jsVersion=2.4.3&loginMethod=1'

auth=`curl -s -A "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.91 Safari/537.36" -e "${loginURL}" -X POST -d "DDDDD=%2C0%2C${userId}%40cccc&upass=${password}&R1=0&R2=0&R3=0&R6=0&para=00&0MKKey=123456" "${loginURL}"`
echo 'success!'
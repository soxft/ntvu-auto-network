## ntvu-auto-network

某学校的校园网自动登录脚本

## Usage Example

```shell
$ go build -o ntvu-auto-network main.go

$ ./ntvu-auto-network -u uname -p pwd -isp 1 -cron "0 */2 * * * *" 
```

## Usage

```shell
Usage of ntvu-auto-network:
  -c string
        配置文件路径 (default "config.yaml")
  -isp int
        1: 移动, 2: 电信, 3: 联通, 4: 校园网 (default 1)
  -u string
        学号: ex:20220000000
  -p string
        密码
  -cron cron表达式
        轮训间隔, 请使用cron表达式, ex: 0 */2 * * * *
        格式: 秒 分 时 日 月 周 (default "0 */2 * * * *")
  -typ string
        操作: login or logout (default "login")
```
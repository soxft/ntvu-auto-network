package main

import (
	"autoLoginNtvu/tool"
	"encoding/base64"
	"errors"
	"flag"
	"github.com/robfig/cron/v3"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	Gateway = "http://10.0.1.52:801/eportal/?c=ACSetting&a=Login&jsVersion=2.4.3"
	JsVer   = "2.4.3"
)

var ConfPath string
var Config tool.Conf
var ISPMap = map[int]string{
	1: ",0,{{UNAME}}@bbbb",
	2: ",0,{{UNAME}}@cccc",
	3: ",0,{{UNAME}}@aaaa",
	4: ",0,{{UNAME}}",
}
var logout bool

func main() {
	flag.IntVar(&Config.ISP, "isp", 1, "1: 移动, 2: 电信, 3: 联通, 4: 校园网")
	flag.StringVar(&Config.Username, "u", "", "学号: ex:20220000000")
	flag.StringVar(&Config.Password, "p", "", "密码")
	flag.StringVar(&Config.Cron, "cron", "0 */2 * * * *", "轮训间隔, 请使用`cron表达式`, ex: 0 */2 * * * *\r\n 格式: 秒 分 时 日 月 周")
	flag.StringVar(&ConfPath, "c", "config.yaml", "配置文件路径")
	flag.IntVar(&Config.RunTyp, "typ", 0, "0: 单次执行, 1: cron")
	flag.BoolVar(&logout, "logout", false, "是否为退出登录, true / false")
	flag.Parse()

	if Cfg, err := tool.ReadConfig(ConfPath); errors.Is(err, tool.ErrFileNotExist) {
		log.Printf("config file not exist, use command line args")
	} else if err != nil {
		log.Fatalf("read config error: %v", err)
	} else {
		Config = Cfg
	}

	if logout {
		doLogout()
		return
	}

	if Config.Username == "" || Config.Password == "" {
		log.Fatalf("username or password is empty")
	}

	_isp := ISPMap[Config.ISP]
	if _isp == "" {
		log.Fatalf("isp not support")
	}

	Config.Username = strings.ReplaceAll(_isp, "{{UNAME}}", Config.Username)

	log.Printf("continue with Identity: DDDDD= %s | upass= %s", Config.Username, Config.Password)

	if Config.RunTyp == 0 {
		doConnect()()
		return
	}

	// start cron job
	cr := cron.New(cron.WithSeconds())

	doConnect()()
	if _, err := cr.AddFunc(Config.Cron, doConnect()); err != nil {
		log.Fatalf("cron add func error: %v", err)
	}

	cr.Start()

	select {}
}

func doConnect() func() {
	return func() {
		log.Println("check network ...")
		if checkNetwork() {
			log.Println("network status ok")
			return
		}
		log.Printf("network status error, try to login ...")
		doVerify(1)
	}
}

func doVerify(retry int) {

	if retry > 3 {
		log.Printf("login error (%d), exected max retries, waiting for next loop", retry)
		return
	}

	_idx := map[string]string{
		"DDDDD": Config.Username,
		"upass": Config.Password,
	}

	client := resty.New().SetRedirectPolicy(resty.FlexibleRedirectPolicy(0))

	var err error
	var resp *resty.Response
	if resp, err = client.R().SetFormData(_idx).Post(Gateway); resp.StatusCode() != 302 {
		log.Printf("login error (%d) | statusCode: %d , err: %v", retry, resp.StatusCode(), err)

		time.Sleep(time.Second * 2)
		doVerify(retry + 1)
	}

	// check
	var ul *url.URL
	if ul, err = url.Parse(resp.Header().Get("Location")); err != nil {
		log.Printf("login failed (%d), unable to parse location: %v", retry, err)
	}
	errMsg := ul.Query().Get("ErrorMsg")
	acLogOut := ul.Query().Get("ACLogOut")
	retCode := ul.Query().Get("RetCode")

	// login Failed
	if acLogOut != "" {
		errMsgByt, err := base64.StdEncoding.DecodeString(errMsg)
		if err != nil {
			errMsgByt = []byte(errMsg)
		}
		log.Printf("Login Failed (%d), ACLogOut: %s,RetCode: %s, ErrMsg: %s", retry, acLogOut, retCode, string(errMsgByt))

		time.Sleep(time.Second * 2)
		doVerify(retry + 1)
		return
	}

	log.Println("login success !")

	time.Sleep(2 * time.Second)
	if checkNetwork() {
		log.Printf("network check pass")
		return
	}

	log.Printf("network check failed. retry now")
	doVerify(1)
}

// generate_204
func checkNetwork() bool {
	client := resty.New().SetTimeout(time.Second * 5)
	if resp, err := client.R().Get("https://api.xsot.cn/generate_204"); err != nil {
		return false
	} else if resp.StatusCode() != 204 {
		return false
	}

	return true
}

func doLogout() {
	client := resty.New().SetRedirectPolicy(resty.FlexibleRedirectPolicy(0)).SetTimeout(time.Second * 2)
	resp, err := client.R().Get("http://10.0.1.52:801/eportal/?c=ACSetting&a=Logout&jsVersion=" + JsVer)
	if resp.StatusCode() != 302 {
		log.Fatalf("logout error: %v", err)
	}

	var ul *url.URL
	if ul, err = url.Parse(resp.Header().Get("Location")); err != nil {
		log.Fatalf("logout error: %v", err)
	}

	acLogOut := ul.Query().Get("ACLogOut")
	if acLogOut != "1" && acLogOut != "2" {
		log.Println("logout failed =>", ul.Path)
	}
	log.Println("logout success")
}

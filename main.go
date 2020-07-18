package main

import (
	"flag"
	"fmt"
	"fyne.io/fyne/widget"
	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"gopkg.in/gomail.v2"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var button *widget.Button

func main()  {

	var starting        string //出发
	var target    		string //到达
	var depDt       	string //出发时间
	var shitf        	string //航班编码

	var help        	bool //命令行帮助


	flag.StringVar(&starting, "s", "", "出发地点")
	flag.StringVar(&target, "t", "", "到达地点")
	flag.StringVar(&depDt, "dt", "", "出发时间")
	flag.StringVar(&shitf, "c", "", "航班编号")
	flag.BoolVar(&help, "h", false, "帮助")

	flag.Parse()

	if help {
		flag.Usage = usage
		flag.Usage()

		os.Exit(0)
	}


	start(starting, target, depDt, shitf)
}



func usage() {
	fmt.Fprintf(os.Stderr, `东航随心飞抢票工具 

使用: dh [-s 出发地] [-t 到达地点] [-dt 出发时间] [-c 航班编号]

选项:
`)
	flag.PrintDefaults()
}

func start(starting string, target string, depDt string,code string) {

	for true {


		html := form(starting, target, depDt)

		if html == "false" {
			now := time.Now()
			dateString := fmt.Sprintf("%d-%d-%d %d:%d:%d",now.Year(),now.Month(),now.Day(),now.Hour(),now.Minute(),now.Second())
			fmt.Println(dateString +" ----->  意外中断 \n")
		}

		doc,_ := goquery.NewDocumentFromReader(strings.NewReader(html))

		doc.Find("#sylvanas_3 .flight").Each(func(i int, s *goquery.Selection) {

			band := s.Find(".title").Text()
			match, _ := regexp.MatchString("("+code+")", band)

			if (match) {
				number,_ :=  s.Find("*[name='lowest']").First().Attr("class")
				match, _ := regexp.MatchString("(noPrice)", number)
				fmt.Println(number)
				if(match) {
					now := time.Now()
					dateString := fmt.Sprintf("%d-%d-%d %d:%d:%d",now.Year(),now.Month(),now.Day(),now.Hour(),now.Minute(),now.Second())
					fmt.Println(dateString +" ----->  查询无余票 \n")
				} else {
					//TODO 实现一个高可达的通知， 把邮件配置文件配置化
					m := gomail.NewMessage()
					m.SetHeader("From", "=?utf-8?b?55uR5o6n5bCP5Yqp5omL?= <@163.com>")
					m.SetHeader("To", "@qq.com")
					m.SetHeader("Subject", "东航监控提示")
					m.SetBody("text/html", "<b>监控到有票！</b> ")

					d := gomail.NewDialer("smtp.163.com", 25, "@163.com", "")

					// Send the email to Bob, Cora and Dan.
					if err := d.DialAndSend(m); err != nil {
						panic(err)
					}
				}
			}

		})
		time.Sleep(30*time.Second)

	}

}


/*
表单操作提取数据
*/
func form (starting string, target string,  depDt string) string{

	rand.Seed(time.Now().UnixNano())

	port := 8*1000 + rand.Intn(9)*100 + rand.Intn(9)*10 + rand.Intn(9)

	fmt.Println("随机chromedriver端口:" + strconv.Itoa(port))

	pwd, _ := os.Getwd()
	seleniumPath := pwd+"/chromedriver"

	opts := []selenium.ServiceOption{

	}
	//selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService(seleniumPath, port, opts...)

	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}

	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}


	chromeCaps := chrome.Capabilities{
		Args: []string{
			"--headless", // 设置Chrome无头模式
			"--no-sandbox",
			//"--blink-settings=imagesEnabled=false",
			//"--disable-gpu",
			"--disable-images",
			//"--proxy-server=116.196.88.52:4444",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36", // 模拟user-agent，防反爬
		},
	}

	/*caps.AddProxy(selenium.Proxy{
		Type:          selenium.Manual,
		SOCKSVersion:     5,
		SOCKS:         "127.0.0.1:24000",
		NoProxy: []string{
			"http://www.ceair.com/upload/2018/9/20190925daxingjichang.jpg",
			"http://www.ceair.com/ad/sytytd/202003/W020200312633151698975.png",
			"http://static-cdn.ceair.com/resource/images/public/map-mark-old.png?v=zh_CN_18459",
		},
	})*/

	caps.AddChrome(chromeCaps)


	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))


	defer  wd.Quit()


	if err != nil {
		panic(err)
	}

	wd.Get("http://lumtest.com/myip.json")
	ip,_ := wd.PageSource()
	fmt.Println(ip)

	if err := wd.Get("http://www.ceair.com"); err != nil {
		fmt.Print(err)
		return "false"
	}




	time.Sleep(10*time.Second)
	//关闭广告
	gg,err := wd.FindElement(selenium.ByClassName, "pop_wrap_close")

	if err != nil {
		fmt.Print(err)
		return "false"
	}

	gg.Click()

	//输入起始地址
	startElement,err := wd.FindElement(selenium.ByID, "label_ID_0")

	if err != nil {
		fmt.Print(err)
		return "false"
	}

	startElement.Clear()
	startElement.SendKeys(starting)

	//点击下拉联想
	time.Sleep(8*time.Second)

	drag,err := wd.FindElement(selenium.ByXPATH, "//body/div/div/div/div/div/div/div/ul/li[1]")
	if err != nil {
		fmt.Print(err)
		return "false"
	}
	drag.Click()


	//输入抵达地址
	time.Sleep(1*time.Second)
	targetElement,err := wd.FindElement(selenium.ByID, "label_ID_1")
	if err != nil {
		return "false"
	}
	targetElement.Clear()
	targetElement.SendKeys(target)


	//点击下拉联想
	time.Sleep(8*time.Second)
	drag,err = wd.FindElement(selenium.ByXPATH, "//div//div//div//div//div//div//div[1]//ul[1]")
	if err != nil {
		fmt.Print(err)
		return "false"
	}
	drag.Click()


	//输入抵达时间

	depDtElement,err := wd.FindElement(selenium.ByID, "depDt")
	if err != nil {
		fmt.Print(err)
		return "false"
	}
	depDtElement.Click()
	depDtElement.Clear()
	wd.ExecuteScript("depDt.value='"+depDt+"'", nil)


	depDtCloseElement,err := wd.FindElement(selenium.ByXPATH, "//html//body")
	if err != nil {
		fmt.Print(err)
		return "false"
	}
	depDtCloseElement.Click()

	//点击提交
	btn,err := wd.FindElement(selenium.ByID, "btn_flight_search")
	if err != nil {
		fmt.Print(err)
		return "false"
	}
	btn.Click()



	//切换窗口
	WindowHandles,_ := wd.WindowHandles()
	if (len(WindowHandles) != 2 ) {
		return  "false"
	}
	wd.SwitchWindow(WindowHandles[1])


	//东航接口慢等它先加载下
	time.Sleep(20*time.Second)
	html,_ := wd.PageSource()

	return html

}



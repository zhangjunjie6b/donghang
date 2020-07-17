package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	url2 "net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

var button *widget.Button

func main()  {


	app := app.NewWithID("io.phpzjj.demo")
	app.SetIcon(theme.FyneLogo())
	w := app.NewWindow("东航周末随心飞监控软件")
	w.Resize(fyne.NewSize(500,600))


	//起点
	starting := widget.NewEntry()
	starting.SetPlaceHolder("start")

	//终点
	target := widget.NewEntry()
	target.SetPlaceHolder("target")


	//起飞时间
	depDt := widget.NewEntry()
	depDt.SetPlaceHolder("depDt (eg:2020-08-02)")


	code := widget.NewEntry()
	code.SetPlaceHolder("code")

	info := widget.NewMultiLineEntry()
	entryLoremIpsumScroller := widget.NewVScrollContainer(info)
	entryLoremIpsumScroller.SetMinSize(fyne.NewSize(500,500))

	link,_ := url2.Parse("http://www.ceair.com/booking/sha-xnn-200801_CNY.html")


	hyperlink := widget.NewHyperlink("click to get URL", link)


	button = widget.NewButton("start", func() {
		info.SetText("Start ->  [" + code.Text + "]")
		go start(starting.Text, target.Text, depDt.Text, code.Text, app, info)
	})


	w.SetContent(widget.NewVBox(
		starting,
		target,
		depDt,
		hyperlink,
		code,
		button,
		widget.NewButton("quit", func() {
			app.Quit()
		}),
		entryLoremIpsumScroller,
	))

	defer app.Quit()

	w.SetFixedSize(true)
	w.ShowAndRun()
}

func start(starting string, target string, depDt string,code string, app fyne.App, info *widget.Entry) {

	button.SetText("ing")
	button.Disable()

	for true {


		html := form(starting, target, depDt)

		if html == "false" {
			now := time.Now()
			dateString := fmt.Sprintf("%d-%d-%d %d:%d:%d",now.Year(),now.Month(),now.Day(),now.Hour(),now.Minute(),now.Second())
			info.SetText( dateString +"    ----->   no! \n" + info.Text)
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
					info.SetText( dateString +"    ----->   no \n" + info.Text)
				} else {
					app.SendNotification(&fyne.Notification{
						Title: "随心飞提示:",
						Content: "有余票",
					})
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
	pwd, _ := os.Getwd()
	seleniumPath := pwd+"/chromedriver"
	const (
		port            = 8080
	)
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
			"--user-agent = Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0.2 Safari/604.4.7", // 模拟user-agent，防反爬
		},
	}

	/*caps.AddProxy(selenium.Proxy{
		Type:          selenium.Manual,
		SOCKSVersion:     5,
		SOCKS:         "116.196.88.52:4444",
		NoProxy: []string{
			"http://www.ceair.com/upload/2018/9/20190925daxingjichang.jpg",
			"http://www.ceair.com/ad/sytytd/202003/W020200312633151698975.png",
			"http://static-cdn.ceair.com/resource/images/public/map-mark-old.png?v=zh_CN_18459",
		},
	})
*/
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
		return "false"
	}




	time.Sleep(10*time.Second)
	//关闭广告
	gg,err := wd.FindElement(selenium.ByClassName, "pop_wrap_close")

	if err != nil {
		return "false"
	}

	gg.Click()

	//输入起始地址
	startElement,err := wd.FindElement(selenium.ByID, "label_ID_0")

	if err != nil {
		return "false"
	}

	startElement.Clear()
	startElement.SendKeys(starting)

	//点击下拉联想
	time.Sleep(8*time.Second)

	drag,err := wd.FindElement(selenium.ByXPATH, "//body/div/div/div/div/div/div/div/ul/li[1]")
	if err != nil {
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
		return "false"
	}
	drag.Click()


	//输入抵达时间

	depDtElement,err := wd.FindElement(selenium.ByID, "depDt")
	if err != nil {
		return "false"
	}
	depDtElement.Click()
	depDtElement.Clear()
	wd.ExecuteScript("depDt.value='"+depDt+"'", nil)


	depDtCloseElement,err := wd.FindElement(selenium.ByXPATH, "//html//body")
	if err != nil {
		return "false"
	}
	depDtCloseElement.Click()

	//点击提交
	btn,err := wd.FindElement(selenium.ByID, "btn_flight_search")
	if err != nil {
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



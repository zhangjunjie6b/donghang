package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/tebeka/selenium"
	url2 "net/url"
	"os"
	"regexp"
	"strings"
	"time"
	"github.com/PuerkitoBio/goquery"
)

var button *widget.Button

func main()  {
	app := app.NewWithID("io.phpzjj.demo")
	app.SetIcon(theme.FyneLogo())
	w := app.NewWindow("东航周末随心飞监控软件")
	w.Resize(fyne.NewSize(500,600))

	url := widget.NewEntry()
	url.SetPlaceHolder("URL")

	code := widget.NewEntry()
	code.SetPlaceHolder("code")

	info := widget.NewMultiLineEntry()
	entryLoremIpsumScroller := widget.NewVScrollContainer(info)
	entryLoremIpsumScroller.SetMinSize(fyne.NewSize(500,500))

	link,_ := url2.Parse("http://www.ceair.com/booking/sha-xnn-200801_CNY.html")


	hyperlink := widget.NewHyperlink("click to get URL", link)


	button = widget.NewButton("start", func() {
		info.SetText("Start ->  [" + code.Text + "]")
		go start(url.Text, code.Text, app, info)
	})


	w.SetContent(widget.NewVBox(
		url,
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

func start( url string,  code string, app fyne.App, info *widget.Entry) {

	button.SetText("ing")
	button.Disable()
	// Start a Selenium WebDriver server instance (if one is not already
	// running).

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

		wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
		if err != nil {
			panic(err)
		}
		defer wd.Quit()

	for true {
		if err := wd.Get(url); err != nil {
			panic(err)
		}
		html, _ := wd.PageSource()

		time.Sleep(30*time.Second)

		doc,_ := goquery.NewDocumentFromReader(strings.NewReader(html))

		doc.Find("#sylvanas_3 .flight").Each(func(i int, s *goquery.Selection) {

			band := s.Find(".title").Text()
			match, _ := regexp.MatchString("("+code+")", band)

			if (match) {
				number,_ :=  s.Find("*[name='lowest']").First().Attr("class")
				match, _ := regexp.MatchString("(noPrice)", number)

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

		time.Sleep(10*time.Minute)

	}

}
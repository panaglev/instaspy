package core

import (
	"fmt"
	"instaspy/src/logger"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
)

type seleniumRemote struct {
	remote selenium.WebDriver
}

// Establish remote connection to
// Selenium image
func EstablishRemote() (*seleniumRemote, error) {
	const op = "core.EstablishRemote"

	/*
		docker run -d -p 4444:4444 -e SE_NODE_SESSION_TIMEOUT=1000 --shm-size="4g" --name selenium-server selenium/standalone-chrome
	*/
	caps := selenium.Capabilities{
		"browserName": "chrome",
		"chromeOptions": map[string]interface{}{ // Some speed improvments, I don't know what are they doing
			"args": []string{ // TODO(?) Disable VNC
				"--headless",
				"--ignore-certificate-errors",
				"--ignore-ssl-errors",
				"--ignore-gpu-blacklist",
				"--use-gl",
				"--no-sandbox",
				"--disable-web-security",
				"--disable-gpu",
			},
		},
	}
	wd, err := selenium.NewRemote(caps, "http://192.168.1.3:4444/wd/hub") // Change to address of container in docker-compose
	if err != nil {
		logger.HandleOpError(op, err)
		return &seleniumRemote{}, err
	}

	return &seleniumRemote{remote: wd}, nil
}

// Correct terminate selenium session
// Still need fix for SIGENV
func (s *seleniumRemote) Quit() {
	const op = "core.Close"

	s.remote.Quit()
}

// Load page -> save source code -> extract history links
func (s *seleniumRemote) Job(username string) ([]string, []string, error) {
	const op = "core.Job"

	pageUrl := fmt.Sprintf("https://instanavigation.com/ru/user-profile/%s", username)
	err := s.remote.Get(pageUrl)
	if err != nil {
		logger.HandleOpError(op, err)
		return []string{}, []string{}, err
	}

	pageSource, err := s.remote.PageSource()
	if err != nil {
		/*
				Problem during parse job at cmd.instaparser.main: Error getting page source at core.Job:
			unexpected alert open: unexpected alert open: {Alert text : Профиль не был найден. Попробуйте перезагрузить страницу!}
			(Session info: chrome=114.0.5735.133)exit status 1
		*/
		s.Quit() // - Quick solution. I guess you have to return control and try one more time
		// Maybe even reccurcive call 5 times, for example.
		logger.HandleOpErrorWithComment(op, err, "Profile might 1. not exists 2. can't be accessed right now")
		return []string{}, []string{}, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(pageSource))
	if err != nil {
		logger.HandleOpError(op, err)
		return []string{}, []string{}, err
	}

	imageLinks := []string{}

	// Find all img elements with the "data-src" attribute within the elements with class "profile-stories-item"
	doc.Find(".profile-stories-item img[data-src]").Each(func(i int, img *goquery.Selection) {
		dataSrc, exists := img.Attr("data-src")
		if exists {
			imageLinks = append(imageLinks, dataSrc)
		}
	})

	/*
		For debug:
		time.Sleep(10 * time.Second)
		pageUrl := fmt.Sprintf("https://iganony.io/profile/%s", username)

		The problem is that I can not have the links for video or even for pictures
		when dealing with "lazyload-wapper " class. Maybe robots.txt
		protecting from selenium so we have to change user-agent
		TODO: add video grabbing

		!!this uses instanavigator!!
		fmt.Println(doc.Find(".story-video video[src]"))
		doc.Find(".story-video video[src]").Each(func(i int, img *goquery.Selection) {
			src, exists := img.Attr("src")
			if exists {
				videoLinks = append(videoLinks, src)
			} else {
				fmt.Println("src attribute not found for image.")
			}
		})

		!!this uses iagony!!
		videoLinks := []string{}
		doc.Find(".lazyload-wrapper ").Each(func(i int, element *goquery.Selection) {
			html, _ := element.Html()
			fmt.Println(html)
			src, exists := element.Find("img, video").Attr("src")
			fmt.Println(src, exists)
			if exists {
				videoLinks = append(videoLinks, src)
			} else {
				fmt.Println("src attribute not found for video.")
			}
		})
		fmt.Println("Video Links:", videoLinks)
	*/

	return imageLinks, nil, nil
}

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

func (s *seleniumRemote) Quit() {
	const op = "core.Close"

	s.remote.Quit()
}

func (s *seleniumRemote) Job(username string) ([]string, []string, error) {
	const op = "core.Job"

	pageUrl := fmt.Sprintf("https://instanavigation.com/ru/user-profile/%s", username)
	err := s.remote.Get(pageUrl)
	if err != nil {
		return []string{}, []string{}, fmt.Errorf("Error getting page by URL at %s: %w", op, err)
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
		return []string{}, []string{}, fmt.Errorf("Error getting page source at %s: %w", op, err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(pageSource))
	if err != nil {
		return []string{}, []string{}, fmt.Errorf("Error getting page info at %s: %w", op, err)
	}

	imageLinks := []string{}

	// Find all img elements with the "data-src" attribute within the elements with class "profile-stories-item"
	doc.Find(".profile-stories-item img[data-src]").Each(func(i int, img *goquery.Selection) {
		dataSrc, exists := img.Attr("data-src")
		if exists {
			imageLinks = append(imageLinks, dataSrc)
		} else {
			fmt.Println("data-src attribute not found for image.")
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

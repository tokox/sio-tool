package szkopul_client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Arapak/sio-tool/cookiejar"

	"github.com/fatih/color"
)

type SzkopulClient struct {
	Jar            *cookiejar.Jar `json:"cookies"`
	Username       string         `json:"handle"`
	Password       string         `json:"password"`
	LastSubmission *Info          `json:"last_submission"`
	host           string
	path           string
	client         *http.Client
}

var Instance *SzkopulClient

func Init(path, host, proxy string) {
	jar, _ := cookiejar.New(nil)
	c := &SzkopulClient{Jar: jar, LastSubmission: nil, path: path, host: host, client: nil}
	if err := c.load(); err != nil {
		color.Red(err.Error())
		color.Green("Create a new session in %v", path)
	}
	Proxy := http.ProxyFromEnvironment
	if len(proxy) > 0 {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			color.Red(err.Error())
			color.Green("Use default proxy from environment")
		} else {
			Proxy = http.ProxyURL(proxyURL)
		}
	}
	c.client = &http.Client{Jar: c.Jar, Transport: &http.Transport{Proxy: Proxy}}
	if err := c.save(); err != nil {
		color.Red(err.Error())
	}
	Instance = c
}

func (c *SzkopulClient) load() (err error) {
	file, err := os.Open(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)

	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, c)
	if err != nil {
		return
	}

	parsed_url, err := url.Parse(c.host)
	if err != nil {
		return
	}

	cookies := c.Jar.Cookies(parsed_url)

	for _, cookie := range cookies {
		if cookie.Name == "lang" {
			cookie.Value = "pl"
		}
	}

	c.Jar.SetCookies(parsed_url, cookies)

	return nil
}

func (c *SzkopulClient) save() (err error) {
	data, err := json.MarshalIndent(c, "", "  ")
	if err == nil {
		err = os.MkdirAll(filepath.Dir(c.path), os.ModePerm)
		if err == nil {
			err = os.WriteFile(c.path, data, 0644)
		}
	}
	if err != nil {
		color.Red("Cannot save session to %v\n%v", c.path, err.Error())
	}
	return
}

package codeforces_client

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

type CodeforcesClient struct {
	Jar            *cookiejar.Jar `json:"cookies"`
	Handle         string         `json:"handle"`
	HandleOrEmail  string         `json:"handle_or_email"`
	Password       string         `json:"password"`
	Ftaa           string         `json:"ftaa"`
	Bfaa           string         `json:"bfaa"`
	LastSubmission *Info          `json:"last_submission"`
	host           string
	proxy          string
	path           string
	client         *http.Client
}

var Instance *CodeforcesClient

func Init(path, host, proxy string) {
	jar, _ := cookiejar.New(nil)
	c := &CodeforcesClient{Jar: jar, LastSubmission: nil, path: path, host: host, proxy: proxy, client: nil}
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

func (c *CodeforcesClient) load() (err error) {
	file, err := os.Open(c.path)
	if err != nil {
		return
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)

	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, c)
}

func (c *CodeforcesClient) save() (err error) {
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

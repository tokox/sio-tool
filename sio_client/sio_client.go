package sio_client

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

type SioInstanceClient int

const (
	Staszic SioInstanceClient = 0
	Mimuw   SioInstanceClient = 1
	Talent  SioInstanceClient = 2
)

type SioClient struct {
	Jar            *cookiejar.Jar `json:"cookies"`
	Username       string         `json:"handle"`
	Password       string         `json:"password"`
	LastSubmission *Info          `json:"last_submission"`
	host           string
	path           string
	client         *http.Client
	instanceClient SioInstanceClient
}

var StaszicInstance *SioClient
var MimuwInstance *SioClient
var TalentInstance *SioClient

func Init(path, host, proxy string, instanceClient SioInstanceClient) {
	jar, _ := cookiejar.New(nil)
	c := &SioClient{Jar: jar, LastSubmission: nil, path: path, host: host, client: nil}
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
	c.instanceClient = instanceClient
	if instanceClient == Staszic {
		StaszicInstance = c
	} else if instanceClient == Mimuw {
		MimuwInstance = c
	} else if instanceClient == Talent {
		TalentInstance = c
	}
}

func (c *SioClient) load() (err error) {
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
		return err
	}

	return json.Unmarshal(bytes, c)
}

func (c *SioClient) save() (err error) {
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

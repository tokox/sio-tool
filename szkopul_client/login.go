package szkopul_client

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"syscall"

	"github.com/Arapak/sio-tool/cookiejar"
	"github.com/Arapak/sio-tool/util"

	"github.com/fatih/color"
	"golang.org/x/term"
)

var ErrorNotLogged = "not logged in"

func findUsername(body []byte) (username string, err error) {
	reg := regexp.MustCompile(`<strong class="username" id="username">([\s\S]+?)</strong>`)
	tmp := reg.FindSubmatch(body)
	if len(tmp) < 2 {
		return "", errors.New(ErrorNotLogged)
	}
	return string(tmp[1]), nil
}

func findCsrf(body []byte) (string, error) {
	reg := regexp.MustCompile(`<input type="hidden" name="csrfmiddlewaretoken" value="(.+?)">`)
	tmp := reg.FindSubmatch(body)
	if len(tmp) < 2 {
		return "", errors.New("cannot find csrf")
	}
	return string(tmp[1]), nil
}

func (c *SzkopulClient) GetCsrf(URL string) (csrf string, err error) {
	body, err := util.GetBody(c.client, URL)
	if err != nil {
		return
	}
	return findCsrf(body)
}

func (c *SzkopulClient) Login() (err error) {
	color.Cyan("Login...\n")

	jar, _ := cookiejar.New(nil)

	c.client.Jar = jar
	csrf, err := c.GetCsrf(c.host + "/login/")
	if err != nil {
		return
	}
	password, err := c.DecryptPassword()
	if err != nil {
		return
	}

	form := url.Values{}
	form.Add("csrfmiddlewaretoken", csrf)
	form.Add("login_view-current_step", "auth")
	form.Add("auth-username", c.Username)
	form.Add("auth-password", password)

	req, err := http.NewRequest("POST", c.host+"/login/", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", c.host+"/login/")
	req.Header.Set("Origin", c.host)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	username, err := findUsername(body)
	if err != nil {
		return err
	}

	c.Username = username
	c.Jar = jar

	color.Green("Succeed!!")
	color.Green("Welcome %v~", c.Username)
	return c.save()
}

func createHash(key string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hasher.Sum(nil)
}

func encrypt(handle, password string) (ret string, err error) {
	block, err := aes.NewCipher(createHash("glhf" + handle + "233"))
	if err != nil {
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}
	text := gcm.Seal(nonce, nonce, []byte(password), nil)
	ret = hex.EncodeToString(text)
	return
}

func decrypt(handle, password string) (ret string, err error) {
	data, err := hex.DecodeString(password)
	if err != nil {
		err = errors.New("cannot decode the password")
		return
	}
	block, err := aes.NewCipher(createHash("glhf" + handle + "233"))
	if err != nil {
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}
	nonceSize := gcm.NonceSize()
	nonce, text := data[:nonceSize], data[nonceSize:]
	plain, err := gcm.Open(nil, nonce, text, nil)
	if err != nil {
		return
	}
	ret = string(plain)
	return
}

func (c *SzkopulClient) DecryptPassword() (string, error) {
	if len(c.Password) == 0 || len(c.Username) == 0 {
		return "", errors.New("you have to configure your username and password by `st config`")
	}
	return decrypt(c.Username, c.Password)
}

func (c *SzkopulClient) ConfigLogin() (err error) {
	if c.Username != "" {
		color.Green("Current user: %v", c.Username)
	}
	color.Cyan("Configure username and password")
	color.Cyan("Note: The password is invisible, just type it correctly.")

	fmt.Printf("username: ")
	username := util.ScanlineTrim()

	password := ""
	if term.IsTerminal(syscall.Stdin) {
		fmt.Printf("password: ")
		bytePassword, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			fmt.Println()
			if err.Error() == "EOF" {
				fmt.Println("Interrupted.")
				return nil
			}
			return err
		}
		password = string(bytePassword)
		fmt.Println()
	} else {
		color.Red("Your terminal does not support the hidden password.")
		fmt.Printf("password: ")
		password = util.Scanline()
	}

	c.Username = username
	c.Password, err = encrypt(username, password)
	if err != nil {
		return
	}
	return c.Login()
}

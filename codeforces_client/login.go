package codeforces_client

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

func genFtaa() string {
	return util.RandString(18)
}

func genBfaa() string {
	return "f1b3f18c715565b589b7823cda7448ce"
}

var ErrorNotLogged = "Not logged in"

func findHandle(body []byte) (string, error) {
	reg := regexp.MustCompile(`handle = "([\s\S]+?)"`)
	tmp := reg.FindSubmatch(body)
	if len(tmp) < 2 {
		return "", errors.New(ErrorNotLogged)
	}
	return string(tmp[1]), nil
}

func findCsrf(body []byte) (string, error) {
	reg := regexp.MustCompile(`csrf='(.+?)'`)
	tmp := reg.FindSubmatch(body)
	if len(tmp) < 2 {
		return "", errors.New("cannot find csrf")
	}
	return string(tmp[1]), nil
}
func AesDecrypt(cipherIn []byte, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData := make([]byte, len(cipherIn))
	blockMode.CryptBlocks(origData, cipherIn)
	return origData, nil
}
func addRCPC(c *CodeforcesClient, body []byte) ([]byte, error) {
	if strings.Contains(string(body), "Redirecting... Please, wait.") {
		reg := regexp.MustCompile(`var a=toNumbers\("([0-9a-f]*)"\),b=toNumbers\("([0-9a-f]*)"\),c=toNumbers\("([0-9a-f]*)"\);`)
		out := reg.FindAllSubmatch(body, -1)
		if len(out) != 1 {
			return nil, fmt.Errorf("cannot find key in body body %v out %+v", string(body), out)
		} else {
			key, _ := hex.DecodeString(string(out[0][1]))
			iv, _ := hex.DecodeString(string(out[0][2]))
			cipherIn, _ := hex.DecodeString(string(out[0][3]))
			cipherOut, err := AesDecrypt(cipherIn, key, iv)
			if err != nil {
				return nil, err
			}
			loginUrl, err := url.Parse(c.host + "/enter")
			if err != nil {
				return nil, err
			}
			c.client.Jar.SetCookies(loginUrl, []*http.Cookie{
				{Name: "RCPC", Value: hex.EncodeToString(cipherOut)},
			})
			body, err := util.GetBody(c.client, c.host+"/enter")
			return body, err
		}
	}
	return body, nil
}

func (c *CodeforcesClient) Login() (err error) {
	color.Cyan("Login %v...\n", c.HandleOrEmail)

	password, err := c.DecryptPassword()
	if err != nil {
		return
	}

	jar, _ := cookiejar.New(nil)

	c.client.Jar = jar
	body, err := util.GetBody(c.client, c.host+"/enter")
	if err != nil {
		return
	}
	body, err = addRCPC(c, body)
	if err != nil {
		return
	}
	csrf, err := findCsrf(body)
	if err != nil {
		return
	}

	ftaa := genFtaa()
	bfaa := genBfaa()

	body, err = util.PostBody(c.client, c.host+"/enter", url.Values{
		"csrf_token":    {csrf},
		"action":        {"enter"},
		"ftaa":          {ftaa},
		"bfaa":          {bfaa},
		"handleOrEmail": {c.HandleOrEmail},
		"password":      {password},
		"_tta":          {"176"},
		"remember":      {"on"},
	})
	if err != nil {
		return
	}

	handle, err := findHandle(body)
	if err != nil {
		return
	}

	c.Ftaa = ftaa
	c.Bfaa = bfaa
	c.Handle = handle
	c.Jar = jar
	color.Green("Succeed!!")
	color.Green("Welcome %v~", handle)
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

func (c *CodeforcesClient) DecryptPassword() (string, error) {
	if len(c.Password) == 0 || len(c.HandleOrEmail) == 0 {
		return "", errors.New("you have to configure your handle and password by `st config`")
	}
	return decrypt(c.HandleOrEmail, c.Password)
}

func (c *CodeforcesClient) ConfigLogin() (err error) {
	if c.Handle != "" {
		color.Green("Current user: %v", c.Handle)
	}
	color.Cyan("Configure handle/email and password")
	color.Cyan("Note: The password is invisible, just type it correctly.")

	fmt.Printf("handle/email: ")
	handleOrEmail := util.ScanlineTrim()

	password := ""
	if term.IsTerminal(int(syscall.Stdin)) {
		fmt.Printf("password: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
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

	c.HandleOrEmail = handleOrEmail
	c.Password, err = encrypt(handleOrEmail, password)
	if err != nil {
		return
	}
	return c.Login()
}

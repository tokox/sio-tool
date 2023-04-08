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
	"sio-tool/util"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"golang.org/x/term"
)

// ErrorNotLogged not logged in
var ErrorNotLogged = "Not logged in"

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

func (c *SzkopulClient) GetUsername(token string) (username string, err error) {
	req, err := http.NewRequest("GET", c.host+"/api/auth_ping", nil)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Token %v", token))

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)

	if !strings.Contains(string(responseBody), "pong ") {
		return "", fmt.Errorf("this token is not valid")
	}

	username = strings.TrimPrefix(string(responseBody), "\"pong ")
	username = strings.TrimSuffix(username, "\"")

	return username, nil
}

// Login codeforces with handler and password
func (c *SzkopulClient) Login() (err error) {
	color.Cyan("Login...\n")

	token, err := c.DecryptToken()
	if err != nil {
		return
	}

	c.Username, err = c.GetUsername(token)
	if err != nil {
		return
	}

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

// DecryptPassword get real password
func (c *SzkopulClient) DecryptToken() (string, error) {
	if len(c.Token) == 0 || len(c.Username) == 0 {
		return "", errors.New("you have to configure your username and password by `st config`")
	}
	return decrypt(c.Username, c.Token)
}

// ConfigLogin configure handle and password
func (c *SzkopulClient) ConfigLogin() (err error) {
	if c.Username != "" {
		color.Green("Current user: %v", c.Username)
	}
	color.Cyan("Configure API token")
	color.Cyan("Note: The token is invisible, just type/paste it correctly.")

	token := ""
	if term.IsTerminal(int(syscall.Stdin)) {
		fmt.Printf("token: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println()
			if err.Error() == "EOF" {
				fmt.Println("Interrupted.")
				return nil
			}
			return err
		}
		token = string(bytePassword)
		fmt.Println()
	} else {
		color.Red("Your terminal does not support the hidden password.")
		fmt.Printf("token: ")
		token = util.Scanline()
	}

	c.Username, err = c.GetUsername(token)
	if err != nil {
		return
	}
	c.Token, err = encrypt(c.Username, token)
	if err != nil {
		return
	}
	return c.Login()
}

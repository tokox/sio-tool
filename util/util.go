package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"
)

const CHA = "abcdefghijklmnopqrstuvwxyz0123456789"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = CHA[rand.Intn(len(CHA))]
	}
	return string(b)
}

func GetBody(client *http.Client, URL string) ([]byte, error) {
	resp, err := client.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func PostBody(client *http.Client, URL string, data url.Values) ([]byte, error) {
	resp, err := client.PostForm(URL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func DebugSave(data interface{}) {
	f, err := os.OpenFile("./tmp/body", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if dataBytes, ok := data.([]byte); ok {
		// Write the slice of bytes to the file
		if _, err := f.Write(dataBytes); err != nil {
			log.Fatal(err)
		}
	} else {
		// Convert the value of data to a string
		dataString := fmt.Sprintf("%v\n\n", data)

		// Write the string to the file
		if _, err := f.Write([]byte(dataString)); err != nil {
			log.Fatal(err)
		}
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func DebugJSON(data interface{}) {
	text, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(text))
}

const colorReset = "\033[0m"

const colorRed = "\033[31m"
const colorGreen = "\033[32m"

func RedString(str string) string {
	return fmt.Sprintf("%v%v%v", colorRed, str, colorReset)
}

func GreenString(str string) string {
	return fmt.Sprintf("%v%v%v", colorGreen, str, colorReset)
}

type Performance struct {
	fetchingStart time.Time
	parsingStart  time.Time

	Fetching time.Duration
	Parsing  time.Duration
}

func (p *Performance) StartFetching() {
	p.fetchingStart = time.Now()
}
func (p *Performance) StartParsing() {
	p.parsingStart = time.Now()
}

func (p *Performance) StopFetching() {
	p.Fetching += time.Since(p.fetchingStart)
}
func (p *Performance) StopParsing() {
	p.Parsing += time.Since(p.parsingStart)
}

func (p *Performance) Parse() string {
	return fmt.Sprintf("Fetching: %v, Parsing: %v", p.Fetching.Round(time.Millisecond).String(), p.Parsing.Round(time.Microsecond*10).String())
}

func AverageTime(t time.Duration, n int) time.Duration {
	if n == 0 {
		return 0
	}
	return time.Duration(int64(t) / int64(n))
}

func LimitNumOfChars(s string, n int) string {
	unicodeSafeString := []rune(s)
	if len(unicodeSafeString) > n {
		return string(unicodeSafeString[:n])
	}
	return string(unicodeSafeString)
}

func PdfToText(body []byte) ([]byte, error) {
	cmd := exec.Command("pdftotext", "-", "-")
	cmd.Stdin = bytes.NewReader(body)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func AddNewLine(body []byte) []byte {
	if !bytes.HasSuffix(body, []byte("\n")) {
		return append(body, byte('\n'))
	}
	return body
}

func GetValue(message string, val *string, required bool) {
	if required {
		if err := survey.AskOne(&survey.Input{Message: message}, val, survey.WithValidator(survey.Required)); err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
	} else {
		if err := survey.AskOne(&survey.Input{Message: message}, val); err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
	}
}

package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// CHA map
const CHA = "abcdefghijklmnopqrstuvwxyz0123456789"

// RandString n is the length. a-z 0-9
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = CHA[rand.Intn(len(CHA))]
	}
	return string(b)
}

// Scanline scan line
func Scanline() string {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}
	fmt.Println("\nInterrupted.")
	os.Exit(1)
	return ""
}

// ScanlineTrim scan line and trim
func ScanlineTrim() string {
	return strings.TrimSpace(Scanline())
}

// ChooseIndex return valid index in [0, maxLen)
func ChooseIndex(maxLen int) int {
	color.Cyan("Please choose one (index): ")
	for {
		index := ScanlineTrim()
		i, err := strconv.Atoi(index)
		if err == nil && i >= 0 && i < maxLen {
			return i
		}
		color.Red("Invalid index! Please try again: ")
	}
}

func Confirm(note string) bool {
	color.Cyan(note)
	tmp := ScanlineTrim()
	if tmp == "y" || tmp == "Y" || tmp == "" {
		return true
	}
	return false
}

// YesOrNo must choose one
func YesOrNo(note string) bool {
	color.Cyan(note)
	for {
		tmp := ScanlineTrim()
		if tmp == "y" || tmp == "Y" {
			return true
		}
		if tmp == "n" || tmp == "N" {
			return false
		}
		color.Red("Invalid input. Please input again: ")
	}
}

// GetBody read body
func GetBody(client *http.Client, URL string) ([]byte, error) {
	resp, err := client.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// PostBody read post body
func PostBody(client *http.Client, URL string, data url.Values) ([]byte, error) {
	resp, err := client.PostForm(URL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// GetJSONBody read json body
func GetJSONBody(client *http.Client, url string) (map[string]interface{}, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var data map[string]interface{}
	if err = decoder.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

// DebugSave write data to temperary file
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

// DebugJSON debug
func DebugJSON(data interface{}) {
	text, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(text))
}

// IsURL returns true if a given string is an url
func IsURL(str string) bool {
	if _, err := url.ParseRequestURI(str); err == nil {
		return true
	}
	return false
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

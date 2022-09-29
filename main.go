package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Data struct {
}

var wg sync.WaitGroup

// var regexXXXX = regexp.MustCompile(`-`)
var regexMd5 = regexp.MustCompile(`\b[a-fA-F0-9]{32}\b`)
var regexSha1 = regexp.MustCompile(`\b[a-fA-F0-9]{40}\b`)
var regexSha256 = regexp.MustCompile(`\b[a-fA-F0-9]{64}\b`)

// r3, _ := regexp.Compile(`[a-fA-F0-9]{64}`)

func main() {
	resp, err := http.Get("https://malshare.com/daily/")
	if err != nil {
		fmt.Println("Erorr", err.Error())
		// log.Fatalln(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// log.Fatalln(err)
		fmt.Println(err)
		return
	}

	var strArray []string
	r, _ := regexp.Compile(`\d{4}-\d{2}-\d{2}`)
	strArray = r.FindAllString(string(body), -1)
	// var str string
	strNewArray := make([]string, len(strArray))
	for i, v := range strArray {
		strNewArray[i] = strings.ReplaceAll(v, "-", "/")
	}
	path := "newmalshare" //make a directory
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	t := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   time.Hour * 2,
			KeepAlive: time.Hour,
		}).Dial,
		// We use ABSURDLY large keys, and should probably not.
		// TLSHandshakeTimeout: 60 * time.Second,
	}
	c := &http.Client{
		Transport: t,
	}
	message := make(chan string, len(strArray))
	for i := 0; i < len(strArray); i++ {
		newPath := filepath.Join("newmalshare", strNewArray[i])

		err := os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
			// log.Println(err)
		}

		wg.Add(1)
		go func(i int, message chan string, c *http.Client, strArray []string) {
			makeRequest(i, message, c, strArray)
			// writeFileMd5(newPath, message)
			// writeFileSha1(newPath, message)
			// writeFileSha256(newPath, message)
			wg.Done()
			// writeFileMd5()

		}(i, message, c, strArray)

		go func(newPath string, message chan string) {
			// direct = fmt.Sprintf("%s/md5.txt", newPath)
			writeFileMd5(newPath, message)
			wg.Done()
		}(newPath, message)
		// go writeFileMd5(message, *regexMd5, fmt.Sprintf("%s/sha1.txt", newPath))
		go func(newPath string, message chan string) {
			writeFileSha1(newPath, message)
			wg.Done()
		}(newPath, message)
		go func(newPath string, message chan string) {
			writeFileSha256(newPath, message)
			wg.Done()
		}(newPath, message)

	}
	wg.Wait()

	// time.Sleep(time.Minute * 5)
}

//function to write to file

func writeFileMd5(newPath string, message chan string) error {
	// r1, _ := regexp.Compile(`\b[a-fA-F0-9]{32}\b`)

	var strArray []string = regexMd5.FindAllString(<-message, -1)
	direct := fmt.Sprintf("%s/md5.txt", newPath)
	//filepath.Join("newmalshare", strNewArray[i])
	//newPath + "/" + "md5.txt"
	if _, err := os.Stat(direct); os.IsNotExist(err) {
		f, err := os.Create(direct)
		if err != nil {
			// log.Fatal(err)
			fmt.Println(err)
			return err
		}
		defer f.Close()

		for i := 0; i < len(strArray); i++ {
			_, err2 := f.WriteString(string(strArray[i]) + "\n")
			if err2 != nil {
				// log.Fatal(err2)
				fmt.Println(err2)
				return err2
			}
		}
	}
	// defer wg.Done()
	return nil
}

func writeFileSha1(newPath string, message chan string) error {
	var strArray []string = regexSha1.FindAllString(<-message, -1)
	direct := fmt.Sprintf("%s/sha1.txt", newPath) //newPath + "/" + "sha1.txt"
	if _, err := os.Stat(direct); os.IsNotExist(err) {
		f, err := os.Create(direct)
		if err != nil {
			// log.Fatal(err)
			fmt.Println(err)

		}
		defer f.Close()

		for i := 0; i < len(strArray); i++ {
			_, err2 := f.WriteString(string(strArray[i]) + "\n")
			if err2 != nil {
				// log.Fatal(err2)
				fmt.Println(err2)
			}
		}
	}
	// defer wg.Done()
	return nil
}

func writeFileSha256(newPath string, message chan string) error {
	// r3, _ := regexp.Compile(`[a-fA-F0-9]{64}`)
	strArray := regexSha256.FindAllString(<-message, -1)
	// direct := newPath + "/" + "sha256.txt"
	direct := fmt.Sprintf("%s/sha256.txt", newPath)
	if _, err := os.Stat(direct); os.IsNotExist(err) {
		f, err := os.Create(direct)
		if err != nil {
			fmt.Println(err)
			// log.Fatal(err)
		}
		defer f.Close()

		for i := 0; i < len(strArray); i++ {
			_, err2 := f.WriteString(string(strArray[i]) + "\n")
			if err2 != nil {
				fmt.Println(err2)
				// log.Fatal(err2)
			}
		}
	}
	// defer wg.Done()
	return nil
}

func makeRequest(i int, message chan string, c *http.Client, strArray []string) {

	resp, err := c.Get(fmt.Sprintf("https://malshare.com/daily/%s/malshare_fileList.%s.all.txt", strArray[i], strArray[i]))
	if err != nil {
		// time.Sleep(time.Second * 10)
		// makeRequest(i, message, c, strArray)
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	message <- string(body)

}

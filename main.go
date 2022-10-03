package main

import (
	"context"
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

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Crawler struct {
	Date   string `json:"date" bson:"date"`
	Md5    string `json:"md5" bson:"md5"`
	Sha1   string `json:"sha1" bson:"sha1"`
	Sha256 string `json:"sha256" bson:"sha256"`
}
type Site struct {
	URL string
}

// type Result struct {
// 	URL     string
// 	message string
// }

var wg sync.WaitGroup

// var regexXXXX = regexp.MustCompile(`-`)
var regexMd5 = regexp.MustCompile(`\b[a-fA-F0-9]{32}\b`)
var regexSha1 = regexp.MustCompile(`\b[a-fA-F0-9]{40}\b`)
var regexSha256 = regexp.MustCompile(`\b[a-fA-F0-9]{64}\b`)

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
			Timeout:   time.Hour * 2 * 10,
			KeepAlive: time.Hour * 10,
		}).Dial,
		// We use ABSURDLY large keys, and should probably not.
		// TLSHandshakeTimeout: 60 * time.Second,
	}
	c := &http.Client{
		Transport: t,
	}
	var urls [5300]string
	for i := 0; i < len(strArray); i++ {
		urls[i] = fmt.Sprintf("https://malshare.com/daily/%s/malshare_fileList.%s.all.txt", strArray[i], strArray[i])
		fmt.Println(urls[i])
	}

	message := make(chan Site, len(strArray)*3)
	result := make(chan string, len(strArray)*3)
	for w := 1; w <= 10; w++ {
		go makeRequest(message, c, urls[:], result)
	}
	for _, url := range urls {
		message <- Site{URL: url}
	}
	close(message)
	for i := 0; i < len(strArray); i++ {
		newPath := filepath.Join("newmalshare", strNewArray[i])
		err := os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
		res := <-result
		go writeFileMd5(newPath, res)
		go writeFileSha1(newPath, res)
		go writeFileSha256(newPath, res)
	}

}
func writeFileMd5(newPath string, res string) error {

	var strArray []string = regexMd5.FindAllString(res, -1)
	direct := fmt.Sprintf("%s/md5.txt", newPath)
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
	return nil
}

func writeFileSha1(newPath string, res string) error {
	var strArray []string = regexSha1.FindAllString(res, -1)
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

func writeFileSha256(newPath string, res string) error {
	// r3, _ := regexp.Compile(`[a-fA-F0-9]{64}`)
	strArray := regexSha256.FindAllString(res, -1)
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

func makeRequest(message <-chan Site, c *http.Client, urls []string, result chan<- string) {
	for site := range message {

		resp, err := c.Get(site.URL)
		if err != nil {
			time.Sleep(time.Second * 10)
			// fmt.Println(err)
			// return
		}
		// fmt.Println(resp.Body)
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		// fmt.Println(j)
		result <- string(body)
	}
	// fmt.Println(message)
	// close(message)
	// }

}
func connectMongoDb() context.Context {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://hoangdznka123:hoangdznka123@cluster0.okc3a4u.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		fmt.Println(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println(err)
	}
	return ctx
}

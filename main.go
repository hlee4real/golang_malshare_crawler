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
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Site struct {
	URL string
}
type Crawler struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Date   string             `json:"date" bson:"date"`
	Md5    []string           `json:"md5" bson:"md5"`
	Sha1   []string           `json:"sha1" bson:"sha1"`
	Sha256 []string           `json:"sha256" bson:"sha256"`
}

type LungLinh struct {
	Value string ``
	Type  string ``
	Date  string ``
}

type APIResponse struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

// Index Value: Unique
// Index Date

// {"value": "9e95d6b47da7b4a86bdbe4c2e0f6efec", "type": "md5", "date": "2020-10-10"}
// {"value": "9e95d6b47da7b4a86bdbe4c2e0f6efef", "type": "md5", "date": "2020-10-10"}
// {"value": "sha111111111", "type": "sha1", "date": "2020-10-10"}

var collection *mongo.Collection
var regexMd5 = regexp.MustCompile(`\b[a-fA-F0-9]{32}\b`)
var regexSha1 = regexp.MustCompile(`\b[a-fA-F0-9]{40}\b`)
var regexSha256 = regexp.MustCompile(`\b[a-fA-F0-9]{64}\b`)

func main() {
	// use config .env
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://hoangdznka123:hoangdznka123@cluster0.okc3a4u.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		fmt.Println(err)
		//
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println(err)
		//
	}
	collection = client.Database("malshare").Collection("crawler")
	// malshareDatabase := client.Database("malshare")
	// collectionMd5 := malshareDatabase.Collection("md5")
	// collectionSha1 := malshareDatabase.Collection("sha1")
	// collectionSha256 := malshareDatabase.Collection("sha256")

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
	}
	c := &http.Client{
		Transport: t,
	}
	var urls [5300]string // 5400 ?
	// fixed size, dynamic size: use append()
	for i := 0; i < len(strArray); i++ {
		urls[i] = fmt.Sprintf("https://malshare.com/daily/%s/malshare_fileList.%s.all.txt", strArray[i], strArray[i])
	}

	message := make(chan Site, len(strArray)*3)
	result := make(chan string, len(strArray)*3)
	for w := 1; w <= 10; w++ {
		go makeRequest(message, c, result)
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
		var md5Data []string = regexMd5.FindAllString(res, -1)
		var md5Result string
		var sha1Result string
		var sha256Result string
		var sha1Data []string = regexSha1.FindAllString(res, -1)
		var sha256Data []string = regexSha256.FindAllString(res, -1)
		for i := 0; i < len(md5Data); i++ {
			md5Result += md5Data[i] + "\n"
		}
		for i := 0; i < len(sha1Data); i++ {
			sha1Result += sha1Data[i] + "\n"
		}
		for i := 0; i < len(sha256Data); i++ {
			sha256Result += sha256Data[i] + "\n"
		}
		// check dupplicate value of each type
		collection.InsertOne(ctx, bson.D{
			{Key: "Value", Value: md5Result},
			{Key: "Type", Value: "md5"},
			{Key: "Date", Value: strNewArray[i]},
		})
		collection.InsertOne(ctx, bson.D{
			{Key: "Value", Value: sha1Result},
			{Key: "Type", Value: "sha1"},
			{Key: "Date", Value: strNewArray[i]},
		})
		collection.InsertOne(ctx, bson.D{
			{Key: "Value", Value: sha256Result},
			{Key: "Type", Value: "sha256"},
			{Key: "Date", Value: strNewArray[i]},
		})
		// if err := writeFileMd5(newPath, res); err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		// go writeFileSha1(newPath, res)
		// go writeFileSha256(newPath, res)
	}

}
func writeFileMd5(newPath string, res string) error {

	var strArray []string = regexMd5.FindAllString(res, -1)
	direct := fmt.Sprintf("%s/md5.txt", newPath)
	if _, err := os.Stat(direct); os.IsNotExist(err) {
		f, err := os.Create(direct)
		if err != nil {
			// log.Fatal(err)
			// fmt.Println(err)
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
			// fmt.Println(err)
			return err

		}
		defer f.Close()

		for i := 0; i < len(strArray); i++ {
			_, err2 := f.WriteString(string(strArray[i]) + "\n")
			if err2 != nil {
				// log.Fatal(err2)
				// fmt.Println(err2)
				return err2
			}
		}
	}
	// fmt.Println(strArray)
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
			// fmt.Println(err)
			return err
			// log.Fatal(err)
		}
		defer f.Close()

		for i := 0; i < len(strArray); i++ {
			_, err2 := f.WriteString(string(strArray[i]) + "\n")
			if err2 != nil {
				// fmt.Println(err2)
				return err2
				// log.Fatal(err2)
			}
		}
	}
	// defer wg.Done()
	return nil
}

func makeRequest(message <-chan Site, c *http.Client, result chan<- string) error {
	for site := range message {
		resp, err := c.Get(site.URL)
		if err != nil {
			time.Sleep(time.Second * 10)
			// fmt.Println(err)
			return err
			// return
		}
		// fmt.Println(resp.Body)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// fmt.Println(err)
			return err
		}
		// fmt.Println(j)
		resp.Body.Close()
		result <- string(body)
		// close(result)
	}
	return nil
}

// {"message": "", "code": 0, "data": {}} 0 --> success, != 0 --> error

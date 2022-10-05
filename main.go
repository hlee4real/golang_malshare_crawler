package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	// . "github.com/gobeam/mongo-go-pagination"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Site struct {
	URL string
}

type LungLinh struct {
	Value string `json:"value" bson:"value"`
	Type  string `json:"type" bson:"type"`
	Date  string `json:"date" bson:"date"`
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

// docker momngo
func main() {
	// use config .env
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println(err)
		return
	}
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// err = client.Connect(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// 	//
	// 	return
	// }
	// defer client.Disconnect(ctx)

	collection = client.Database("newmalshare").Collection("crawler")

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
	// t := &http.Transport{
	// 	Dial: (&net.Dialer{
	// 		Timeout:   time.Hour * 2 * 10,
	// 		KeepAlive: time.Hour * 10,
	// 	}).Dial,
	// }
	// c := &http.Client{
	// 	Transport: t,
	// }
	// var urls = make([]string, len(strArray))
	// for i := 0; i < len(strArray); i++ {
	// 	urls[i] = fmt.Sprintf("https://malshare.com/daily/%s/malshare_fileList.%s.all.txt", strArray[i], strArray[i])
	// }

	// message := make(chan Site, len(strArray)*3)
	// result := make(chan string, len(strArray)*3)
	// for w := 1; w <= 10; w++ {
	// 	go makeRequest(message, c, result)
	// }
	// for _, url := range urls {
	// 	message <- Site{URL: url}
	// }
	// close(message)
	// for i := 0; i < len(strArray); i++ {
	// 	newPath := filepath.Join("newmalshare", strNewArray[i])
	// 	err := os.MkdirAll(newPath, os.ModePerm)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	res := <-result
	// 	var md5Data []string = regexMd5.FindAllString(res, -1)
	// 	// var md5Result string
	// 	// var sha1Result string
	// 	// var sha256Result string
	// 	var sha1Data []string = regexSha1.FindAllString(res, -1)
	// 	var sha256Data []string = regexSha256.FindAllString(res, -1)
	// 	for i := 0; i < len(md5Data); i++ {
	// 		collection.InsertOne(ctx, bson.D{
	// 			{Key: "Value", Value: md5Data[i]},
	// 			{Key: "Type", Value: "md5"},
	// 			{Key: "Date", Value: strNewArray[i]},
	// 		})
	// 	}
	// 	for i := 0; i < len(sha1Data); i++ {
	// 		collection.InsertOne(ctx, bson.D{
	// 			{Key: "Value", Value: sha1Data[i]},
	// 			{Key: "Type", Value: "sha1"},
	// 			{Key: "Date", Value: strNewArray[i]},
	// 		})
	// 	}
	// 	for i := 0; i < len(sha256Data); i++ {
	// 		collection.InsertOne(ctx, bson.D{
	// 			{Key: "Value", Value: sha256Data[i]},
	// 			{Key: "Type", Value: "sha256"},
	// 			{Key: "Date", Value: strNewArray[i]},
	// 		})
	// 	}
	// collection.InsertOne(ctx, bson.D{
	// 	{Key: "Value", Value: md5Result},
	// 	{Key: "Type", Value: "md5"},
	// 	{Key: "Date", Value: strNewArray[i]},
	// })
	// collection.InsertOne(ctx, bson.D{
	// 	{Key: "Value", Value: sha1Result},
	// 	{Key: "Type", Value: "sha1"},
	// 	{Key: "Date", Value: strNewArray[i]},
	// })
	// collection.InsertOne(ctx, bson.D{
	// 	{Key: "Value", Value: sha256Result},
	// 	{Key: "Type", Value: "sha256"},
	// 	{Key: "Date", Value: strNewArray[i]},
	// })
	// if err := writeFileMd5(newPath, res); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// go writeFileSha1(newPath, res)
	// go writeFileSha256(newPath, res)
	// }
	router := gin.Default()
	router.GET("/lunglinh", getAllLungLinh())
	router.GET("/lunglinh/:value", getALungLinh())
	router.POST("/lunglinh", createLungLinh())
	router.PUT("/lunglinh/:value", updateLungLinh())
	router.DELETE("/lunglinh/:value", deleteLungLinh())
	router.Run()

}
func writeFileMd5(newPath string, res string) error {

	var strArray []string = regexMd5.FindAllString(res, -1)
	direct := fmt.Sprintf("%s/md5.txt", newPath)
	if _, err := os.Stat(direct); os.IsNotExist(err) {
		f, err := os.Create(direct)
		if err != nil {
			return err
		}
		defer f.Close()

		for i := 0; i < len(strArray); i++ {
			_, err2 := f.WriteString(string(strArray[i]) + "\n")
			if err2 != nil {
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
			return err

		}
		defer f.Close()

		for i := 0; i < len(strArray); i++ {
			_, err2 := f.WriteString(string(strArray[i]) + "\n")
			if err2 != nil {
				return err2
			}
		}
	}

	return nil
}

func writeFileSha256(newPath string, res string) error {
	strArray := regexSha256.FindAllString(res, -1)
	direct := fmt.Sprintf("%s/sha256.txt", newPath)
	if _, err := os.Stat(direct); os.IsNotExist(err) {
		f, err := os.Create(direct)
		if err != nil {
			return err
		}
		defer f.Close()

		for i := 0; i < len(strArray); i++ {
			_, err2 := f.WriteString(string(strArray[i]) + "\n")
			if err2 != nil {
				return err2
			}
		}
	}
	return nil
}

func makeRequest(message <-chan Site, c *http.Client, result chan<- string) error {
	for site := range message {
		resp, err := c.Get(site.URL)
		if err != nil {
			time.Sleep(time.Second * 10)
			return err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
		result <- string(body)
	}
	return nil
}

// {"message": "", "code": 0, "data": {}} 0 --> success, != 0 --> error

func getAllLungLinh() gin.HandlerFunc {
	// getAll --> paginate --> limit & offset, page & perPage

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var lunglinhs []LungLinh

		defer cancel()
		results, err := collection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error", Code: 1, Data: nil})
			return
		}
		defer results.Close(ctx)
		var page int = 1
		var perPage int64 = 10
		findOptions := options.Find()
		findOptions.SetSkip((int64(page) - 1) * perPage)
		findOptions.SetLimit(perPage)
		cursor, _ := collection.Find(ctx, bson.M{}, findOptions)
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var lunglinh LungLinh
			cursor.Decode(&lunglinh)
			lunglinhs = append(lunglinhs, lunglinh)
		}

		c.JSON(http.StatusOK, APIResponse{Message: "Success", Code: 0, Data: lunglinhs})
	}
}

func createLungLinh() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var lunglinh LungLinh
		defer cancel()
		if err := c.BindJSON(&lunglinh); err != nil {
			c.JSON(http.StatusBadRequest, APIResponse{Message: "Error", Code: 1, Data: nil})
			return
		}
		newLungLinh := LungLinh{
			Value: lunglinh.Value,
			Type:  lunglinh.Type,
			Date:  lunglinh.Date,
		}
		result, err := collection.InsertOne(ctx, newLungLinh)
		if err != nil {
			c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error", Code: 1, Data: nil})
			//
			return
		}
		c.JSON(http.StatusCreated, APIResponse{Message: "Success", Code: 0, Data: result})
	}
}
func getALungLinh() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		value := c.Param("value")
		var lunglinh LungLinh
		defer cancel()

		err := collection.FindOne(ctx, bson.M{"Value": value}).Decode(&lunglinh)
		if err != nil {
			c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error", Code: 1, Data: nil})
			return
		}
		c.JSON(http.StatusOK, APIResponse{Message: "Success", Code: 0, Data: lunglinh})
	}
}

func updateLungLinh() gin.HandlerFunc {
	// by value
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		value := c.Param("value")
		var lunglinh LungLinh
		defer cancel()
		if err := c.BindJSON(&lunglinh); err != nil {
			c.JSON(http.StatusBadRequest, APIResponse{Message: "Error", Code: 1, Data: nil})
			return
		}
		updateLungLinh := LungLinh{
			Value: lunglinh.Value,
			Type:  lunglinh.Type,
			Date:  lunglinh.Date,
		}
		result, err := collection.UpdateOne(ctx, bson.M{"Value": value}, bson.M{"$set": updateLungLinh})
		if err != nil {
			c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error", Code: 1, Data: nil})
			return
		}
		var updatedLungLinh LungLinh
		if result.MatchedCount == 1 {
			err := collection.FindOne(ctx, bson.M{"Value": value}).Decode(&updatedLungLinh)
			if err != nil {
				c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error", Code: 1, Data: nil})
				return
			}
		}
		c.JSON(http.StatusOK, APIResponse{Message: "Success", Code: 0, Data: updatedLungLinh})
	}
}
func deleteLungLinh() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		value := c.Param("value")
		defer cancel()
		result, err := collection.DeleteOne(ctx, bson.M{"Value": value})
		if err != nil {
			c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error", Code: 1, Data: nil})
			return
		}
		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, APIResponse{Message: "Not found", Code: 1, Data: nil})
			return
		}
		c.JSON(http.StatusOK, APIResponse{Message: "Success", Code: 0, Data: result})
	}
}

// write benchmark func for find
// func BenchmarkFind(b *testing.B) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	for i := 0; i < b.N; i++ {
// 		//find row md5 in collection
// 		results, err := collection.Find(ctx, bson.M{"Type": "md5"})
// 		if err != nil {
// 			b.Errorf("Error: %v", err)
// 		}
// 		defer results.Close(ctx)
// 	}
// }

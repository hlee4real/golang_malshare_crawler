package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Data struct {
}

func main() {
	resp, err := http.Get("https://malshare.com/daily/")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var strArray []string
	r, _ := regexp.Compile(`\d{4}-\d{2}-\d{2}`)
	strArray = r.FindAllString(string(body), -1)
	var str string
	strNewArray := make([]string, len(strArray))
	for i, v := range strArray {
		strNewArray[i] = regexp.MustCompile(`-`).ReplaceAllString(v, "/")
	}
	path := "malshare" //make a directory
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	for i := 0; i < len(strArray); i++ {
		//make small directory for year, month, date
		newPath := "malshare" + "/" + strNewArray[i]
		err := os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
		//request to malshare with date
		str = "https://malshare.com/daily/" + strArray[i] + "/malshare_fileList." + strArray[i] + ".all.txt"
		resp, err := http.Get(str)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		//write to file
		r1, _ := regexp.Compile(`\b[a-fA-F0-9]{32}\b`)
		var strArray1 []string = r1.FindAllString(string(body), -1)
		// fmt.Println(strArray1)
		directMd5 := newPath + "/" + "md5.txt"
		writeFile(str, directMd5, strArray1)

		r2, _ := regexp.Compile(`\b[a-fA-F0-9]{40}\b`)
		var strArray2 []string = r2.FindAllString(string(body), -1)
		// fmt.Println(strArray2)
		directSha1 := newPath + "/" + "sha1.txt"
		writeFile(str, directSha1, strArray2)

		r3, _ := regexp.Compile(`[a-fA-F0-9]{64}`)
		var strArray3 []string = r3.FindAllString(string(body), -1)
		// fmt.Println(strArray3)
		directSha256 := newPath + "/" + "sha256.txt"
		writeFile(str, directSha256, strArray3)
	}
}

//function to write to file

func writeFile(str string, direct string, strArray []string) {
	f, err := os.Create(direct)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	for i := 0; i < len(strArray); i++ {
		_, err2 := f.WriteString(string(strArray[i]) + "\n")
		if err2 != nil {
			log.Fatal(err2)
		}
	}
}

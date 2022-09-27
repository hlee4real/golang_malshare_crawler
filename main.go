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
		f, err := os.Create(directMd5)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()
		for j := 0; j < len(strArray1); j++ {
			_, err2 := f.WriteString(string(strArray1[j]) + "\n")
			if err2 != nil {
				log.Fatal(err2)
			}
		}
		r2, _ := regexp.Compile(`\b[a-fA-F0-9]{40}\b`)
		var strArray2 []string = r2.FindAllString(string(body), -1)
		// fmt.Println(strArray2)
		directSha1 := newPath + "/" + "sha1.txt"
		f1, err3 := os.Create(directSha1)
		if err3 != nil {
			log.Fatal(err3)
		}

		defer f1.Close()
		for j := 0; j < len(strArray2); j++ {
			_, err2 := f1.WriteString(string(strArray2[j]) + "\n")
			if err2 != nil {
				log.Fatal(err2)
			}
		}

		r3, _ := regexp.Compile(`[a-fA-F0-9]{64}`)
		var strArray3 []string = r3.FindAllString(string(body), -1)
		// fmt.Println(strArray3)
		directSha256 := newPath + "/" + "sha256.txt"
		f3, err4 := os.Create(directSha256)
		if err4 != nil {
			log.Fatal(err4)
		}

		defer f3.Close()
		for j := 0; j < len(strArray3); j++ {
			_, err2 := f3.WriteString(string(strArray3[j]) + "\n")
			if err2 != nil {
				log.Fatal(err2)
			}
		}
	}

	// r1, _ := regexp.Compile(`[a-fA-F0-9]{64}`)
	// var strArray1 []string = r1.FindAllString(string(body1), -1)
	// fmt.Println(strArray1)

	// r2, _ := regexp.Compile(`\b[a-fA-F0-9]{32}\b`)
	// var strArray2 []string = r2.FindAllString(string(body1), -1)
	// fmt.Println(strArray2)

	// r3, _ := regexp.Compile(`\b[a-fA-F0-9]{40}\b`)
	// var strArray3 []string = r3.FindAllString(string(body1), -1)
	// fmt.Println(strArray3)

}

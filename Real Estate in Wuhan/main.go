package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

func httpget(url string) (result string, err error){
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()

	// Loop read the page data and transfer
	buf := make([]byte,4096)
	for {
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}

		result += string(buf[:n])
	}
	return
}

func Savefile(idx int, Name, Location, Link [][]string) {
	/*f, err := os.Create("Page"+strconv.Itoa(idx)+".csv")
	if err != nil {
		fmt.Println("Savefile error",err)
		return
	}
	defer f.Close()

	n := len(Name)
	f.WriteString("Real Estate Name"+"\t\t\t"+"Location"+"\t\t\t"+"House Area"+"\t\t\t"+"Website Link\n")
	for i := 0; i < n; i ++ {
		f.WriteString(Name[i][1]+"\t\t\t"+Location[i][1]+"\t\t\t"+Area[i][1]+"\t\t\t"+Link[i][1]+"\n")
	}*/

	n := len(Name)
	f,err := os.Create("Page"+strconv.Itoa(idx)+".csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(f)
	head := []string{"Name","Location","Website Link"}
	w.Write(head)

	for i:=0;i<n-1;i++ {
		str := []string{Name[i][1],Location[i][1],Link[i][1]}
		w.Write(str)
		w.Flush()
	}
}

func spider(i int, page chan int) {
	url := "https://wh.fang.anjuke.com/loupan/all/p" + strconv.Itoa(i) + "/"
	result, err := httpget(url)
	if err != nil {
		fmt.Println("httpget error",err)
		return
	}
	//fmt.Println("Result =",result)

	// write the regular expression

	// The real estate name
	ret1 := regexp.MustCompile(`<span class="items-name">(.*?)</span>`)
	name := ret1.FindAllStringSubmatch(result, -1)

	// Location
	ret2 := regexp.MustCompile(`&nbsp;]&nbsp;(.*?)</span>`)
	location := ret2.FindAllStringSubmatch(result, -1)

	// Area
	// ret3 := regexp.MustCompile(`<span class="building-area">建筑面积：(.*?)</span>`)
	// area := ret3.FindAllStringSubmatch(result, -1)

	//Link
	ret4 := regexp.MustCompile(`<a class="tags-wrap" href="(.*?)" soj=`)
	link := ret4.FindAllStringSubmatch(result, -1)

	// Save the obtained information into a file
	Savefile(i, name, location, link)
	page <- i    // synchronize with main go
}

func crawler(start, end int) {
	fmt.Printf("Crawling from page %d to page %d...\n", start, end)

	page := make(chan int)
	// Loop for each page
	for i:=start; i<=end; i++ {
		go spider(i, page)
	}

	for i:=start;i <=end; i++ {
		fmt.Printf("Page %d crawled successfully!\n",<-page)
	}
}

func main() {
	var start,end int
	fmt.Println("Please enter the start page number (>= 1):")
	fmt.Scan(&start)
	fmt.Println("Please enter the end page number (>= 1):")
	fmt.Scan(&end)
	t1 := time.Now()
	crawler(start,end)
	T := time.Since(t1)
	fmt.Println("This crawling process eventually took",T)
}


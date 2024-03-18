package main

import (
	"encoding/csv"
	"fmt"
	"github.com/antchfx/htmlquery"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func SaveData(results [][]string, filename string) {
	os.Mkdir("spider_data", os.ModePerm)
	csvFile, err := os.Create(fmt.Sprintf("spider_data/%s.csv", filename))
	if err != nil {
		panic(err)
	}
	csvFile.WriteString("\xEF\xBB\xBF")
	writer := csv.NewWriter(csvFile)
	writer.Write([]string{"电影", "描述", "评分", "评论"})
	writer.WriteAll(results)
	writer.Flush()
}

func Spider(filename string) {

	// 基础信息
	var results [][]string
	index_page := 1

	log.Println("Start Crawling at ", time.Now().Format("2024-03-07 15:04:05"))

	// 1.发送请求
	for num := 0; num < 250; num = num + 25 {
		// 1.初始化
		client := http.Client{}
		url := fmt.Sprintf("https://movie.douban.com/top250?start=%d&filter=", num)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("req error : ", err)
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 SLBrowser/9.0.3.1311 SLBChan/10")
		req.Header.Set("Cookie", "bid=1M5ISxZjwWQ; _pk_id.100001.4cf6=a0f14f60273ff6bf.1709822045.; __utmz=30149280.1709822045.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); __utmz=223695111.1709822045.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); __yadk_uid=u80cpoEZg2tQRBgH31CKykaqf5j21aCc; __utma=30149280.1560904285.1709822045.1709822045.1710079316.2; __utma=223695111.46849392.1709822045.1709822045.1710079316.2")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("请求失败，原因如下: ", err)
		}
		defer resp.Body.Close()

		// 2.解析网页
		root, err := htmlquery.Parse(resp.Body)
		if err != nil {
			fmt.Println("解析失败，原因如下: ", err)
		}

		title_len := htmlquery.Find(root, "//*[@class='grid_view']/li")
		sparling_num := len(title_len)
		log.Printf("\n正在爬取第 %d 页", index_page)
		log.Printf("\n第 %d 页爬取到了 %d 条", index_page, sparling_num)
		//log.Print("第 %d 页 爬取到了 %d 条title", index_page, len(title_len))
		for time_num := 1; time_num <= sparling_num; time_num++ {

			// 初始化xpath语法
			titlexpathExpr := fmt.Sprintf("//*[@class='grid_view']/li[%d]//div[@class='hd']//span[@class='title'][1]/text()", time_num)
			descriptionxpathExpr := fmt.Sprintf("//*[@class='grid_view']/li[%d]//div[@class='bd']/p[1]/text()", time_num)
			scorexpathExpr := fmt.Sprintf("//*[@class='grid_view']/li[%d]//div[@class='bd']/div[1]/span[@class='rating_num'][1]text()", time_num)
			reviewxpahtExpr := fmt.Sprintf("//*[@class='grid_view']/li[%d]//div[@class='bd']//span[@class='inq'][1]/text()", time_num)

			title := htmlquery.Find(root, titlexpathExpr)
			descrpiton := htmlquery.Find(root, descriptionxpathExpr)
			score := htmlquery.Find(root, scorexpathExpr)
			review := htmlquery.Find(root, reviewxpahtExpr)

			var review_text string
			if len(review) > 0 {
				review_text = htmlquery.InnerText(review[0])
				// 使用 review_text
			} else {
				review_text = "none"
			}

			title_text := htmlquery.InnerText(title[0])
			descrpiton_text_demo1 := htmlquery.InnerText(descrpiton[0])
			descrpiton_text_demo2 := strings.ReplaceAll(descrpiton_text_demo1, " ", "")
			descrpiton_text := strings.ReplaceAll(descrpiton_text_demo2, "\n", "")
			score_text := htmlquery.InnerText(score[0])

			tmp := []string{title_text, descrpiton_text, score_text, review_text}
			results = append(results, tmp)
		}

		index_page = index_page + 1
	}
	log.Println("Finish Crawling at ", time.Now().Format("2024-03-07 15:04:05"))
	log.Println("正在保存数据")
	SaveData(results, filename)

}

func main() {
	var filename string
	fmt.Println("请输入你要保存的文件名")
	fmt.Scan(&filename)
	Spider(filename)
	fmt.Println("爬取完毕，请到spider_data对应文件夹下进行查看")
}

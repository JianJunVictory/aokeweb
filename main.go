package main

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"log"
	"os"
	"strings"
)

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

var (
	filename = "matchData.xls"
	url      = "http://www.okooo.com/jingcai/"
)
var fileToSave *os.File

func init() {
	fileToSave, _ = os.Create(filename)
	fmt.Println("初始化文件")

}
func DoLottery() {
	fmt.Println("开始爬取澳客网足彩信息")
	aokeBody, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	aokeBody.Find(".cont").Each(func(i int, s *goquery.Selection) {
		cont_1 := s.Find(".cont_1 .riqi a span").Eq(0).Text()
		newdata := ConvertToString(cont_1, "gbk", "utf-8")
		newdata = strings.Replace(newdata, "聽", "", -1)
		riqi := []rune(newdata) //日期

		dateInfo := string(riqi[:len(riqi)-3]) //年月日
		xingqi := string(riqi[len(riqi)-3:])   //星期

		touzhu := s.Find(".touzhu").Find(".touzhu_1")
		touzhu.Each(func(j int, sj *goquery.Selection) {
			var gameInfo []string
			var r1 []string
			var r2 []string
			var r3 []string
			var r4 []string
			var r5 []string
			var s1 []string
			var s2 []string
			var s3 []string
			var s4 []string
			var s5 []string

			num := sj.Find(".xulie").Text() //序列值

			aa := sj.Find(".zhum").Eq(0).Text() //主场名字
			ss := ConvertToString(aa, "gbk", "utf-8")

			bb := sj.Find(".zhum").Eq(1).Text() //客场名称
			tt := ConvertToString(bb, "gbk", "utf-8")

			gameInfo = append(gameInfo, ("序列:" + num), ss, tt, ("   " + dateInfo + "-----" + xingqi))

			tmpUrl := fmt.Sprintf("%s?action=more&LotteryNo=%s&MatchOrder=%s", url, dateInfo, caculteDay(xingqi)+num)

			tmpbody, _ := goquery.NewDocument(tmpUrl)
			tmpbody.Find(".pingd").Each(func(k int, sk *goquery.Selection) {
				var bifenname string
				var bifen string
				sk.Find(".peilv").Each(func(p int, sp *goquery.Selection) {
					bifenname = sp.Text()
					if k <= 12 {
						bifenname = strings.Replace(bifenname, "-", "--", 1)
						r1 = append(r1, bifenname)
					} else if k > 12 && k <= 17 {
						bifenname = strings.Replace(bifenname, "-", "--", 1)
						r2 = append(r2, bifenname)
					} else if k > 17 && k <= 30 {
						bifenname = strings.Replace(bifenname, "-", "--", 1)
						r3 = append(r3, bifenname)
					} else if k > 30 && k <= 38 {
						bifenname = strings.Replace(bifenname, "-", "--", 1)
						r4 = append(r4, bifenname)
					} else {
						bifenname = strings.Replace(bifenname, "-", "--", 1)
						r5 = append(r5, bifenname)
					}
				})
				sk.Find(".peilv_1").Each(func(q int, sq *goquery.Selection) {
					bifen = sq.Text()
					if k <= 12 {
						s1 = append(s1, bifen)
					} else if k > 12 && k <= 17 {
						s2 = append(s2, bifen)
					} else if k > 17 && k <= 30 {
						s3 = append(s3, bifen)
					} else if k > 30 && k <= 38 {
						s4 = append(s4, bifen)
					} else {
						s5 = append(s5, bifen)
					}
				})

			})
			fmt.Println(gameInfo)
			fmt.Println(r1)
			fmt.Println(s1)
			fmt.Println(r2)
			fmt.Println(s2)
			fmt.Println(r3)
			fmt.Println(s3)
			fmt.Println(r4)
			fmt.Println(s4)
			fmt.Println(r5)
			fmt.Println(s5)
			doEccel(fileToSave, gameInfo, r1, r2, r3, r4, r5, s1, s2, s3, s4, s5)
		})

	})

}

func caculteDay(xingqi string) string {

	switch xingqi {
	case "星期一":
		return "1"
	case "星期二":
		return "2"
	case "星期三":
		return "3"
	case "星期四":
		return "4"
	case "星期五":
		return "5"
	case "星期六":
		return "6"
	default:
		return "7"
	}
}
func doEccel(f *os.File, gameInfo, r1, r2, r3, r4, r5, s1, s2, s3, s4, s5 []string) {
	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w := csv.NewWriter(f)
	w.Write(gameInfo)
	w.Write(r1)
	w.Write(s1)
	w.Write(r2)
	w.Write(s2)
	w.Write(r3)
	w.Write(s3)
	w.Write(r4)
	w.Write(s4)
	w.Write(r5)
	w.Write(s5)
	w.Write([]string{""})
	w.Flush()
}
func main() {
	DoLottery()
	fileToSave.Close()
	log.Fatal("抓取结束！")
}

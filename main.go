package main

import (
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"qiming/utils"
	"time"
)

var (
	xing    = "李"
	genders = "2" //1男孩，2女孩
	//公历生日，年月日时分
	year  = "2016"
	month = "10"
	day   = "5"
)

func main() {

	start := time.Now().Unix()
	fmt.Printf("start time:%d", start)
	//proxy,如果不用，可以直接注释tr中的proxy选项
	proxy, _ := url.Parse("http://10.80.18.1:3128")
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 10, //超时时间
	}

	datas := make(chan NameInfo, 10)

	go syncNameInfoToFile(datas)

	waterNames := getWaterWords()
	cys := getChangyongWords()
	for _, water := range waterNames {
		for _, cy := range cys {
			if strings.Trim(string(water), " ") != "" && strings.Trim(string(cy), " ") != "" {
				nameInfo := meimingteng(client, string(water)+string(cy))
				if nameInfo != nil {
					datas <- nameInfo.(NameInfo)
				}

				nameInfo1 := meimingteng(client, string(cy)+string(water))
				if nameInfo1 != nil {
					datas <- nameInfo1.(NameInfo)
				}
			}
		}
	}
	//close the channel
	close(datas)
	end := time.Now().Unix()
	fmt.Printf("end time:%d,cost time:%d", end, end-start)
}

func syncNameInfoToFile(datas <-chan NameInfo) {
	now := time.Now().Unix()
	filename := fmt.Sprintf("result_%d.txt", now)
	file, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	file.WriteString("名字  文化 五行 生肖 五格 分数\n")
	for data := range datas {
		file.WriteString(data.name + " " + strconv.Itoa(data.wenhua) + " " + strconv.Itoa(data.wuxing) + " " + strconv.Itoa(data.shengxiao) + " " + strconv.Itoa(data.wuge) + " " + strconv.Itoa(data.scope) + "\n")
	}
}

//获取名字，可以自己组合
func getName(water string) string {
	return "馥" + water
}

//获取水属性的汉字列表
func getWaterWords() string {
	//水属性的字
	b, _ := ioutil.ReadFile("water")
	words := strings.Replace(string(b), " ", "", -1)
	words = strings.Replace(words, "\n", "", -1)
	return words
}

func getChangyongWords() string {
	b, _ := ioutil.ReadFile("changyong")
	words := strings.Replace(string(b), " ", "", -1)
	words = strings.Replace(words, "\n", "", -1)
	return words
}

func meimingteng(client *http.Client, name string) interface{} {
	data := url.Values{}
	//必须要有，反爬虫
	data.Set("__EVENTTARGET", "ctl00$ContentPlaceHolder1$InputBasicInfo1$btNext")
	data.Set("__EVENTARGUMENT", "")
	data.Set("__VIEWSTATE", "/wEPDwULLTEyNjU5OTUwOTBkGAEFHl9fQ29udHJvbHNSZXF1aXJlUG9zdEJhY2tLZXlfXxYeBTtjdGwwMCRDb250ZW50UGxhY2VIb2xkZXIxJElucHV0QmFzaWNJbmZvMSRyYlNwZWNpZnlCaXJ0aGRheQU9Y3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRJbnB1dEJhc2ljSW5mbzEkcmJTcGVjaWZ5TGluQ2hhblFpbgU9Y3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRJbnB1dEJhc2ljSW5mbzEkcmJTcGVjaWZ5TGluQ2hhblFpbgU+Y3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRJbnB1dEJhc2ljSW5mbzEkcmJOb3RTcGVjaWZ5QmlydGhkYXkFPmN0bDAwJENvbnRlbnRQbGFjZUhvbGRlcjEkSW5wdXRCYXNpY0luZm8xJHJiTm90U3BlY2lmeUJpcnRoZGF5BTFjdGwwMCRDb250ZW50UGxhY2VIb2xkZXIxJElucHV0QmFzaWNJbmZvMSRyYlNvbGFyBTFjdGwwMCRDb250ZW50UGxhY2VIb2xkZXIxJElucHV0QmFzaWNJbmZvMSRyYkx1bmFyBTFjdGwwMCRDb250ZW50UGxhY2VIb2xkZXIxJElucHV0QmFzaWNJbmZvMSRyYkx1bmFyBTdjdGwwMCRDb250ZW50UGxhY2VIb2xkZXIxJElucHV0QmFzaWNJbmZvMSRjYklzTGVhcE1vbnRoBTFjdGwwMCRDb250ZW50UGxhY2VIb2xkZXIxJElucHV0QmFzaWNJbmZvMSRsYnROb25lBTFjdGwwMCRDb250ZW50UGxhY2VIb2xkZXIxJElucHV0QmFzaWNJbmZvMSRsYnROb25lBTJjdGwwMCRDb250ZW50UGxhY2VIb2xkZXIxJElucHV0QmFzaWNJbmZvMSRyYnRMdW5ZdQUyY3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRJbnB1dEJhc2ljSW5mbzEkcmJ0THVuWXUFNGN0bDAwJENvbnRlbnRQbGFjZUhvbGRlcjEkSW5wdXRCYXNpY0luZm8xJHJidFNoaUppbmcFNGN0bDAwJENvbnRlbnRQbGFjZUhvbGRlcjEkSW5wdXRCYXNpY0luZm8xJHJidFNoaUppbmcFMWN0bDAwJENvbnRlbnRQbGFjZUhvbGRlcjEkSW5wdXRCYXNpY0luZm8xJHJidFBvZW0FMWN0bDAwJENvbnRlbnRQbGFjZUhvbGRlcjEkSW5wdXRCYXNpY0luZm8xJHJidFBvZW0FMmN0bDAwJENvbnRlbnRQbGFjZUhvbGRlcjEkSW5wdXRCYXNpY0luZm8xJHJidElkaW9tBTJjdGwwMCRDb250ZW50UGxhY2VIb2xkZXIxJElucHV0QmFzaWNJbmZvMSRyYnRJZGlvbQU6Y3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRJbnB1dEJhc2ljSW5mbzEkY2JsUGVyc29uYWxpdHkkMAU6Y3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRJbnB1dEJhc2ljSW5mbzEkY2JsUGVyc29uYWxpdHkkMQU6Y3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRJbnB1dEJhc2ljSW5mbzEkY2JsUGVyc29uYWxpdHkkMgU6Y3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRJbnB1dEJhc2ljSW5mbzEkY2JsUGVyc29uYWxpdHkkMwU6Y3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRJbnB1dEJhc2ljSW5mbzEkY2JsUGVyc29uYWxpdHkkNAU6Y3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRJbnB1dEJhc2ljSW5mbzEkY2JsUGVyc29uYWxpdHkkNQU6Y3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRJbnB1dEJhc2ljSW5mbzEkY2JsUGVyc29uYWxpdHkkNgU6Y3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRJbnB1dEJhc2ljSW5mbzEkY2JsUGVyc29uYWxpdHkkNwU6Y3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRJbnB1dEJhc2ljSW5mbzEkY2JsUGVyc29uYWxpdHkkOAU6Y3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRJbnB1dEJhc2ljSW5mbzEkY2JsUGVyc29uYWxpdHkkOQU6Y3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRJbnB1dEJhc2ljSW5mbzEkY2JsUGVyc29uYWxpdHkkOc+RBppac3/CXaY8AJwjzwaxBvxk")
	data.Set("__VIEWSTATEGENERATOR", "9F5AD4C7")
	data.Set("__EVENTVALIDATION", "/wEWuwICkof8EwK19uvYAgKlkY7lBAKBmOYKAoCY5goCqcPP0AcCqMPP0AcCq8PP0AcCnd3BDAKVuaPyDgKV99WYCAL9kPD0BgKmsv/0AQKL2c+pBQKMlpPbCAKMloemAQKMluuBCgKMlt/sAgKMlsO3CwKMlreTBAKMlpv+DAKMls+WCgKMlrPyAgLhr8HqBQLhr7W2DgLhr5mRBwLhr438DwLhr/HHCALhr+WiAQLhr8mNCgLhr73pAgLhr+EBAuGv1ewIAvq448ULAvq416AEAvq4u4wNAvq4r9cFAvq4k7IOAvq4h50HAvq46/gPAvq438MIAvq4g/wFAvq498cOAt/RhLABAt/R6JsKAt/R3OYCAt/RwMELAt/RtK0EAt/RmIgNAt/RjNMFAt/R8L4OAt/RpNcLAt/RiLIEArDrpqsHArDrivYPArDr/tEIArDr4rwBArDr1ocKArDruuMCArDrrs4LArDrkqkEArDrxsEBArDrqq0KApWEuIYNApWErOEFApWEkMwOApWEhJcHApWE6PIPApWE3N0IApWEwLgBApWEtIQKApWE2LwHApWEzAcC7p3a8AIC7p3O2wsC7p2ypwQC7p2mgg0C7p2K7QUC7p3+yA4C7p3ikwcC7p3W/g8C7p36lw0C7p3u8gUCw7b86wgCw7bgtgECw7bUkQoCw7a4/QICw7as2AsCw7aQowQCw7aEjg0Cw7bo6QUCw7acggMCw7aA7QsC9Ny8qgEC9Nyg9QkC9NyU0AIC9Nz4uwsC9NzshgQC9NzQ4QwC9NzEzAUC9NyoqA4C9NzcwAsC9NzAqwQCyfXehAcCyfXC7w8CyfW2ywgCyfWalgECyfWO8QkCyfXy3AICyfXmpwsCyfXKggQCyfX+uwECyfXihgoCjZaL8A8CjZb/2wgCjZbjpgECjZbXgQoCjZa77QICjZavyAsCjZaTkwQCjZaH/gwCjZarlwoCjZaf8gIC5q+t6wUC5q+Rtg4C5q+FkQcC5q/p/A8C5q/dxwgC5q/BogEC5q+1jgoC5q+Z6QIC5q/NAQLmr7HtCAL7uM/FCwL7uLOhBAL7uKeMDQLyr4XtDALzr4XtDALwr4XtDALxr4XtDAL2r4XtDAL3r4XtDAL0r4XtDALlr4XtDALqr4XtDALyr8XuDALyr8nuDALyr83uDAKVwsetDgKUwsetDgKXwsetDgKWwsetDgKRwsetDgKQwsetDgKTwsetDgKCwsetDgKNwsetDgKVwoeuDgKVwouuDgKVwo+uDgKVwrOuDgKVwreuDgKVwruuDgKVwr+uDgKVwqOuDgKVwuetDgKVwuutDgKUwoeuDgKUwouuDgKUwo+uDgKUwrOuDgKUwreuDgKUwruuDgKUwr+uDgKUwqOuDgKUwuetDgKUwuutDgKXwoeuDgKXwouuDgKArtmFCAKTrpWGCAKMrpWGCAKNrpWGCAKOrpWGCAKPrpWGCAKIrpWGCAKJrpWGCAKKrpWGCAKbrpWGCAKUrpWGCAKMrtWFCAKMrtmFCAKMrt2FCAKMruGFCAKMruWFCAKMrumFCAKMru2FCAKMrvGFCAKMrrWGCAKMrrmGCAKNrtWFCAKNrtmFCAKNrt2FCAKNruGFCALCsOLeDALdsOLeDALcsOLeDALfsOLeDALesOLeDALZsOLeDALYsOLeDALbsOLeDALKsOLeDALFsOLeDALdsKLdDALdsK7dDALdsKrdDALdsJbdDALdsJLdDALdsJ7dDALdsJrdDALdsIbdDALdsMLeDALdsM7eDALcsKLdDALcsK7dDALcsKrdDALcsJbdDALcsJLdDALcsJ7dDALcsJrdDALcsIbdDALcsMLeDALcsM7eDALfsKLdDALfsK7dDALfsKrdDALfsJbdDALfsJLdDALfsJ7dDALfsJrdDALfsIbdDALfsMLeDALfsM7eDALesKLdDALesK7dDALesKrdDALesJbdDALesJLdDALesJ7dDALesJrdDALesIbdDALesMLeDALesM7eDALZsKLdDALZsK7dDALZsKrdDALZsJbdDALZsJLdDALZsJ7dDALZsJrdDALZsIbdDALZsMLeDALZsM7eDALq3rerBALCuMjbDwKi2vn0BQK1gMHLAgKy5bvtBQK236G8AQLV9/naAgKgqtmeCgLVnOXeBwL62evxBgL22aPyBgL32aPyBgL02aPyBgL12aPyBgLy2aPyBgLz2aPyBgLw2aPyBgLh2aPyBgLu2aPyBgL22ePxBgL22e/xBgL22evxBgL22dfxBgL22dPxBgL22dvxBgL22cfxBgL62e/xBgL/l8bNAwL/l8rNAwL/l77NAwL/l8LNAwL/l9bNAwL/l9rNAwL/l87NAwL/l9LNAwL/l6bNAwL/l6rNAwL9lKKdCgKWk7qbCgL1g6aRDwLXw/fkCALvyL/mDALev/OZDgLL5ZPJCAKi/N2EAQKD5eu3CwKgl5eyCgKeoZaXAwKd8qmzBALwr6/EDwLS6rKCBxK3OW282UF90M07q6v0IUMlKeiv")
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$tbXing", xing)
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$tbMingWords", name)
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$ddlGenders", genders)
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$SPECIFY_BIRHDAY", "rbSpecifyBirthday")
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$CalendarType", "rbSolar")
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$ddlYear", year)
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$ddlMonth", month)
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$ddlDay", day)
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$ddlHour", "13")
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$ddlMinute", "40")
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$tbCountry", "中国")
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$tbProvince", "")
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$tbCity", "")
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$tbOtherHopes", "")
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$ddlCareer", "-2")
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$tbFather", "")
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$tbMother", "")
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$tbAvoidWords", "")
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$tbAvoidSimpParts", "")
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$LoginAnywhere1$tbUserName", "")
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$LoginAnywhere1$tbPwd", "")
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$LoginAnywhere1$tbVCode", "")
	data.Set("ctl00$ContentPlaceHolder1$InputBasicInfo1$LoginAnywhere1$loginParam", "2")

	req, _ := http.NewRequest("POST", "https://www.meimingteng.com/Naming/Default.aspx?Tag=4", strings.NewReader(data.Encode()))

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,fr;q=0.7,tr;q=0.6,zh-TW;q=0.5")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "Params=%26Xing%3d%e6%9d%8e%26Gender%3d1%26Year%3d2018%26Month%3d3%26Day%3d13%26Hour%3d13%26Minute%3d40%26IsSolarCalendar%3d1%26IsLeapMonth%3d0%26NameType%3d2%26ReiterativeLocution%3d0%26Location%3d%e4%b8%ad%e5%9b%bd++%26Career%3d-2%26Personality%3d%26Father%3d%26Mother%3d%26SpecifiedName%3d%26SpecifiedNameIndex%3d0%26OtherHopes%3d%26AvoidWords%3d%26AvoidSimpParts%3d%26SpecifiedMing1SimpParts%3d%26SpecifiedMing2SimpParts%3d%26SpecifiedMing1Stroke%3d%26SpecifiedMing2Stroke%3d%26Tag%3d4%7c2%26LinChanQi%3dFalse%26NamingByCategoryCategoryID%3d-1%26SM1S%3d%26SM2S%3d%26SM1T%3d%26SM2T%3d%26SM1M%3d%26SM2M%3d%26RN%3d%26SpecifiedMing1Spell%3d%26SpecifiedMing2Spell%3d%26SM1Y%3d%26SM2Y%3d%26FA%3d%e4%b8%8b%e5%8d%88++++%e6%98%a5%e5%ad%a3++%e6%ad%a3%e6%9c%88%26LOCATION_COUNTY%3d%e4%b8%ad%e5%9b%bd%26LOCATION_PROVINCE%3d%26LOCATION_CITY%3d%26MING_WORDS%3d%e9%93%ad%e6%ba%a5; ASP.NET_SessionId=tgy22x55l3jdvk2ci4gjpxru; mmtsuser=1; ckcookie=chcookie; HELLO_USER=1; Hm_lvt_637e96da78d1c6c8f8a218c811dea5fb=1521614202; Hm_lpvt_637e96da78d1c6c8f8a218c811dea5fb=1521614588; qrcode=1")
	req.Header.Set("Host", "www.meimingteng.com")
	req.Header.Set("Origin", "https://www.meimingteng.com")
	req.Header.Set("Referer", "https://www.meimingteng.com/Naming/Default.aspx?Tag=4")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 8_10_2) AppleWebKit/456.20 (KHTML, like Gecko) Chrome/45.0.14.25 Safari/567.24")

	//反安全策略，随机IP
	ip := utils.RandomIp()
	req.Header.Set("CLIENT-IP", ip)
	req.Header.Set("X-FORWARDED-FOR", ip)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	doc, _ := goquery.NewDocumentFromResponse(resp)

	//获取文化、五行、生肖、五格分数
	wenhua := doc.Find("#ctl00_ContentPlaceHolder1_ShowNameDetails1_lbNameScore > font:nth-child(5) > b").Text()
	wuxing := doc.Find("#bdAppSummDiv > table:nth-child(7) > tbody > tr > td > font:nth-child(5) > b").Text()
	shengxiao := doc.Find("#bdAppSummDiv > table:nth-child(7) > tbody > tr > td > font:nth-child(9) > b").Text()
	wuge := doc.Find("#bdAppSummDiv > table:nth-child(7) > tbody > tr > td > font:nth-child(13) > b").Text()

	//fmt.Println(xing+name, wenhua, wuxing, shengxiao, wuge)

	wenhuaI, _ := strconv.Atoi(wenhua)
	wuxingI, _ := strconv.Atoi(wuxing)
	shengxiaoI, _ := strconv.Atoi(shengxiao)
	wugeI, _ := strconv.Atoi(wuge)

	return NameInfo{name: xing + name, wenhua: wenhuaI, wuxing: wuxingI, shengxiao: shengxiaoI, wuge: wugeI, scope: (wenhuaI + wuxingI + shengxiaoI + wugeI) / 4}
}

type NameInfo struct {
	name      string
	wenhua    int
	wuxing    int
	shengxiao int
	wuge      int
	scope     int
}

//排序使用
type NameInfoSlice []NameInfo

func (p NameInfoSlice) Len() int           { return len(p) }
func (p NameInfoSlice) Less(i, j int) bool { return p[i].scope > p[j].scope }
func (p NameInfoSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func appendFile(data string) {
	file, _ := os.OpenFile("result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	file.WriteString(data)
	defer file.Close()
}

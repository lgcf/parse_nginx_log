package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var help = flag.Bool("h", false, "使用帮助")
var timepos = flag.Int("timepos", 0, "时间所在的位置")
var urlpos = flag.Int("urlpos", 6, "URL所在的位置")
var statuspos = flag.Int("statuspos", 8, "HTTP状态码所在位置")
var maxtime = flag.Int("maxtime", 1000, "单位ms 过滤大于maxtime的请求")
var maxnum = flag.Int("maxnum", 50, "显示结果的条数")
var statuscode = flag.Int("code", 500, "http错误码")
var slowflag = flag.Bool("slowflag", true, "true:汇总所有URL请求中执行时间超过maxtime的URL，false：汇总所有URL中500错误的页面")

type LogMsgBef struct {
	TotalTime float64
	Total     int
	Err500Num int
}

type LogMsg struct {
	URL       string
	TotalTime float64
	AvgTime   float64
	Total     int
	Err500Num int
}
type LogMsgSlice []LogMsg

func (s LogMsgSlice) Len() int      { return len(s) }
func (s LogMsgSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByAvgTime struct{ LogMsgSlice }

func (s ByAvgTime) Less(i, j int) bool { return s.LogMsgSlice[i].AvgTime > s.LogMsgSlice[j].AvgTime }

type ByTotal struct{ LogMsgSlice }

func (s ByTotal) Less(i, j int) bool { return s.LogMsgSlice[i].Total > s.LogMsgSlice[j].Total }

type ByErr500Num struct{ LogMsgSlice }

func (s ByErr500Num) Less(i, j int) bool {
	return s.LogMsgSlice[i].Err500Num > s.LogMsgSlice[j].Err500Num
}

func readline(logfile string) {
	t_start := time.Now()
	f, err := os.Open(logfile)
	defer f.Close()
	if err != nil {
		fmt.Println("打开文件失败：" + logfile)
		os.Exit(1)
	}

	lineNum := 0

	ireader := bufio.NewReader(f)

	ret := make(map[string]LogMsgBef)

	for {
		lineNum = lineNum + 1
		line, _, errMsg := ireader.ReadLine()

		if errMsg == io.EOF {
			t_end := time.Now()
			t := t_end.Sub(t_start)
			fmt.Printf("解析完成;用时：%v 行数：%v\n", t, lineNum)
			break
		}

		content := string(line)

		strtime, strurl, status := parseline(content)
		if strurl == "" {
			fmt.Println(content)
			continue
		}
		tmp, ok := ret[strurl]
		err500 := 0

		if status == *statuscode {
			err500 = 1
		}

		logmsgbef := new(LogMsgBef)
		if ok {
			if *slowflag {
				logmsgbef.TotalTime = tmp.TotalTime + strtime
			}
			logmsgbef.Total = tmp.Total + 1
			logmsgbef.Err500Num = tmp.Err500Num + err500
		} else {
			if *slowflag {
				logmsgbef.TotalTime = strtime
			}
			logmsgbef.Total = 1
			logmsgbef.Err500Num = err500
		}
		//tmparr := [3]string{strconv.FormatFloat(tmptime, 'f', 6, 64), strconv.Itoa(tmpnum), strconv.Itoa(tmpstatus)}
		ret[strurl] = *logmsgbef
	}

	var logmsgslice LogMsgSlice
	for k, v := range ret {
		logmsg := new(LogMsg)
		logmsg.URL = k
		logmsg.Total = v.Total
		if *slowflag {
			logmsg.TotalTime = v.TotalTime
			logmsg.AvgTime = v.TotalTime / float64(v.Total)
		}
		logmsg.Err500Num = v.Err500Num

		if *slowflag {
			t_maxtime := float64(*maxtime)
			if logmsg.AvgTime >= t_maxtime/1000.0 {
				logmsgslice = append(logmsgslice, *logmsg)
			}
		} else {
			logmsgslice = append(logmsgslice, *logmsg)
		}
	}
	if *slowflag {
		fmt.Println("\n慢请求TOP", *maxnum)
		sort.Sort(ByAvgTime{logmsgslice})
		printLogMsg(logmsgslice, *maxnum)
	} else {
		fmt.Println("\n", *statuscode, "错误页面TOP", *maxnum)
		sort.Sort(ByErr500Num{logmsgslice})
		printErr500LogMsg(logmsgslice, *maxnum)
	}
}

func printLogMsg(ss []LogMsg, snum int) {
	fmt.Println("\n平均时长 \t总次数 \t\t500次数 \tURL地址")
	for i, o := range ss {
		//fmt.Println(o.URL, o.AvgTime, o.TotalTime, o.Total, o.Err500Num)
		fmt.Printf("%.3f \t\t%d \t\t%d \t\t%s\n", o.AvgTime, o.Total, o.Err500Num, o.URL)
		if i > snum {
			break
		}
	}
}

func printErr500LogMsg(ss []LogMsg, snum int) {
	fmt.Println("\n次数\t\tURL地址")
	for i, o := range ss {
		if o.Err500Num > 0 {
			fmt.Printf("%d\t\t%s\n", o.Err500Num, o.URL)
		}
		if i > snum {
			break
		}
	}
}

func parseline(line string) (float64, string, int) {
	line = strings.Replace(line, ", ", ",", -1)
	line = strings.Replace(line, " +", "+", -1)
	arr := strings.Split(line, " ")
	if len(arr) < 9 {
		return 0, "", 0
	}
	tmptime := 0.0
	if *slowflag {
		tmptime, _ = strconv.ParseFloat(arr[*timepos], 32)
	}
	tmpstatus, _ := strconv.Atoi(arr[*statuspos])
	return tmptime, arr[*urlpos], tmpstatus
}

func main() {
	//var logfile = "./maccess.20161129"
	flag.Parse()
	if *help {
		fmt.Println("\n解析Ngxin日志\n")
		flag.PrintDefaults()
		os.Exit(0)
	}

	//	fmt.Println(*timepos, *urlpos, *statuspos)
	var logfile = flag.Arg(0)
	if logfile == "" {
		fmt.Println("请输入要解析文件的路径")
		os.Exit(0)
	}

	//	fmt.Println(logfile)

	readline(logfile)
}

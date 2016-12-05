# parse_nginx_log
Nginx日志解析
参数说明：
    -code int
          http错误码 (default 500)
      -h	使用帮助
      -maxnum int
          显示结果的条数 (default 50)
      -maxtime int
          单位ms 过滤大于maxtime的请求 (default 1000)
      -slowflag
          true:汇总所有URL请求中执行时间超过maxtime的URL，false：汇总所有URL中500错误的页面 (default true)
      -statuspos int
          HTTP状态码所在位置 (default 8)
      -timepos int
          时间所在的位置
      -urlpos int
          URL所在的位置 (default 6)
          
使用方法
  nginx 日志如下：
  0.032 119.6.226.151, 61.55.167.208 - - [30/Nov/2016:00:00:00 +0800] "GET / HTTP/1.0" 301 0 "-" "Mozilla/5.0 (Linux;u;Android 4.2.2;zh-cn;) AppleWebKit/534.46 (KHTML,like Gecko) Version/5.1 Mobile Safari/10600.6.3 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)"
  
  则
  
  时间(timepos) 0.032 pos为 0
  URL地址(urlpos)： /  pos为 6
  HTTP状态码(statuspos): 301 pos 为 8
  
  
  查找请求时间大于500ms的前10个URL地址：
  ./parse_nginx_log -urlpos=6 -timepos=0 -statuspos=8 -maxnum=10 -maxtime=500 access.log
  
  查找HTTP状态码为500的前10个URL地址：
  ./parse_nginx_log -urlpos=6 -statuspos=8 -maxnum=10 -slowflag=false  access.log
  
  
  
  

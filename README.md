# CLI
* [GUI版本](https://github.com/zhangjunjie6b/donghang) 传送门

```cassandraql

东航随心飞抢票工具

使用: dh [-s 出发地] [-t 到达地点] [-dt 出发时间] [-c 航班编号]

选项:
  -c string
    	航班编号
  -dt string
    	出发时间
  -h	帮助
  -s string
    	出发地点
  -t string
    	到达地点

```

# 前言
> 东方航空推出3322周末随心飞后，机票异常难买。这个项目就是因为我也买了这个套餐，需要一点取巧的方法来捡漏刷票而存在的，这样可以说是一劳半年内都可以用。

声明 
> 我也知道代码有诸多不完善的地方诸如
* 一个main文件没有解耦开来
* 邮件通知不可配置必须自己编译，通知未设置频率
* 各种兼容问题，没有封装成docker包
> 就上述问题我持着知道但是不会去改的态度，因为这是一个具有时效性的非团队软件。

#原理

> 其实我们的目的就一个，就是监控东航官网出票系统，看有没有余票

官网查询地址 

```
http://www.ceair.com/booking/sha-xnn-200801_CNY.html
``` 

我们可以发现，航班数据其实是进行了一个异步请求的，并且这个异步请求里面埋了一下暗装。它的`CURL`格式如下

```cassandraql
curl $'http://www.ceair.com/otabooking/flight-search\u0021doFlightSearch.shtml' \
  -H 'Connection: keep-alive' \
  -H 'Accept: application/json, text/javascript, */*; q=0.01' \
  -H 'X-Requested-With: XMLHttpRequest' \
  -H 'User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1' \
  -H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' \
  -H 'Origin: http://www.ceair.com' \
  -H 'Referer: http://www.ceair.com/booking/sha-xnn-200801_CNY.html' \
  -H 'Accept-Language: zh-CN,zh;q=0.9' \
  -H 'Cookie: gr_user_id=316cac51-2692-45cd-a39e-b70248150dd8; grwng_uid=399b9fae-32ca-4fcd-a02b-2204cbe3ab2b; s_fid=73AB8873AD8160C8-25BD91BA291AE489; smidV2=20200618160625fe2625aafe9e7f1be5d6388b4253ecf000c2501c131082fb0; es_login_user=ZHANG%252FJUNJIE; _fbp=fb.1.1592467719875.1625792994; pt_71d4c6a5=uid=2XPdIyg/7hcAc6xJVDxZNg&nid=0&vid=LijfdB1Jcu/OiSL1ChfTDQ&vn=2&pvn=7&sact=1592467754060&to_flag=0&pl=T/Sc0xCX-mv523FP0TDGZw*pt*1592467746520; language=zh_CN; Webtrends=d644d8a7.5a8574dd256c2; _ga=GA1.2.5901444.1592467805; 84bb15efa4e13721_gr_last_sent_cs1=34A0F183FA98F1222CE52BBDA098BBF4; ssxmod_itna=Qq+xnQG=i=ExBiDzOImPY5GQE4BKND70tmYY00QxBMb4iNDnD8x7YDvmIh4pn4cjYnqNLet0mxqhFhoeuPezj3rfa8KicGXzD84i7DKqibDCqD1D3qDkWm7FZ8DADi3DEDDm2Dmqi8DITtDAfLlDDfcDAUV0qGw9MdDGdI/4wr574t9nLj0D0PoKGorpA5=zYnLYie=0Bh5KRvHq2GTZYfmBGFNYOxTlEWTQmGij74=eD===; ssxmod_itna2=Qq+xnQG=i=ExBiDzOImPY5GQE4BKND70tmYY00D8T1olxGNq=qGaKlKfHzbT4hxAP4c4adIrwzebnYxGqdzeFR87efr84IbHKaoFE2/t=7Fvd4qIlFRzkFkkaWWpfxu9uctxTv27dYca9DKhcjWGVxuo1xqK44wK06exOETQoCuhdK7GHrGtxbor9eqwO83rSlorOUm=xZ0fqDWa6P7jkKeUNjwi0SrO0K7kSjLa4OC0zOQfoHTiOzaiDHFFxZB70IZXiVAHHVBn0RB7KVWOKu6OiyMnhLMctuzixN9cxN=rc5NnguPVnvV9uZGKSYDbGImfhKzwxo+KfPnn2bQbiA71A9IjYK67FeYMb7Yu7FeOoxBo6YreorU+3cc/8O1DudPY==YR+bWfYmnO3Pq3vb=quPhE4du/FPr6hqcBMGFG3ATofAA8tkFV+9zGR/RW4WxDKMxzYac0D1YtqYdvDYAkDa0jG3tzGDdba064L8dATqqBrX=dvcqdD49ph/ZKu7NNOrK2DzYd1BNg00NtRnR3McYPGWe7O77DKIu2cz6c2hY7BG+8iqEhrKDcyY552qD7=DYFk=DKFGDD; ecrmWebtrends=124.79.118.230.1594005387905; _gid=GA1.2.1121419070.1594005373; user_cookie=true; ceairWebType=new; JSESSIONID=ukmOH+wWEr13SvauKe5kztE8.laputaServer6; 84bb15efa4e13721_gr_session_id=f2c59e6b-fd4c-4778-9cde-8698dca37a4a; 84bb15efa4e13721_gr_last_sent_sid_with_cs1=f2c59e6b-fd4c-4778-9cde-8698dca37a4a; 84bb15efa4e13721_gr_cs1=34A0F183FA98F1222CE52BBDA098BBF4; 84bb15efa4e13721_gr_session_id_f2c59e6b-fd4c-4778-9cde-8698dca37a4a=true; _gat=1; _gat_UA-80008755-11=1' \
  --data-raw '_=47c796f0c0f311ea931a0db120f1f5cc&searchCond={"adtCount":1,"chdCount":0,"infCount":0,"currency":"CNY","tripType":"OW","recommend":false,"reselect":"","page":"0","sortType":"a","sortExec":"a","seriesid":"47c796f0c0f311ea931a0db120f1f5cc","segmentList":[{"deptCd":"SHA","arrCd":"XNN","deptDt":"2020-08-01","deptAirport":"","arrAirport":"","deptCdTxt":"上海","arrCdTxt":"西宁","deptCityCode":"SHA","arrCityCode":"XNN"}],"version":"A.1.0"}' \
  --compressed \
  --insecure
```

好的告辞，不想在找暗装上花费太多的时间，因此我走了第二条路 selenium + chromedriver


# 安装

> cli 版本主要是用于服务器上，因此安装分2部

## 安装 linux chrome


1. 下载 [chrome](https://www.google.cn/intl/zh-CN/chrome/) 最下方找到 `其他平台` 下载Linux版
2. rpm -ivh [下载文件名]

或者直接yum安装

```cassandraql
yum install google-chrome-stable.x86_64
```

## 找到对应版本的chromedriver

[镜像包下载](http://npm.taobao.org/mirrors/chromedriver/)
[版本对应查看](https://sites.google.com/a/chromium.org/chromedriver/downloads)


## 运行

```cassandraql
nohup ./main.x -c MU9501 -dt 2020-08-08 -s 上海 -t 潮汕 >> log.txt &
```

# 技术调用

* [selenium的GO库](https://godoc.org/github.com/tebeka/selenium#Capabilities.AddLogging)








# web-crawler
GOALNG实现的一个分布式爬虫， 用Elasticsearch保存爬到的东西，可以在多台机器上爬取。它们之间通过rpc通信。
//根目录是coding-180,如果不用这个根目录需要在修改所有的import路径。
git clone https://github.com/SweetWhen/web-crawler.git -C coding-180

1）先启动Elasticsearch，并通过web页面登录确认Elasticsearch在容器中启动正常
#虚拟机中跑elasticsearch经常出现内存不足，不能够启动JVM，下面这个命令可以绕过
docker run  --name esl2 -e ES_JAVA_OPTS="-Xms512m -Xmx512m" -p 9200:9200 elasticsearch

2）启动frontend web服务器
设置查询数据库web工作路径，
/opt/mygo/src/coding-180/crawler/frontend/starter.go文件的
const webRoot = "/opt/mygo/src/coding-180/crawler/frontend/view"指定了web跟目录。
然后启动 go run starter.go 默认是监听8888端口，启动之后可以web访问一下

3）启动itemsaver 模块（可在不同主机上启动）：
root@myubuntu:server# go run itemsaver.go -port 1234
2018/11/09 22:28:42 Listening on :1234

4）启动worker模块（可在不同主机上启动）：
root@myubuntu:server# pwd
/opt/mygo/src/coding-180/crawler_distributed/worker/server
root@myubuntu:server# go run worker.go --help
Usage of /tmp/go-build003730233/b001/exe/worker:
  -port int
    	the port for me to listen on
exit status 2
root@myubuntu:server# go run worker.go -port 9000
2018/11/09 22:30:36 Listening on :9000

5）启动engine模块，并同时指定itemsaver和worker地址：
go run main.go --help

add hellow.

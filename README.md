goscws
======

goscws是scws分词的go语言绑定，源代码只有一个scws.go文件

安装方法
========

首先要安装一下scws的库

* https://github.com/hightman/scws

按里面的方法安装一下，然后去

* http://www.xunsearch.com/scws/

下载字典和规则

例子
====

    package main

    import (
        "fmt"
        "github.com/bylevel/goscws"
    )

    func main() {
        scws := goscws.Scws{}
        scws.New()
        scws.SetDict("/Volumes/dev/dev/c/scws-1.2.2/dict.utf8.xdb", goscws.SCWS_XDICT_XDB)
        scws.SetRule("/Volumes/dev/dev/c/scws-1.2.2/etc/rules.utf8.ini")
        scws.SetCharset("utf8")
        scws.SetIgnore(1)
        scws.SetMulti(goscws.SCWS_MULTI_SHORT & goscws.SCWS_MULTI_DUALITY & goscws.SCWS_MULTI_ZMAIN)
        scws.SendText("这只是一个测试程序")
        for scws.Next() {
            fmt.Println(scws.GetRes())
        }
        scws.Free()
    }

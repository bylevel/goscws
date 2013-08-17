// 单线程版本的测试程序
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
	text := "这只是一个测试程序"
	scws.SendText(text)
	for scws.Next() {
		fmt.Println(scws.GetRes())
	}
	// 循环一万次测试性能
	for i := 0; i < 10000; i++ {
		scws.SendText(text)
		for scws.Next() {

		}
		if i%1000 == 0 {
			fmt.Println(i)
		}
	}
	fmt.Println("按回车键结束")
	var scan string
	fmt.Scanf("%s", &scan)
	scws.Free()
	fmt.Println("完成")
}

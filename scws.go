// goscws是C语言版本的scws分词库的go封装，使之更易于在Go程序中使用。
package goscws

/*
#include <stdio.h>
#include <stdlib.h>
#include <scws/scws.h>
#define SCWS_PREFIX     "/usr/local"
#cgo LDFLAGS: -lscws
*/
import "C"
import (
	"errors"
	"unsafe"
)

const (
	SCWS_XDICT_XDB     = C.SCWS_XDICT_XDB     // 表示直接读取 xdb 文件
	SCWS_XDICT_MEM     = C.SCWS_XDICT_MEM     // 表示将 xdb 文件全部加载到内存中，以 XTree 结构存放，可用异或结合另外2个使用。
	SCWS_XDICT_TXT     = C.SCWS_XDICT_TXT     // 表示要读取的词典文件是文本格式，可以和后2项结合用
	SCWS_MULTI_SHORT   = C.SCWS_MULTI_SHORT   // 短词
	SCWS_MULTI_DUALITY = C.SCWS_MULTI_DUALITY // 二元（将相邻的2个单字组合成一个词）
	SCWS_MULTI_ZMAIN   = C.SCWS_MULTI_ZMAIN   // 重要单字
	SCWS_MULTI_ZALL    = C.SCWS_MULTI_ZALL    // 全部单字
)

// scws主结构体
// typedef struct scws_st {
//   struct scws_st *p;
//   xdict_t d; // 词典指针，可检测是否为 NULL 来判断是否加载成功
//   rule_t r; // 规则集指针，可检测是否为 NULL 来判断是否加载成功
//   unsigned char *mblen;
//   unsigned int mode;
//   unsigned char *txt;
//   int len;
//   int off;
//   int wend;
//   scws_res_t res0; // scws_res_t 解释见后面
//   scws_res_t res1;
//   word_t **wmap;
//   struct scws_zchar *zmap;
// } scws_st, *scws_t;
type Scws struct {
	s    C.scws_t     // scws的C对象
	res  C.scws_res_t // C语言 分词结果集
	cur  C.scws_res_t
	text string //分词的内容
	rs   Res    //分词结果
}

// 分词结果
type Res struct {
	String string  //分词的结果
	Attr   string  //词性
	Idf    float64 //idf值
}

//  分配或初始化与 scws 系列操作的 scws_st 对象。该函数将自动分配、初始化、并返回新对象的指针。 只能通过调用 scws_free() 来释放该对象
func (s *Scws) New() (err error) {
	s.s = C.scws_new()
	if s.s == nil {
		err = errors.New("内存不足")
	}
	return
}

// 释放对象
func (s *Scws) Free() {
	C.scws_free(s.s)
}

// 设定字符编码，支持gbk和utf8两种，默认是gbk
func (s *Scws) SetCharset(charset string) {
	C.scws_set_charset(s.s, C.CString(charset))
}

// 清除并设定当前 scws 操作所有的词典文件
//
// 参数 fpath 词典的文件路径，词典格式是 XDB或TXT 格式。
// 参数 mode 有3种值，参见 AddDict。
// 返回值 成功返回 0，失败返回 -1。
// 注意 若此前 scws 句柄已经加载过词典，则此调用会先释放已经加载的全部词典。和 AddDict 的区别在于会覆盖已有词典。
func (s *Scws) SetDict(fpath string, mode C.int) (err error) {
	if int(C.scws_set_dict(s.s, C.CString(fpath), mode)) == -1 {
		err = errors.New("加载词典失败")
	}
	return
}

// 设定规则集文件
// 参数 fpath 规则集文件的路径。若此前 scws 句柄已经加载过规则集，则此调用会先释放已经加载的规则集。
// 错误 加载失败，scws_t 结构中的 r 元素为 NULL，即通过 s->r == NULL 与否来判断加载的失败与成功。
// 注意 规则集定义了一些新词自动识别规则，包括常见的人名、地区、数字年代等。规则编写方法另行参考其它部分。
func (s *Scws) SetRule(fpath string) (err error) {
	C.scws_set_rule(s.s, C.CString(fpath))
	if s.s.r == nil {
		err = errors.New("加载失败")
	}
	return
}

// 设定分词执行时是否执行针对长词复合切分。（例：“中国人”分为“中国”、“人”、“中国人”）。

// 参数 mode 复合分词法的级别，缺省不复合分词。取值由下面几个常量异或组合：

// SCWS_MULTI_SHORT 短词
// SCWS_MULTI_DUALITY 二元（将相邻的2个单字组合成一个词）
// SCWS_MULTI_ZMAIN 重要单字
// SCWS_MULTI_ZALL 全部单字
func (s *Scws) SetMulti(mode C.int) {
	C.scws_set_multi(s.s, mode)
}

// 设定分词结果是否忽略所有的标点等特殊符号（不会忽略\r和\n）。
//
// 参数 yes 1 表示忽略，0 表示不忽略，缺省情况为不忽略。
func (s *Scws) SetIgnore(yes int) {
	C.scws_set_ignore(s.s, C.int(yes))
}

// C.scws_send_text(s, text, C.int(len(C.GoString(text))))
// void scws_send_text(scws_t s, const char *text, int len) 设定要切分的文本数据。
// 参数 text 文本字符串指针。
// 参数 len 文本的长度。
// 注意 该函数可安全用于二进制数据，不会因为字符串中包括 \0 而停止切分。 这个函数应在 scws_get_result() 和 scws_get_tops() 之前调用。
// scws 结构内部维护着该字符串的指针和相应的偏移及长度，连续调用后会覆盖之前的设定；故不应在多次的 scws_get_result 循环中再调用 scws_send_text() 以免出错。
func (s *Scws) SendText(text string) {
	s.text = text
	C.scws_send_text(s.s, C.CString(text), C.int(len(text)))
}

// 获取下一个分词结果
func (s *Scws) Next() bool {
	// cur为nil说明一个res已经结束
	if s.cur == nil {
		C.scws_free_result(s.res)
		s.res = C.scws_get_result(s.s)
		s.cur = s.res
	}
	// res为nil代表分词已结束
	if s.res == nil {
		return false
	}
	// 生成分词结果
	s.rs = Res{}
	s.rs.String = s.text[s.cur.off : int(s.cur.off)+int(s.cur.len)]
	// 将C语言的char数组转成字符串
	goArray := make([]byte, len(s.cur.attr))
	p := uintptr(unsafe.Pointer(&s.cur.attr[0]))
	for i := 0; i < len(s.cur.attr); i++ {
		j := *(*byte)(unsafe.Pointer(p))
		goArray[i] = j
		p += unsafe.Sizeof(j)
	}
	s.rs.Attr = string(goArray)
	s.rs.Idf = float64(s.cur.idf)
	s.cur = s.cur.next
	return true
}

// 返回结果
func (s *Scws) GetRes() Res {
	return s.rs
}

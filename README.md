# iconvs

就地转换文件编码

## 1.安装
```
go install github.com/fsgo/iconvs@latest
```

## 2.使用

### 2.1 参数
```
# iconvs -help
usage: iconvs [flags] [files ...]
  -f string
    	from encoding
  -l	list all available encodings
  -t string
    	to encoding
  -w	write file (default true)
```

参数 `-f` 和 `-t` ：
1. 忽略大小写和`-`（"UTF-8" 和 utf8 等效）
2. 支持的编码可以使用 `iconvs -l` 查看（见下列 2.3）。

### 2.2 示例

转换一个文件:
```
iconvs -f gbk -t utf8 gbk.txt
```

转换多个文件:
```
iconvs -f gbk -t utf8 gbk.txt a.txt
```

转换目录里所有文件:
```
iconvs -f gbk -t utf8 dir/
```


转换目录里特定文件:
```
iconvs -f gbk -t utf8 dir/*.c
```

若有文件转换失败:
1. 当前文件会跳过，并打印日志，继续处理其他文件
2. 程序 exit code =2

### 2.3 支持的编码
```
# iconvs -l
```

```
Big5
EUC-JP
EUC-KR
GB18030
GBK
HZ-GB2312
IBMCodePage037
IBMCodePage1047
IBMCodePage1140
IBMCodePage437
IBMCodePage850
IBMCodePage852
IBMCodePage855
IBMCodePage860
IBMCodePage862
IBMCodePage863
IBMCodePage865
IBMCodePage866
ISO-2022-JP
ISO-8859-6E
ISO-8859-6I
ISO-8859-8E
ISO-8859-8I
ISO8859-1
ISO8859-10
ISO8859-13
ISO8859-14
ISO8859-15
ISO8859-16
ISO8859-2
ISO8859-3
ISO8859-4
ISO8859-5
ISO8859-6
ISO8859-7
ISO8859-8
ISO8859-9
KOI8-R
KOI8-U
Macintosh
MacintoshCyrillic
ShiftJIS
UTF-16
UTF-16BE
UTF-16LE
UTF-32
UTF-32BE
UTF-32LE
UTF-8
Windows1250
Windows1251
Windows1252
Windows1253
Windows1254
Windows1255
Windows1256
Windows1257
Windows1258
Windows874
WindowsCodePage858
```
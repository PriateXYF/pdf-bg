# PDF-BG

> 一个简单的 PDF 背景颜色修改工具

## 参数说明

```
Usage of ./pdfbg:
  -c string
        需要设置的背景颜色 (default "2D2D2D")
  -f string
        需要处理的 PDF 文件 / 文件夹 (default "test.pdf")
  -o string
        输出的 PDF 文件名 (default "auto")
  -p string
        处理的 PDF 文件页码(如: 1,2,3 ) (default "all")
  -r    遍历模式，处理文件夹下所有 PDF 文件
```

## DEMO

* 修改单个 PDF 文件背景色

```bash
./pdfbg -f test.pdf -c FFFFFF
```

* 遍历当前文件夹下所有 PDF 文件

```bash
./pdfbg -f . -c FFFFFF -r
```

* 指定输出文件名（仅在非遍历模式下有效）

```bash
./pdfbg -f test.pdf -c FFFFFF -o test2.pdf
```

* 指定修改页码（目前有点小BUG）

```bash
./pdfbg -f test.pdf -c FFFFFF -p 1,3
```
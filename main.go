package main

import (
    "pdfbg/common/msg"
	"github.com/signintech/gopdf"
	fpdi "github.com/phpdave11/gofpdi"
    "os"
	"fmt"
    "flag"
    "strconv"
    "path"
    "path/filepath"
    "strings"
    "bufio"
)

var (
    flagColor string
    flagFile string
    flagPages string
    flagIsTraverse bool
    flagOutput string
    dealPages []int
    R, G, B int
)

// 判断文件夹路径是否存在
func pathExists(name string) (bool, error) {
    _, err := os.Stat(name)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

func flagInit() error {
    flag.StringVar(&flagColor, "c", "2D2D2D", "需要设置的背景颜色")
    flag.StringVar(&flagFile, "f", "test.pdf", "需要处理的 PDF 文件 / 文件夹")
    flag.StringVar(&flagOutput, "o", "auto", "输出的 PDF 文件名")
    flag.StringVar(&flagPages, "p", "all", "处理的 PDF 文件页码(如: 1,2,3 )")
    flag.BoolVar(&flagIsTraverse, "r", false, "遍历模式，处理文件夹下所有 PDF 文件")
    flag.Parse()
    err := checkFlag()
    if err != nil {
        msg.Fail("%s，程序退出", err)
        return err
    }
    return nil
}

func str2rgb(color_str string) (red, green, blue int, err error) {
    color64, err := strconv.ParseInt(color_str, 16, 32)
    if err != nil {
        return
    }
    color32 := int(color64) //类型强转
    red = color32 >> 16
    green = (color32 & 0x00FF00) >> 8
    blue = color32 & 0x0000FF
    return
}

func checkFlag() (err error) {
    msg.Info("正在检测参数")
    // 检测颜色参数
    if len(flagColor) != 6 {
        return fmt.Errorf("颜色参数设置错误")
    } else {
        R, G, B, err = str2rgb(flagColor)
        if err != nil {
            return fmt.Errorf("颜色参数设置错误")
        }
    }
    // 检测文件参数
    isExist, _ := pathExists(flagFile)
    if isExist == false {
        return fmt.Errorf("文件或文件夹 %s 不存在", flagFile)
    }
    if strings.ToLower(path.Ext(flagFile)) != ".pdf" && !flagIsTraverse {
        return fmt.Errorf("只支持处理 PDF 文件")
    }
    // 检测页面参数
    if flagPages != "all" {
        pages := strings.Split(flagPages, ",")
        for _, page := range pages {
            ipage, err := strconv.Atoi(page)
            if err != nil {
                return fmt.Errorf("页面参数错误")
            }
            dealPages = append(dealPages, ipage)
        }
    }
    // 检测输出文件参数
    if flagOutput != "auto" && strings.ToLower(path.Ext(flagOutput)) != ".pdf" {
        return fmt.Errorf("只支持输出 PDF 文件")
    }

    info, err := os.Stat(flagFile)
    if err != nil {
        return err
    }
    // 遍历模式下检测
    if flagIsTraverse {
        if flagOutput != "auto" {
            msg.Warn("遍历模式下无法手动指定输出文件名")
        }
        // 检测 -f 参数是否为文件夹
        if !info.IsDir() {
            return fmt.Errorf("遍历模式下 -f 参数只能是文件夹")
        }
    } else {
        if info.IsDir() {
            return fmt.Errorf("非遍历模式下 -f 参数只能是 PDF 文件")
        }
    }
    msg.Good("参数检测完成")
    return nil
}


// 获取原文件的第一页尺寸
func getFirstPageSize(sourceFile string)(width, height float64){
    importer := fpdi.NewImporter()
    importer.SetSourceFile(sourceFile)
    pageSizes := importer.GetPageSizes()
    width = pageSizes[1]["/MediaBox"]["w"]
    height = pageSizes[1]["/MediaBox"]["h"]
    return
}

// 判断数组中是否包含某元素
func has(array []int, val int) bool {
    for _, key := range array {
        if val == key {
            return true
        }
    }
    return false
}

// 处理 PDF 背景颜色
func renderPDF(sourceFile string) (string, error) {
    pdf := gopdf.GoPdf{}
    width, height := getFirstPageSize(sourceFile)
    pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: width, H: height}})
    importer := fpdi.NewImporter()
    importer.SetSourceFile(sourceFile)
    pageSizes := importer.GetPageSizes()
    for pageNo := 1; pageNo <= len(pageSizes); pageNo++ {
        width := pageSizes[pageNo]["/MediaBox"]["w"]
        height := pageSizes[pageNo]["/MediaBox"]["h"]
        pdf.AddPage()
        if flagPages != "all" {
            if has(dealPages, pageNo) {
                pdf.SetFillColor(uint8(R), uint8(G), uint8(B))
            }
        } else {
            pdf.SetFillColor(uint8(R), uint8(G), uint8(B))
        }
        pdf.RectFromUpperLeftWithStyle(0, 0, width, height, "FD")
        tpl1 := pdf.ImportPage(sourceFile, pageNo, "/MediaBox")
        pdf.UseImportedTemplate(tpl1, 0, 0, width, height)
    }
    outputFile := flagOutput
    // 如果未指定输出名称
    if flagOutput == "auto" {
        fileNameWithoutExt := filepath.Base(sourceFile)
        extIndex := strings.LastIndex(fileNameWithoutExt, ".") // 查找最后一个点号的位置
        if extIndex != -1 {
            outputFile = fileNameWithoutExt[:extIndex] // 去除后缀部分
        } else {
            outputFile = fileNameWithoutExt // 没有后缀则直接返回完整文件名
        }
        outputFile = fmt.Sprintf("%s-addbg.pdf", outputFile)
    }
    isExist, _ := pathExists(outputFile)
    if isExist {
        msg.Warn("文件 %s 已存在，将覆盖原文件，是否继续？[按下回车继续]", outputFile)
        reader := bufio.NewReader(os.Stdin)
        input, err := reader.ReadString('\n')
        if err != nil {
            return "", err
        } else if input == "\r\n" || input == "\n" {
            pdf.WritePdf(outputFile)
        } else {
            return "", fmt.Errorf("取消覆盖文件")
        }
    } else {
        pdf.WritePdf(outputFile)
    }
    return outputFile, nil
    
}

func main() {
    err := flagInit()
    if err != nil {
        return
    }
    // 如果处理文件夹下所有文件
    var files []string
    if flagIsTraverse {
        err := filepath.Walk(flagFile, func(fpath string, info os.FileInfo, err error) error {
            if strings.ToLower(path.Ext(fpath)) == ".pdf" {
                files = append(files, fpath)
            }
            
            return nil
        })
        if err != nil {
            panic(err)
        }
        msg.Warn("即将处理目录及子目录下的 %d 个 PDF 文件 [按下回车继续]", len(files))
        reader := bufio.NewReader(os.Stdin)
        input, err := reader.ReadString('\n')
        if err != nil {
            return
        } else if input == "\r\n" || input == "\n" {
            for _, file := range files {
                msg.Info("正在处理 %s 文件", file)
                outputFile, err := renderPDF(file)
                if err != nil {
                    msg.Fail("出现错误，%v", err)
                }else{
                    msg.Good("处理完成，输出文件 : %s", outputFile)
                }
            }
        }
        
    // 处理单个文件
    } else{
        msg.Info("正在处理 %s 文件", flagFile)
        outputFile, err := renderPDF(flagFile)
        if err != nil {
            msg.Fail("程序终止，%v", err)
        }else{
            msg.Good("处理完成，输出文件 : %s", outputFile)
        }
    }
}
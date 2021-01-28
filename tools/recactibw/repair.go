package recactibw

import (
	"errors"
	"fmt"
	"github.com/beevik/etree"
	"log"
	"mytools/utils"
	"path/filepath"
	"strconv"
	"strings"
)

func RepairWorkerFunc(p Repair, bar utils.MyBarChan, f func(*etree.Document, Repair, utils.MyBarChan) error) error {
	_, filename := utils.GetFileName(p.Filepath, "")
	desc := fmt.Sprintf("读取文件 %v", filename)
	bar.OptChan <- utils.MyBarOptChan{
		Desc:  desc,
		Total: 2,
	}
	bar.AddChan <- 1
	ext := filepath.Ext(p.Filepath)
	if ext != ".rrd" {
		e := fmt.Sprintf("%v 文件不是rrd文件", p.Filepath)
		return errors.New(e)
	}
	context, err := GetRrdDumpXml(p.Filepath)
	if err != nil {
		e := fmt.Sprintf("%v 文件转换失败: %v", p.Filepath, err)
		return errors.New(e)
	}
	//读取xml
	doc := etree.NewDocument()
	if err := doc.ReadFromString(context); err != nil {
		e := fmt.Sprintf("%v 内容解析失败: %v", p.Filepath, err)
		return errors.New(e)
	}
	bar.AddChan <- 1
	err = f(doc, p, bar)
	if err != nil {
		return err
	}
	return nil
}

func Print(doc *etree.Document, p Repair, bar utils.MyBarChan) error {
	root := doc.SelectElement("rrd")
	stepInt64, lastRecordInt64 := GetLastRecord(root)
	for i := p.StartTime; i <= p.EndTime; i = i + stepInt64 {
		var in, out string
		for _, rra := range root.SelectElements("rra") {
			database := rra.SelectElement("database")
			row := database.SelectElements("row")
			startRecordInt64 := GetStartRecordPerDatabase(row, stepInt64, lastRecordInt64)
			rewriteIndex := (i - startRecordInt64) / stepInt64
			err := verifyTime(p.Filepath, i, i, startRecordInt64, lastRecordInt64)
			if err == nil {
				in = row[rewriteIndex].SelectElements("v")[0].Text()
				out = row[rewriteIndex].SelectElements("v")[1].Text()
			}
		}
		fmt.Printf("%v %v in: %v, out: %v\n",
			p.Filepath, utils.TimestampToString(i), in, out)
	}
	return nil
}

func Modify(doc *etree.Document, p Repair, bar utils.MyBarChan) error {
	root := doc.SelectElement("rrd")
	stepInt64, lastRecordInt64 := GetLastRecord(root)
	_, filename := utils.GetFileName(p.Filepath, "")

	desc := fmt.Sprintf("处理文件 %v", filename)
	bar.OptChan <- utils.MyBarOptChan{
		Desc:  desc,
		Total: len(root.SelectElements("rra")),
	}
	for _, rra := range root.SelectElements("rra") {
		bar.AddChan <- 1
		database := rra.SelectElement("database")
		row := database.SelectElements("row")
		for i := p.StartTime; i <= p.EndTime; i = i + stepInt64 {

			startRecordInt64 := GetStartRecordPerDatabase(row, stepInt64, lastRecordInt64)
			//判断修改时间点
			referTime := i - p.Interval
			err := verifyTime(p.Filepath, i, referTime, startRecordInt64, lastRecordInt64)
			if err != nil {
				//if p.Details {
				//	log.Println(err)
				//}
				continue
			}
			//获取参考时间值
			referIndex := (referTime - startRecordInt64) / stepInt64
			referIn := row[referIndex].SelectElements("v")[0].Text()
			referOut := row[referIndex].SelectElements("v")[1].Text()
			//重新数据

			referIn, referOut, err = ModifyPercent(p.Filepath, referIn, referOut, p.Percent)
			if err != nil {
				//if p.Details{
				//	log.Println(err)
				//}
				continue
			}
			rewriteIndex := (i - startRecordInt64) / stepInt64
			row[rewriteIndex].SelectElements("v")[0].SetText(referIn)
			row[rewriteIndex].SelectElements("v")[1].SetText(referOut)
		}
	}
	desc = fmt.Sprintf("新做文件 %v", filename)
	bar.OptChan <- utils.MyBarOptChan{
		Desc:  desc,
		Total: 2,
	}
	bar.AddChan <- 1
	err := MakeNewRrd(p.Filepath, doc)
	if err != nil {
		return err
	}
	bar.AddChan <- 1
	return nil
}

//按照百分百进行数据修改
func ModifyPercent(filepath, referIn, referOut string, percent float64) (string, string, error) {
	if strings.ContainsAny(referIn, "NaN") && strings.ContainsAny(referOut, "NaN") {
		return "NaN", "NaN", nil
	}
	var referInFloat64, referOutFloat64 float64
	_, _ = fmt.Sscanf(referIn, "%e", &referInFloat64)
	//if err != nil || 1 != n{
	//	e := fmt.Sprintf("%v %v转换出错", filepath, referIn)
	//	return "NaN", "NaN", errors.New(e)
	//}
	_, _ = fmt.Sscanf(referOut, "%e", &referOutFloat64)
	//if err != nil || 1 != n{
	//	e := fmt.Sprintf("%v %v转换出错", filepath, referOut)
	//	return "NaN", "NaN", errors.New(e)
	//}

	referIn = strconv.FormatFloat(referInFloat64*percent, 'E', -1, 64)
	referOut = strconv.FormatFloat(referOutFloat64*percent, 'E', -1, 64)
	return referIn, referOut, nil

}

func verifyTime(filepath string, i, referTime, startRecordInt64, lastRecordInt64 int64) error {
	if i > lastRecordInt64 {
		e := fmt.Sprintf("%v 修改时间 %v 大于最后采样时间 %v，无法修改",
			filepath, utils.TimestampToString(i), utils.TimestampToString(lastRecordInt64))
		return errors.New(e)
	}
	if i < startRecordInt64 {
		e := fmt.Sprintf("%v 修改时间 %v 小于开始采样时间 %v，无法修改",
			filepath, utils.TimestampToString(i), utils.TimestampToString(startRecordInt64))
		return errors.New(e)
	}
	if referTime < startRecordInt64 {
		e := fmt.Sprintf("%v 参考时间 %v 小于采样开始时间 %v，无法修改",
			filepath, utils.TimestampToString(referTime), utils.TimestampToString(startRecordInt64))
		return errors.New(e)
	}
	return nil
}

//直接获取dump出来的信息
func GetRrdDumpXml(filename string) (string, error) {
	out, err := utils.Bash("rrdtool", "dump", filename)
	if err != nil {
		return "", err
	}
	return out, nil
}

// 备份rrd文件，生成新的rrd
func MakeNewRrd(filename string, doc *etree.Document) error {
	_, suffix := utils.GetFileName(filename, ".rrd")
	//bakRrdName := filename + ".bak"
	xmlfile := "/tmp/" + suffix + ".xml"

	err := doc.WriteToFile(xmlfile)
	if err != nil {
		e := fmt.Sprintf("%v文件重新写入失败", xmlfile)
		return errors.New(e)
	}

	//_, err = utils.Bash("cp", "-rf", filename, bakRrdName)
	//if err != nil {
	//	e := fmt.Sprintf("%v文件备份失败: %v", filename, err)
	//	return errors.New(e)
	//}

	_, err = utils.Bash("rrdtool", "restore", "-f", xmlfile, filename)
	if err != nil {
		e := fmt.Sprintf("%v新文件生成失败: %v", filename, err)
		return errors.New(e)
	}

	return nil
}

// 获取每个database的开始时间
func GetStartRecordPerDatabase(row []*etree.Element, stepInt64, lastRecordInt64 int64) int64 {
	recordSize := len(row) - 1
	recordSizeString := strconv.Itoa(recordSize)
	recordSizeInt64, err := strconv.ParseInt(strings.TrimSpace(recordSizeString), 10, 64)
	if err != nil {
		log.Println(err)
	}
	startRecordInt64 := lastRecordInt64 - (stepInt64 * recordSizeInt64)
	return startRecordInt64
}

// 获取监控step和最后一次记录时间
func GetLastRecord(root *etree.Element) (int64, int64) {
	step := root.SelectElement("step").Text()
	stepInt64, err := strconv.ParseInt(strings.TrimSpace(step), 10, 64)
	if err != nil {
		log.Println(err)
	}
	//获取最后一次更新时间
	lastupdate := root.SelectElement("lastupdate").Text()
	lastupdateInt64, err := strconv.ParseInt(strings.TrimSpace(lastupdate), 10, 64)
	if err != nil {
		log.Println(err)
	}
	lastRecordInt64 := lastupdateInt64 / stepInt64 * stepInt64
	return stepInt64, lastRecordInt64
}

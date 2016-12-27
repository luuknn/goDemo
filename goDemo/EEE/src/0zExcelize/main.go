package main

import (
	"ToolExcelize"
	"fmt"
	"os"
)

func main() {

	//创建xlsx文件
	xlsx := excelize.CreateFile()
	xlsx.NewSheet(2, "Sheet2")
	xlsx.NewSheet(3, "Sheet3")
	xlsx.SetCellInt("Sheet2", "A2", 10)
	xlsx.SetCellStr("Sheet3", "B20", "Hello")
	err := xlsx.WriteTo("D:/test/abc.xlsx")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("创建成功")
	}
	/*
		//对现有的xlsx文件进行写入操作
		xlsx, err := excelize.OpenFile("D:/test/abc.xlsx")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		xlsx.SetCellValue("Sheet1", "A1", 110)
		xlsx.SetCellValue("Sheet1", "B1", "Hello")
		xlsx.NewSheet(4, "GQY")
		xlsx.SetCellStr("Sheet4", "B2", "hi")
		xlsx.SetActiveSheet(2)
		err = xlsx.Save()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}*/

	/*
		//对xlsx文件 进行 读取操作
		xlsx, err := excelize.OpenFile("D:/test/abc.xlsx")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cell := xlsx.GetCellValue("Sheet4", "B2")
		fmt.Println("hello" + cell + "world")
	*/
}

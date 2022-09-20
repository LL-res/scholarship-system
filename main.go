package main

import (
	"github.com/flopp/go-findfont"
	"os"
	"scholarship/UI"
	"scholarship/data"
	"strings"
)

func init() {
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		//fmt.Println(path)
		//楷体:simkai.ttf
		//黑体:simhei.ttf
		if strings.Contains(path, "simkai.ttf") {
			//fmt.Println(path)
			os.Setenv("FYNE_FONT", path) // 设置环境变量  // 取消环境变量 os.Unsetenv("FYNE_FONT")
			break
		}
	}
	//fmt.Println("=============")
}

func main() {
	data.NewStudents()
	UI.APP()
}

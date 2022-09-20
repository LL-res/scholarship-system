package UI

import (
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"html/template"
	"image/color"
	"net/url"
	"scholarship/data"
	"sort"
	"strconv"
	"strings"
	"time"
)

func APP() {
	myApp := app.New()
	NewLoginPage(myApp)
	myApp.Run()
}

func NewLoginPage(myApp fyne.App) {
	loginPage := myApp.NewWindow("奖学金分级评测系统")
	loginPage.SetIcon(theme.HomeIcon())
	loginPage.CenterOnScreen()
	loginPage.Resize(fyne.NewSize(500, 600))
	username := widget.NewEntry()
	password := widget.NewPasswordEntry()
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "账号", Widget: username},
			{Text: "密码", Widget: password}},
		OnSubmit: func() {
			id, _ := strconv.Atoi(username.Text)
			if s, ok := data.DataBase[id]; !ok || s.Password != password.Text {
				NewErrLogin(myApp)
			} else {
				NewCodePage(myApp, id)
				loginPage.Close()

			}
		},
		SubmitText: "登录",
	}
	logo := canvas.NewImageFromFile("./source/badge.png")
	logo.SetMinSize(fyne.NewSize(140, 140))
	badge := container.New(layout.NewCenterLayout(), logo)
	cform := container.NewAdaptiveGrid(3, layout.NewSpacer(), form, layout.NewSpacer())
	bigbox := container.NewVBox(layout.NewSpacer(), badge, layout.NewSpacer(), cform, layout.NewSpacer())
	colorBox := container.NewMax(canvas.NewRectangle(color.RGBA{
		R: 204,
		G: 204,
		B: 255,
		A: 215,
	}), bigbox)
	loginPage.SetContent(colorBox)
	loginPage.Show()
}
func NewCodePage(myApp fyne.App, id int) {
	codePage := myApp.NewWindow("奖学金分级评测系统")
	codePage.CenterOnScreen()
	codePage.SetIcon(theme.InfoIcon())
	codePage.Resize(fyne.NewSize(500, 600))
	var level string
	switch data.DataBase[id].LastLevel {
	case 1:
		level = "一等"
	case 2:
		level = "二等"
	case 3:
		level = "三等"
	}
	infoBox := widget.NewForm([]*widget.FormItem{
		{Text: "姓名\t:", Widget: widget.NewLabel(data.DataBase[id].Name)},
		{Text: "学号\t:", Widget: widget.NewLabel(strconv.Itoa(id))},
		{Text: "曾获奖学金级别:", Widget: widget.NewLabel(level)},
		{Text: "必修课成绩", Widget: widget.NewLabel(strings.Join(data.DataBase[id].GetStringCompulsory(), " "))},
		{Text: "选修课成绩", Widget: widget.NewLabel(strings.Join(data.DataBase[id].GetStringSelective(), " "))},
		{Text: "课业加权成绩:", Widget: widget.NewLabel(data.DataBase[id].GetFinalScore())},
	}...)
	avator := canvas.NewImageFromFile("./source/avator.png") //"C:\\Users\\LL\\Pictures\\Camera Roll\\kuku.jpg")
	avator.SetMinSize(fyne.NewSize(180, 120))
	upBox := container.NewHBox(infoBox, layout.NewSpacer(), avator, layout.NewSpacer())

	score1 := binding.NewString()
	score2 := binding.NewString()
	score3 := binding.NewString()

	codeForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "科研成绩", Widget: widget.NewLabelWithData(score1)},
			{Text: "社会工作", Widget: widget.NewLabelWithData(score2)},
			{Text: "文体竞赛", Widget: widget.NewLabelWithData(score3)},
		},
		OnSubmit: func() {
			s1, _ := score1.Get()
			s2, _ := score2.Get()
			s3, _ := score3.Get()
			NewFinalPage(myApp, id, s1, s2, s3)
			codePage.Close()
		},
		SubmitText: "提交",
	}
	tabs := container.NewAppTabs(
		container.NewTabItem("\t科研成绩\t\t", widget.NewRadioGroup(chooseItem(1), func(s string) {
			score1.Set(indexSearch(chooseItem(1), s))
		})),
		container.NewTabItem("\t社会工作\t\t", widget.NewRadioGroup(chooseItem(2), func(s string) {
			score2.Set(indexSearch(chooseItem(2), s))
		})),
		container.NewTabItem("\t文体竞赛\t\t", widget.NewRadioGroup(chooseItem(3), func(s string) {
			score3.Set(indexSearch(chooseItem(3), s))
		})),
		container.NewTabItem("\t总结\t\t", codeForm),
	)
	pic := canvas.NewImageFromFile("./source/mainbuilding.jpg")
	pic.SetMinSize(fyne.NewSize(150, 120))
	pic.Translucency = 0.4
	codeBigBox := container.NewVBox(widget.NewToolbar(
		widget.NewToolbarAction(theme.NavigateBackIcon(), func() {
			NewLoginPage(myApp)
			codePage.Close()
		})), upBox, tabs, pic)
	colorBox := container.NewMax(canvas.NewRectangle(color.RGBA{
		R: 204,
		G: 255,
		B: 229,
		A: 205,
	}), codeBigBox)
	codePage.SetContent(colorBox)
	codePage.Show()
}
func NewErrLogin(myApp fyne.App) {
	errPage := myApp.NewWindow("错误")
	errPage.CenterOnScreen()
	errPage.SetIcon(theme.ErrorIcon())
	errPage.Resize(fyne.NewSize(200, 100))
	text := canvas.NewText("用户名或密码错误", color.Black)
	ctext := container.NewAdaptiveGrid(3, layout.NewSpacer(), text, layout.NewSpacer())
	cbutton := container.NewAdaptiveGrid(3, layout.NewSpacer(), widget.NewButtonWithIcon("确定", theme.ErrorIcon(), func() {
		errPage.Close()
	}), layout.NewSpacer())
	bigbox := container.NewVBox(ctext, layout.NewSpacer(), cbutton, layout.NewSpacer())
	errPage.SetContent(bigbox)
	errPage.Show()
}
func NewFinalPage(myApp fyne.App, id int, s1, s2, s3 string) {
	finalPage := myApp.NewWindow("成绩排名")
	finalPage.SetIcon(theme.RadioButtonIcon())
	finalPage.Resize(fyne.NewSize(500, 600))
	finalPage.CenterOnScreen()
	infinite := widget.NewProgressBarInfinite()
	finalPage.SetContent(container.NewVBox(layout.NewSpacer(), widget.NewLabel("结果计算中"), infinite, layout.NewSpacer()))
	finalPage.Show()
	time.Sleep(3 * time.Second)
	go func() {
		var rank int
		stuPool := make([]*data.Student, 0)
		func(i int) {
			for _, v := range data.DataBase {
				if v.LastLevel == i {
					stuPool = append(stuPool, v)
				}
			}
		}(data.DataBase[id].LastLevel)
		s1, _ := strconv.Atoi(s1)
		s2, _ := strconv.Atoi(s2)
		s3, _ := strconv.Atoi(s3)
		showScore, _ := strconv.ParseFloat(data.DataBase[id].GetFinalScore(), 64)
		showScore = showScore*0.6 + 0.3*float64(s1) + 0.05*float64(s2) + 0.05*float64(s3)
		sort.Slice(stuPool, func(i, j int) bool {
			var (
				si float64
				sj float64
			)
			if stuPool[i] == data.DataBase[id] {
				si = showScore
			} else {
				si = stuPool[i].GetShowScore()
			}
			if stuPool[j] == data.DataBase[id] {
				sj = showScore
			} else {
				sj = stuPool[j].GetShowScore()
			}
			return si > sj
		})
		for i, v := range stuPool {
			if v.Id == id {
				rank = i + 1
				break
			}
		}
		PoolRank := strconv.Itoa(rank) + "/" + strconv.Itoa(len(stuPool))
		finalPage.SetIcon(theme.ConfirmIcon())
		text := generateText(id, PoolRank, fmt.Sprintf("%.1f", showScore), s1, s2, s3)
		links := generateLink(s1, s2, s3)
		box := container.NewVBox(
			widget.NewToolbar(
				widget.NewToolbarAction(theme.NavigateBackIcon(), func() {
					NewCodePage(myApp, id)
					finalPage.Close()
				}),
				widget.NewToolbarSpacer(),
				widget.NewToolbarAction(theme.HomeIcon(), func() {
					NewLoginPage(myApp)
					finalPage.Close()
				})),
			widget.NewRichTextFromMarkdown(text),
			container.NewHBox(layout.NewSpacer(), container.NewVBox(links...), layout.NewSpacer(), layout.NewSpacer()),
		)
		finalPage.SetContent(box)
		finalPage.Show()
	}()

}
func chooseItem(x int) []string {
	switch x {
	case 1:
		return []string{
			"9分 导师十分满意，有成型科研成果",
			"7分 导师很是满意，有快完成的成果",
			"5分 导师较为满意，有进行中的科研",
			"3分 导师一般满意，无任何科研成果",
			"1分 导师极不满意，无任何科研成果",
		}
	case 2:
		return []string{
			"9分 担任重要学生工作，经常参与社会实践",
			"7分 担任普通学生工作，经常参与社会实践",
			"5分 担任普通学生工作，偶尔参与社会实践",
			"3分 没有任何学生工作，较少参与社会实践",
			"1分 没有任何学生工作，从未参与社会实践",
		}
	case 3:
		return []string{
			"9分 已取得文体或竞赛的国家级奖项",
			"7分 已取得文体或竞赛的省市级奖项",
			"5分 已取得文体或竞赛的校内级奖项",
			"3分 认为自己即将取得相关竞赛奖项",
			"1分 未作出任何文体贡献或取得奖项",
		}
	default:
		res := make([]string, 10)
		for i := range res {
			res[i] = strconv.Itoa(i)
		}
		return res
	}
}
func indexSearch(s []string, target string) string {
	for i, v := range s {
		if v == target {
			return strconv.Itoa(9 - i*2)
		}
	}
	return ""
}
func generateText(id int, PoolRank, ShowScore string, s1, s2, s3 int) string {
	type showStu struct {
		Name           string
		Compulsory     string
		Selective      string
		ShowScore      string
		PoolRank       string
		LastLevel      string
		Prospect       string
		AdviceKeep     string
		Advice1        string
		Advice2        string
		Advice3        string
		Recommendation string
	}
	var stu showStu
	stu.Name = data.DataBase[id].Name
	stu.Compulsory = strings.Join(data.DataBase[id].GetStringCompulsory(), " ")
	stu.Selective = strings.Join(data.DataBase[id].GetStringSelective(), " ")
	stu.ShowScore = ShowScore
	stu.PoolRank = PoolRank
	r, _ := strconv.Atoi(strings.Split(PoolRank, "/")[0])
	switch data.DataBase[id].LastLevel {
	case 1:
		if r > data.FirstClass-data.FtoS {
			stu.Prospect = "降档"
		} else {
			stu.Prospect = "保持不变"
		}
		stu.LastLevel = "一等"
	case 2:
		if r <= data.FtoS {
			stu.Prospect = "升档"
		} else if r > data.SecondClass-data.StoT {
			stu.Prospect = "降档"
		} else {
			stu.Prospect = "保持不变"
		}
		stu.LastLevel = "二等"
	case 3:
		if r <= data.StoT {
			stu.Prospect = "升档"
		} else {
			stu.Prospect = "保持不变"
		}
		stu.LastLevel = "三等"
	}
	num := '1'
	if s1 <= 3 {
		stu.Advice1 = "##  ⁣⁣⁣⁣" + string(num) + ". 建议与导师加强沟通，致力科研"
		num++
	}
	if s2 <= 3 {
		stu.Advice2 = "##  ⁣⁣⁣⁣" + string(num) + ". 建议日常做社会实践，关心集体"
		num++
	}
	if s3 <= 3 {
		stu.Advice3 = "##  ⁣⁣⁣⁣" + string(num) + ". 建议多参加文体竞赛，贡献力量"
		num++
	}
	if s1 > 3 && s2 > 3 && s3 > 3 {
		stu.AdviceKeep = "继续努力"
		stu.Recommendation = "[大连理工大学校内活动网](https://huodong.dlut.edu.cn/)"
	}
	buf := new(bytes.Buffer)
	t, _ := template.ParseFiles("./source/finalPage.txt")
	t.Execute(buf, stu)

	return buf.String()
}
func generateLink(s1, s2, s3 int) []fyne.CanvasObject {
	ans := []fyne.CanvasObject{}
	if s1 <= 3 {
		u1, _ := url.Parse("https://huodong.dlut.edu.cn/info/1192/7705.htm")
		u2, _ := url.Parse("https://huodong.dlut.edu.cn/info/1192/7704.htm")
		ans = append(ans, widget.NewHyperlink("非单调数据系数的线性二次平均场博弈问题", u1))
		ans = append(ans, widget.NewHyperlink("大规模社交网络上的“回声室”效应和社会极化现象", u2))

	}
	if s2 <= 3 {
		u1, _ := url.Parse("https://huodong.dlut.edu.cn/info/1044/7534.htm")
		u2, _ := url.Parse("https://huodong.dlut.edu.cn/info/1202/7467.htm")
		ans = append(ans, widget.NewHyperlink("政策试验与变革管理", u1))
		ans = append(ans, widget.NewHyperlink("学会收纳，重塑自我，开拓思路", u2))
	}
	if s3 <= 3 {
		u1, _ := url.Parse("https://huodong.dlut.edu.cn/info/1073/6245.htm")
		u2, _ := url.Parse("https://huodong.dlut.edu.cn/info/1073/6246.htm")
		ans = append(ans, widget.NewHyperlink("“乐赞芳华”庆祝建校72周年师生音乐会", u1))
		ans = append(ans, widget.NewHyperlink("大连理工大学喜迎建党百年研究生文艺晚会", u2))
	}
	if s1 > 3 && s2 > 3 && s3 > 3 {
		u1, _ := url.Parse("https://huodong.dlut.edu.cn/")
		ans = append(ans, widget.NewHyperlink("大连理工大学校内活动网", u1))
	}

	return ans
}

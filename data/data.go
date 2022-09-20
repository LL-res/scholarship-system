package data

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	classNum    = 5
	std         = 10.0
	mean        = 75.0
	studentNum  = 500
	password    = "88888888"
	baseId      = 20174001
	FirstClass  = 100
	SecondClass = 200
	ThirdClass  = 200
	FtoS        = 20
	StoT        = 40
)

var (
	DataBase map[int]*Student
)

type Student struct {
	Id         int
	LastLevel  int
	Password   string
	Name       string
	Compulsory []int
	Selective  []int
}

func NewStudent(id int, g func() int) *Student {
	rand.Seed(time.Now().UnixNano())
	compulsory := make([]int, classNum)
	selective := make([]int, classNum)
	for i := range compulsory {
		compulsory[i] = int(rand.NormFloat64()*std + mean)
		for compulsory[i] >= 100 || compulsory[i] < 0 {
			compulsory[i] = int(rand.NormFloat64()*std + mean)
		}
	}
	for i := range selective {
		selective[i] = int(rand.NormFloat64()*std + mean)
		for selective[i] >= 100 || selective[i] < 0 {
			selective[i] = int(rand.NormFloat64()*std + mean)
		}
	}
	f, err := os.Open(".\\data\\name.txt")
	if err != nil {
		log.Println("can not open the name file :", err)
		return nil
	}
	namesBytes, _ := io.ReadAll(f)
	names := string(namesBytes)
	names = strings.TrimSpace(names)
	nameRunes := []rune(names)
	return &Student{
		Id:         id,
		Name:       string(nameRunes[(id-baseId)*3 : (id-baseId)*3+3]),
		LastLevel:  g(),
		Compulsory: compulsory,
		Selective:  selective,
		Password:   password,
	}
}

func NewStudents() {
	DataBase = make(map[int]*Student)
	g := NewLevelGenerator(FirstClass, SecondClass, ThirdClass)
	for i := 0; i < studentNum; i++ {
		DataBase[baseId+i] = NewStudent(baseId+i, g)
	}
}
func (s *Student) GetStringCompulsory() []string {
	ans := []string{}
	for _, v := range s.Compulsory {
		ans = append(ans, strconv.Itoa(v))
	}
	return ans
}
func (s *Student) GetStringSelective() []string {
	ans := []string{}
	for _, v := range s.Selective {
		ans = append(ans, strconv.Itoa(v))
	}
	return ans
}
func (s *Student) GetFinalScore() string {
	ss := make([]int, 0)
	cc := make([]int, 0)
	minCompulsory := s.Compulsory[0]
	minSelective := s.Selective[0]
	for i := range s.Selective {
		minCompulsory = min(s.Compulsory[i], minCompulsory)
		minSelective = min(s.Selective[i], minSelective)
	}
	for i := range s.Selective {
		if s.Selective[i] == minSelective {
			ss = append(s.Selective[:i], s.Selective[i+1:]...)
			break
		}
	}
	for i := range s.Compulsory {
		if s.Compulsory[i] == minCompulsory {
			cc = append(s.Compulsory[:i], s.Selective[i+1:]...)
			break
		}
	}
	ans := 0.0
	for i := range ss {
		ans += float64(3*cc[i]) / float64(20)
		ans += float64(ss[i]) / float64(10)
	}
	return fmt.Sprintf("%.1f", ans*0.1)

}
func min(x, y int) int {
	if x > y {
		return y
	} else {
		return x
	}
}
func NewLevelGenerator(m1, m2, m3 int) func() int {
	if m1+m2+m3 != studentNum {
		log.Fatalln("wrong parameters")
	}
	var (
		first  int
		second int
		third  int
	)
	return func() int {
		level := rand.Intn(3) + 1
		for level == 1 && first == m1 || level == 2 && second == m2 || level == 3 && third == m3 {
			level = rand.Intn(3) + 1
		}
		switch level {
		case 1:
			first++
		case 2:
			second++
		case 3:
			third++
		}
		return level
	}
}
func Count() {
	var m1, m2, m3 int
	for _, v := range DataBase {
		switch v.LastLevel {
		case 1:
			m1++
		case 2:
			m2++
		case 3:
			m3++
		}
	}
	println(m1, m2, m3)
}
func (s *Student) GetShowScore() float64 {
	fs := s.GetFinalScore()
	ffs, _ := strconv.ParseFloat(fs, 64)
	ffs = ffs * 0.6
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 3; i++ {
		switch i {
		case 0:
			ffs += 0.3 * float64(2*rand.Intn(5)+1)
		default:
			ffs += 0.05 * float64(2*rand.Intn(5)+1)
		}
	}
	return ffs
}

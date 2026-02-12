package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
)

type ProjectWork struct {
	name       string
	nameOfWork string
	date       time.Time
}

func main() {

	fmt.Println("========Обработка Тем Работ студентов========")
	fmt.Print("Введите строку в формате: имя студента(кавычки), название темы(кавычки), дата выдачи(гггг.мм.дд)\n")
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	//fmt.Println(input)
	split := regexp.MustCompile(`"[^"]*"|\d{4}\.\d{2}\.\d{2}`).FindAllString(input, -1)
	for i := 0; i < len(split); i++ {
		fmt.Println(split[i])
	}
	date, _ := time.Parse("2006.01.02", split[2])

	work := ProjectWork{
		name:       split[0],
		nameOfWork: split[1],
		date:       date,
	}
	fmt.Println(work.name + " " + work.nameOfWork + " " + work.date.Format("2006.01.02"))
	fmt.Println(work)
}

//"Новокрещенов Александр Денисович", "Инфобез", 2025.02.12

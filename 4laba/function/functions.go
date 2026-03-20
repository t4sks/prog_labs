package function

import (
	"auiapp/model"
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var lineRegex = regexp.MustCompile(`^"([^"]+)"\s+"([^"]+)"\s+(\d{4}\.\d{2}\.\d{2})\s+([a-f])$`)

func ParsingFile(data []byte) []model.ProjectWork {

	var works []model.ProjectWork

	scanner := bufio.NewScanner(bytes.NewReader(data))

	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		work := ParcingLine(line)

		if work.Name != "" {
			works = append(works, work)
		}
	}

	return works
}

func ParcingLine(line string) model.ProjectWork {

	split := lineRegex.FindStringSubmatch(line)

	res := model.ProjectWork{}

	if len(split) == 5 {
		res.Name = split[1]
		res.NameOfWork = split[2]
		res.Date, _ = time.Parse("2006.01.02", split[3])
		res.Type = split[4]
	} else {
		log.Println("Некорректная строка:", line)
		return model.ProjectWork{}
	}
	return res
}

func ObjectTobytes(worksToLine []model.ProjectWork) []byte {
	var b bytes.Buffer
	for _, work := range worksToLine {
		line := fmt.Sprintf(
			"\"%s\" \"%s\" %s %s",
			work.Name,
			work.NameOfWork,
			work.Date.Format("2006.01.02"),
			work.Type,
		)
		parsed := ParcingLine(line)
		if parsed.Name == "" {
			continue
		}
		fmt.Fprintln(&b, line)
	}
	return b.Bytes()
}

func ReadExecCommandFile(works []model.ProjectWork, commandData []byte) []model.ProjectWork {
	scanner := bufio.NewScanner(bytes.NewReader(commandData))
	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "ADD") {
			args := strings.TrimSpace(strings.TrimPrefix(line, "ADD"))
			works = addCommand(works, args)
		} else if strings.HasPrefix(line, "REM") {
			args := strings.TrimSpace(strings.TrimPrefix(line, "REM"))
			works = remCommand(works, args)
		} else if strings.HasPrefix(line, "SAVE") {
			args := strings.TrimSpace(strings.TrimPrefix(line, "SAVE"))
			saveCommand(works, args)
		}

	}
	return works
}

func addCommand(works []model.ProjectWork, arguments string) []model.ProjectWork {
	args := strings.Split(arguments, ";")
	if len(args) != 4 {
		return works
	}
	for i := range args {
		args[i] = strings.TrimSpace(args[i])
	}
	date, _ := time.Parse("2006.01.02", args[2])

	newWork := model.ProjectWork{
		Name:       args[0],
		NameOfWork: args[1],
		Date:       date,
		Type:       args[3],
	}
	works = append(works, newWork)
	return works
}

func remCommand(works []model.ProjectWork, arguments string) []model.ProjectWork {
	args := strings.Fields(arguments)
	if len(args) != 3 {
		return works
	}
	field := args[0]
	condition := args[1]
	value := strings.Join(args[2:], " ")
	result := make([]model.ProjectWork, 0)

	for _, work := range works {
		if !checkCondition(work, field, condition, value) {
			result = append(result, work)
		}
	}

	return result
}

func checkCondition(work model.ProjectWork, field, condition, value string) bool {
	switch field {
	case "Name":
		return condition == "=" && value == work.Name
	case "NameOfWork":
		return condition == "=" && value == work.NameOfWork
	case "Date":
		date, _ := time.Parse("2006.01.02", value)
		return condition == "=" && work.Date.Equal(date)
	case "Type":
		return condition == "=" && value == work.Type
	}
	return false
}

func saveCommand(works []model.ProjectWork, filename string) {
	file, _ := os.Create(filename)
	defer file.Close()
	for _, work := range works {
		fmt.Fprintf(file, "\"%s\" \"%s\" %s %s\n",
			work.Name,
			work.NameOfWork,
			work.Date.Format("2006.01.02"),
			work.Type)
	}
}

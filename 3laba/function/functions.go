package function

import (
	"auiapp/model"
	"bufio"
	"bytes"
	"log"
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

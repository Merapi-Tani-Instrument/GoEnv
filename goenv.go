package goenv

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strings"
)

const (
	TAB_SIZE = 2
)

type GoEnv map[string]string

func pop(alist *[]string) {
	f := len(*alist)
	*alist = (*alist)[:f-1]
}

func Parse(fs fs.File) (map[string]string, error) {
	res := make(map[string]string)
	s, err := fs.Stat()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(fs)
	if s.Size() >= bufio.MaxScanTokenSize {
		maxSize := int(s.Size())
		buf := make([]byte, maxSize)
		scanner.Buffer(buf, maxSize)
	}

	var output strings.Builder
	var prop []string
	for scanner.Scan() {
		line := scanner.Text()
		ignore := regexp.MustCompile(`^\s*#`)
		arrayLine := regexp.MustCompile(`^\s*-\s`).FindString(line)

		arrayIndex := 0
		if arrayLine == "" {
			arrayIndex = -1
		}
		if ignore.FindString(line) != "" || !strings.Contains(line, ":") {
			if output.Len() > 0 {
				output.WriteString("\n")
			}
			fmt.Println("ignore: ", ignore)
			continue
		}
		tabReg := ""
		for i := 0; i < TAB_SIZE; i++ {
			tabReg += "\\s"
		}
		index := 0
		if ll := regexp.MustCompile(fmt.Sprintf("(%s)", tabReg)).FindAllString(strings.Split(line, ":")[0], -1); ll != nil {
			index = len(ll)
		}
		resultProps := regexp.MustCompile(`(.+)(?:.)`).FindStringSubmatch(line)
		if len(resultProps) < 2 {
			continue
		}

		if index == 0 {
			prop = nil
			prop = append(prop, resultProps[1])
		} else {
			propName := ""
			if arrayLine != "" {
				arrayIndex++
			} else {
				propName = strings.TrimSpace(resultProps[1])
			}
			for arrayIndex < 0 && prop != nil && index < len(prop) {
				pop(&prop)
			}
			if arrayIndex < 0 {
				prop = append(prop, propName)
			}
		}
		var value []string = nil
		if arrayIndex < 0 {
			value = regexp.MustCompile(`:(.+)`).FindStringSubmatch(line)
		} else {
			value = regexp.MustCompile(`(?<=-\s).+`).FindStringSubmatch(line)
		}
		if value != nil && len(value) > 1 {
			idx := ""
			if arrayIndex >= 0 {
				idx = fmt.Sprintf("[%d]", arrayIndex)
			}
			res[fmt.Sprintf("%s%s", strings.Join(prop, "."), idx)] = strings.TrimSpace(value[1])
		}
	}
	fmt.Println("res ", res)
	return res, nil
}

func (g GoEnv) SetAsEnvironment() {
	for k, v := range g {
		if _, ok := os.LookupEnv(k); !ok {
			os.Setenv(k, v)
		}
	}
}

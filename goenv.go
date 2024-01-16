package goenv

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"reflect"
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

func Parse(fs fs.File) (GoEnv, error) {
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
			continue
		}
		tabReg := ""
		for i := 0; i < TAB_SIZE; i++ {
			tabReg += `\s`
		}
		index := 0
		if ll := regexp.MustCompile(fmt.Sprintf("(%s)", tabReg)).FindAllString(strings.Split(line, ":")[0], -1); ll != nil {
			index = len(ll)
		}
		resultProps := regexp.MustCompile(`(.+):[\s.]?`).FindStringSubmatch(line)
		if len(resultProps) < 2 {
			continue
		}
		if index == 0 {
			prop = nil
			prop = append(prop, strings.TrimSpace(resultProps[1]))
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
			key := fmt.Sprintf("%s%s", strings.Join(prop, "."), idx)
			if v, ok := os.LookupEnv(key); ok {
				res[key] = v
			} else {
				res[key] = strings.TrimSpace(value[1])
			}
		}
	}
	return GoEnv(res), nil
}

func (g GoEnv) SetAsEnvironment() {
	for k, v := range g {
		if _, ok := os.LookupEnv(k); !ok {
			os.Setenv(k, v)
		}
	}
}

func (g GoEnv) Inject(param ...any) {
	for _, p := range param {
		t := reflect.TypeOf(p)
		if t.Kind() != reflect.Pointer {
			continue
		}
		pE := t.Elem()
		if pE.Kind() != reflect.Struct {
			continue
		}
		rV := reflect.ValueOf(p).Elem()
		pnf := pE.NumField()
		for i := 0; i < pnf; i++ {
			sF := pE.Field(i)
			pT, ok := sF.Tag.Lookup("env")
			if !ok {
				continue
			}
			if sF.Type.Kind() == reflect.String {
				if eV, ok := g[pT]; ok {
					rV.Field(i).SetString(eV)
				}
			} else if sF.Type.Kind() == reflect.Bool {
				if eV, ok := g[pT]; ok {
					rV.Field(i).SetBool(eV == "true")
				}
			}
		}
	}
}

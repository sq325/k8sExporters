package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// RunCmd run cmd and transfer output to bufio.Scanner
// Return: bufio.Scanner and if output is empty
func RunCmd(cmd *exec.Cmd) (*bufio.Scanner, bool) {
	os.Setenv("LANG", "C")
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil { // 开始执行cmd
		fmt.Println(err)
	}
	var buf1, buf2 bytes.Buffer
	buf := io.MultiWriter(&buf1, &buf2)
	io.Copy(buf, stdout)
	content, _ := io.ReadAll(&buf2)
	return bufio.NewScanner(&buf1), strings.Count(string(content), "\n") < 1
}

func IsEleInStringSlice(strSlice []string, ele string) bool {
	if len(strSlice) == 0 {
		return false
	}
	tmp := make(map[string]struct{})
	for _, e := range strSlice {
		tmp[e] = struct{}{}
	}
	_, isIn := tmp[ele]
	return isIn
}

func GetTextFCmd(cmd *exec.Cmd) string {
	var textBuild strings.Builder
	fmt.Println(cmd.String())
	scanner, _ := RunCmd(cmd)
	for scanner.Scan() {
		line := scanner.Text()
		textBuild.WriteString(line + "\n")
	}
	return textBuild.String()
}

func Utf8ToGbk(str string) string {
	reader := transform.NewReader(strings.NewReader(str), simplifiedchinese.GBK.NewEncoder())
	bys, _ := io.ReadAll(reader)
	return string(bys)
}

func JsonStrToMap(jsonStr string) (map[string]string, error) {
	m := make(map[string]string)
	if jsonStr == "" {
		return nil, errors.New("jsonStr is empty string")
	}

	// fix {...key:value ,} condition
	pattern_endwithComma := regexp.MustCompile(`\s*,\s*\}`)
	if pattern_endwithComma.MatchString(jsonStr) {
		jsonStr = pattern_endwithComma.ReplaceAllString(jsonStr, "}")
	}

	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return nil, fmt.Errorf("JsonStrToMap Err: %s.\n%s", jsonStr, err.Error())
	}
	return m, nil
}

func MapToStr(m map[string]string) (string, error) {
	if m == nil {
		return "", nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	str := string(b[:])
	return str, nil
}

package helper

import (
	"fmt"
	"strconv"
	"strings"
)

type codeFrameLevel = uint8

const (
	codeFrameWarn codeFrameLevel = iota + 1
	codeFrameError
)

// 生成 n 个长度的空白字符
func genIndentStr(prefixWidth int, width int, lastLine string) string {
	lastLineChars := []rune(lastLine)
	str := strings.Builder{}

	i := 0
	for i < prefixWidth {
		str.WriteByte(' ')
		i++
	}

	i = 0
	for i < width {
		if i < len(lastLineChars) && lastLineChars[i] == '\t' {
			str.WriteByte('\t')
		} else {
			str.WriteByte(' ')
		}
		i++
	}
	return str.String()
}

// 返回无符号整数的位数
func getUintWidth(value int) int {
	width := 1
	n := value
	for n >= 10 {
		width++
		n /= 10
	}
	return width
}

// 返回固定宽度的数字字符
func getFixedWidthStr(value string, width int, padChar byte) string {
	result := strings.Builder{}
	gap := width - len(value)
	if gap >= 0 {
		for i := 0; i < gap; i++ {
			result.WriteByte(padChar)
		}
		result.WriteString(value)
	} else {
		gap = -gap
		chars := []rune(value)
		for i, ch := range chars {
			if i >= gap {
				result.WriteString(string(ch))
			}
		}
	}
	return result.String()
}

// 返回在源代码中的行列信息
func getSourcePosition(source *string, index int) (line int, column int) {
	line = 1
	column = 1
	for i, ch := range *source {
		if i == index {
			return
		}
		if ch == '\n' {
			line++
			column = 1
		} else {
			column++
		}
	}
	return
}

// 打印代码帧信息，返回目标位置的行列信息
func printCodeFrame(source []rune, pos int, message string, level codeFrameLevel) (targetLine int, targetColumn int) {
	input := string(source)
	beforeLines := make([]string, 0, 3)
	afterLines := make([]string, 0, 3)

	// 分割提示信息的前后代码片段（打印目标位置，上面3行，下面2行）
	lines := strings.Split(input, "\n")
	targetLine, targetColumn = getSourcePosition(&input, pos)

	min := targetLine - 3
	max := targetLine + 2
	if min < 0 {
		min = 0
	}
	if max > len(lines) {
		max = len(lines)
	}
	for i := min; i < max; i++ {
		if i < targetLine {
			beforeLines = append(beforeLines, lines[i])
		} else {
			afterLines = append(afterLines, lines[i])
		}
	}

	// 行号数字长度 (最大行号为目标位置加上下面的2行)
	lineNoWidth := getUintWidth(targetLine + 2)
	lineNo := min

	// 打印提示信息前面代码
	for _, rawLine := range beforeLines {
		lineNo++
		formatLineNo := getFixedWidthStr(strconv.Itoa(lineNo), lineNoWidth, ' ')

		head := fmt.Sprintf("%s | ", formatLineNo)
		fmt.Print(head)
		fmt.Println(rawLine)
	}

	// 打印提示信息（需预留行号空白位置）
	formatMsg := strings.Builder{}
	lastLine := ""
	if len(beforeLines) > 0 {
		lastLine = beforeLines[len(beforeLines)-1]
	}
	formatMsg.WriteString(genIndentStr(lineNoWidth+3, targetColumn-1, lastLine))
	formatMsg.WriteString("^ ")
	formatMsg.WriteString(message)
	fmt.Println(formatMsg.String())

	// 打印提示信息后面代码
	for _, rawLine := range afterLines {
		lineNo++
		formatLineNo := getFixedWidthStr(strconv.Itoa(lineNo), lineNoWidth, ' ')

		head := fmt.Sprintf("%s | ", formatLineNo)
		fmt.Print(head)
		fmt.Println(rawLine)
	}

	return
}

// PrintWarnFrame 打印警告代码帧信息
func PrintWarnFrame(source []rune, pos int, message string) (int, int) {
	return printCodeFrame(source, pos, message, codeFrameWarn)
}

// PrintErrorFrame 打印错误代码帧信息
func PrintErrorFrame(source []rune, pos int, message string) (int, int) {
	return printCodeFrame(source, pos, message, codeFrameError)
}

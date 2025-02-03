package main

import (
	"Assembler/code"
	"Assembler/parser"
	"Assembler/symbol"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	// 커맨드라인 인자 확인
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <input.asm>")
	}

	filename := os.Args[1]

	// 파일 열기
	asmFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer asmFile.Close()

	// 1. .asm 확장자 → .hack으로 교체
	ext := filepath.Ext(filename) // .asm 추출
	base := strings.TrimSuffix(filename, ext)
	outputFile := base + ".hack"

	hackFile, err := os.Create(outputFile)
	if err != nil {
		log.Fatal("파일 생성 실패:", err)
	}
	defer hackFile.Close() // 함수 종료 시 파일 닫기

	// 1st loop, set label to address into symbol table
	st := scanSymbol(asmFile)
	asmFile.Seek(0, io.SeekStart) // 파일 위치 초기화

	// 1st loop, translate line by line into binary code
	lineNumber := 0
	scanner := bufio.NewScanner(asmFile)
	for scanner.Scan() {
		line := scanner.Text()

		preprocessed := preprocess(line) // 주석/공백 제거
		if preprocessed == "" {
			continue
		}

		var instType parser.InstructionType
		instType, _ = parser.ParseLine(preprocessed)
		// fmt.Printf("%s -> %s: ", preprocessed, instType.String())

		var binCode string
		switch instType {
		case parser.A_INSTRUCTION:
			remaining := preprocessed[1:]
			// 숫자인지 확인
			if num, err := strconv.Atoi(remaining); err == nil {
				// 숫자이면 16bit binary code로 변환
				binCode, _ = code.TranslateAInstruction(num)
			} else {
				// 숫자가 아니면 심볼 테이블에서 찾기
				if st.Contains(remaining) {
					address := st.GetAddress(remaining)
					binCode, _ = code.TranslateAInstruction(address)
				} else {
					// 심볼 테이블에 없으면 에러 처리 (또는 새로운 변수 할당)
					fmt.Printf("Error: Symbol '%s' not found.\n", remaining)
				}
			}
		case parser.C_INSTRUCTION:
			comp, dest, jump, err := parser.ParseCComponents(preprocessed)
			if err != nil {
				fmt.Printf("Error parsing C-instruction: %v\n", err)
				continue
			}
			binCode, _ = code.TranslateCInstruction(comp, dest, jump)
		case parser.LABEL:
			fmt.Println("LABEL")
			continue
		case parser.UNKNOWN:
			fmt.Println("UNKNOWN")
			continue
		}

		// 2. 각 라인 쓰기
		fmt.Fprintln(hackFile, binCode) // 코드 + 줄바꿈
		fmt.Printf("%04d: %-15s %s\n", lineNumber, preprocessed, binCode)
		lineNumber++
	}

	// 스캐너 에러 체크
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error during file scan: %v", err)
	}
}

func preprocess(rawLine string) string {
	// 주석 제거
	if i := strings.Index(rawLine, "//"); i != -1 {
		rawLine = rawLine[:i]
	}
	// 공백 제거
	return strings.TrimSpace(rawLine)
}

func scanSymbol(inputFile *os.File) *symbol.SymbolTable {
	st := symbol.NewSymbolTable()
	romAddr := 0
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()

		preprocessed := preprocess(line) // 주석/공백 제거
		if preprocessed == "" {
			continue
		}

		switch parser.CommandType(preprocessed) {
		case parser.A_INSTRUCTION, parser.C_INSTRUCTION:
			romAddr++
		case parser.LABEL:
			labelName := parser.GetSymbol(preprocessed)
			fmt.Println(labelName)
			st.AddEntry(labelName, romAddr)
		}
	}

	return st
}

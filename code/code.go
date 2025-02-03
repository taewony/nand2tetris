package code

import "fmt"

// A-Instruction → 16비트 2진수 변환
func TranslateAInstruction(value int) (string, error) {
	if value < 0 || value > 32767 {
		return "", fmt.Errorf("주소 범위 초과: %d (0~32767 허용)", value)
	}
	return fmt.Sprintf("0%015b", value), nil
}

// C-Instruction → 16비트 2진수 변환
func TranslateCInstruction(comp, dest, jump string) (string, error) {
	// Validate comp field
	compCode, ok := compTable[comp]
	if !ok {
		return "", fmt.Errorf("잘못된 comp 필드: %q", comp)
	}

	// Validate dest field
	destCode, ok := destTable[dest]
	if !ok {
		return "", fmt.Errorf("잘못된 dest 필드: %q", dest)
	}

	// Validate jump field
	jumpCode, ok := jumpTable[jump]
	if !ok {
		return "", fmt.Errorf("잘못된 jump 필드: %q", jump)
	}

	// Additional validation for empty fields
	if comp == "" {
		return "", fmt.Errorf("comp 필드는 비워둘 수 없습니다")
	}

	return "111" + compCode + destCode + jumpCode, nil
}

var compTable = map[string]string{
	// a=0 (ALU 연산)
	"0":   "0101010", // 0
	"1":   "0111111", // 1
	"-1":  "0111010", // -1
	"D":   "0001100", // D
	"A":   "0110000", // A
	"!D":  "0001101", // !D
	"!A":  "0110001", // !A
	"-D":  "0001111", // -D
	"-A":  "0110011", // -A
	"D+1": "0011111", // D+1
	"A+1": "0110111", // A+1
	"D-1": "0001110", // D-1
	"A-1": "0110010", // A-1
	"D+A": "0000010", // D+A
	"D-A": "0010011", // D-A
	"A-D": "0000111", // A-D
	"D&A": "0000000", // D&A
	"D|A": "0010101", // D|A

	// a=1 (M 연산, A 대신 M 사용)
	"M":   "1110000", // M
	"!M":  "1110001", // !M
	"-M":  "1110011", // -M
	"M+1": "1110111", // M+1
	"M-1": "1110010", // M-1
	"D+M": "1000010", // D+M
	"D-M": "1010011", // D-M
	"M-D": "1000111", // M-D
	"D&M": "1000000", // D&M
	"D|M": "1010101", // D|M
}

var destTable = map[string]string{
	"":    "000", // 값 저장 안 함
	"M":   "001", // M 레지스터
	"D":   "010", // D 레지스터
	"MD":  "011", // M과 D
	"A":   "100", // A 레지스터
	"AM":  "101", // A와 M
	"AD":  "110", // A와 D
	"AMD": "111", // A, M, D
}

var jumpTable = map[string]string{
	"":    "000", // 점프 없음
	"JGT": "001", // > 0
	"JEQ": "010", // == 0
	"JGE": "011", // >= 0
	"JLT": "100", // < 0
	"JNE": "101", // != 0
	"JLE": "110", // <= 0
	"JMP": "111", // 무조건 점프
}

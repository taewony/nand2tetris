package parser

import (
	"fmt"
	"strings"
)

// InstructionType - 어셈블리 명령어 유형을 나타내는 사용자 정의 타입
type InstructionType int

const (
	A_INSTRUCTION InstructionType = iota // @value 형태 (예: @42)
	C_INSTRUCTION                        // dest=comp;jump 형태 (예: D=M+1;JGT)
	LABEL                                // (LOOP) 형태의 심볼 정의
	UNKNOWN                              // 알 수 없는 명령어 (에러 처리용)
)

// String() 메서드 추가 (디버깅용)
func (it InstructionType) String() string {
	return [...]string{"A_INSTRUCTION", "C_INSTRUCTION", "LABEL", "UNKNOWN"}[it]
}

func CommandType(text string) InstructionType {
	if text[0] == '@' {
		// @Xxx
		return A_INSTRUCTION
	}
	if text[0] == '(' {
		// (LABEL)
		return LABEL
	}
	// dest=comp;jump
	// p.analyze()
	return C_INSTRUCTION
}

func GetSymbol(text string) string {
	if strings.HasPrefix(text, "@") {
		return strings.TrimPrefix(text, "@")
	} else {
		return strings.Trim(text, "()")
	}
}

// ParseLine은 어셈블리 라인을 파싱하는 함수 (골격)
func ParseLine(line string) (InstructionType, error) {
	var instType InstructionType
	switch {
	case strings.HasPrefix(line, "@"):
		instType = A_INSTRUCTION
	case strings.HasPrefix(line, "("):
		instType = LABEL
	default:
		instType = C_INSTRUCTION
	}

	return instType, nil
}

// C-Instruction 컴포넌트 분해 (예: "D=M+1;JGT" → comp="M+1", dest="D", jump="JGT")
func ParseCComponents(line string) (comp, dest, jump string, err error) {
	if line == "" {
		return "", "", "", fmt.Errorf("empty instruction")
	}

	// dest 부분 추출 ( "=" 가 있는 경우)
	if strings.Contains(line, "=") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return "", "", "", fmt.Errorf("invalid destination format")
		}
		dest = parts[0]
		if dest == "" {
			return "", "", "", fmt.Errorf("empty destination")
		}
		line = parts[1]
	}

	// jump 부분 추출 ( ";" 가 있는 경우)
	if strings.Contains(line, ";") {
		parts := strings.SplitN(line, ";", 2)
		if len(parts) != 2 {
			return "", "", "", fmt.Errorf("invalid jump format")
		}
		comp = parts[0]
		jump = parts[1]
		if jump == "" {
			return "", "", "", fmt.Errorf("empty jump field")
		}
	} else {
		comp = line
	}

	if comp == "" {
		return "", "", "", fmt.Errorf("empty computation field")
	}

	return comp, dest, jump, nil
}

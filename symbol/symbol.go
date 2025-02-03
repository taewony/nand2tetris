package symbol

// SymbolTable - 심볼 테이블 구조체
type SymbolTable struct {
	table                map[string]int // 심볼과 주소를 매핑하는 맵
	nextAvailableAddress int            // 다음 사용 가능한 변수 주소
}

// AddEntry - 심볼과 주소를 테이블에 추가
func (st *SymbolTable) AddEntry(symbol string, address int) {
	st.table[symbol] = address
}

// Contains - 심볼이 테이블에 있는지 확인
func (st *SymbolTable) Contains(symbol string) bool {
	_, exists := st.table[symbol]
	return exists
}

// GetAddress - 심볼에 대응하는 주소 반환
func (st *SymbolTable) GetAddress(symbol string) int {
	return st.table[symbol]
}

// AddVariable - 새로운 변수를 테이블에 추가 (주소 자동 할당)
func (st *SymbolTable) AddVariable(symbol string) int {
	if !st.Contains(symbol) {
		st.AddEntry(symbol, st.nextAvailableAddress)
		st.nextAvailableAddress++
	}
	return st.GetAddress(symbol)
}

func NewSymbolTable() *SymbolTable {
	st := &SymbolTable{
		table: map[string]int{
			"SP":     0,
			"LCL":    1,
			"ARG":    2,
			"THIS":   3,
			"THAT":   4,
			"R0":     0,
			"R1":     1,
			"R2":     2,
			"R3":     3,
			"R4":     4,
			"R5":     5,
			"R6":     6,
			"R7":     7,
			"R8":     8,
			"R9":     9,
			"R10":    10,
			"R11":    11,
			"R12":    12,
			"R13":    13,
			"R14":    14,
			"R15":    15,
			"SCREEN": 16384,
			"KBD":    24576,
		},
		nextAvailableAddress: 16,
	}
	return st
}

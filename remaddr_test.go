package parser

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRemaddr(t *testing.T) {
	input := `ADDR ASSOC_ID HB_ACT RTO MAX_PATH_RTX REM_ADDR_RTX START STATE
127.0.0.10  69 1 1000 5 0 0 2
127.0.0.20  69 1 3000 5 0 0 3
127.0.0.1  68 1 1000 5 0 0 2
127.0.0.2  68 1 3000 5 0 0 2
`
	eps, err := ParseRemaddr(bufio.NewScanner(strings.NewReader(input)))
	assert.NoError(t, err)
	assert.EqualValues(t, &Remaddr{
		Addr:       "127.0.0.10",
		AssocID:    69,
		HbAct:      1,
		RTO:        1000,
		MaxPathRtx: 5,
		RemAddrRtx: 0,
		Start:      0,
		State:      2,
	}, eps[0])
	assert.EqualValues(t, &Remaddr{
		Addr:       "127.0.0.20",
		AssocID:    69,
		HbAct:      1,
		RTO:        3000,
		MaxPathRtx: 5,
		RemAddrRtx: 0,
		Start:      0,
		State:      3,
	}, eps[1])
	assert.EqualValues(t, &Remaddr{
		Addr:       "127.0.0.1",
		AssocID:    68,
		HbAct:      1,
		RTO:        1000,
		MaxPathRtx: 5,
		RemAddrRtx: 0,
		Start:      0,
		State:      2,
	}, eps[2])
	assert.EqualValues(t, &Remaddr{
		Addr:       "127.0.0.2",
		AssocID:    68,
		HbAct:      1,
		RTO:        3000,
		MaxPathRtx: 5,
		RemAddrRtx: 0,
		Start:      0,
		State:      2,
	}, eps[3])
}

func TestParseRemaddr_WithInvalidInput(t *testing.T) {
	input := `ADDR ASSOC_ID HB_ACT RTO MAX_PATH_RTX REM_ADDR_RTX START STATE
`
	_, err := ParseRemaddr(bufio.NewScanner(strings.NewReader(input)), true)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidRemaddrFormat)
	assert.Contains(t, err.Error(), "ASSOC_ID")
}

func TestParseRemaddr_WithInsufficientItems(t *testing.T) {
	input := `ADDR ASSOC_ID HB_ACT RTO MAX_PATH_RTX REM_ADDR_RTX START STATE
127.0.0.10  69 1 1000 5 0 0
`
	_, err := ParseRemaddr(bufio.NewScanner(strings.NewReader(input)))
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInsufficientNumberOfRemaddrItems)
	assert.Contains(t, err.Error(), "at line #2")
}

func ExampleParseRemaddr() {
	input := `ADDR ASSOC_ID HB_ACT RTO MAX_PATH_RTX REM_ADDR_RTX START STATE
127.0.0.10  69 1 1000 5 0 0 2
127.0.0.20  69 1 3000 5 0 0 3
127.0.0.1  68 1 1000 5 0 0 2
127.0.0.2  68 1 3000 5 0 0 2
`
	assocs, err := ParseRemaddr(bufio.NewScanner(strings.NewReader(input)))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", assocs)
}

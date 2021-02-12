package parser

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEPS(t *testing.T) {
	input := `ENDPT     SOCK   STY SST HBKT LPORT   UID INODE LADDRS
0        0 2   10  24   12345     0 227065 127.0.0.1
0        0 2   10  16   54321     0 232851 127.0.0.3
`
	eps, err := ParseEPS(bufio.NewScanner(strings.NewReader(input)))
	assert.NoError(t, err)
	assert.EqualValues(t, &EPS{
		Endpt:  0,
		Sock:   0,
		Sty:    2,
		Sst:    10,
		Hbkt:   24,
		LPort:  12345,
		Uid:    0,
		Inode:  227065,
		LAddrs: []string{"127.0.0.1"},
	}, eps[0])
	assert.EqualValues(t, &EPS{
		Endpt:  0,
		Sock:   0,
		Sty:    2,
		Sst:    10,
		Hbkt:   16,
		LPort:  54321,
		Uid:    0,
		Inode:  232851,
		LAddrs: []string{"127.0.0.3"},
	}, eps[1])
}

func TestParseEPS_MultipleLAddress(t *testing.T) {
	input := `ENDPT     SOCK   STY SST HBKT LPORT   UID INODE LADDRS
0        0 2   10  24   12345     0 227065 127.0.0.1 127.0.0.2
0        0 2   10  16   54321     0 232851 127.0.0.3 127.0.0.4
`
	eps, err := ParseEPS(bufio.NewScanner(strings.NewReader(input)))
	assert.NoError(t, err)
	assert.EqualValues(t, &EPS{
		Endpt:  0,
		Sock:   0,
		Sty:    2,
		Sst:    10,
		Hbkt:   24,
		LPort:  12345,
		Uid:    0,
		Inode:  227065,
		LAddrs: []string{"127.0.0.1", "127.0.0.2"},
	}, eps[0])
	assert.EqualValues(t, &EPS{
		Endpt:  0,
		Sock:   0,
		Sty:    2,
		Sst:    10,
		Hbkt:   16,
		LPort:  54321,
		Uid:    0,
		Inode:  232851,
		LAddrs: []string{"127.0.0.3", "127.0.0.4"},
	}, eps[1])
}

func TestParseEPS_WithInvalidInput(t *testing.T) {
	input := `ENDPT     SOCK   STY SST HBKT LPORT   UID INODE LADDRS
`
	_, err := ParseEPS(bufio.NewScanner(strings.NewReader(input)), true)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidEPSFormat)
	assert.Contains(t, err.Error(), "ENDPT")
}

func TestParseEPS_WithInsufficientItems(t *testing.T) {
	input := `ENDPT     SOCK   STY SST HBKT LPORT   UID INODE LADDRS
0        0 2   10  24   12345     0 227065
`
	_, err := ParseEPS(bufio.NewScanner(strings.NewReader(input)))
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInsufficientNumberOfEPSItems)
	assert.Contains(t, err.Error(), "at line #2")
}

func ExampleParseEPS() {
	input := `ENDPT     SOCK   STY SST HBKT LPORT   UID INODE LADDRS
0        0 2   10  24   12345     0 227065 127.0.0.1
0        0 2   10  16   54321     0 232851 127.0.0.3
`

	assocs, err := ParseEPS(bufio.NewScanner(strings.NewReader(input)))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", assocs)
}

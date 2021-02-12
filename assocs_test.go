package parser

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseAssocs(t *testing.T) {
	input := `ASSOC     SOCK   STY SST ST HBKT ASSOC-ID TX_QUEUE RX_QUEUE UID INODE LPORT RPORT LADDRS <-> RADDRS HBINT INS OUTS MAXRT T1X T2X RTXC wmema wmemq sndbuf rcvbuf
     0        0 2   1   3  0      60        0      496       0 188897 12345 54321  127.0.0.1 <-> *127.0.0.2     30000 65535 65535   10    0    0        0        1        0   212992   212992
     0        0 2   1   3  0      59        0        0       0 189472 54321 12345  127.0.0.2 <-> *127.0.0.1     30000 65535 65535   10    0    0        0        1        0   212992   212992

`
	assocs, err := ParseAssocs(bufio.NewScanner(strings.NewReader(input)))
	assert.NoError(t, err)
	assert.EqualValues(t, &Assoc{
		Assoc:   0,
		Sock:    0,
		Sty:     2,
		Sst:     1,
		St:      3,
		Hbkt:    0,
		AssocId: 60,
		TxQueue: 0,
		RxQueue: 496,
		Uid:     0,
		Inode:   188897,
		Lport:   12345,
		Rport:   54321,
		LAddrs:  []string{"127.0.0.1"},
		RAddrs:  []string{"127.0.0.2"},
		Hbint:   30000,
		Ins:     65535,
		Outs:    65535,
		Maxrt:   10,
		T1x:     0,
		T2x:     0,
		Rtxc:    0,
		Wmema:   1,
		Wmemq:   0,
		Sndbuf:  212992,
		Rcvbuf:  212992,
	}, assocs[0])
	assert.EqualValues(t, &Assoc{
		Assoc:   0,
		Sock:    0,
		Sty:     2,
		Sst:     1,
		St:      3,
		Hbkt:    0,
		AssocId: 59,
		TxQueue: 0,
		RxQueue: 0,
		Uid:     0,
		Inode:   189472,
		Lport:   54321,
		Rport:   12345,
		LAddrs:  []string{"127.0.0.2"},
		RAddrs:  []string{"127.0.0.1"},
		Hbint:   30000,
		Ins:     65535,
		Outs:    65535,
		Maxrt:   10,
		T1x:     0,
		T2x:     0,
		Rtxc:    0,
		Wmema:   1,
		Wmemq:   0,
		Sndbuf:  212992,
		Rcvbuf:  212992,
	}, assocs[1])
}

func Test_ParseAssocsWithMultiple(t *testing.T) {
	input := `ASSOC     SOCK   STY SST ST HBKT ASSOC-ID TX_QUEUE RX_QUEUE UID INODE LPORT RPORT LADDRS <-> RADDRS HBINT INS OUTS MAXRT T1X T2X RTXC wmema wmemq sndbuf rcvbuf
	0        0 2   1   3  0      62        0        0       0 212110 54321 12345  127.0.0.10 127.0.0.20 <-> *127.0.0.1 127.0.0.2    30000 65535 65535   10    0    0        0        1        0   212992   212992
	0        0 2   1   3  0      63        0        0       0 212095 12345 54321  127.0.0.1 127.0.0.2 <-> *127.0.0.10 127.0.0.20    30000 65535 65535   10    0    0        0        1        0   212992   212992
`
	assocs, err := ParseAssocs(bufio.NewScanner(strings.NewReader(input)))
	assert.NoError(t, err)
	assert.EqualValues(t, &Assoc{
		Assoc:   0,
		Sock:    0,
		Sty:     2,
		Sst:     1,
		St:      3,
		Hbkt:    0,
		AssocId: 62,
		TxQueue: 0,
		RxQueue: 0,
		Uid:     0,
		Inode:   212110,
		Lport:   54321,
		Rport:   12345,
		LAddrs:  []string{"127.0.0.10", "127.0.0.20"},
		RAddrs:  []string{"127.0.0.1", "127.0.0.2"},
		Hbint:   30000,
		Ins:     65535,
		Outs:    65535,
		Maxrt:   10,
		T1x:     0,
		T2x:     0,
		Rtxc:    0,
		Wmema:   1,
		Wmemq:   0,
		Sndbuf:  212992, Rcvbuf: 212992}, assocs[0])
	assert.EqualValues(t, &Assoc{
		Assoc:   0,
		Sock:    0,
		Sty:     2,
		Sst:     1,
		St:      3,
		Hbkt:    0,
		AssocId: 63,
		TxQueue: 0,
		RxQueue: 0,
		Uid:     0,
		Inode:   212095,
		Lport:   12345,
		Rport:   54321,
		LAddrs:  []string{"127.0.0.1", "127.0.0.2"},
		RAddrs:  []string{"127.0.0.10", "127.0.0.20"},
		Hbint:   30000,
		Ins:     65535,
		Outs:    65535,
		Maxrt:   10,
		T1x:     0,
		T2x:     0,
		Rtxc:    0,
		Wmema:   1,
		Wmemq:   0,
		Sndbuf:  212992,
		Rcvbuf:  212992,
	}, assocs[1])
}

func Test_ParseInvalidAssocs(t *testing.T) {
	input := `ASSOC     SOCK   STY SST ST HBKT ASSOC-ID TX_QUEUE RX_QUEUE UID INODE LPORT RPORT LADDRS <-> RADDRS HBINT INS OUTS MAXRT T1X T2X RTXC wmema wmemq sndbuf rcvbuf`

	_, err := ParseAssocs(bufio.NewScanner(strings.NewReader(input)), true)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidAssocsFormat)
	assert.Contains(t, err.Error(), "ASSOC")
}

func Test_ParseInsufficientItems(t *testing.T) {
	input := `ASSOC     SOCK   STY SST ST HBKT ASSOC-ID TX_QUEUE RX_QUEUE UID INODE LPORT RPORT LADDRS <-> RADDRS HBINT INS OUTS MAXRT T1X T2X RTXC wmema wmemq sndbuf rcvbuf
     0        0 2   1   3  0      60        0      496       0 188897 12345 54321  127.0.0.1 <-> *127.0.0.2     30000 65535 65535   10    0    0        0        1        0   212992
`
	_, err := ParseAssocs(bufio.NewScanner(strings.NewReader(input)))
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInsufficientNumberOfAssocItems)
	assert.Contains(t, err.Error(), "at line #2")
}

func ExampleParseAssocs() {
	input := `ASSOC     SOCK   STY SST ST HBKT ASSOC-ID TX_QUEUE RX_QUEUE UID INODE LPORT RPORT LADDRS <-> RADDRS HBINT INS OUTS MAXRT T1X T2X RTXC wmema wmemq sndbuf rcvbuf
     0        0 2   1   3  0      60        0      496       0 188897 12345 54321  127.0.0.1 <-> *127.0.0.2     30000 65535 65535   10    0    0        0        1        0   212992   212992
     0        0 2   1   3  0      59        0        0       0 189472 54321 12345  127.0.0.2 <-> *127.0.0.1     30000 65535 65535   10    0    0        0        1        0   212992   212992
`

	assocs, err := ParseAssocs(bufio.NewScanner(strings.NewReader(input)))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", assocs)
}

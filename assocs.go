package parser

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInsufficientNumberOfAssocItems = errors.New("insufficient number of assoc items on a line")
	ErrInvalidAssocsFormat            = errors.New("invalid assocs format")
)

// Assoc represents the structure of SCTP assoc.
type Assoc struct {
	Assoc   uint64
	Sock    uint64
	Sty     int64
	Sst     int64
	St      int64
	Hbkt    int64
	AssocId int64
	TxQueue int64
	RxQueue int64
	Uid     uint64
	Inode   uint64
	LPort   int64
	RPort   int64
	LAddrs  []string
	RAddrs  []string
	Hbint   uint64
	Ins     int64
	Outs    int64
	Maxrt   int64
	T1x     int64
	T2x     int64
	Rtxc    int64
	Wmema   int64
	Wmemq   int64
	Sndbuf  int64
	Rcvbuf  int64
}

// ParseAssocs parses SCTP assocs contents; for example the contents of `/proc/net/sctp/assocs` file.
//
// - input: the contents of SCTP assocs
// - noHeader: this has to be true if the input doesn't have a header line (default: `false`)
//
// example input:
// ```
//  ASSOC     SOCK   STY SST ST HBKT ASSOC-ID TX_QUEUE RX_QUEUE UID INODE LPORT RPORT LADDRS <-> RADDRS HBINT INS OUTS MAXRT T1X T2X RTXC wmema wmemq sndbuf rcvbuf
//       0        0 2   1   3  0      60        0      496       0 188897 12345 54321  127.0.0.1 <-> *127.0.0.2     30000 65535 65535   10    0    0        0        1        0   212992   212992
//       0        0 2   1   3  0      59        0        0       0 189472 54321 12345  127.0.0.2 <-> *127.0.0.1     30000 65535 65535   10    0    0        0        1        0   212992   212992
// ```
func ParseAssocs(input *bufio.Scanner, noHeader ...bool) ([]*Assoc, error) {
	// implementation memo:
	// https://github.com/torvalds/linux/blob/dcc0b49040c70ad827a7f3d58a21b01fdb14e749/net/sctp/proc.c#L243

	lineNum := 1
	if len(noHeader) <= 0 || !noHeader[0] {
		input.Scan() // skip a header line
		lineNum++
	}

	assocs := make([]*Assoc, 0)

	for input.Scan() {
		if input.Text() == "" {
			lineNum++
			continue
		}

		leaves := spacesRe.Split(strings.TrimSpace(input.Text()), -1)
		leavesLen := len(leaves)
		if leavesLen < 27 {
			return nil, fmt.Errorf("at line #%d: %w", lineNum, ErrInsufficientNumberOfAssocItems)
		}

		assoc, err := strconv.ParseUint(leaves[0], 16, 64)
		if err != nil {
			return nil, fmt.Errorf("ASSOC at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		sock, err := strconv.ParseUint(leaves[1], 16, 64)
		if err != nil {
			return nil, fmt.Errorf("SOCK at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		sty, err := strconv.ParseInt(leaves[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("STY at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		sst, err := strconv.ParseInt(leaves[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("SST at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		st, err := strconv.ParseInt(leaves[4], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("ST at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		hbkt, err := strconv.ParseInt(leaves[5], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("HBKT at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		assocId, err := strconv.ParseInt(leaves[6], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("ASSOC-ID at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		txQueue, err := strconv.ParseInt(leaves[7], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("TX_QUEUE at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		rxQueue, err := strconv.ParseInt(leaves[8], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("RX_QUEUE at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		uid, err := strconv.ParseUint(leaves[9], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("UID at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		inode, err := strconv.ParseUint(leaves[10], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("INODE at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		lport, err := strconv.ParseInt(leaves[11], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("LPORT at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		rport, err := strconv.ParseInt(leaves[12], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("RPORT at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}

		cur := 13
		laddrs := make([]string, 0, 1)
		for {
			if cur >= leavesLen {
				return nil, fmt.Errorf("there is no separater ('<->') for laddr and raddr: %w", ErrInvalidAssocsFormat)
			}

			leaf := leaves[cur]
			cur++
			if leaf == "<->" {
				break
			}

			laddrs = append(laddrs, leaf)
		}

		endCursorForRaddrs := leavesLen - 11
		raddrs := make([]string, 0, 1)
		for {
			if cur >= endCursorForRaddrs {
				break
			}
			raddrs = append(raddrs, strings.Trim(leaves[cur], "*"))
			cur++
		}

		hbint, err := strconv.ParseUint(leaves[leavesLen-11], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("HBINT at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		ins, err := strconv.ParseInt(leaves[leavesLen-10], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("INS at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		outs, err := strconv.ParseInt(leaves[leavesLen-9], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("OUTS at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		maxrt, err := strconv.ParseInt(leaves[leavesLen-8], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("MAXRT at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		t1x, err := strconv.ParseInt(leaves[leavesLen-7], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("T1X at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		t2x, err := strconv.ParseInt(leaves[leavesLen-6], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("T2X at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		rtxc, err := strconv.ParseInt(leaves[leavesLen-5], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("RTXC at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		wmema, err := strconv.ParseInt(leaves[leavesLen-4], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("wmema at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		wmemq, err := strconv.ParseInt(leaves[leavesLen-3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("wmemq at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		sndbuf, err := strconv.ParseInt(leaves[leavesLen-2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("sndbuf at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}
		rcvbuf, err := strconv.ParseInt(leaves[leavesLen-1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("rcvbuf at line #%d: %w", lineNum, ErrInvalidAssocsFormat)
		}

		lineNum++
		assocs = append(assocs, &Assoc{
			Assoc:   assoc,
			Sock:    sock,
			Sty:     sty,
			Sst:     sst,
			St:      st,
			Hbkt:    hbkt,
			AssocId: assocId,
			TxQueue: txQueue,
			RxQueue: rxQueue,
			Uid:     uid,
			Inode:   inode,
			LPort:   lport,
			RPort:   rport,
			LAddrs:  laddrs,
			RAddrs:  raddrs,
			Hbint:   hbint,
			Ins:     ins,
			Outs:    outs,
			Maxrt:   maxrt,
			T1x:     t1x,
			T2x:     t2x,
			Rtxc:    rtxc,
			Wmema:   wmema,
			Wmemq:   wmemq,
			Sndbuf:  sndbuf,
			Rcvbuf:  rcvbuf,
		})
	}

	return assocs, nil
}

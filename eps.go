package parser

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInsufficientNumberOfEPSItems = errors.New("insufficient number of EPS items on a line")
	ErrInvalidEPSFormat             = errors.New("invalid EPS format")
)

// EPS represents the structure of SCTP EPS.
type EPS struct {
	Endpt  uint64
	Sock   uint64
	Sty    int64
	Sst    int64
	Hbkt   int64
	LPort  int64
	Uid    uint64
	Inode  uint64
	LAddrs []string
}

// ParseEPS parses SCTP EPS contents; for example the contents of `/proc/net/sctp/eps` file.
//
// - input: the contents of SCTP EPS
// - noHeader: this has to be true if the input doesn't have a header line (default: `false`)
//
// example input:
// ```
// ENDPT     SOCK   STY SST HBKT LPORT   UID INODE LADDRS
// 0        0 2   10  24   12345     0 227065 127.0.0.1
// 0        0 2   10  16   54321     0 232851 127.0.0.3
// ```
func ParseEPS(input *bufio.Scanner, noHeader ...bool) ([]*EPS, error) {
	// implementation memo:
	// https://github.com/torvalds/linux/blob/dcc0b49040c70ad827a7f3d58a21b01fdb14e749/net/sctp/proc.c#L179

	lineNum := 1
	if len(noHeader) <= 0 || !noHeader[0] {
		input.Scan() // skip a header line
		lineNum++
	}

	epses := make([]*EPS, 0)

	for input.Scan() {
		if input.Text() == "" {
			lineNum++
			continue
		}

		leaves := spacesRe.Split(strings.TrimSpace(input.Text()), -1)
		if len(leaves) < 9 {
			return nil, fmt.Errorf("at line #%d: %w", lineNum, ErrInsufficientNumberOfEPSItems)
		}

		endpt, err := strconv.ParseUint(leaves[0], 16, 64)
		if err != nil {
			return nil, fmt.Errorf("ENDPT at line #%d: %w", lineNum, ErrInvalidEPSFormat)
		}
		sock, err := strconv.ParseUint(leaves[1], 16, 64)
		if err != nil {
			return nil, fmt.Errorf("SOCK at line #%d: %w", lineNum, ErrInvalidEPSFormat)
		}
		sty, err := strconv.ParseInt(leaves[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("STY at line #%d: %w", lineNum, ErrInvalidEPSFormat)
		}
		sst, err := strconv.ParseInt(leaves[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("SST at line #%d: %w", lineNum, ErrInvalidEPSFormat)
		}
		hbkt, err := strconv.ParseInt(leaves[4], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("HBKT at line #%d: %w", lineNum, ErrInvalidEPSFormat)
		}
		lport, err := strconv.ParseInt(leaves[5], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("LPORT at line #%d: %w", lineNum, ErrInvalidEPSFormat)
		}
		uid, err := strconv.ParseUint(leaves[6], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("UID at line #%d: %w", lineNum, ErrInvalidEPSFormat)
		}
		inode, err := strconv.ParseUint(leaves[7], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("INODE at line #%d: %w", lineNum, ErrInvalidEPSFormat)
		}
		laddrs := leaves[8:]

		lineNum++
		epses = append(epses, &EPS{
			Endpt:  endpt,
			Sock:   sock,
			Sty:    sty,
			Sst:    sst,
			Hbkt:   hbkt,
			LPort:  lport,
			Uid:    uid,
			Inode:  inode,
			LAddrs: laddrs,
		})
	}

	return epses, nil
}

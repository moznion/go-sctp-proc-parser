package parser

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInsufficientNumberOfRemaddrItems = errors.New("insufficient number of remaddr items on a line")
	ErrInvalidRemaddrFormat             = errors.New("invalid Remaddr format")
)

// Remaddr represents the structure of SCTP remaddr.
type Remaddr struct {
	Addr       string
	AssocID    int64
	HbAct      int64
	RTO        uint64
	MaxPathRtx int64
	RemAddrRtx int64
	Start      int64
	State      int64
}

// ParseRemaddr parses SCTP remaddr contents; for example the contents of `/proc/net/sctp/remaddr` file.
//
// - input: the contents of SCTP remaddr
// - noHeader: this has to be true if the input doesn't have a header line (default: `false`)
//
// example input:
// ```
// ADDR ASSOC_ID HB_ACT RTO MAX_PATH_RTX REM_ADDR_RTX START STATE
// 127.0.0.10  69 1 1000 5 0 0 2
// 127.0.0.20  69 1 3000 5 0 0 3
// 127.0.0.1  68 1 1000 5 0 0 2
// 127.0.0.2  68 1 3000 5 0 0 2
// ```
func ParseRemaddr(input *bufio.Scanner, noHeader ...bool) ([]*Remaddr, error) {
	// implementation memo:
	// https://github.com/torvalds/linux/blob/dcc0b49040c70ad827a7f3d58a21b01fdb14e749/net/sctp/proc.c#L302

	lineNum := 1
	if len(noHeader) <= 0 || !noHeader[0] {
		input.Scan() // skip a header line
		lineNum++
	}

	remaddrs := make([]*Remaddr, 0)

	for input.Scan() {
		if input.Text() == "" {
			lineNum++
			continue
		}

		leaves := spacesRe.Split(strings.TrimSpace(input.Text()), -1)
		if len(leaves) < 8 {
			return nil, fmt.Errorf("at line #%d: %w", lineNum, ErrInsufficientNumberOfRemaddrItems)
		}

		addr := leaves[0]
		assocID, err := strconv.ParseInt(leaves[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("ASSOC_ID at line #%d: %w", lineNum, ErrInvalidRemaddrFormat)
		}
		hbAct, err := strconv.ParseInt(leaves[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("HB_ACT at line #%d: %w", lineNum, ErrInvalidRemaddrFormat)
		}
		rto, err := strconv.ParseUint(leaves[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("RTO at line #%d: %w", lineNum, ErrInvalidRemaddrFormat)
		}
		maxPathRtx, err := strconv.ParseInt(leaves[4], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("MAX_PATH_RTX at line #%d: %w", lineNum, ErrInvalidRemaddrFormat)
		}
		remAddrRtx, err := strconv.ParseInt(leaves[5], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("REM_ADDR_RTX at line #%d: %w", lineNum, ErrInvalidRemaddrFormat)
		}
		start, err := strconv.ParseInt(leaves[6], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("START at line #%d: %w", lineNum, ErrInvalidRemaddrFormat)
		}
		state, err := strconv.ParseInt(leaves[7], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("STATE at line #%d: %w", lineNum, ErrInvalidRemaddrFormat)
		}

		lineNum++
		remaddrs = append(remaddrs, &Remaddr{
			Addr:       addr,
			AssocID:    assocID,
			HbAct:      hbAct,
			RTO:        rto,
			MaxPathRtx: maxPathRtx,
			RemAddrRtx: remAddrRtx,
			Start:      start,
			State:      state,
		})
	}

	return remaddrs, nil
}

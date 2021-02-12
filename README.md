# go-sctp-proc-parser ![Go](https://github.com/moznion/go-sctp-proc-parser/workflows/Go/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/moznion/go-sctp-proc-parser.svg)](https://pkg.go.dev/github.com/moznion/go-sctp-proc-parser)

A parser for `/proc/net/sctp/*` files.

## Synopsis

### Parse `/proc/net/sctp/assocs`

```go
import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"github.com/moznion/go-sctp-proc-parser"
)

func main() {
	input := `ASSOC     SOCK   STY SST ST HBKT ASSOC-ID TX_QUEUE RX_QUEUE UID INODE LPORT RPORT LADDRS <-> RADDRS HBINT INS OUTS MAXRT T1X T2X RTXC wmema wmemq sndbuf rcvbuf
     0        0 2   1   3  0      60        0      496       0 188897 12345 54321  127.0.0.1 <-> *127.0.0.2     30000 65535 65535   10    0    0        0        1        0   212992   212992
     0        0 2   1   3  0      59        0        0       0 189472 54321 12345  127.0.0.2 <-> *127.0.0.1     30000 65535 65535   10    0    0        0        1        0   212992   212992

`
	assocs, err := parser.ParseAssocs(bufio.NewScanner(strings.NewReader(input)))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", assocs)
}
```


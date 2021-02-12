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

### Parse `/proc/net/sctp/eps`

```go
import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"github.com/moznion/go-sctp-proc-parser"
)

func main() {
	input := `ENDPT     SOCK   STY SST HBKT LPORT   UID INODE LADDRS
0        0 2   10  24   12345     0 227065 127.0.0.1
0        0 2   10  16   54321     0 232851 127.0.0.3
`
	assocs, err := parser.ParseEPS(bufio.NewScanner(strings.NewReader(input)))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", assocs)
```

### Parse `/proc/net/sctp/remaddr`

```go
import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"github.com/moznion/go-sctp-proc-parser"
)

func main() {
	input := `ADDR ASSOC_ID HB_ACT RTO MAX_PATH_RTX REM_ADDR_RTX START STATE
127.0.0.10  69 1 1000 5 0 0 2
127.0.0.20  69 1 3000 5 0 0 3
127.0.0.1  68 1 1000 5 0 0 2
127.0.0.2  68 1 3000 5 0 0 2
`
	assocs, err := parser.ParseRemaddr(bufio.NewScanner(strings.NewReader(input)))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", assocs)
```


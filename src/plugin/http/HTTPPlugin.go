package plugin

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/google/gopacket"
)

const (
	// Port : Default port to listen to
	Port = 80
	// Version : Plugin version
	Version = "0.1"
)

const (
	// CmdPort : Input Port from Command line
	CmdPort = "-p"
)

// HTTPPlugin :
type HTTPPlugin struct {
	port    int
	version string
}

var hp *HTTPPlugin

// NewInstance : Instantiate
func NewInstance() *HTTPPlugin {
	if hp == nil {
		hp = &HTTPPlugin{
			port:    Port,
			version: Version,
		}
	}
	return hp
}

// ResolveStream : Resolve application layer with HTTP protocol
func (m *HTTPPlugin) ResolveStream(net, transport gopacket.Flow, buf io.Reader) {
	bio := bufio.NewReader(buf)
	for {
		req, err := http.ReadRequest(bio)

		if err == io.EOF {
			return
		} else if err != nil {
			continue
		} else {
			req.ParseForm()
			msg := fmt.Sprintf("[%v] [%v%v] [%v]\n", req.Method, req.Host, req.URL.String(), req.Form.Encode())

			log.Printf(msg)

			req.Body.Close()
		}
	}
}

// BPFFilter : Setup BPF
func (m *HTTPPlugin) BPFFilter() string {
	return fmt.Sprintf("tcp and port %d", m.port)
}

// Version : Get version as string
func (m *HTTPPlugin) Version() string {
	return Version
}

// SetFlag : Set Flag for HTTP Protocol
func (m *HTTPPlugin) SetFlag(flg []string) {
	c := len(flg)

	if c == 0 {
		return
	}

	if c>>1 == 0 {
		fmt.Println("ERR : HTTP Number of parameters")
		os.Exit(1)
	}

	for i := 0; i < c; i = i + 2 {
		key := flg[i]
		val := flg[i+1]

		switch key {
		case CmdPort:
			port, err := strconv.Atoi(val)
			m.port = port
			if err != nil {
				panic("ERR : port")
			}
			if port < 0 || port > 65535 {
				panic("ERR : port(0-65535)")
			}
			break
		default:
			panic("ERR : HTTP's params")
		}
	}
}

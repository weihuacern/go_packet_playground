package core

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

// ProtocolStreamFactory : A implementation of StreamFactory
// Interface defined in : https://github.com/google/gopacket/blob/v1.1.17/tcpassembly/assembly.go#L186-L191
type ProtocolStreamFactory struct {
	pcapHandler *PCAPHandler
}

// ProtocolStream : An aux struct
type ProtocolStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
}

// New : Entrypoint to parse stream
func (m *ProtocolStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	// Init struct
	stm := &ProtocolStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}

	// Decode packet
	go m.pcapHandler.Plugin.ResolveStream(net, transport, &(stm.r))
	return &(stm.r)
}

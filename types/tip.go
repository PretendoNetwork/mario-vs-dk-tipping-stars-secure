package types

import (
	"github.com/PretendoNetwork/nex-go"
)

type Tip struct {
	DataStoreProfileID uint64
	MiiversePostID     string
	Stars              uint32
}

// ExtractFromStream extracts a Tip from a stream
func (tip *Tip) ExtractFromStream(stream *nex.StreamIn) {
	tip.DataStoreProfileID = stream.ReadUInt64BE()
	tip.MiiversePostID = string(stream.ReadBytesNext(32))
	tip.Stars = stream.ReadUInt32BE()
}

func NewTip() *Tip {
	return &Tip{}
}

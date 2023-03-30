package types

import (
	"github.com/PretendoNetwork/nex-go"
	"modernc.org/mathutil"
)

type TipBucket struct {
	Platform uint16
	Tips     []*Tip
}

// ExtractFromStream extracts a TipBucket from a stream
func (tipBucket *TipBucket) ExtractFromStream(stream *nex.StreamIn) {
	_ = stream.ReadBytesNext(4) // * Magic "DSTB"
	_ = stream.ReadBytesNext(2) // * Version (always 1)
	tipBucket.Platform = stream.ReadUInt16LE()

	tipCount := mathutil.ClampUint32(stream.ReadUInt32BE(), 0, 5)

	tipBucket.Tips = make([]*Tip, 0, tipCount)

	for i := 0; i < int(tipCount); i++ {
		tip := NewTip()
		tip.ExtractFromStream(stream)

		tipBucket.Tips = append(tipBucket.Tips, tip)
	}
}

func NewTipBucket() *TipBucket {
	return &TipBucket{}
}

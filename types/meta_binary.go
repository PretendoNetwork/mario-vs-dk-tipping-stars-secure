package types

type MetaBinary struct {
	DataID            uint32
	OwnerPID          uint32
	Name              string
	DataType          uint16
	Buffer            []byte
	Permission        uint8
	DeletePermission  uint8
	Flag              uint32
	Period            uint16
	Tags              []string
	PersistenceSlotID uint16
	ExtraData         []string
}

func NewMetaBinary() *MetaBinary {
	return &MetaBinary{}
}

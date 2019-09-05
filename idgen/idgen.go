package idgen

import (
	"fmt"
	"strconv"
	"time"
)

const (
	bitsNodeID = 8  // Up to 256 nodes supported
	bitsSeqNum = 14 // Up to 16384 unique sequence numbers per millisecond

	maxNodeID = -1 ^ (-1 << bitsNodeID)
	maxSeqNum = -1 ^ (-1 << bitsSeqNum)

	bitsToShiftTS     = bitsNodeID + bitsSeqNum
	bitsToShiftNodeID = bitsSeqNum
)

type ID uint64

func (id ID) ToBinary() string {
	return strconv.FormatInt(int64(id), 2)
}

type GenTS func() uint64

type IDGen struct {
	tsPrev     uint64
	nodeID     uint64
	seqNumPrev uint64
	genTS      GenTS
}

func genTSNow() uint64 {
	ts := time.Now().UnixNano() / 1000000 // In milliseconds
	return uint64(ts)
}

func (idGen *IDGen) getSeqNumNext(tsNext uint64) uint64 {
	seqNumNext := idGen.seqNumPrev + 1

	// Reset if max reached
	if (seqNumNext & maxSeqNum) == 0 {
		return 0
	}

	// Reset if new timestamp
	if tsNext > idGen.tsPrev {
		return 0
	}

	return seqNumNext
}

// GenID generates a 64-bit pseudo ID consisting of timestamp, node ID and sequence number.
func (idGen *IDGen) GenID() ID {
	tsNext := idGen.genTS()
	defer func() { idGen.tsPrev = tsNext }()

	seqNumNext := idGen.getSeqNumNext(tsNext)
	defer func() { idGen.seqNumPrev = seqNumNext }()

	generated := tsNext<<bitsToShiftTS | idGen.nodeID<<bitsToShiftNodeID | seqNumNext
	return ID(generated)
}

func NewIDGen(nodeID uint64) (*IDGen, error) {
	// Validate node ID
	if nodeID > maxNodeID {
		return nil, fmt.Errorf("Node ID too large: %d (max: %d)", nodeID, maxNodeID)
	}

	return &IDGen{
		tsPrev:     0,
		nodeID:     nodeID,
		seqNumPrev: 0,
		genTS:      genTSNow,
	}, nil
}

package idgen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidNodeIDs(t *testing.T) {
	nodeIDs := []uint64{0, 255} // List of valid node IDs

	for _, nodeID := range nodeIDs {
		idGen, err := NewIDGen(nodeID)
		assert.NotNil(t, idGen)
		assert.Nil(t, err)
	}
}

func TestInvalidNodeIDs(t *testing.T) {
	nodeIDs := []uint64{256, 257, 258} // List of invalid node IDs (> 255)

	for _, nodeID := range nodeIDs {
		idGen, err := NewIDGen(nodeID)
		assert.Nil(t, idGen)
		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("Node ID too large: %d (max: %d)", nodeID, 255))
	}
}

func TestSeqNumInit(t *testing.T) {
	idGen := &IDGen{
		tsPrev:     0,
		nodeID:     0,
		seqNumPrev: 0,
	}

	tsNext := genTSNow()
	assert.True(t, tsNext > 0)
	seqNumNext := idGen.getSeqNumNext(tsNext)
	assert.Equal(t, uint64(0), seqNumNext)
}

func TestSeqNumMaxReached(t *testing.T) {
	idGen := &IDGen{
		tsPrev:     0,
		nodeID:     0,
		seqNumPrev: 16383,
	}

	tsNext := uint64(0)
	seqNumNext := idGen.getSeqNumNext(tsNext)
	assert.Equal(t, uint64(0), seqNumNext)
}

func TestSeqNumIncremented(t *testing.T) {
	ts := genTSNow()
	idGen := &IDGen{
		tsPrev:     ts,
		nodeID:     0,
		seqNumPrev: 0,
	}

	tsNext := ts
	seqNumNext := idGen.getSeqNumNext(tsNext)
	assert.Equal(t, uint64(1), seqNumNext)
}

func TestIDsSameTS(t *testing.T) {
	// Timestamp unchanged forcing sequence number to increment
	var genTS GenTS = func() uint64 {
		return uint64(1)
	}

	idGen := &IDGen{
		tsPrev:     0,
		nodeID:     5,
		seqNumPrev: 0,
		genTS:      genTS,
	}

	expected := []ID{
		4276224, // 1 << (8 + 14) | 5 << 14 | 0
		4276225, // 1 << (8 + 14) | 5 << 14 | 1
		4276226, // 1 << (8 + 14) | 5 << 14 | 2
		4276227, // 1 << (8 + 14) | 5 << 14 | 3
		4276228, // 1 << (8 + 14) | 5 << 14 | 4
	}

	for _, value := range expected {
		id := idGen.GenID()
		assert.Equal(t, ID(value), id)
	}
}

func TestIDsDifferentTSs(t *testing.T) {
	// Timestamp incrementing forcing sequence number to reset
	var ts uint64 = 1
	var genTS GenTS = func() uint64 {
		defer func() { ts += 1 }()
		return ts
	}

	idGen := &IDGen{
		tsPrev:     0,
		nodeID:     5,
		seqNumPrev: 0,
		genTS:      genTS,
	}

	expected := []ID{
		4276224,  // 1 << (8 + 14) | 5 << 14 | 0
		8470528,  // 2 << (8 + 14) | 5 << 14 | 0
		12664832, // 3 << (8 + 14) | 5 << 14 | 0
		16859136, // 4 << (8 + 14) | 5 << 14 | 0
		21053440, // 5 << (8 + 14) | 5 << 14 | 0
	}

	for _, value := range expected {
		id := idGen.GenID()
		assert.Equal(t, ID(value), id)
	}
}

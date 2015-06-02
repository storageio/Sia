package consensus

import (
	"testing"

	"github.com/NebulousLabs/Sia/types"
)

// TestEarliestChildTimestamp probes the earliestChildTimestamp method of the
// block node type.
func TestEarliestChildTimestamp(t *testing.T) {
	// Check the earliest timestamp generated when the block node has no
	// parent.
	bn1 := &blockNode{block: types.Block{Timestamp: 1}}
	if bn1.earliestChildTimestamp() != 1 {
		t.Error("earliest child timestamp has been calculated incorrectly.")
	}

	// Set up a series of targets, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
	bn2 := &blockNode{block: types.Block{Timestamp: 2}, parent: bn1}
	bn3 := &blockNode{block: types.Block{Timestamp: 3}, parent: bn2}
	bn4 := &blockNode{block: types.Block{Timestamp: 4}, parent: bn3}
	bn5 := &blockNode{block: types.Block{Timestamp: 5}, parent: bn4}
	bn6 := &blockNode{block: types.Block{Timestamp: 6}, parent: bn5}
	bn7 := &blockNode{block: types.Block{Timestamp: 7}, parent: bn6}
	bn8 := &blockNode{block: types.Block{Timestamp: 8}, parent: bn7}
	bn9 := &blockNode{block: types.Block{Timestamp: 9}, parent: bn8}
	bn10 := &blockNode{block: types.Block{Timestamp: 10}, parent: bn9}
	bn11 := &blockNode{block: types.Block{Timestamp: 11}, parent: bn10}

	// Median should be '1' for bn6.
	if bn6.earliestChildTimestamp() != 1 {
		t.Error("incorrect child timestamp")
	}
	// Median should be '2' for bn7.
	if bn7.earliestChildTimestamp() != 2 {
		t.Error("incorrect child timestamp")
	}
	// Median should be '6' for bn11.
	if bn11.earliestChildTimestamp() != 6 {
		t.Error("incorrect child timestamp")
	}

	// Mix up the sorting:
	//           7, 5, 5, 2, 3, 9, 12, 1, 8, 6, 14
	// sorted11: 1, 2, 3, 5, 5, 6, 7, 8, 9, 12, 14
	// sorted10: 1, 2, 3, 5, 5, 6, 7, 7, 8, 9, 12
	// sorted9:  1, 2, 3, 5, 5, 7, 7, 7, 8, 9, 12
	bn1.block.Timestamp = 7
	bn2.block.Timestamp = 5
	bn3.block.Timestamp = 5
	bn4.block.Timestamp = 2
	bn5.block.Timestamp = 3
	bn6.block.Timestamp = 9
	bn7.block.Timestamp = 12
	bn8.block.Timestamp = 1
	bn9.block.Timestamp = 8
	bn10.block.Timestamp = 6
	bn11.block.Timestamp = 14

	// Median of bn11 should be '6'.
	if bn11.earliestChildTimestamp() != 6 {
		t.Error("incorrect child timestamp")
	}
	// Median of bn10 should be '6'.
	if bn10.earliestChildTimestamp() != 6 {
		t.Error("incorrect child timestamp")
	}
	// Median of bn9 should be '7'.
	if bn9.earliestChildTimestamp() != 7 {
		t.Error("incorrect child timestamp")
	}
}

// TestHeavierThan probes the heavierThan method of the blockNode.
func TestHeavierThan(t *testing.T) {
	// Create a light node.
	bnLight := new(blockNode)
	bnLight.depth[0] = 64
	bnLight.childTarget[0] = 200

	// Create a node that's heavier, but not enough to beat the surpass
	// threshold.
	bnMiddle := new(blockNode)
	bnMiddle.depth[0] = 60
	bnMiddle.childTarget[0] = 200

	// Create a node that's heavy enough to break the surpass threshold.
	bnHeavy := new(blockNode)
	bnHeavy.depth[0] = 16
	bnHeavy.childTarget[0] = 200

	// bnLight should not be heavier than bnHeavy.
	if bnLight.heavierThan(bnHeavy) {
		t.Error("light heavier than heavy")
	}
	// bnLight should not be heavier than middle.
	if bnLight.heavierThan(bnMiddle) {
		t.Error("light heavier than middle")
	}
	// bnLight should not be heavier than itself.
	if bnLight.heavierThan(bnLight) {
		t.Error("light heavier than itself")
	}

	// bnMiddle should not be heavier than bnLight.
	if bnMiddle.heavierThan(bnLight) {
		t.Error("middle heaver than light - surpass threshold should not have been broken")
	}
	// bnHeavy should be heaver than bnLight.
	if !bnHeavy.heavierThan(bnLight) {
		t.Error("heavy is not heavier than light")
	}
	// bnHeavy should be heavier than bnMiddle.
	if !bnHeavy.heavierThan(bnMiddle) {
		t.Error("heavy is not heavier than middle")
	}
}

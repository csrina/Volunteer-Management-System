package main

import (
	"testing"
	"time"

	_ "github.com/lib/pq"
)

/* TestBlocksIUS ...
 * Ensure block can be inserted into DB,
 * that the inserted block can be updated,
 * and that the block can be retrieved.
 */
func TestBlocksIUS(t *testing.T) {
	err := startDb()
	defer db.Close() // defer teardown

	block := TimeBlock{
		Start:    time.Now(),
		End:      time.Now(),
		Room:     1,
		Modifier: 1,
		Note:     "note"}

	block.End = block.End.Add(24000)

	// test insertion
	err = block.insertBlock()
	if err != nil {
		t.Fail()
		t.Log("Failed on insertBlock\n", err)
	}

	// test retrieval after insert
	endD := block.End
	endD.Add(50000)
	blocksGot, err := getBlocks(block.Start, endD)
	if err != nil || len(blocksGot) == 0 {
		t.Fail()
		t.Log("Failed to retrieve inserted block\n", err)
	}

	// test updating
	block.End.Add(5000)
	err = block.updateBlock()
	if err != nil {
		t.Fail()
		t.Log("Failed to update the block", err)
	}

	// test retrieval post update
	blocksGot, err = getBlocks(block.Start, endD)
	if err != nil || len(blocksGot) == 0 {
		t.Fail()
		t.Log("Failed to retrieve inserted block")
	}
	t.Log("IUS test successful")
}

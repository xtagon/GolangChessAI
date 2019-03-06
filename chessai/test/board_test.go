package test

import (
	"ChessAI3/chessai/board"
	"ChessAI3/chessai/board/color"
	"ChessAI3/chessai/board/util"
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestBoardMove(t *testing.T) {
	board2 := board.Board{}
	board2.SetPiece(util.End, &board.Rook{})
	board2.SetPiece(util.Start, &board.Rook{})
	startPiece := board2.GetPiece(util.Start)
	startPiece.SetPosition(util.End)
	board.MakeMove(&board.Move{
		Start: util.Start,
		End:   util.End,
	}, &board2)
	assert.Nil(t, board2.GetPiece(util.Start))
	assert.Equal(t, startPiece, board2.GetPiece(util.End))
	assert.Equal(t, util.End, board2.GetPiece(util.End).GetPosition())
}

func TestBoardMoveClear(t *testing.T) {
	board2 := board.Board{}
	assert.Panics(t, func() {
		for i := 0; i < 3; i++ {
			board.MakeMove(&board.Move{
				Start: util.Start,
				End:   util.End,
			}, &board2)
		}
	})
}

func TestBoardFlags(t *testing.T) {
	board2 := board.Board{}
	for i := 0; i < board.FlagRightRookMoved; i++ {
		for c := 0; c < 2; c++ {
			assert.False(t, board2.GetFlag(byte(i), byte(c)))
		}
	}
	for i := 0; i < board.FlagRightRookMoved; i++ {
		for c := 0; c < 2; c++ {
			board2.SetFlag(byte(i), byte(c), true)
			assert.True(t, board2.GetFlag(byte(i), byte(c)))
			for i2 := 0; i2 < board.FlagRightRookMoved; i2++ {
				for c2 := 0; c2 < 2; c2++ {
					if i != i2 && c != c2 {
						assert.False(t, board2.GetFlag(byte(i2), byte(c2)))
					}
				}
			}
			board2.SetFlag(byte(i), byte(c), false)
			assert.False(t, board2.GetFlag(byte(i), byte(c)))
		}
	}
}

func TestBoardSetAndCopy(t *testing.T) {
	bo1 := board.Board{}
	bo2 := board.Board{}
	bo1.ResetDefault()
	bo1.SetFlag(board.FlagCastled, color.Black, true)
	bo1.SetFlag(board.FlagRightRookMoved, color.Black, true)
	bo1.SetFlag(board.FlagRightRookMoved, color.White, true)
	bo1.SetFlag(board.FlagLeftRookMoved, color.White, true)
	assert.False(t, bo1.Equals(&bo2))
	assert.False(t, bo2.Equals(&bo1))
	bo2 = *bo1.Copy()
	assert.True(t, bo1.Equals(&bo2))
	assert.True(t, bo2.Equals(&bo1))
}

func TestBoardResetDefault(t *testing.T) {
	bo1 := board.Board{}
	bo2 := board.Board{}
	bo1.ResetDefault()
	bo2.ResetDefaultSlow()
	assert.True(t, bo1.Equals(&bo2))
	assert.True(t, bo2.Equals(&bo1))
	bo2.SetFlag(board.FlagCastled, color.Black, true)
	assert.False(t, bo1.Equals(&bo2))
	assert.False(t, bo2.Equals(&bo1))
}

func TestBoardHash(t *testing.T) {
	bo1 := board.Board{}
	bo2 := board.Board{}
	bo1.ResetDefault()
	bo1.SetFlag(board.FlagCastled, color.Black, true)
	bo1.SetFlag(board.FlagRightRookMoved, color.Black, true)
	bo1.SetFlag(board.FlagRightRookMoved, color.White, true)
	bo1.SetFlag(board.FlagLeftRookMoved, color.White, true)
	assert.False(t, reflect.DeepEqual(bo1.Hash(), bo2.Hash()))
	bo2 = *bo1.Copy()
	assert.True(t, reflect.DeepEqual(bo1.Hash(), bo2.Hash()))
}

func TestBoardHashLookupParallel(t *testing.T) {
	const (
		NumThreads = 8
		NumOps     = 1000
	)
	scoreMap := util.NewConcurrentScoreMap()

	done := make([]chan int, NumThreads)
	for tIdx := 0; tIdx < NumThreads; tIdx++ {
		done[tIdx] = make(chan int)
		go func(thread int) {
			bo1 := board.Board{}
			bo1.TestRandGen = rand.New(rand.NewSource(time.Now().UnixNano() + int64(thread)))
			numStores := 0
			for i := 0; i < NumOps; i++ {
				bo1.RandomizeIllegal()
				hash := bo1.Hash()
				r := bo1.TestRandGen.Uint32()
				_, ok := scoreMap.Read(&hash)
				if !ok {
					scoreMap.Store(&hash, r)
					score, _ := scoreMap.Read(&hash)
					assert.Equal(t, r, score)
					numStores++
				}
			}
			done[thread] <- numStores
		}(tIdx)
	}
	start := time.Now()
	totalNumStores := 0
	for tIdx := 0; tIdx < NumThreads; tIdx++ {
		totalNumStores += <-done[tIdx]
	}
	duration := time.Now().Sub(start)
	timePerOp := duration.Nanoseconds() / int64(totalNumStores)
	pSuccess := 100.0 * float64(totalNumStores) / (NumOps * NumThreads)
	log.Printf("Parallel randomize,hash,write,read %d ops with %d us/loop. %.1f%% stores successful (%d)\n",
		NumOps*NumThreads, timePerOp, pSuccess, totalNumStores)
	//scoreMap.PrintMetrics()
}

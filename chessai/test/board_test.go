package test

import (
	"ChessAI3/chessai/board"
	"ChessAI3/chessai/board/color"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"reflect"
	"testing"
)

var start = board.Location{Row: 2, Col: 5}
var end = board.Location{Row: 4, Col: 5}

func TestBoardMove(t *testing.T) {
	board2 := board.Board{}
	board2.SetPiece(end, &board.Rook{})
	board2.SetPiece(start, &board.Rook{})
	startPiece := board2.GetPiece(start)
	startPiece.SetPosition(end)
	board.MakeMove(&board.Move{
		Start: start,
		End:   end,
	}, &board2)
	assert.Nil(t, board2.GetPiece(start))
	assert.Equal(t, startPiece, board2.GetPiece(end))
	assert.Equal(t, end, board2.GetPiece(end).GetPosition())
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
	assert.False(t, bo1.Hash() == bo2.Hash())
	assert.False(t, reflect.DeepEqual(bo1.Hash(), bo2.Hash()))
	bo2 = *bo1.Copy()
	assert.True(t, bo1.Hash() == bo2.Hash())
	assert.True(t, reflect.DeepEqual(bo1.Hash(), bo2.Hash()))
}

func BenchmarkCopy(b *testing.B) {
	board2 := board.Board{}
	bNew := board2.Copy()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bNew = board2.Copy()
	}
	b.StopTimer()
	bNew.Copy()
}

func BenchmarkSetPiece(b *testing.B) {
	board2 := board.Board{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		board2.SetPiece(start, &board.Rook{})
	}
}

func BenchmarkGetPiece(b *testing.B) {
	board2 := board.Board{}
	board2.SetPiece(start, &board.Rook{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		board2.GetPiece(start)
	}
}

func BenchmarkBoardMove(b *testing.B) {
	board2 := board.Board{}
	board2.SetPiece(end, &board.Rook{})
	board2.SetPiece(start, &board.Rook{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		board.MakeMove(&board.Move{
			Start: start,
			End:   end,
		}, &board2)
	}
}

func BenchmarkBoardHash(b *testing.B) {
	bo1 := board.Board{}
	bo2 := board.Board{}
	bo1.ResetDefault()
	bo2.ResetDefaultSlow()
	b.ResetTimer()
	for i := 0; i < b.N/2; i++ {
		bo2.Hash()
		bo1.Hash()
	}
}

func BenchmarkBoardResetDefault(b *testing.B) {
	board2 := board.Board{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		board2.ResetDefault()
	}
}

func BenchmarkBoardEquals(b *testing.B) {
	bo1 := board.Board{}
	bo2 := board.Board{}
	bo1.ResetDefault()
	bo2.ResetDefaultSlow()
	b.ResetTimer()
	for i := 0; i < b.N/2; i++ {
		bo1.Equals(&bo2)
		bo2.Equals(&bo1)
	}
}

func BenchmarkBoardHashLookup(b *testing.B) {
	var scoreMap = make(map[uint64]map[uint64]map[uint64]map[uint64]map[byte]uint32)
	bo1 := board.Board{}
	bo1.ResetDefault()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash := bo1.Hash()
		idx := make([]uint64, 4)
		for x := 0; x < 32; x += 8 {
			idx[x/8] = binary.BigEndian.Uint64(hash[x : x+8])
		}

		_, ok := scoreMap[idx[0]]
		if !ok {
			scoreMap[idx[0]] = make(map[uint64]map[uint64]map[uint64]map[byte]uint32)
		}
		_, ok = scoreMap[idx[0]][idx[1]]
		if !ok {
			scoreMap[idx[0]][idx[1]] = make(map[uint64]map[uint64]map[byte]uint32)
		}
		_, ok = scoreMap[idx[0]][idx[1]][idx[2]]
		if !ok {
			scoreMap[idx[0]][idx[1]][idx[2]] = make(map[uint64]map[byte]uint32)
		}
		_, ok = scoreMap[idx[0]][idx[1]][idx[2]][idx[3]]
		if !ok {
			scoreMap[idx[0]][idx[1]][idx[2]][idx[3]] = make(map[byte]uint32)
		}

		scoreMap[idx[0]][idx[1]][idx[2]][idx[3]][hash[32]] = rand.Uint32()

		b.StopTimer()
		bo1.RandomizeIllegal()
		b.StartTimer()
	}
}

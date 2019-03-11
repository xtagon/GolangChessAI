package board

import (
	"ChessAI3/chessai/board/piece"
)

type Bishop struct {
	Location Location
	Color    byte
}

func (r *Bishop) GetChar() rune {
	return piece.BishopChar
}

func (r *Bishop) GetPieceType() byte {
	return piece.BishopType
}

func (r *Bishop) GetColor() byte {
	return r.Color
}

func (r *Bishop) SetColor(color byte) {
	r.Color = color
}

func (r *Bishop) SetPosition(loc Location) {
	r.Location.Set(loc)
}

func (r *Bishop) GetPosition() Location {
	return r.Location
}

/**
 * Explores a board using canMove, a function that determines how much to explore.
 */
func (r *Bishop) ExploreMoves(board *Board,
	canMove func(pieceColor byte, l Location, b *Board) (validMove bool, checkNext bool)) *[]Move {
	var moves []Move
	for i := 0; i < 4; i++ {
		l := r.GetPosition()
		for l.InBounds() {
			if i == 0 {
				l = l.Add(RightUpMove)
			} else if i == 1 {
				l = l.Add(RightDownMove)
			} else if i == 2 {
				l = l.Add(LeftUpMove)
			} else if i == 3 {
				l = l.Add(LeftDownMove)
			}
			validMove, checkNext := canMove(r.GetColor(), l, board)
			if validMove {
				moves = append(moves, Move{r.GetPosition(), l})
			}
			if !checkNext {
				break
			}
		}
	}
	return &moves
}

func (r *Bishop) GetMoves(board *Board) *[]Move {
	return r.ExploreMoves(board, CheckLocationForPiece)
}

func (r *Bishop) GetAttackableMoves(board *Board) *[]Move {
	return r.ExploreMoves(board, CheckLocationForAttackability)
}

func (r *Bishop) Move(m *Move, b *Board) {}

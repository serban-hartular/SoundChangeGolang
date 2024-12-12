package main

import (
	"fmt"
	"slices"
	"strings"
)

func string_distance(s SymStr, t SymStr) int {
	m := len(s)
	n := len(t)
	v0 := make([]int, n+1)
	v1 := make([]int, n+1)
	for i := range v0 {
		v0[i] = i
	}
	for i := 0; i < m; i++ {
		v1[0] = i + 1
		for j := 0; j < n; j++ {
			delCost := v0[j+1] + 1
			insCost := v1[j] + 1
			var substCost int
			if s[i] == t[j] {
				substCost = v0[j]
			} else {
				substCost = v0[j] + 1
			}
			v1[j+1] = min(delCost, insCost, substCost)
		}
		v0, v1 = v1, v0 // swap
	}
	return v0[n]
}

type Position struct {
	X, Y int
}

var NO_POS = Position{-1, -1}

const BOS = "#"
const EOS = "#"

type Edit struct {
	s_in  SymStr
	s_out SymStr
	from  Position
	score int
}

func (chg Edit) String() string {
	s_in, s_out := chg.s_in, chg.s_out
	if s_in.Empty() {
		s_in = SingleSymbol("0")
	}
	if s_out.Empty() {
		s_out = SingleSymbol("0")
	}
	return fmt.Sprintf("%s > %s", s_in, s_out)
}

func (chg Edit) NoChange() bool {
	return chg.s_in.equals(chg.s_out)
}

func INS(ss string, from Position, score int) Edit {
	return Edit{EmptySymStr(), SingleSymbol(ss), from, score}
}

func DEL(ss string, from Position, score int) Edit {
	return Edit{SingleSymbol(ss), EmptySymStr(), from, score}
}

func SUB(ss1, ss2 string, from Position, score int) Edit {
	return Edit{SingleSymbol(ss1), SingleSymbol(ss2), from, score}
}

func (chg Edit) isINS() bool { return chg.s_in.Empty() && !chg.s_out.Empty() }
func (chg Edit) isDEL() bool { return !chg.s_in.Empty() && chg.s_out.Empty() }
func (chg Edit) isSUB() bool {
	return !chg.s_in.equals(chg.s_out) && !chg.s_in.Empty() && !chg.s_out.Empty()
}

type MatrixCell struct {
	changes []Edit
	score   int
}

func NewMatrixCell(changes []Edit) MatrixCell {
	cell := MatrixCell{changes, 0}
	if len(changes) > 0 { // set cell score to be minimum score
		cell.score = changes[0].score
		for _, chg := range changes {
			if chg.score < cell.score {
				cell.score = chg.score
			}
		}
	}
	return cell
}

type ChangeMatrix [][]MatrixCell

func NewChangeMatrix(x_len, y_len int) ChangeMatrix {
	matrix := make(ChangeMatrix, x_len)
	for i := range matrix {
		matrix[i] = make([]MatrixCell, y_len)
	}
	return matrix
}

func (mm ChangeMatrix) xlen() int {
	return len(mm)
}

func (mm ChangeMatrix) ylen() int {
	return len(mm[0])
}

func (mm ChangeMatrix) get(pos Position) *MatrixCell {
	return &mm[pos.X][pos.Y]
}

func add_BOS_EOS(ss SymStr) SymStr {
	//	return NewSymStr(BOS + ss.String() + EOS)
	return SymStrConcat(SingleSymbol(BOS), ss, SingleSymbol(EOS))
}

func GenerateChangeMatrix(w1, w2 SymStr) ChangeMatrix {
	w1, w2 = add_BOS_EOS(w1), add_BOS_EOS(w2)
	matrix := NewChangeMatrix(len(w1), len(w2))

	//Initialize first row and first column
	for x := 1; x < matrix.xlen(); x++ {
		matrix[x][0] = NewMatrixCell([]Edit{DEL(w1[x], Position{x - 1, 0}, x)})
	}
	for y := 1; y < matrix.ylen(); y++ {
		matrix[0][y] = NewMatrixCell([]Edit{DEL(w2[y], Position{0, y - 1}, y)})
	}
	//Populate matrix
	for y := 1; y < matrix.ylen(); y++ {
		for x := 1; x < matrix.xlen(); x++ {
			var substCost int
			if w1[x] == w2[y] {
				substCost = 0
			} else {
				substCost = 1
			}
			options := []Edit{
				DEL(w1[x], Position{x - 1, y}, matrix[x-1][y].score+1),
				INS(w2[y], Position{x, y - 1}, matrix[x][y-1].score+1),
				SUB(w1[x], w2[y], Position{x - 1, y - 1}, matrix[x-1][y-1].score+substCost),
			}
			min_score := min(options[0].score, options[1].score, options[2].score)
			matrix[x][y].score = min_score
			for _, o := range options {
				if o.score == min_score {
					matrix[x][y].changes = append(matrix[x][y].changes, o)
				}
			}
		}
	}
	return matrix
}

type EditSequence []Edit

func (chgSeq EditSequence) String() string {
	chgStrings := make([]string, len(chgSeq))
	for i, chg := range chgSeq {
		chgStrings[i] = chg.String()
	}
	return "[" + strings.Join(chgStrings, ", ") + "]"
}

func find_change_sequences(matrix ChangeMatrix, pos Position) []EditSequence {
	if pos == NO_POS {
		pos = Position{matrix.xlen() - 1, matrix.ylen() - 1}
	}
	orig := Position{0, 0}
	if pos == orig {
		// return []ChangeSequence{ChangeSequence{SUB(BOS, BOS, NO_POS, 0)}}
		return []EditSequence{{SUB(BOS, BOS, NO_POS, 0)}}
	}
	cell := matrix.get(pos)
	mod_paths := make([]EditSequence, 0)
	for _, mod := range cell.changes {
		new_paths := find_change_sequences(matrix, mod.from)
		for i := range new_paths {
			new_paths[i] = append(new_paths[i], mod)
		}
		mod_paths = append(mod_paths, new_paths...)
	}
	return mod_paths
}

func word_pair_change_sequences(w1 SymStr, w2 SymStr) []EditSequence {
	matrix := GenerateChangeMatrix(w1, w2)
	return find_change_sequences(matrix, NO_POS)
}

func wordPairChangeSequencesAll(w1 SymStr, w2 SymStr) []EditSequence {
	original_seqs := word_pair_change_sequences(w1, w2)
	all_seqs := make([]EditSequence, len(original_seqs))
	copy(all_seqs, original_seqs)
	for _, seq := range original_seqs {
		versions := changeSequenceGetVersions(seq)
		all_seqs = append(all_seqs, versions...)
	}
	return all_seqs
}

func changeSequenceGetVersions(c_seq EditSequence) []EditSequence {
	//substitution + deletion -> merger (a>0, j>e -> aj>e, eg. Caesar -> Cesar)
	//substitution + insertion -> expansion (e>j, 0>a -> e>ja, e.g. erba -> iarba)
	// substitution + substitution -> transposition if c1.in == c2.out and c1.out == c2.in (e>r, r>e -> er > re, e.g. per -> pre)
	//...or vice-versa
	versions := make([]EditSequence, 0, len(c_seq)-1)
	for i := 0; i < len(c_seq)-1; i++ {
		current := c_seq[i]
		next := c_seq[i+1]
		if !current.NoChange() && !next.NoChange() && !(current.isDEL() && next.isDEL()) &&
			(!(current.isSUB() && next.isSUB()) ||
				(current.s_in.equals(next.s_out) && current.s_out.equals(next.s_in))) {
			new_chg := Edit{slices.Concat(current.s_in, next.s_in), slices.Concat(current.s_out, next.s_out),
				current.from, 1}
			new_chg_slice := EditSequence{new_chg}
			new_seq := slices.Concat(c_seq[:i], new_chg_slice, c_seq[i+2:])
			versions = append(versions, new_seq)
		}
	}
	return versions
}

// def find_change_sequences(matrix : ModMatrix, pos : Position = None) -> List[ChangeSequence]:
//     if pos is None:
//         pos = (matrix.rows-1, matrix.cols-1)
//     if pos == (0, 0):
//         return [[Transition(BOS, BOS)]] # beginning of string
//     mod_paths : List[ChangeSequence] = []
//     for mod in matrix.at(pos).modifications:
//         new_paths = find_change_sequences(matrix, mod.from_pos)
//         for path in new_paths:
//             # path.append((pos[0], mod))
//             path.append(Transition(mod.op.d_in, mod.op.d_out))
//         mod_paths.extend(new_paths)
//     return mod_paths

// def simple_string_distance(s : str, t : str) -> int:
//     m = len(s)
//     n = len(t)
//     v0 = list(range(n+1))
//     v1 = [0]*(n+1)
//     for i in range(m):
//         v1[0] = i + 1
//         for j in range(n):
//             # // calculating costs for A[i + 1][j + 1]
//             deletionCost = v0[j + 1] + 1
//             insertionCost = v1[j] + 1
//             substitutionCost = v0[j] if s[i] == t[j] else v0[j]+1
//             v1[j + 1] = min(deletionCost, insertionCost, substitutionCost)
//         v0, v1 = v1, v0 # swap
//     return v0[n]

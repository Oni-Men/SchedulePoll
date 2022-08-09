package poll

import (
	"github.com/google/uuid"
)

type Column struct {
	ID      uuid.UUID
	content string
	voters  map[Voter]bool
}

func CreateColumn(content string) *Column {
	col := Column{
		ID:      uuid.New(),
		content: content,
		voters:  make(map[Voter]bool),
	}

	return &col
}

func (col *Column) GetContent() string {
	return col.content
}

func (col *Column) AppendVoter(v Voter) {
	// If the key didn't exist. the flag will be false, zero-value of bool
	// キーが存在しないとき、flagはゼロ値であるfalseになる
	flag := col.voters[v]
	if flag {
		return
	}
	col.voters[v] = true
}

func (col *Column) RemoveVoter(v Voter) {
	//キーが存在しないかnilとき、何もしない
	delete(col.voters, v)
}

func (col *Column) VoteCount() int {
	return len(col.voters)
}

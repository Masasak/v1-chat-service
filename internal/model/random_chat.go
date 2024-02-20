package model

import (
	"github.com/google/uuid"
)

// RevealStage indicates the status of random chat's reveal stage.
type RevealStage uint8

const (
	// RevealStageInit indicates nobody hasn't voted for the reveal.
	RevealStageInit RevealStage = iota

	// RevealStageWaiting indicates at least 1 person voted for the reveal.
	RevealStageWaiting

	// RevealStageWaiting indicates all people voted for the reveal.
	RevealStageRevealed
)

// RevealVote is user's vote to the reveal.
type RevealVote struct {
	UserID    uuid.UUID
	CreatedAt uuid.UUID
}

// RandomChat is chat generated randomly. it includes chat in it.
type RandomChat struct {
	Chat
	Stage       RevealStage
	RevealVotes []RevealVote
	aliases     map[uuid.UUID]string
}

// SetAliases sets alias hashmap into given hashamp.
// Once it is set, it should be immutable.
func (rc *RandomChat) SetAliases(aliases map[uuid.UUID]string) {
	tmp := make(map[uuid.UUID]string, len(aliases))
	for k, v := range aliases {
		tmp[k] = v
	}
	rc.aliases = tmp
}

// GetAlias maps userID into specified alias.
func (rc *RandomChat) GetAlias(userID uuid.UUID) string {
	alias, ok := rc.aliases[userID]
	if !ok {
		// TODO: Change this to make it fit into requirements.
		return "Unknown"
	}
	return alias
}

// CanVote checks wheter user can vote for reveal in the chat.
func (rc *RandomChat) CanVote(userID uuid.UUID) bool {
	if rc.IsVoteFinished() {
		return false
	}
	for _, vote := range rc.RevealVotes {
		if vote.UserID == userID {
			return false
		}
	}
	return true
}

func (rc *RandomChat) IsVoteFinished() bool {
	return len(rc.UserIDs) == len(rc.RevealVotes)
}

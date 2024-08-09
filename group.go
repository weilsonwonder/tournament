package tournament

// Group is constructed in accordance to: https://www.ittf.com/wp-content/uploads/2021/01/HTR-2021-final.pdf, page 11
type Group interface {
	// TotalRounds required to complete the entire group match
	TotalRounds() int
	// TotalMatches total number of actual matches, minus null-matches (matches without opponent)
	TotalMatches() int
	// MatchesPerRound total number of actual matches in each round, minus null-matches (matches without opponent)
	MatchesPerRound() int
	// GetMatches for any round [1, totalRounds], includes null-matches
	GetMatches(round int) []Match
	// GetOpponents for any player by index (0-index based), and return list of opponents' index (based on initial players list, also need check for no-opponent scenario)
	GetOpponents(playerIdx int) []PlayerIndex
}

// Match represents 2 player indexes in a match, if one of the player index is invalid, this is a null match
type Match struct {
	Player1, Player2 PlayerIndex
}

// PlayerIndex represents index value of players based on original list
type PlayerIndex int

type groupImpl struct {
	players         int // number of players in a match
	seats           int // number of seats in a match, always even number (includes a null seat)
	rounds          int // total valid rounds (exclude null seat) in a group tournament
	matches         int // total valid matches (exclude null seat) in a group tournament
	matchesPerRound int // number of valid matches (exclude null seat) per round
}

// NewGroup returns a group containing all players in order to perform matching
func NewGroup(players int) Group {
	if players < 2 {
		panic("require at least 2 players")
	}

	group := groupImpl{
		players: players,
		seats:   players + players%2,
		rounds:  players + players%2 - 1,
		matches: players * (players - 1) / 2,
	}
	group.matchesPerRound = group.matches / group.rounds

	return group
}

func (z groupImpl) TotalRounds() int {
	return z.rounds
}

func (z groupImpl) TotalMatches() int {
	return z.matches
}

func (z groupImpl) MatchesPerRound() int {
	return z.matchesPerRound
}

func (z groupImpl) GetMatches(round int) []Match {
	round = max(0, round-1)%z.rounds + 1 // ensure round is [1, totalRounds]
	moves := z.rounds - round
	pairs := make([]Match, z.matchesPerRound+z.players%2) // include null-match where needed

	// keep first position 0 fixed, start from position 1
	pairs[0].Player1 = 0
	for i := 1; i < z.seats; i++ {
		var newPos int
		if i%2 == 0 { // even pos, +2 per move
			newPos = i + 2*moves
			if newPos >= z.seats {
				newPos = 2*z.seats - newPos - 1
				if newPos < 0 {
					newPos = -(newPos - 1)
				}
			}
		} else { // odd pos, -2 per move
			newPos = i - 2*moves
			if newPos < 0 {
				newPos = -(newPos - 1)
				if newPos >= z.seats {
					newPos = 2*z.seats - newPos - 1
				}
			}
		}

		// prepare player index
		playerIdx := PlayerIndex(newPos)
		if newPos == z.players {
			// this is a null seat
			playerIdx = -1
		}

		// update pair
		if i%2 == 0 {
			pairs[i/2].Player1 = playerIdx
		} else {
			pairs[i/2].Player2 = playerIdx
		}
	}

	return pairs
}

func (z groupImpl) GetOpponents(playerIdx int) []PlayerIndex {
	// validity check
	if playerIdx < 0 || z.players <= playerIdx {
		return nil // invalid player index, no opponents
	}

	// for playerIdx 0, easy scenario
	opponents := make([]PlayerIndex, z.rounds)
	if playerIdx == 0 {
		opponentIdx := playerIdx + 1
		for round := z.rounds - 1; round >= 0; round-- {
			// assign opponent if available
			if opponentIdx < z.players {
				opponents[round] = PlayerIndex(opponentIdx)
			} else {
				opponents[round] = -1
			}

			// get next opponentIdx
			if opponentIdx%2 == 0 { // even pos, +2
				opponentIdx += 2
				if opponentIdx >= z.seats {
					opponentIdx = z.seats - 1
				}
			} else { // odd pos, -2
				opponentIdx -= 2
				if opponentIdx < 0 {
					opponentIdx = 2
				}
			}
		}
	} else {
		getOpponentIdx := func(idx int) int {
			if idx%2 == 0 {
				return idx + 1
			}
			return idx - 1
		}
		opponentIdx := getOpponentIdx(playerIdx)
		for round := z.rounds - 1; round >= 0; round-- {
			// assign opponent if available
			if opponentIdx < z.players {
				opponents[round] = PlayerIndex(opponentIdx)
			} else {
				opponents[round] = -1
			}

			// move playerIdx (clockwise instead)
			if playerIdx%2 == 0 { // even pos, -2
				playerIdx -= 2
				if playerIdx <= 0 {
					playerIdx = 1
				}
			} else { // odd pos, +2
				playerIdx += 2
				if playerIdx >= z.seats {
					playerIdx = z.seats - 2
				}
			}

			// update with roundOpponentIdx
			roundOpponentIdx := getOpponentIdx(playerIdx)

			// reverse engineer opponentIdx
			moves := z.rounds - round
			if roundOpponentIdx == 0 { // special case, 0 index player do not move
				opponentIdx = 0
			} else if roundOpponentIdx%2 == 0 { // even pos, +2 per move
				opponentIdx = roundOpponentIdx + 2*moves
				if opponentIdx >= z.seats {
					opponentIdx = 2*z.seats - opponentIdx - 1
					if opponentIdx < 0 {
						opponentIdx = -(opponentIdx - 1)
					}
				}
			} else { // odd pos, -2 per move
				opponentIdx = roundOpponentIdx - 2*moves
				if opponentIdx < 0 {
					opponentIdx = -(opponentIdx - 1)
					if opponentIdx >= z.seats {
						opponentIdx = 2*z.seats - opponentIdx - 1
					}
				}
			}
		}
	}

	return opponents
}

// IsValid returns true if both players are valid
func (z Match) IsValid() bool {
	return z.Player1.IsValid() && z.Player2.IsValid()
}

// IsValid returns true for non-negative values
func (z PlayerIndex) IsValid() bool {
	return z >= 0
}

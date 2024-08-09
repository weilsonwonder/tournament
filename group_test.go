package tournament

import (
	"runtime"
	"testing"
)

type GroupTestCase struct {
	TotalPlayers    int
	TotalRounds     int
	TotalMatches    int
	MatchesPerRound int
	Matches         [][]Match
	Opponents       [][]PlayerIndex
}

var groupTestCases = []GroupTestCase{
	{
		// this is from ittf handbook
		TotalPlayers:    8,
		TotalRounds:     7,
		TotalMatches:    28,
		MatchesPerRound: 4,
		Matches: [][]Match{
			makematch(0, 3, 1, 5, 2, 7, 4, 6),
			makematch(0, 5, 3, 7, 1, 6, 2, 4),
			makematch(0, 7, 5, 6, 3, 4, 1, 2),
			makematch(0, 6, 7, 4, 5, 2, 3, 1),
			makematch(0, 4, 6, 2, 7, 1, 5, 3),
			makematch(0, 2, 4, 1, 6, 3, 7, 5),
			makematch(0, 1, 2, 3, 4, 5, 6, 7),
		},
		Opponents: [][]PlayerIndex{
			{3, 5, 7, 6, 4, 2, 1},
			{5, 6, 2, 3, 7, 4, 0},
			{7, 4, 1, 5, 6, 0, 3},
			{0, 7, 4, 1, 5, 6, 2},
			{6, 2, 3, 7, 0, 1, 5},
			{1, 0, 6, 2, 3, 7, 4},
			{4, 1, 5, 0, 2, 3, 7},
			{2, 3, 0, 4, 1, 5, 6},
		},
	},
	{
		TotalPlayers:    7,
		TotalRounds:     7,
		TotalMatches:    21,
		MatchesPerRound: 3,
		Matches: [][]Match{
			makematch(0, 3, 1, 5, 2, -1, 4, 6),
			makematch(0, 5, 3, -1, 1, 6, 2, 4),
			makematch(0, -1, 5, 6, 3, 4, 1, 2),
			makematch(0, 6, -1, 4, 5, 2, 3, 1),
			makematch(0, 4, 6, 2, -1, 1, 5, 3),
			makematch(0, 2, 4, 1, 6, 3, -1, 5),
			makematch(0, 1, 2, 3, 4, 5, 6, -1),
		},
		Opponents: [][]PlayerIndex{
			{3, 5, -1, 6, 4, 2, 1},
			{5, 6, 2, 3, -1, 4, 0},
			{-1, 4, 1, 5, 6, 0, 3},
			{0, -1, 4, 1, 5, 6, 2},
			{6, 2, 3, -1, 0, 1, 5},
			{1, 0, 6, 2, 3, -1, 4},
			{4, 1, 5, 0, 2, 3, -1},
		},
	},
	{
		TotalPlayers:    6,
		TotalRounds:     5,
		TotalMatches:    15,
		MatchesPerRound: 3,
		Matches: [][]Match{
			makematch(0, 3, 1, 5, 2, 4),
			makematch(0, 5, 3, 4, 1, 2),
			makematch(0, 4, 5, 2, 3, 1),
			makematch(0, 2, 4, 1, 5, 3),
			makematch(0, 1, 2, 3, 4, 5),
		},
		Opponents: [][]PlayerIndex{
			{3, 5, 4, 2, 1},
			{5, 2, 3, 4, 0},
			{4, 1, 5, 0, 3},
			{0, 4, 1, 5, 2},
			{2, 3, 0, 1, 5},
			{1, 0, 2, 3, 4},
		},
	},
	{
		TotalPlayers:    5,
		TotalRounds:     5,
		TotalMatches:    10,
		MatchesPerRound: 2,
		Matches: [][]Match{
			makematch(0, 3, 1, -1, 2, 4),
			makematch(0, -1, 3, 4, 1, 2),
			makematch(0, 4, -1, 2, 3, 1),
			makematch(0, 2, 4, 1, -1, 3),
			makematch(0, 1, 2, 3, 4, -1),
		},
		Opponents: [][]PlayerIndex{
			{3, -1, 4, 2, 1},
			{-1, 2, 3, 4, 0},
			{4, 1, -1, 0, 3},
			{0, 4, 1, -1, 2},
			{2, 3, 0, 1, -1},
		},
	},
	{
		TotalPlayers:    4,
		TotalRounds:     3,
		TotalMatches:    6,
		MatchesPerRound: 2,
		Matches: [][]Match{
			makematch(0, 3, 1, 2),
			makematch(0, 2, 3, 1),
			makematch(0, 1, 2, 3),
		},
		Opponents: [][]PlayerIndex{
			{3, 2, 1},
			{2, 3, 0},
			{1, 0, 3},
			{0, 1, 2},
		},
	},
	{
		TotalPlayers:    3,
		TotalRounds:     3,
		TotalMatches:    3,
		MatchesPerRound: 1,
		Matches: [][]Match{
			makematch(0, -1, 1, 2),
			makematch(0, 2, -1, 1),
			makematch(0, 1, 2, -1),
		},
		Opponents: [][]PlayerIndex{
			{-1, 2, 1},
			{2, -1, 0},
			{1, 0, -1},
		},
	},
	{
		TotalPlayers:    2,
		TotalRounds:     1,
		TotalMatches:    1,
		MatchesPerRound: 1,
		Matches: [][]Match{
			makematch(0, 1),
		},
		Opponents: [][]PlayerIndex{
			{1},
			{0},
		},
	},
}

func TestGroup(t *testing.T) {
	for i, tc := range groupTestCases {
		group := NewGroup(tc.TotalPlayers)
		if group.TotalRounds() != tc.TotalRounds {
			t.Errorf("tc_%d: wrong total rounds: %d expected %d", i, group.TotalRounds(), tc.TotalRounds)
		}
		if group.TotalMatches() != tc.TotalMatches {
			t.Errorf("tc_%d: wrong total matches: %d expected %d", i, group.TotalMatches(), tc.TotalMatches)
		}
		if group.MatchesPerRound() != tc.MatchesPerRound {
			t.Errorf("tc_%d: wrong matches per round: %d expected %d", i, group.MatchesPerRound(), tc.MatchesPerRound)
		}
		for roundIdx, tcMatches := range tc.Matches {
			round := roundIdx + 1
			//matches := group.(*groupImpl[int]).getMatchesStable(round) // can swap to stable if necessary
			matches := group.GetMatches(round)
			if len(matches) != len(tcMatches) {
				t.Errorf("tc_%d_round_%d: invalid number of matches: %d expected %d", i, round, len(matches), len(tcMatches))
			}
			for x, match := range matches {
				tcMatch := tcMatches[x]
				if !cmpMatch(match, tcMatch) {
					t.Errorf("tc_%d_round_%d_match_%d: unexpected pairing: %v expected %v", i, round, x, match, tcMatch)
				}
			}
		}
		for playerIdx, tcOpponents := range tc.Opponents {
			opponents := group.GetOpponents(playerIdx)
			if len(opponents) != len(tcOpponents) {
				t.Errorf("tc_%d_playeridx_%d: unexpected number of opponents: %v expected %v", i, playerIdx, len(opponents), len(tcOpponents))
			}
			for x, opponent := range opponents {
				tcOpponent := tcOpponents[x]
				if !cmpOpponent(opponent, tcOpponent) {
					t.Errorf("tc_%d_playeridx_%d_opponentidx_%d: unexpected opponents: %v expected %v", i, playerIdx, x, opponent, tcOpponent)
				}
			}
		}
	}

	// just to print and for debugging or visualization purposes
	//tc := groupTestCases[0]
	//group := NewGroup(tc.TotalPlayers)
	//for roundIdx := range tc.Matches {
	//	round := roundIdx + 1
	//	matches := group.GetMatches(round)
	//	printMatches(t, round, matches)
	//}
	//for playerIdx, player := range group.Players() {
	//	opponents := group.GetOpponents(playerIdx)
	//	t.Log(player, "|", opponents)
	//}
}

func BenchmarkBigGroup(b *testing.B) {
	const totalPlayers = 10_000

	// print base info
	group := NewGroup(totalPlayers)
	b.Log("total rounds:", group.TotalRounds())
	b.Log("total matches:", group.TotalMatches())
	b.Log("matches per round:", group.MatchesPerRound())

	// begin test
	b.Run("get a single match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = group.GetMatches((i % group.TotalRounds()) + 1)
		}
		b.ReportAllocs()
	})
	b.Run("get a single match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = group.GetMatches((i % group.TotalRounds()) + 1)
		}
		b.ReportAllocs()
	})
	b.Run("get all matches", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for round := 1; round <= group.TotalRounds(); round++ {
				_ = group.GetMatches(round)
			}
		}
		b.ReportAllocs()
	})
	b.Run("get a player's opponents", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = group.GetOpponents(i % 2) // alternate between 0 and 1
		}
		b.ReportAllocs()
	})
	b.Run("get all players' opponents", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for playerIdx := 0; playerIdx < totalPlayers; playerIdx++ {
				_ = group.GetOpponents(playerIdx)
			}
		}
		b.ReportAllocs()
	})
}

func TestGroupMemory(t *testing.T) {
	tc := groupTestCases[0]
	group := NewGroup(tc.TotalPlayers)
	debugMemory(t, "[even-matches]", func() {
		for roundIdx := range tc.Matches {
			_ = group.GetMatches(roundIdx + 1)
		}
	})
	debugMemory(t, "[even-opponents]", func() {
		for playerIdx := 0; playerIdx < tc.TotalPlayers; playerIdx++ {
			_ = group.GetOpponents(playerIdx)
		}
	})

	tc = groupTestCases[2]
	group = NewGroup(tc.TotalPlayers)
	debugMemory(t, "[odd-matches]", func() {
		for roundIdx := range tc.Matches {
			_ = group.GetMatches(roundIdx + 1)
		}
	})
	debugMemory(t, "[odd-opponents]", func() {
		for playerIdx := 0; playerIdx < tc.TotalPlayers; playerIdx++ {
			_ = group.GetOpponents(playerIdx)
		}
	})
}

func BenchmarkGroup(b *testing.B) {
	tc := groupTestCases[0]

	b.Run("new group", func(b *testing.B) {
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			_ = NewGroup(tc.TotalPlayers)
		}
		b.ReportAllocs()
	})

	group := NewGroup(tc.TotalPlayers)

	b.Run("get matches (even) single", func(b *testing.B) {
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			_ = group.GetMatches(0)
		}
		b.ReportAllocs()
	})
	b.Run("get matches (even) all rounds", func(b *testing.B) {
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			for roundIdx := range tc.Matches {
				_ = group.GetMatches(roundIdx + 1)
			}
		}
		b.ReportAllocs()
	})
	b.Run("get matches (even) single opponents", func(b *testing.B) {
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			_ = group.GetOpponents(0)
		}
		b.ReportAllocs()
	})
	b.Run("get matches (even) all opponents", func(b *testing.B) {
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			for playerIdx := 0; playerIdx < tc.TotalPlayers; playerIdx++ {
				_ = group.GetOpponents(playerIdx)
			}
		}
		b.ReportAllocs()
	})

	tc = groupTestCases[2]
	group = NewGroup(tc.TotalPlayers)

	b.Run("get matches (odd) single", func(b *testing.B) {
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			_ = group.GetMatches(0)
		}
		b.ReportAllocs()
	})
	b.Run("get matches (odd) all rounds", func(b *testing.B) {
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			for roundIdx := range tc.Matches {
				_ = group.GetMatches(roundIdx + 1)
			}
		}
		b.ReportAllocs()
	})
	b.Run("get matches (odd) single opponents", func(b *testing.B) {
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			_ = group.GetOpponents(0)
		}
		b.ReportAllocs()
	})
	b.Run("get matches (odd) all opponents", func(b *testing.B) {
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			for playerIdx := 0; playerIdx < tc.TotalPlayers; playerIdx++ {
				_ = group.GetOpponents(playerIdx)
			}
		}
		b.ReportAllocs()
	})
}

func cmpMatch(a, b Match) bool {
	return a.Player1 == b.Player1 && a.Player2 == b.Player2
}

func cmpOpponent(a, b PlayerIndex) bool {
	return a == b
}

func makematch(indexes ...int) []Match {
	matches := make([]Match, 0, len(indexes)/2)
	for i := 0; i < len(indexes); i += 2 {
		matches = append(matches, Match{
			Player1: PlayerIndex(indexes[i]),
			Player2: PlayerIndex(indexes[i+1]),
		})
	}
	return matches
}

func printMatches(t *testing.T, round int, matches []Match) {
	t.Log("round:", round)
	for _, pair := range matches {
		t.Log("\t", pair.Player1, "v", pair.Player2)
	}
}

func debugMemory(t *testing.T, name string, doWork func()) {
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	doWork()
	runtime.ReadMemStats(&m2)
	t.Log(name, "total:", m2.TotalAlloc-m1.TotalAlloc)
	t.Log(name, "mallocs:", m2.Mallocs-m1.Mallocs)
}

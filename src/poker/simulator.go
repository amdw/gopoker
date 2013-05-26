package poker

type Simulator struct {
	HandCount           int
	WinCount            int
	OurClassCounts      []int
	OpponentClassCounts []int
}

func (s *Simulator) SimulateHoldem(yourCards, tableCards []Card, players, handsToPlay int) {
	s.HandCount = handsToPlay
	s.WinCount = 0
	s.OurClassCounts = make([]int, MAX_HANDCLASS)
	s.OpponentClassCounts = make([]int, MAX_HANDCLASS)

	p := NewPack()
	for i := 0; i < handsToPlay; i++ {
		p.shuffleFixing(tableCards, yourCards)
		won, ourLevel, opponentLevel := p.SimulateOneHoldemHand(players)
		if won {
			s.WinCount++
		}
		s.OurClassCounts[ourLevel.Class]++
		s.OpponentClassCounts[opponentLevel.Class]++
	}
}

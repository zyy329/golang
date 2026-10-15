package ordern

type OrderType int32

//goland:noinspection GoUnusedConst
const (
	BetType       OrderType = 1
	AwardType     OrderType = 2
	CancelBetType OrderType = 3
)

package status

type RoomStatus string

const (
	Waiting     = RoomStatus("waiting")
	Setup       = RoomStatus("setup")
	GameStarted = RoomStatus("game_started")
	Shuffling   = RoomStatus("shuffling")
	Picking     = RoomStatus("picking")
	Finished    = RoomStatus("finished")
)

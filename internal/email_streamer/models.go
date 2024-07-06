package emailstreamer

type Streamer struct {
	ID         uint `gorm:"primaryKey"`
	Topic      string
	Partition  int
	LastOffset uint
}

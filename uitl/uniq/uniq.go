package uniq

var (
	num = make(chan uint64)
)

func init() {
	go func() {
		for i := uint64(0); ; i++ {
			num <- i
		}
	}()
}

func GetUniq() uint64 {
	return <-num
}

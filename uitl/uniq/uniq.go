package uniq

var (
	num = make(chan int)
)

func init() {
	go func() {
		for i := 0; ; i++ {
			num <- i
		}
	}()
}

func GetUniq() int {
	return <-num
}

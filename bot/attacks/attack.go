package attacks

type Attack interface {
	Name() string
	Send(host string, port int, seconds int, size int, threads int)
	Stop(host string, port int)
}

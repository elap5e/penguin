package rpc

type Client interface {
	Close() error
	Go(serviceMethod string, args *Args, reply *Reply, done chan *Call) *Call
	Call(serviceMethod string, args *Args, reply *Reply) error

	GetNextSeq() int32
	GetAppID() int32
	SetAppID(id int32)
}

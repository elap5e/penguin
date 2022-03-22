package rpc

type Call struct {
	ServiceMethod string
	Seq           uint64
	Args          *Args
	Reply         *Reply
	Error         error
	Done          chan *Call
}

func (call *Call) done() {
	select {
	case call.Done <- call:
		// ok
	default:
		// We don't want to block here. It is the caller's responsibility to make
		// sure the channel has enough buffer space. See comment in Go().
	}
}

type Args struct {
	Seq    int32
	AppID  int32
	FixID  int32
	Buffer []byte
}

type Reply struct {
	Seq int32
}

type Request struct {
	ServiceMethod string
	Seq           uint64
}

type Response struct {
	ServiceMethod string
	Seq           uint64
}

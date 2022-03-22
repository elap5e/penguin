package msf

import (
	"math/rand"
	"sync/atomic"

	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

func NewClient() rpc.Client {
	return &client{
		rs:  rpc.NewSender(),
		seq: rand.Int31n(100000),
	}
}

type client struct {
	rs  rpc.Sender
	seq int32
}

func (c *client) Close() error {
	return c.rs.Close()
}

func (c *client) Go(serviceMethod string, args *rpc.Args, reply *rpc.Reply, done chan *rpc.Call) *rpc.Call {
	return c.rs.Go(serviceMethod, args, reply, done)
}

func (c *client) Call(serviceMethod string, args *rpc.Args, reply *rpc.Reply) error {
	call := <-c.Go(serviceMethod, args, reply, make(chan *rpc.Call, 1)).Done
	return call.Error
}

func (c *client) GetNextSeq() int32 {
	seq := atomic.AddInt32(&c.seq, 1)
	if seq > 1000000 {
		c.seq = rand.Int31n(100000) + 60000
	}
	return seq
}

func (c *client) GetAppID() int32 {
	return -1
}

func (c *client) SetAppID(id int32) {
}

var _ rpc.Client = (*client)(nil)

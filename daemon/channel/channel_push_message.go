// Copyright 2022 Elapse and contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package channel

import (
	"encoding/json"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin/daemon/channel/pb"
	pbmsg "github.com/elap5e/penguin/daemon/message/pb"
	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

func (m *Manager) handlePushMessage(reply *rpc.Reply) (*rpc.Args, error) {
	push := pb.MsgPush_MsgOnlinePush{}
	if err := proto.Unmarshal(reply.Payload, &push); err != nil {
		return nil, err
	}
	for _, msg := range push.GetMsgs() {
		typ, sub := msg.GetHead().GetContentHead().GetMsgType(), msg.GetHead().GetContentHead().GetSubType()
		if typ == 3840 {
			_ = m.OnRecvChannelMessage(reply.Uin, msg)
		} else if typ == 3841 {
			dumpUnknown(typ, sub, msg)
			var got *pbmsg.IMMsgBody_CommonElem
			for _, elem := range msg.GetBody().GetRichText().GetElems() {
				if elem := elem.GetCommonElem(); elem != nil {
					got = elem
					break
				}
			}

			var body pb.ChannelService_EventBody
			if err := proto.Unmarshal(got.GetPbElem(), &body); err != nil {
				log.Error("failed to unmarshal channel event: %v", err)
				continue
			}
			p, _ := json.Marshal(&body)
			log.Warn("unknown msg type:%d sub_type:%d body:%s", typ, sub, p)
		} else if typ == 3842 {
			dumpUnknown(typ, sub, msg)
		} else {
			dumpUnknown(typ, sub, msg)
		}
	}
	return nil, nil
}

func dumpUnknown(typ, sub uint64, msg *pb.Common_Msg) {
	p, _ := json.Marshal(msg)
	log.Warn("unknown msg type:%d sub_type:%d msg:%s", typ, sub, p)
}

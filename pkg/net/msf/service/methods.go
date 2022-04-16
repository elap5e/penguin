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

package service

import "fmt"

const (
	MethodHeartbeatAlive = "Heartbeat.Alive"
)

const (
	MethodAuthName2Uin = "wtlogin.name2uin"
	MethodAuthSignIn   = "wtlogin.login"
	MethodAuthSignInA2 = "wtlogin.exchange_emp"
)

const (
	MethodChatGetGroups     = "friendlist.GetTroopListReqV2"
	MethodChatGetGroupUsers = "friendlist.getTroopMemberList"
	MethodChatGetGroupUser  = "group_member_card.get_group_member_card_info"
)

const (
	MethodChannelGetMessage  = "trpc.group_pro.synclogic.SyncLogic.GetChannelMsg"
	MethodChannelSendMessage = "MsgProxy.SendMsg"

	MethodChannelSyncFirstView = "trpc.group_pro.synclogic.SyncLogic.SyncFirstView"
	MethodChannelPushFirstView = "trpc.group_pro.synclogic.SyncLogic.PushFirstView"
	MethodChannelPushMessage   = "MsgPush.PushGroupProMsg"
)

const (
	MethodContactDeleteContact   = "friendlist.delFriend"
	MethodContactGetContacts     = "friendlist.getFriendGroupList"
	MethodContactSetContactGroup = "friendlist.SetGroupReq"
)

const (
	MethodMessageDeleteMessage = "MessageSvc.PbDeleteMsg"
	MethodMessageGetMessage    = "MessageSvc.PbGetMsg"
	MethodMessageSendMessage   = "MessageSvc.PbSendMsg"
	MethodMessageRecallMessage = "PbMessageSvc.PbMsgWithDraw"

	MethodMessagePushNotify = "MessageSvc.PushNotify"
	MethodMessagePushReaded = "MessageSvc.PushReaded"

	MethodMessageChatDownloadImage = "ImgStore.GroupPicDown"
	MethodMessageChatUploadImage   = "ImgStore.GroupPicUp"
	MethodMessageUserDownloadImage = "LongConn.OffPicDown"
	MethodMessageUserUploadImage   = "LongConn.OffPicUp"
)

const (
	MethodServiceRegister            = "StatSvc.register"
	MethodServiceSetStatusFromClient = "StatSvc.SetStatusFromClient"

	MethodServiceConfigPushDomain   = "ConfigPushSvc.PushDomain"
	MethodServiceConfigPushRequest  = "ConfigPushSvc.PushReq"
	MethodServiceConfigPushResponse = "ConfigPushSvc.PushResp"

	MethodServiceOnlinePushUserMessage   = "OnlinePush.PbC2CMsgSync"
	MethodServiceOnlinePushChatMessage   = "OnlinePush.PbPushGroupMsg"
	MethodServiceOnlinePushChatSerivce   = "OnlinePush.PbPushTransMsg"
	MethodServiceOnlinePushRequest       = "OnlinePush.ReqPush"
	MethodServiceOnlinePushResponse      = "OnlinePush.RespPush"
	MethodServiceOnlinePushTicketExpired = "OnlinePush.SidTicketExpired"
)

func MethodOidbSvcTrpcTcp(cmd, typ uint32) string {
	return fmt.Sprintf("OidbSvcTrpcTcp.0x%x_%d", cmd, typ)
}

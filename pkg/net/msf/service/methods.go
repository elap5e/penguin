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

import (
	"fmt"
)

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

	MethodChannelDownloadPhoto = "ImgStore.QQMeetPicDown"
	MethodChannelUploadPhoto   = "ImgStore.QQMeetPicUp"
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

	MethodMessageChatDownloadMedia = "MultiMsg.ApplyDown"
	MethodMessageChatUploadMedia   = "MultiMsg.ApplyUp"

	MethodMessageChatDownloadPhoto = "ImgStore.GroupPicDown"
	MethodMessageChatUploadPhoto   = "ImgStore.GroupPicUp"
	MethodMessageUserDownloadPhoto = "LongConn.OffPicDown"
	MethodMessageUserUploadPhoto   = "LongConn.OffPicUp"

	MethodMessageChatDownloadVideo = "PttCenterSvr.GroupShortVideoDownReq"
	MethodMessageChatUploadVideo   = "PttCenterSvr.GroupShortVideoUpReq"
	MethodMessageUserDownloadVideo = "PttCenterSvr.ShortVideoDownReq"
	MethodMessageUserForwardVideo  = "PttCenterSvr.ShortVideoRetweetReq"
	MethodMessageUserUploadVideo   = "PttCenterSvr.ShortVideoUpReq"

	MethodMessageChatDownloadVoice = "PttStore.GroupPttDown"
	MethodMessageChatUploadVoice   = "PttStore.GroupPttUp"

	MethodMessageChatDownloadVideoSticker = "OidbSvcTrpcTcp.0x10dd_1"
	MethodMessageChatUploadVideoSticker   = "OidbSvcTrpcTcp.0x10dd_0"

	_ = "MultiVideo.s2c"
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

	MethodServicePushForceOffline = "MessageSvc.PushForceOffline"
	MethodServicePushLoginNotify  = "StatSvc.SvcReqMSFLoginNotify"
)

func MethodOidbSvcTrpcTcp(cmd, svc uint32) string {
	return fmt.Sprintf("OidbSvcTrpcTcp.0x%x_%d", cmd, svc)
}

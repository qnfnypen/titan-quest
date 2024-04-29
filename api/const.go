package api

const (
	MissionIdConnectWallet int64 = iota + 1001
	MissionIdFollowTwitter
	MissionIdRetweet
	MissionIdLikeTwitter
	MissionIdJoinDiscord
	MissionIdJoinTelegram
	MissionIdBindingKOL
	MissionIdVisitOfficialWebsite
	MissionIdVisitReferrerPage
)

const (
	MissionIdQuoteTweet int64 = iota + 1106
	MissionIdPostTweet
	MissionIdInviteFriendsToDiscord
)

const (
	MissionTypeBasic int32 = iota + 1
	MissionTypeDaily
	MissionTypeWeekly
)

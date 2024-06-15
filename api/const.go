package api

const (
	MissionIdBindTitanWallet int64 = iota + 1000
	MissionIdConnectWallet
	MissionIdFollowTwitter
	MissionIdRetweet
	MissionIdLikeTwitter
	MissionIdJoinDiscord
	MissionIdJoinTelegramGroup
	MissionIdBindingKOL
	MissionIdVisitOfficialWebsite
	MissionIdVisitReferrerPage
	MissionIdBrowsOfficialWebSite
	MissionIdBecomeVolunteer // 成为预备志愿者
	MissionIdJoinDCVolunteerChannel
	MissionIdJoinTelegramVolunteerGroup
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

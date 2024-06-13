package api

const (
	MissionIdConnectWallet int64 = iota + 1001
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
	MissionIdJoinSpecifyDCChannel
	MissionIdJoinSpecifyDCVolunteerChannel
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

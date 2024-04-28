package discord

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gnasnik/titan-quest/core/dao"
	"github.com/golang-module/carbon/v2"
	"log"
)

var (
	inviteCounterKey = "gm::discord::invitecounter::%s::%s"
	recordsKey       = "gm::discord::inviterecords"
	membersKey       = "gm::discord::members"
)

func RunDiscordBot(token string) {
	fmt.Println("==>", token)

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("create discord cmd: %v", err)
		return
	}

	ctx := context.Background()

	invitesBefore := make(map[string][]*discordgo.Invite)

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as %s %s\n", r.User, r.User.ID)

		var guildIds []string
		for _, guild := range r.Guilds {
			guildIds = append(guildIds, guild.ID)

			invites, err := s.GuildInvites(guild.ID)
			if err != nil {
				log.Println(err)
				continue
			}

			invitesBefore[guild.ID] = invites
		}

		go UpdateAllMemberTask(ctx, s, guildIds)

	})

	s.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
		log.Printf("User %s join", m.User.ID)

		_, err := dao.RedisCache.SAdd(ctx, membersKey, m.User.ID).Result()
		if err != nil {
			log.Printf("redis sadd: %s %v\n", membersKey, err)
			return
		}

		guildInvitesBefore := invitesBefore[m.GuildID]
		log.Printf("invs before: %d\n", len(guildInvitesBefore))

		invitesAfter, err := s.GuildInvites(m.GuildID)
		if err != nil {
			log.Print("GuildInvites", err)
			return
		}

		log.Printf("invs after: %d\n", len(invitesAfter))

		exist, err := dao.RedisCache.SIsMember(ctx, recordsKey, m.User.ID).Result()
		if err != nil {
			log.Printf("redis sis member: %v\n", err)
			return
		}

		if exist {
			log.Printf("User %s have joined before, the number of invitations will no longer be counted", m.User.ID)
			return
		}

		for _, newInvite := range invitesAfter {
			var isNewInvite bool

			retInv := findInviteByCode(guildInvitesBefore, newInvite.Code)
			if retInv == nil {
				// first newInvite
				if newInvite.Uses == 1 {
					isNewInvite = true
				}
				continue
			}

			if retInv.Uses < newInvite.Uses {
				isNewInvite = true
				fmt.Println("newInvite code:", newInvite.Code)
				fmt.Printf("inviter: %s %s\n", newInvite.Inviter, newInvite.Inviter.ID)
			}

			if isNewInvite {
				_, err := dao.RedisCache.SAdd(ctx, recordsKey, m.User.ID).Result()
				if err != nil {
					log.Printf("redis sadd: %s %v\n", recordsKey, err)
					continue
				}

				startOfWeek := carbon.Now().StartOfWeek().String()
				_, err = dao.RedisCache.Incr(ctx, fmt.Sprintf(inviteCounterKey, newInvite.Inviter.ID, startOfWeek)).Result()
				if err != nil {
					log.Printf("redis incr: %v\n", err)
				}

				return
			}
		}
	})

	s.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
		log.Printf("User %s left", m.User.ID)

		_, err := dao.RedisCache.SRem(ctx, membersKey, m.User.ID).Result()
		if err != nil {
			log.Printf("redis srem: %s %v\n", membersKey, err)
			return
		}
	})

	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()
}

func findInviteByCode(invites []*discordgo.Invite, code string) *discordgo.Invite {
	for _, invite := range invites {
		if invite.Code == code {
			return invite
		}
	}
	return nil
}

func UpdateAllMemberTask(ctx context.Context, s *discordgo.Session, guildIds []string) {
	log.Println("start update all member")

	for _, guild := range guildIds {

		log.Printf("guild: %s members to redis\n", guild)

		var after string

		for {
			finish, last, err := loadMembers(ctx, s, guild, after)
			if err != nil {
				log.Printf("load members: %v\n", err)
			}

			after = last

			if finish {
				break
			}
		}
	}

	log.Println("update all member done")
}

func loadMembers(ctx context.Context, s *discordgo.Session, guild string, after string) (bool, string, error) {
	members, err := s.GuildMembers(guild, after, 1000)
	if err != nil {
		log.Println(err)
		return true, "", err
	}

	if len(members) == 0 {
		log.Println("finish")
		return true, "", nil
	}

	var memberIds []interface{}
	for _, mem := range members {
		memberIds = append(memberIds, mem.User.ID)
	}

	_, err = dao.RedisCache.SAdd(ctx, membersKey, memberIds...).Result()
	if err != nil {
		log.Printf("redis sadd: %s %v\n", recordsKey, err)
		return true, "", err
	}

	last := members[len(members)-1].User.ID
	log.Printf("loading %d members to redis", len(members))

	return false, last, nil
}

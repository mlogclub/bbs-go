package migrations

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/repositories"

	"errors"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

type levelSeed struct {
	Level   int
	NeedExp int
	Title   string
	Status  int
}

type badgeSeed struct {
	Name        string
	Title       string
	Description string
	Icon        string
	SortNo      int
	Status      int
}

const (
	badgeIconNewcomer     = "/res/images/badges/badge_newcomer.svg"
	badgeIconFirstPost    = "/res/images/badges/badge_first_post.svg"
	badgeIconFirstComment = "/res/images/badges/badge_first_comment.svg"
	badgeIconAuthor       = "/res/images/badges/badge_author.svg"
	badgeIconHelper       = "/res/images/badges/badge_helper.svg"
	badgeIconStreak7      = "/res/images/badges/badge_streak_7.svg"
	badgeIconStreak30     = "/res/images/badges/badge_streak_30.svg"
	badgeIconVeteran      = "/res/images/badges/badge_veteran.svg"
)

type taskConfigSeed struct {
	GroupName   constants.TaskGroup
	EventType   string
	Title       string
	Description string

	Score  int
	Exp    int
	Badge  string // Badge.Name
	Period int
	SortNo int
	Status int

	MaxFinishCount int
	EventCount     int

	BtnName   string
	ActionUrl string

	StartTime int64
	EndTime   int64
}

type taskSystemSeed struct {
	Levels []levelSeed
	Badges []badgeSeed
	Tasks  []taskConfigSeed
}

func migrate_init_task_data() error {
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		now := dates.NowTimestamp()
		seed := taskSeedForLanguage()

		if err := validateLevelSeeds(seed.Levels); err != nil {
			return err
		}

		for _, s := range seed.Levels {
			existing := repositories.LevelConfigRepository.Take(tx, "level = ?", s.Level)
			if existing != nil {
				existing.NeedExp = s.NeedExp
				existing.Title = s.Title
				existing.Status = s.Status
				existing.UpdateTime = now
				if err := repositories.LevelConfigRepository.Update(tx, existing); err != nil {
					return err
				}
				continue
			}

			levelCfg := &models.LevelConfig{
				Level:      s.Level,
				NeedExp:    s.NeedExp,
				Title:      s.Title,
				Status:     s.Status,
				CreateTime: now,
				UpdateTime: now,
			}
			if err := repositories.LevelConfigRepository.Create(tx, levelCfg); err != nil {
				return err
			}
		}

		badgeIDMap := make(map[string]int64, len(seed.Badges))
		for _, s := range seed.Badges {
			existing := repositories.BadgeRepository.Take(tx, "name = ?", s.Name)
			if existing != nil {
				existing.Title = s.Title
				existing.Description = s.Description
				existing.Icon = s.Icon
				existing.SortNo = s.SortNo
				existing.Status = s.Status
				existing.UpdateTime = now
				if err := repositories.BadgeRepository.Update(tx, existing); err != nil {
					return err
				}
				badgeIDMap[s.Name] = existing.Id
				continue
			}

			badge := &models.Badge{
				Name:        s.Name,
				Title:       s.Title,
				Description: s.Description,
				Icon:        s.Icon,
				SortNo:      s.SortNo,
				Status:      s.Status,
				CreateTime:  now,
				UpdateTime:  now,
			}
			if err := repositories.BadgeRepository.Create(tx, badge); err != nil {
				return err
			}
			badgeIDMap[s.Name] = badge.Id
		}

		for _, s := range seed.Tasks {
			var badgeId int64
			if s.Badge != "" {
				badgeId = badgeIDMap[s.Badge]
				if badgeId == 0 {
					return errors.New("badge not found: " + s.Badge)
				}
			}

			task := &models.TaskConfig{
				EventType:      s.EventType,
				Title:          s.Title,
				Description:    s.Description,
				Score:          s.Score,
				Exp:            s.Exp,
				BadgeId:        badgeId,
				Period:         s.Period,
				MaxFinishCount: s.MaxFinishCount,
				EventCount:     s.EventCount,
				BtnName:        s.BtnName,
				ActionUrl:      s.ActionUrl,
				SortNo:         s.SortNo,
				StartTime:      s.StartTime,
				EndTime:        s.EndTime,
				Status:         s.Status,
				CreateTime:     now,
				UpdateTime:     now,
				GroupName:      s.GroupName,
			}
			if err := repositories.TaskConfigRepository.Create(tx, task); err != nil {
				return err
			}
		}

		return nil
	})
}

func taskSeedForLanguage() taskSystemSeed {
	lang := config.Instance.Language
	if !lang.IsValid() {
		lang = config.DefaultLanguage
	}
	if lang == config.LanguageEnUS {
		return taskSystemSeed{
			Levels: defaultLevelsEn(),
			Badges: defaultBadgesEn(),
			Tasks:  defaultTasksEn(),
		}
	}
	return taskSystemSeed{
		Levels: defaultLevelsZh(),
		Badges: defaultBadgesZh(),
		Tasks:  defaultTasksZh(),
	}
}

func validateLevelSeeds(levels []levelSeed) error {
	if len(levels) == 0 {
		return errors.New("level seeds empty")
	}
	for i, l := range levels {
		expectedLevel := i + 1
		if l.Level != expectedLevel {
			return errors.New("level must start from 1 and be continuous")
		}
		if i == 0 {
			continue
		}
		if l.NeedExp <= levels[i-1].NeedExp {
			return errors.New("needExp must be strictly increasing")
		}
	}
	return nil
}

func defaultLevelsEn() []levelSeed {
	return []levelSeed{
		{Level: 1, NeedExp: 0, Title: "Level 1", Status: constants.StatusOk},
		{Level: 2, NeedExp: 60, Title: "Level 2", Status: constants.StatusOk},
		{Level: 3, NeedExp: 150, Title: "Level 3", Status: constants.StatusOk},
		{Level: 4, NeedExp: 300, Title: "Level 4", Status: constants.StatusOk},
		{Level: 5, NeedExp: 500, Title: "Level 5", Status: constants.StatusOk},
		{Level: 6, NeedExp: 750, Title: "Level 6", Status: constants.StatusOk},
		{Level: 7, NeedExp: 1050, Title: "Level 7", Status: constants.StatusOk},
		{Level: 8, NeedExp: 1400, Title: "Level 8", Status: constants.StatusOk},
		{Level: 9, NeedExp: 1800, Title: "Level 9", Status: constants.StatusOk},
		{Level: 10, NeedExp: 2250, Title: "Level 10", Status: constants.StatusOk},
		{Level: 11, NeedExp: 2750, Title: "Level 11", Status: constants.StatusOk},
		{Level: 12, NeedExp: 3300, Title: "Level 12", Status: constants.StatusOk},
		{Level: 13, NeedExp: 3900, Title: "Level 13", Status: constants.StatusOk},
		{Level: 14, NeedExp: 4550, Title: "Level 14", Status: constants.StatusOk},
		{Level: 15, NeedExp: 5250, Title: "Level 15", Status: constants.StatusOk},
		{Level: 16, NeedExp: 6000, Title: "Level 16", Status: constants.StatusOk},
		{Level: 17, NeedExp: 6800, Title: "Level 17", Status: constants.StatusOk},
		{Level: 18, NeedExp: 7650, Title: "Level 18", Status: constants.StatusOk},
		{Level: 19, NeedExp: 8550, Title: "Level 19", Status: constants.StatusOk},
		{Level: 20, NeedExp: 9500, Title: "Level 20", Status: constants.StatusOk},
	}
}

func defaultLevelsZh() []levelSeed {
	return []levelSeed{
		{Level: 1, NeedExp: 0, Title: "等级 1", Status: constants.StatusOk},
		{Level: 2, NeedExp: 60, Title: "等级 2", Status: constants.StatusOk},
		{Level: 3, NeedExp: 150, Title: "等级 3", Status: constants.StatusOk},
		{Level: 4, NeedExp: 300, Title: "等级 4", Status: constants.StatusOk},
		{Level: 5, NeedExp: 500, Title: "等级 5", Status: constants.StatusOk},
		{Level: 6, NeedExp: 750, Title: "等级 6", Status: constants.StatusOk},
		{Level: 7, NeedExp: 1050, Title: "等级 7", Status: constants.StatusOk},
		{Level: 8, NeedExp: 1400, Title: "等级 8", Status: constants.StatusOk},
		{Level: 9, NeedExp: 1800, Title: "等级 9", Status: constants.StatusOk},
		{Level: 10, NeedExp: 2250, Title: "等级 10", Status: constants.StatusOk},
		{Level: 11, NeedExp: 2750, Title: "等级 11", Status: constants.StatusOk},
		{Level: 12, NeedExp: 3300, Title: "等级 12", Status: constants.StatusOk},
		{Level: 13, NeedExp: 3900, Title: "等级 13", Status: constants.StatusOk},
		{Level: 14, NeedExp: 4550, Title: "等级 14", Status: constants.StatusOk},
		{Level: 15, NeedExp: 5250, Title: "等级 15", Status: constants.StatusOk},
		{Level: 16, NeedExp: 6000, Title: "等级 16", Status: constants.StatusOk},
		{Level: 17, NeedExp: 6800, Title: "等级 17", Status: constants.StatusOk},
		{Level: 18, NeedExp: 7650, Title: "等级 18", Status: constants.StatusOk},
		{Level: 19, NeedExp: 8550, Title: "等级 19", Status: constants.StatusOk},
		{Level: 20, NeedExp: 9500, Title: "等级 20", Status: constants.StatusOk},
	}
}

func defaultBadgesEn() []badgeSeed {
	return []badgeSeed{
		{Name: string(constants.BadgeNameNewcomer), Title: "Newcomer", Description: "First check-in completed", Icon: badgeIconNewcomer, SortNo: 10, Status: constants.StatusOk},
		{Name: string(constants.BadgeNameFirstPost), Title: "First Post", Description: "Created the first post", Icon: badgeIconFirstPost, SortNo: 20, Status: constants.StatusOk},
		{Name: string(constants.BadgeNameFirstComment), Title: "First Comment", Description: "Created the first comment", Icon: badgeIconFirstComment, SortNo: 30, Status: constants.StatusOk},
		{Name: string(constants.BadgeNameAuthor), Title: "Author", Description: "Created 10 posts", Icon: badgeIconAuthor, SortNo: 40, Status: constants.StatusOk},
		{Name: string(constants.BadgeNameHelper), Title: "Helper", Description: "Created 50 comments", Icon: badgeIconHelper, SortNo: 50, Status: constants.StatusOk},
		{Name: string(constants.BadgeNameStreak7), Title: "7-day Streak", Description: "Checked in for 7 days", Icon: badgeIconStreak7, SortNo: 60, Status: constants.StatusOk},
		{Name: string(constants.BadgeNameStreak30), Title: "30-day Streak", Description: "Checked in for 30 days", Icon: badgeIconStreak30, SortNo: 70, Status: constants.StatusOk},
		{Name: string(constants.BadgeNameVeteran), Title: "Veteran", Description: "Reached level 10", Icon: badgeIconVeteran, SortNo: 80, Status: constants.StatusOk},
	}
}

func defaultBadgesZh() []badgeSeed {
	return []badgeSeed{
		{Name: string(constants.BadgeNameNewcomer), Title: "新人报到", Description: "完成首次签到", Icon: badgeIconNewcomer, SortNo: 10, Status: constants.StatusOk},
		{Name: string(constants.BadgeNameFirstPost), Title: "初来乍到", Description: "发布第一篇帖子", Icon: badgeIconFirstPost, SortNo: 20, Status: constants.StatusOk},
		{Name: string(constants.BadgeNameFirstComment), Title: "初次发言", Description: "发表第一条评论", Icon: badgeIconFirstComment, SortNo: 30, Status: constants.StatusOk},
		{Name: string(constants.BadgeNameAuthor), Title: "内容创作者", Description: "累计发帖 10 篇", Icon: badgeIconAuthor, SortNo: 40, Status: constants.StatusOk},
		{Name: string(constants.BadgeNameHelper), Title: "热心助人", Description: "累计评论 50 次", Icon: badgeIconHelper, SortNo: 50, Status: constants.StatusOk},
		{Name: string(constants.BadgeNameStreak7), Title: "坚持一周", Description: "累计签到 7 天", Icon: badgeIconStreak7, SortNo: 60, Status: constants.StatusOk},
		{Name: string(constants.BadgeNameStreak30), Title: "月度打卡王", Description: "累计签到 30 天", Icon: badgeIconStreak30, SortNo: 70, Status: constants.StatusOk},
		{Name: string(constants.BadgeNameVeteran), Title: "资深玩家", Description: "达到等级 10", Icon: badgeIconVeteran, SortNo: 80, Status: constants.StatusOk},
	}
}

func defaultTasksEn() []taskConfigSeed {
	return []taskConfigSeed{
		// newbie (one-time)
		{GroupName: constants.TaskGroupNewbie, Title: "First check-in", Description: "Sign in for the first time to get started and earn your Newcomer badge. A small step that kicks off your journey here.", EventType: "checkin", EventCount: 1, MaxFinishCount: 1, Period: 0, Score: 10, Exp: 20, Badge: string(constants.BadgeNameNewcomer), SortNo: 10, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupNewbie, Title: "First post", Description: "Share your first topic with the community. Your voice matters—unlock the First Post badge and let others see what you have to say.", EventType: "topic.create", EventCount: 1, MaxFinishCount: 1, Period: 0, Score: 20, Exp: 40, Badge: string(constants.BadgeNameFirstPost), SortNo: 20, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupNewbie, Title: "First comment", Description: "Leave your first comment on a post. Join the conversation and connect with other members—rewards include the First Comment badge.", EventType: "comment.create", EventCount: 1, MaxFinishCount: 1, Period: 0, Score: 10, Exp: 25, Badge: string(constants.BadgeNameFirstComment), SortNo: 30, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupNewbie, Title: "Follow 1 user", Description: "Follow someone you find interesting. You'll see their posts in your feed and stay in the loop with the community.", EventType: "follow.create", EventCount: 1, MaxFinishCount: 1, Period: 0, Score: 5, Exp: 10, SortNo: 40, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupNewbie, Title: "Favorite 1 item", Description: "Save a post or topic to your favorites. Build your personal collection and revisit great content anytime.", EventType: "favorite.create", EventCount: 1, MaxFinishCount: 1, Period: 0, Score: 5, Exp: 10, SortNo: 50, Status: constants.StatusOk},

		// daily
		{GroupName: constants.TaskGroupDaily, Title: "Daily login", Description: "Log in to the site today. A quick visit keeps you connected and earns you daily rewards.", EventType: "user.login", EventCount: 1, MaxFinishCount: 1, Period: 1, Score: 2, Exp: 5, SortNo: 110, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupDaily, Title: "Daily check-in", Description: "Complete today's check-in to build your streak. Consistency pays off with extra points and experience.", EventType: "checkin", EventCount: 1, MaxFinishCount: 1, Period: 1, Score: 5, Exp: 10, SortNo: 120, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupDaily, Title: "Create 1 post", Description: "Publish one new topic today. Share a question, an idea, or a find—every post enriches the community.", EventType: "topic.create", EventCount: 1, MaxFinishCount: 1, Period: 1, Score: 8, Exp: 15, SortNo: 130, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupDaily, Title: "Comment 3 times", Description: "Leave at least 3 comments on posts today. Your feedback and discussion keep the community lively.", EventType: "comment.create", EventCount: 3, MaxFinishCount: 1, Period: 1, Score: 6, Exp: 12, SortNo: 140, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupDaily, Title: "Like 10 times", Description: "Give 10 likes to posts or comments you enjoy today. A simple way to show appreciation and support creators.", EventType: "like.create", EventCount: 10, MaxFinishCount: 1, Period: 1, Score: 4, Exp: 8, SortNo: 150, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupDaily, Title: "Favorite 2 items", Description: "Add 2 posts or topics to your favorites today. Curate your reading list and never lose track of good content.", EventType: "favorite.create", EventCount: 2, MaxFinishCount: 1, Period: 1, Score: 3, Exp: 6, SortNo: 160, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupDaily, Title: "Follow 1 user", Description: "Follow one new user today. Expand your feed and discover more voices that match your interests.", EventType: "follow.create", EventCount: 1, MaxFinishCount: 1, Period: 1, Score: 2, Exp: 4, SortNo: 170, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupDaily, Title: "Comment 10 times", Description: "Leave 10 comments today. Active participants earn bonus rewards—dive into discussions and share your thoughts.", EventType: "comment.create", EventCount: 10, MaxFinishCount: 1, Period: 1, Score: 10, Exp: 20, SortNo: 180, Status: constants.StatusOk},

		// achievement (one-time)
		{GroupName: constants.TaskGroupAchievement, Title: "10 posts", Description: "Publish 10 topics in total. Reach this milestone to earn the Author badge and be recognized as an active contributor.", EventType: "topic.create", EventCount: 10, MaxFinishCount: 1, Period: 0, Score: 50, Exp: 120, Badge: string(constants.BadgeNameAuthor), SortNo: 210, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupAchievement, Title: "50 comments", Description: "Reach 50 total comments. Your engagement helps the community thrive—unlock the Helper badge as a thank you.", EventType: "comment.create", EventCount: 50, MaxFinishCount: 1, Period: 0, Score: 50, Exp: 120, Badge: string(constants.BadgeNameHelper), SortNo: 220, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupAchievement, Title: "7 check-ins", Description: "Complete 7 check-ins in total. Build a week-long habit and earn the 7-day Streak badge for your dedication.", EventType: "checkin", EventCount: 7, MaxFinishCount: 1, Period: 0, Score: 30, Exp: 80, Badge: string(constants.BadgeNameStreak7), SortNo: 230, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupAchievement, Title: "30 check-ins", Description: "Complete 30 check-ins in total. A month of consistency earns the 30-day Streak badge and substantial rewards.", EventType: "checkin", EventCount: 30, MaxFinishCount: 1, Period: 0, Score: 120, Exp: 300, Badge: string(constants.BadgeNameStreak30), SortNo: 240, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupAchievement, Title: "100 likes", Description: "Give 100 likes in total. Spread positivity across the community and get rewarded for supporting others' work.", EventType: "like.create", EventCount: 100, MaxFinishCount: 1, Period: 0, Score: 30, Exp: 80, SortNo: 250, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupAchievement, Title: "Follow 20 users", Description: "Follow 20 users in total. Grow your network and fill your feed with diverse, interesting content.", EventType: "follow.create", EventCount: 20, MaxFinishCount: 1, Period: 0, Score: 30, Exp: 80, SortNo: 260, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupAchievement, Title: "Reach level 10", Description: "Reach level 10 to earn the Veteran badge and mark your steady growth in the community.", EventType: constants.TaskEventTypeLevel10, EventCount: 1, MaxFinishCount: 1, Period: 0, Score: 80, Exp: 200, Badge: string(constants.BadgeNameVeteran), SortNo: 270, Status: constants.StatusOk},
	}
}

func defaultTasksZh() []taskConfigSeed {
	return []taskConfigSeed{
		// 新手任务（一次性）
		{GroupName: constants.TaskGroupNewbie, Title: "完成首次签到", Description: "完成第一次签到即可迈出第一步，并获得「新人报到」徽章。小小一步，开启你在本站的旅程。", EventType: "checkin", EventCount: 1, MaxFinishCount: 1, Period: 0, Score: 10, Exp: 20, Badge: string(constants.BadgeNameNewcomer), SortNo: 10, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupNewbie, Title: "发布第一篇帖子", Description: "在社区发布你的第一篇主题。你的声音很重要——完成即可解锁「初来乍到」徽章，让更多人看到你的分享。", EventType: "topic.create", EventCount: 1, MaxFinishCount: 1, Period: 0, Score: 20, Exp: 40, Badge: string(constants.BadgeNameFirstPost), SortNo: 20, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupNewbie, Title: "发表第一条评论", Description: "在任意帖子下留下你的第一条评论，参与讨论、结识同好，完成即可获得「初次发言」徽章。", EventType: "comment.create", EventCount: 1, MaxFinishCount: 1, Period: 0, Score: 10, Exp: 25, Badge: string(constants.BadgeNameFirstComment), SortNo: 30, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupNewbie, Title: "关注 1 位用户", Description: "关注一位你感兴趣的用户，他们的动态会出现在你的时间线中，让你与社区保持同步。", EventType: "follow.create", EventCount: 1, MaxFinishCount: 1, Period: 0, Score: 5, Exp: 10, SortNo: 40, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupNewbie, Title: "收藏 1 篇内容", Description: "将一篇帖子或主题加入收藏，建立你的个人收藏夹，随时回看优质内容。", EventType: "favorite.create", EventCount: 1, MaxFinishCount: 1, Period: 0, Score: 5, Exp: 10, SortNo: 50, Status: constants.StatusOk},

		// 每日任务
		{GroupName: constants.TaskGroupDaily, Title: "每日登录", Description: "今日访问并登录站点即可完成。保持每日登录，领取每日奖励，与社区保持联系。", EventType: "user.login", EventCount: 1, MaxFinishCount: 1, Period: 1, Score: 2, Exp: 5, SortNo: 110, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupDaily, Title: "每日签到", Description: "完成今日签到，延续你的签到 streak。坚持签到可获得更多积分与经验值。", EventType: "checkin", EventCount: 1, MaxFinishCount: 1, Period: 1, Score: 5, Exp: 10, SortNo: 120, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupDaily, Title: "发 1 篇帖子", Description: "今日发布一篇新主题。可以是一则提问、一个想法或一次分享——每篇帖子都在丰富社区。", EventType: "topic.create", EventCount: 1, MaxFinishCount: 1, Period: 1, Score: 8, Exp: 15, SortNo: 130, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupDaily, Title: "评论 3 次", Description: "今日在帖子下至少发表 3 条评论。你的回复与讨论能让社区更有活力。", EventType: "comment.create", EventCount: 3, MaxFinishCount: 1, Period: 1, Score: 6, Exp: 12, SortNo: 140, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupDaily, Title: "点赞 10 次", Description: "今日为你喜欢的帖子或评论点赞 10 次。用点赞表达认可，支持创作者。", EventType: "like.create", EventCount: 10, MaxFinishCount: 1, Period: 1, Score: 4, Exp: 8, SortNo: 150, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupDaily, Title: "收藏 2 次", Description: "今日将 2 篇帖子或主题加入收藏，整理你的阅读清单，不错过好内容。", EventType: "favorite.create", EventCount: 2, MaxFinishCount: 1, Period: 1, Score: 3, Exp: 6, SortNo: 160, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupDaily, Title: "关注 1 人", Description: "今日新关注一位用户，扩展你的时间线，发现更多志趣相投的创作者。", EventType: "follow.create", EventCount: 1, MaxFinishCount: 1, Period: 1, Score: 2, Exp: 4, SortNo: 170, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupDaily, Title: "评论 10 次", Description: "今日发表 10 条评论。积极参与讨论的用户可获得额外每日奖励。", EventType: "comment.create", EventCount: 10, MaxFinishCount: 1, Period: 1, Score: 10, Exp: 20, SortNo: 180, Status: constants.StatusOk},

		// 成就任务（一次性）
		{GroupName: constants.TaskGroupAchievement, Title: "累计发帖 10 篇", Description: "累计发布 10 篇主题即可达成。解锁「内容创作者」徽章，成为被认可的活跃贡献者。", EventType: "topic.create", EventCount: 10, MaxFinishCount: 1, Period: 0, Score: 50, Exp: 120, Badge: string(constants.BadgeNameAuthor), SortNo: 210, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupAchievement, Title: "累计评论 50 次", Description: "累计发表 50 条评论。你的参与让社区更有活力，达成即可获得「热心助人」徽章。", EventType: "comment.create", EventCount: 50, MaxFinishCount: 1, Period: 0, Score: 50, Exp: 120, Badge: string(constants.BadgeNameHelper), SortNo: 220, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupAchievement, Title: "累计签到 7 天", Description: "累计完成 7 次签到。养成一周签到习惯，即可获得「坚持一周」徽章。", EventType: "checkin", EventCount: 7, MaxFinishCount: 1, Period: 0, Score: 30, Exp: 80, Badge: string(constants.BadgeNameStreak7), SortNo: 230, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupAchievement, Title: "累计签到 30 天", Description: "累计完成 30 次签到。坚持一个月即可获得「月度打卡王」徽章及丰厚奖励。", EventType: "checkin", EventCount: 30, MaxFinishCount: 1, Period: 0, Score: 120, Exp: 300, Badge: string(constants.BadgeNameStreak30), SortNo: 240, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupAchievement, Title: "累计点赞 100 次", Description: "累计为帖子或评论点赞 100 次。用点赞传递认可，支持他人创作，即可领取奖励。", EventType: "like.create", EventCount: 100, MaxFinishCount: 1, Period: 0, Score: 30, Exp: 80, SortNo: 250, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupAchievement, Title: "累计关注 20 人", Description: "累计关注 20 位用户。扩大你的关注列表，让时间线充满更多有趣、多元的内容。", EventType: "follow.create", EventCount: 20, MaxFinishCount: 1, Period: 0, Score: 30, Exp: 80, SortNo: 260, Status: constants.StatusOk},
		{GroupName: constants.TaskGroupAchievement, Title: "达到等级 10", Description: "达到等级 10 即可解锁「资深玩家」徽章，见证你在社区的成长。", EventType: constants.TaskEventTypeLevel10, EventCount: 1, MaxFinishCount: 1, Period: 0, Score: 80, Exp: 200, Badge: string(constants.BadgeNameVeteran), SortNo: 270, Status: constants.StatusOk},
	}
}

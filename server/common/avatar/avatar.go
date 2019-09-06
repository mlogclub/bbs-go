package avatar

var DefaultAvatars = []string{
	"https://file.mlog.club/avatar/club_default_avatar1.png",
	"https://file.mlog.club/avatar/club_default_avatar2.png",
	"https://file.mlog.club/avatar/club_default_avatar3.png",
	"https://file.mlog.club/avatar/club_default_avatar4.png",
	"https://file.mlog.club/avatar/club_default_avatar5.png",
	"https://file.mlog.club/avatar/club_default_avatar6.png",
}

// 获取默认头像
func GetDefaultAvatar(id int64) string {
	if id <= 0 {
		return DefaultAvatars[0]
	}
	i := int(id) % len(DefaultAvatars)
	return DefaultAvatars[i]
}

// 是否是默认头像
func IsDefaultAvatar(avatar string) bool {
	if len(avatar) == 0 {
		return true
	}
	for _, a := range DefaultAvatars {
		if a == avatar {
			return true
		}
	}
	return false
}

package avatar

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"strconv"
	"strings"

	"github.com/issue9/identicon"
)

const DefaultAvatar = "https://file.mlog.club/avatar/club_default_avatar1.png"

var (
	err                   error
	avatarBackgroundColor = color.Transparent // 头像背景色
	avatarFrontColors     []color.Color       // 头像前景色
	avatarFrontHexColors  = []string{
		"#FFC1C1", "#FFC125", "#FFC0CB", "#FFBBFF",
		"#FFB90F", "#FFB6C1", "#FFB5C5", "#FFAEB9",
		"#FFA54F", "#FFA500", "#FFA07A", "#FF8C69",
		"#FF8C00", "#FF83FA", "#FF82AB", "#FF8247",
		"#FF7F50", "#FF7F24", "#FF7F00", "#FF7256",
		"#FF6EB4", "#FF6A6A", "#FF69B4", "#FF6347",
		"#FF4500", "#FF4040", "#FF3E96", "#FF34B3",
		"#FF3030", "#FF1493", "#FF00FF", "#FF0000",
	}
	identiconIns *identicon.Identicon
)

func init() {
	for _, hexColor := range avatarFrontHexColors {
		c, _ := colorToRGB(hexColor)
		avatarFrontColors = append(avatarFrontColors, *c)
	}
	identiconIns, _ = identicon.New(300, color.Transparent, avatarFrontColors...)
}

func Generate(userId int64) ([]byte, error) {
	buf := new(bytes.Buffer)
	img := GenerateAvatar(userId)
	if err = png.Encode(buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 生成默认头像
func GenerateAvatar(userId int64) image.Image {
	return identiconIns.Make([]byte(strconv.FormatInt(userId, 10)))
}

func colorToRGB(colorStr string) (*color.RGBA, error) {
	colorStr = strings.TrimPrefix(colorStr, "#")
	color64, err := strconv.ParseInt(colorStr, 16, 32)
	if err != nil {
		return nil, err
	}
	colorInt := int(color64)
	r, g, b := colorInt>>16, (colorInt&0x00FF00)>>8, colorInt&0x0000FF
	return &color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, nil
}

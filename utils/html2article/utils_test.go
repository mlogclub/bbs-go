package html2article

import (
	"testing"
	"time"
)

func TestCompress(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test0",
			args: args{
				str: "test ",
			},
			want: "test ",
		},
		{
			name: "test1",
			args: args{
				str: " test ",
			},
			want: " test ",
		},
		{
			name: "test2",
			args: args{
				str: "test 2  ",
			},
			want: "test 2 ",
		},
		{
			name: "test3",
			args: args{
				str: "test 3  \n    ",
			},
			want: "test 3 ",
		},
		{
			name: "test4",
			args: args{
				str: "test4",
			},
			want: "test4",
		},
		{
			name: "test5",
			args: args{
				str: "test5  test5  \n test   ",
			},
			want: "test5 test5 test ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Compress(tt.args.str); got != tt.want {
				t.Errorf("Compress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDistance(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		maxValue int

		want int
		ok   bool
	}{
		{"1", "abc", "ab", 10, 1, true},
		{"2", "abc", "abd", 10, 1, true},
		{"3", "ab", "abcef", 10, 3, true},
		{"4", "我爱中国", "我是中国人", 10, 2, true},
		{"5", "abcdefg", "hijklmnxfdfd", 3, (1 << 31) - 1, false},
		{"6", "36氪-创业学院教授：创业者都应该学习亚马逊，敢于颠覆市场 | 创投圈", "36氪-创业学院教授：创业者都应该学习亚马逊，敢于颠覆市场", 100, 6, true},
		{"7", "提示", "36氪-创业学院教授：创业者都应该学习亚马逊，敢于颠覆市场 | 创投圈", 117, 35, true},
		{"8", "提示", "36氪-创业学院教授：创业者都应该学习亚马逊，敢于颠覆市场 | 创投圈", 11, 35, false},
		{"9", "本文作者：三川 2017-01-20 18:42 导语： Torch 的新生还是终结？", "Facebook 发布开源框架 PyTorch orch 终于被移植到 Python 生态圈 | 雷锋网", 15, 100, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := distanceExit(tt.a, tt.b, tt.maxValue)
			if ok != tt.ok {
				t.Errorf("Not ok equal %v, want %v", ok, tt.ok)
			}
			if ok {
				if got != tt.want {
					t.Errorf("distanceExit() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestDiffstr(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string

		want int
	}{
		{"1", "abc", "ab", 1},
		{"2", "abc", "abd", 2},
		{"3", "ab", "abcef", 3},
		{"4", "我的站长之路 丨 伊成Blog", "伊成Blog", 8},
		{"5", "我的站长之路 丨 伊成Blog", "我的站长之路", 8},
		{"6", "abc", "abef", 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := diffString(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("diffString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTime(t *testing.T) {
	tests := []struct {
		name string
		a    string
		want int64
	}{
		{"1", "fdaf5小时前 ggagg", time.Now().Add(-(5-8)*time.Hour).Unix() / 3600 * 3600},
		{"2", "hgha3天前fdsa", (time.Now().Add(-3*time.Hour*24 + time.Hour*8).Unix()) / int64(24*3600) * int64(24*3600)},
		{"3", "2017-02-14 05:48", 1487051280},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTime(tt.a); got != tt.want {
				t.Errorf("Time = %v, want %v", got, tt.want)
			}
		})
	}
}

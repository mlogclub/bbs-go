package iplocator

import "testing"

func TestIpLocationFromRegion(t *testing.T) {
	tests := []struct {
		name   string
		region string
		want   string
	}{
		{
			name:   "old_ipv4_format",
			region: "中国|0|湖北省|武汉市|电信",
			want:   "湖北省",
		},
		{
			name:   "new_ipv6_format_cn",
			region: "中国|广东省|深圳市|电信|CN",
			want:   "广东省",
		},
		{
			name:   "new_ipv6_format_foreign",
			region: "Australia|Queensland|Brisbane|0|AU",
			want:   "Queensland",
		},
		{
			name:   "fallback_to_nation",
			region: "中国|0|0|0|0",
			want:   "中国",
		},
		{
			name:   "invalid",
			region: "",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ipLocationFromRegion(tt.region)
			if got != tt.want {
				t.Fatalf("ipLocationFromRegion()=%q want=%q", got, tt.want)
			}
		})
	}
}

package iplocator

import (
	"strings"
	"sync"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"
	"golang.org/x/exp/slog"
)

var (
	once     sync.Once
	searcher *xdb.Searcher
)

func InitIpLocator(dbPath string) {
	once.Do(func() {
		if strs.IsBlank(dbPath) {
			dbPath = "ip2region.xdb"
		}
		start := dates.NowTimestamp()
		data, err := xdb.LoadContentFromFile(dbPath)
		if err != nil {
			slog.Error("failed to load content", slog.Any("path", dbPath), slog.Any("err", err))
			return
		}
		if searcher, err = xdb.NewWithBuffer(data); err != nil {
			slog.Error("failed to create searcher with content", slog.Any("err", err))
			return
		}
		slog.Info("Load ip2region.xdb success", slog.Any("elapsed", dates.NowTimestamp()-start))
	})
}

func Search(ip string) string {
	if searcher == nil || strs.IsBlank(ip) {
		return ""
	}
	region, _ := searcher.SearchByStr(ip)
	return region
}

func IpLocation(ip string) string {
	region := Search(ip) // eg. 中国|0|湖北省|武汉市|电信
	if strs.IsBlank(region) {
		return ""
	}
	ss := strings.Split(region, "|")
	if len(ss) != 5 {
		return ""
	}
	var (
		nation   = ss[0]
		province = ss[2]
	)
	if strs.IsNotBlank(province) && province != "0" {
		return province
	}
	if strs.IsNotBlank(nation) && nation != "0" {
		return nation
	}
	return ""
}

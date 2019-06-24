package html2article

import (
	"math"

	"golang.org/x/net/html"
)

type Info struct {
	TextCount     int
	LinkTextCount int
	TagCount      int
	LinkTagCount  int
	LeafList      []int
	Density       float64
	DensitySum    float64
	Pcount        int
	InputCount    int
	ImageCount    int

	Data  string
	avg   float64
	score float64

	node *html.Node
}

func NewInfo() *Info {
	return &Info{}
}

func (info *Info) CalScore(sn_sum, swn_sum float64) {
	sn := countSn(info.Data)
	swn := countStopWords(info.Data)
	info.avg = info.getAvg()
	info.score = info.DensitySum * math.Log(info.avg) * math.Log10(float64(info.Pcount+2)) * (float64(sn)/sn_sum + 2) * (float64(swn)/swn_sum + 2)

	// return
	// if info.score >= 0 {
	// 	c := attr(info.node, "class")
	// 	if c == "" {
	// 		c = attr(info.node, "id")
	// 	}
	// 	if true {
	// 		println("class:", c, info.score, info.DensitySum, math.Log(info.avg), math.Log10(float64(info.Pcount+2)), (float64(sn)/sn_sum + 2), (float64(swn)/swn_sum + 2), sn, sn_sum)
	// 	}
	// }
}

func (info *Info) getAvg() float64 {
	if len(info.LeafList) == 0 {
		return 0
	}
	flen := float64(len(info.LeafList))
	sum := 0
	for _, l := range info.LeafList {
		sum += l
	}
	var sum2 float64 = 0
	avg := float64(sum) / flen
	for _, l := range info.LeafList {
		sum2 += (avg - float64(l)) * (avg - float64(l))
	}
	return math.Sqrt(sum2/flen + 1.0)
}

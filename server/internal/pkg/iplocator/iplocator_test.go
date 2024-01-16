package iplocator_test

import (
	"bbs-go/internal/pkg/iplocator"
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	iplocator.InitIpLocator("/data/ip2region.xdb")
	ip := "47.52.26.78"
	fmt.Println(iplocator.Search(ip))
	fmt.Println(iplocator.IpLocation(ip))
}

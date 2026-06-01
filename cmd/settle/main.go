package main

import (
	"bbs-go/internal/server"
	"bbs-go/internal/services/heatpoints"
	"fmt"
	"os"

	_ "bbs-go/internal/services/eventhandler"
)

func main() {
	server.Init()

	fmt.Println("=== 开始手动结算 ===")
	if err := heatpoints.Settlement.SettleAll(); err != nil {
		fmt.Println("结算失败:", err)
		os.Exit(1)
	}
	fmt.Println("=== 结算完成 ===")
}

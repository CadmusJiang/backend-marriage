package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	fmt.Println("=== 检查启动标志 ===")
	
	// 定义标志
	forceReseed := flag.Bool("force-reseed", false, "If true, truncate and reseed core tables on startup")
	rebuildDatabase := flag.Bool("rebuild-db", false, "If true, drop and recreate all tables with city codes")
	
	// 解析标志
	flag.Parse()
	
	fmt.Printf("force-reseed: %v\n", *forceReseed)
	fmt.Printf("rebuild-db: %v\n", *rebuildDatabase)
	
	fmt.Println("\n=== 环境变量检查 ===")
	fmt.Printf("REBUILD_DB: %s\n", os.Getenv("REBUILD_DB"))
	fmt.Printf("FORCE_RESEED: %s\n", os.Getenv("FORCE_RESEED"))
	
	fmt.Println("\n=== 命令行参数 ===")
	fmt.Printf("命令行参数数量: %d\n", len(os.Args))
	for i, arg := range os.Args {
		fmt.Printf("  [%d] %s\n", i, arg)
	}
	
	fmt.Println("\n=== 标志默认值 ===")
	fmt.Printf("force-reseed 默认值: %v\n", flag.Lookup("force-reseed").DefValue)
	fmt.Printf("rebuild-db 默认值: %v\n", flag.Lookup("rebuild-db").DefValue)
}

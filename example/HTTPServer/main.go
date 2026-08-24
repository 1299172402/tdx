package main

import (
	"flag"
	"log"
	"strings"

	"github.com/injoyai/tdx/extend/httpserver"
)

func main() {
	// 定义命令行参数
	var (
		addr     string
		hostsStr string
		poolSize int
	)
	flag.StringVar(&addr, "addr", ":8080", "监听地址，例如 :8080 或 0.0.0.0:8080")
	flag.StringVar(&hostsStr, "hosts", "103.221.142.73", "标准行情服务器地址，多个用逗号分隔，例如 host1,host2,host3")
	flag.IntVar(&poolSize, "pool", 4, "标准连接池大小")
	flag.Parse()

	// 处理 addr：如果是纯数字端口（不含冒号），自动加上 ":"
	if addr != "" && !strings.Contains(addr, ":") {
		addr = ":" + addr
	}

	// 解析 hosts 字符串为切片
	var hosts []string
	if hostsStr != "" {
		hosts = strings.Split(hostsStr, ",")
		// 去除可能的空格
		for i, h := range hosts {
			hosts[i] = strings.TrimSpace(h)
		}
	}

	// 创建 HTTP 服务
	s, err := httpserver.New(
		httpserver.WithAddr(addr),
		httpserver.WithHosts(hosts...), // 可变参数展开
		httpserver.WithPoolSize(poolSize),
	)
	if err != nil {
		log.Fatalf("创建服务失败: %v", err)
	}

	// 启动服务
	log.Printf("服务启动，监听地址: %s，使用 hosts: %v，连接池大小: %d\n", addr, hosts, poolSize)
	if err := s.Run(); err != nil {
		log.Fatalf("服务运行异常: %v", err)
	}
}

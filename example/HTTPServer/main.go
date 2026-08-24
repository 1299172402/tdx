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
		corsStr  string
	)
	flag.StringVar(&addr, "addr", ":8080", "监听地址，例如 :8080 或 0.0.0.0:8080")
	flag.StringVar(&hostsStr, "hosts", "103.221.142.73", "标准行情服务器地址，多个用逗号分隔，例如 host1,host2,host3")
	flag.IntVar(&poolSize, "pool", 4, "标准连接池大小")
	flag.StringVar(&corsStr, "cors", "", "允许的跨域来源，多个用逗号分隔；传 * 表示允许所有来源（默认不开启 CORS）")
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
	opts := []httpserver.Option{
		httpserver.WithAddr(addr),
		httpserver.WithHosts(hosts...), // 可变参数展开
		httpserver.WithPoolSize(poolSize),
	}
	// 解析 cors 参数：传 * 表示允许所有来源，否则作为具体来源列表
	if corsStr != "" {
		if corsStr == "*" {
			opts = append(opts, httpserver.WithCORS())
		} else {
			var origins []string
			for _, o := range strings.Split(corsStr, ",") {
				if o = strings.TrimSpace(o); o != "" {
					origins = append(origins, o)
				}
			}
			opts = append(opts, httpserver.WithCORS(origins...))
		}
	}

	s, err := httpserver.New(opts...)
	if err != nil {
		log.Fatalf("创建服务失败: %v", err)
	}

	// 启动服务
	log.Printf("服务启动，监听地址: %s，使用 hosts: %v，连接池大小: %d\n", addr, hosts, poolSize)
	if err := s.Run(); err != nil {
		log.Fatalf("服务运行异常: %v", err)
	}
}

package utils

import (
	"bufio"
	"context"
	"fmt"
	"go-study/cache"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var ips = []string{"20.205.243.166", "140.82.114.4"}

func GithubFlush() {
	// 创建channel用于接收IP和延迟结果,错误[3]interface{ip,where_from|0_host&1_dns,err}
	resultsCh := make(chan [3]interface{}, 10)
	var wg sync.WaitGroup
	RCli := cache.RCli()
	ctx := context.Background()
	// 清除排行榜
	RCli.ZDiff(ctx, "github_latency")
	// 从hosts文件读取IP地址
	wg.Add(1)
	go func() {
		defer wg.Done()
		readHostsFile("/etc/hosts", "github.com", resultsCh)
	}()
	// 解析DNS获取IP地址
	wg.Add(1)
	go func() {
		defer wg.Done()
		resolveDNS("github.com", resultsCh)
	}()
	// 自定义ip
	wg.Add(1)
	go func() {
		defer wg.Done()
		resultsCh <- [3]interface{}{nil, 1, nil}
	}()
	// 开始接收和处理结果

	go func() {
		for result := range resultsCh {
			wg.Add(1)
			fmt.Println("add 1 worker")
			fmt.Println("result:", result)
			ipStr := ToString(result[0])
			latency, _ := testLatency(ipStr, &wg)

			if latency != 0 {
				RCli.ZAdd(ctx, "github_latency", &redis.Z{
					Score:  latency,
					Member: ipStr,
				})
			}
		}
	}()

	// 等待所有goroutines完成
	wg.Wait()

	// 获取延迟最小的一条...默认只有五条
	ranking, err := RCli.ZRangeWithScores(ctx, "github_latency", 0, 10).Result()
	if err != nil {
		fmt.Println("ZRangeWithScores error:", err)
	}
	fmt.Printf("Player: %s, Score: %f\n", ranking[0].Member, ranking[0].Score)
}

// 从hosts文件读取指定域名的IP地址并发送到resultsCh
func readHostsFile(hostsFilePath string, domain string, resultsCh chan<- [3]interface{}) {
	file, err := os.Open(hostsFilePath)
	if err != nil {
		resultsCh <- [3]interface{}{nil, 0, err}
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 2 && parts[1] == domain {
			ip := net.ParseIP(parts[0])
			if ip != nil {
				resultsCh <- [3]interface{}{ip, 0, err}
				// testLatencyAndSendResult(ip.String(), resultsCh)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		resultsCh <- [3]interface{}{nil, 0, err}
	}
}

// 解析DNS获取IP地址并发送到resultsCh
func resolveDNS(domain string, resultsCh chan<- [3]interface{}) {
	addrs, err := net.LookupIP(domain)
	if err != nil {
		resultsCh <- [3]interface{}{nil, 1, err}
		return
	}

	for _, addr := range addrs {
		fmt.Println(addr)
		// resultsCh <- [3]interface{}{addr, 1, nil}
		// testLatencyAndSendResult(addr.String(), resultsCh)
	}
}

// 测试IP地址的网络延迟并将结果发送到resultsCh
func testLatencyAndSendResult(ip string, resultsCh chan<- [3]interface{}) {
	latency, err := testLatency(ip, &sync.WaitGroup{})
	if err != nil {
		resultsCh <- [3]interface{}{ip, 0, err}
		return
	}
	resultsCh <- [3]interface{}{ip, latency, nil}
}

// 测试IP地址的网络延迟
func testLatency(ip string, wg *sync.WaitGroup) (float64, error) {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", ip+":80", 3*time.Second)
	if err != nil {
		return 0, err
	}
	fmt.Println("ip:", ip)
	defer func() {
		conn.Close()
		wg.Done()
	}()
	duration := time.Since(start).Seconds() * 1000
	return duration, nil
}

package task

import (
	"net"
	"testing"
)

// 界面文本框按行填、支持 IP:端口 与 # 备注，
// 这里确认内核解析这些写法都不翻车（坏行 panic 是内核的老行为，不在本测范围）
func TestLoadIPRangesFromText(t *testing.T) {
	IPText = "1.1.1.1:2053\n2.2.2.2:443,3.3.3.3;4.4.4.4 # 备注\n8.8.8.8\n"
	PortMapping = make(map[string]int)
	defer func() { IPText = ""; PortMapping = make(map[string]int) }()

	got := loadIPRanges()
	seen := make(map[string]bool, len(got))
	for _, ip := range got {
		seen[ip.String()] = true
	}
	for _, want := range []string{"1.1.1.1", "2.2.2.2", "3.3.3.3", "4.4.4.4", "8.8.8.8"} {
		if !seen[want] {
			t.Errorf("缺 IP %s，解析结果: %v", want, got)
		}
	}
	// IP:端口 的映射要记上
	if PortMapping["1.1.1.1"] != 2053 {
		t.Errorf("1.1.1.1 的端口映射应为 2053，实际 %d", PortMapping["1.1.1.1"])
	}
	_ = net.IPv4
}

// 备注里带逗号 / 分号：# 到行尾整段是注释，不该被当成 IP 分隔符切出脏值
func TestLoadIPRangesFromTextNoteWithSeparator(t *testing.T) {
	IPText = "1.1.1.1:2053 # 北京,上海;香港\n2.2.2.2\n"
	PortMapping = make(map[string]int)
	defer func() { IPText = ""; PortMapping = make(map[string]int) }()

	got := loadIPRanges()
	seen := make(map[string]bool, len(got))
	for _, ip := range got {
		seen[ip.String()] = true
	}
	for _, want := range []string{"1.1.1.1", "2.2.2.2"} {
		if !seen[want] {
			t.Errorf("缺 IP %s，解析结果: %v", want, got)
		}
	}
	// 备注里的「北京」「上海」「香港」必须没有被当成 IP 段解析
	if len(got) != 2 {
		t.Errorf("应只解析出 2 个 IP，实际 %d 个: %v", len(got), got)
	}
}

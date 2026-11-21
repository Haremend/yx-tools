package app

import (
	"testing"
)

func TestRetainedResults(t *testing.T) {
	t.Run("只测延迟时原样保留", func(t *testing.T) {
		rs := []Result{{IP: "1.1.1.1", Port: 443, Delay: 30}, {IP: "2.2.2.2", Port: 443, Delay: 50}}
		if got := retainedResults(rs, nil, true); len(got) != 2 {
			t.Fatalf("want 2, got %d", len(got))
		}
	})

	t.Run("保留所有非 0 速度节点并按速度降序", func(t *testing.T) {
		// 内核子集只给达标的 1 条；流式全量里有 3 条非 0、1 条 0 速度
		rs := []Result{{IP: "3.3.3.3", Port: 443, Speed: 50}}
		streamed := []Result{
			{IP: "3.3.3.3", Port: 443, Speed: 50},
			{IP: "1.1.1.1", Port: 443, Speed: 20},   // 不达标但非 0，要保留
			{IP: "2.2.2.2", Port: 443, Speed: 0},    // 0 速度，丢掉
			{IP: "4.4.4.4", Port: 2053, Speed: 5.5}, // 反代端口，要保留
		}
		got := retainedResults(rs, streamed, false)
		if len(got) != 3 {
			t.Fatalf("want 3, got %d: %+v", len(got), got)
		}
		if got[0].IP != "3.3.3.3" || got[1].IP != "1.1.1.1" || got[2].IP != "4.4.4.4" {
			t.Fatalf("应按速度降序: %+v", got)
		}
	})

	t.Run("同 IP 不同端口不算重复", func(t *testing.T) {
		rs := []Result{}
		streamed := []Result{
			{IP: "1.1.1.1", Port: 443, Speed: 10},
			{IP: "1.1.1.1", Port: 2053, Speed: 20},
		}
		if got := retainedResults(rs, streamed, false); len(got) != 2 {
			t.Fatalf("want 2, got %d", len(got))
		}
	})

	t.Run("流式全量为空时退回内核结果", func(t *testing.T) {
		rs := []Result{{IP: "9.9.9.9", Port: 443, Speed: 1}}
		if got := retainedResults(rs, nil, false); len(got) != 1 || got[0].IP != "9.9.9.9" {
			t.Fatalf("want kernel result, got %+v", got)
		}
	})
}

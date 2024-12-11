package test

import (
	"fmt"
	"github.com/EarlyZhao/hotDetect"
	"testing"
	"time"
)

func callback(list []hotDetect.TopItem) {
	// 处理top-k的热点
	for _, item := range list {
		fmt.Println("count:", item.Freq, ", key:", item.Key)
	}
}

func TestHot(t *testing.T) {
	detectSomething := hotDetect.NewDetect(hotDetect.DefualtConf("test", callback))
	detectSomething.Record("1")
	detectSomething.Record("2")
	detectSomething.Record("3")
	detectSomething.Record("4")
	detectSomething.Record("5")
	detectSomething.Record("1")
	detectSomething.Record("1")
	detectSomething.Record("1")
	detectSomething.Record("1")
	detectSomething.Record("1")
	detectSomething.Record("1")
	detectSomething.Record("1")
	detectSomething.Record("1")
	detectSomething.Record("1")
	detectSomething.Record("1")
	detectSomething.Record("1")
	detectSomething.Record("1")
	detectSomething.Record("1")
	detectSomething.Record("1")
	detectSomething.Record("1")
	detectSomething.Record("1")
	detectSomething.Record("1")

	time.Sleep(time.Minute)
}

package stream

import (
	"fmt"
	"net/url"
	"time"
)

// endpoint
const (
	BinanceSocketUrl = "stream.binance.com:9443"
)

const (
	PingPeriod = 5 * time.Minute
)

// pattern enums
type StreamNameType string

const (
	KlinePattern                 StreamNameType = "%s@kline_%s"
	AggTradePattern              StreamNameType = "%s@aggTrade"
	AllBookTickersPattern        StreamNameType = "!bookTicker"
	AllMiniTickersPattern        StreamNameType = "!miniTicker@arr"
	AllRollingWindowStatsPattern StreamNameType = "!ticker_%s@arr"
	AllTickerPattern             StreamNameType = "!ticker@arr"
	DiffDepthPattern             StreamNameType = "%s@depth"
	BookTickerPattern            StreamNameType = "%s@bookTicker"
	MiniTickerPattern            StreamNameType = "%s@miniTicker"
	RollingWindowStatsPattern    StreamNameType = "%s@ticker_%s"
	TickerPattern                StreamNameType = "%s@ticker"
	PartialBookDepthPattern      StreamNameType = "%s@depth%d"
	TradePattern                 StreamNameType = "%s@trade"
)

type Stream interface {
	Call(url url.URL, c chan interface{})
}

func CreateUrl(tag string) string {
	return "wss://" + BinanceSocketUrl + "/ws/" + tag
}

func CreateStreamName(pattern string, tags ...interface{}) string {
	return fmt.Sprintf(pattern, tags...)
}

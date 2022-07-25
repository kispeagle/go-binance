package stream

import (
	"time"

	"encoding/json"

	"github.com/gorilla/websocket"
	logger "github.com/kispeagle/go-binance/logs"
)

type TradeService struct {
	Symbol string
}

func NewTradeService(symbol string) TradeService {
	s := TradeService{symbol}
	return s
}

func (s TradeService) Call(url string, c chan interface{}) error {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logger.Log.Error(err)
		c <- nil
		return err
	}
	defer conn.Close()

	ticker := time.NewTicker(PingPeriod)

	for {
		select {
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(60))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				logger.Log.Error(err)
				c <- nil
				return err
			}
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				logger.Log.Error(err)
				c <- nil
				return err
			}

			var data Trade
			err = json.Unmarshal(msg, &data)
			if err != nil {
				logger.Log.Error(err)
				c <- nil
				return err
			}
			c <- data
		}
	}
}

type Trade struct {
	E  string `json:"e"`
	E0 int    `json:"E"`
	S  string `json:"s"`
	T  int    `json:"t"`
	P  string `json:"p"`
	Q  string `json:"q"`
	B  int    `json:"b"`
	A  int    `json:"a"`
	T0 int    `json:"T"`
	M  bool   `json:"m"`
	M0 bool   `json:"M"`
}

package stream

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	logger "github.com/kispeagle/go-binance/logs"
)

type AggTradeService struct {
	Symbol string
}

func NewAggTradeService(symbol string) AggTradeService {
	s := AggTradeService{symbol}
	return s
}

func (s AggTradeService) Call(url string, c chan interface{}) error {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logger.Log.Error(err)
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
				return err
			}
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				logger.Log.Error(err)
				return err
			}

			var data AggTrade
			err = json.Unmarshal(msg, &data)
			if err != nil {
				logger.Log.Error(err)
				return err
			}
			c <- data
		}
	}

}

type AggTrade struct {
	E  string `json:"e"`
	E0 int    `json:"E"`
	S  string `json:"s"`
	A  int    `json:"a"`
	P  string `json:"p"`
	Q  string `json:"q"`
	F  int    `json:"f"`
	L  int    `json:"l"`
	T  int    `json:"T"`
	M  bool   `json:"m"`
	M0 bool   `json:"M"`
}

package stream

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	logger "github.com/kispeagle/go-binance/logs"
)

type AllTickersService struct {
	Symbol string
}

func NewAllTickersService(symbol string) AllTickersService {
	s := AllTickersService{symbol}
	return s
}

func (s AllTickersService) Call(url string, c chan interface{}) error {

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

			var data TickerList
			err = json.Unmarshal(msg, &data.Data)
			if err != nil {
				logger.Log.Error(err)
				c <- nil
				return err
			}
			c <- data
		}
	}
}

type TickerList struct {
	Data []Ticker
}

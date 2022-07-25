package stream

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	logger "github.com/kispeagle/go-binance/logs"
)

type RollingWindowStatsService struct {
	Symbol string
}

func NewRollingWindowStatsService(symbol string) RollingWindowStatsService {
	s := RollingWindowStatsService{symbol}
	return s
}

func (s RollingWindowStatsService) Call(url string, c chan interface{}) error {
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

			var data RollingWindowStats
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

type RollingWindowStats struct {
	E  string `json:"e"`
	E0 int    `json:"E"`
	S  string `json:"s"`
	P  string `json:"p"`
	P0 string `json:"P"`
	O  string `json:"o"`
	H  string `json:"h"`
	L  string `json:"l"`
	C  string `json:"c"`
	W  string `json:"w"`
	V  string `json:"v"`
	Q  string `json:"q"`
	O0 int    `json:"O"`
	C0 int    `json:"C"`
	F  int    `json:"F"`
	L0 int    `json:"L"`
	N  int    `json:"n"`
}

package stream

import (
	"encoding/json"
	"strconv"
	"time"

	logger "github.com/kispeagle/go-binance/logs"

	"github.com/gorilla/websocket"
)

type KlineService struct {
	Symbol string
}

func NewKlineService(symbol string) KlineService {
	klineService := KlineService{Symbol: symbol}
	return klineService

}

func (s KlineService) Call(url string, c chan interface{}) error {

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
		// ping check connection
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(60))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				logger.Log.Error(err)
				c <- nil
				return err
			}
		default:
			// read data from socket
			_, msg, err := conn.ReadMessage()
			if err != nil {
				logger.Log.Error(err)
				c <- nil
				return err
			}

			// parse data
			var data Kline
			err = json.Unmarshal(msg, &data)
			if err != nil {
				logger.Log.Error(err)
				c <- nil
				return err
			}
			c <- data
		}
	}

	return nil
}

func processUnixtime(t string) time.Time {
	i, _ := strconv.ParseInt(t, 10, 64)
	i /= 1000
	return time.Unix(i, 0)
}

type Kline struct {
	E  string `json:"e"`
	E0 int    `json:"E"`
	S  string `json:"s"`
	K  struct {
		T  int    `json:"t"`
		T0 int    `json:"T"`
		S  string `json:"s"`
		I  string `json:"i"`
		F  int    `json:"f"`
		L  int    `json:"L"`
		O  string `json:"o"`
		C  string `json:"c"`
		H  string `json:"h"`
		L0 string `json:"l"`
		V  string `json:"v"`
		N  int    `json:"n"`
		X  bool   `json:"x"`
		Q  string `json:"q"`
		V0 string `json:"V"`
		Q0 string `json:"Q"`
		B  string `json:"B"`
	} `json:"k"`
}

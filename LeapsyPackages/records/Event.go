package records

import (
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Event - 事件紀錄
type Event struct {
	Time    time.Time     `form:"time" json:"time"`
	Device  string        `form:"device" json:"device"`
	Method  string        `form:"method" json:"method"`
	URL     string        `form:"url" json:"url"`
	Status  int           `form:"status" json:"status"`
	Results []interface{} `form:"results" json:"results"`
	IP      string        `form:"ip" json:"ip"`
	GPS     string        `form:"gps" json:"gps"`
	Header  http.Header   `json:"header"`
	Input   interface{}   `form:"input" json:"input"`
	Output  interface{}   `form:"output" json:"output"`
}

// PrimitiveM - 轉成primitive.M
/*
 * @return primitive.M returnPrimitiveM 回傳結果
 */
func (event Event) PrimitiveM() (returnPrimitiveM primitive.M) {

	data, _ := bson.Marshal(&event)
	bson.Unmarshal(data, &returnPrimitiveM)

	return
}

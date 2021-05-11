package message

import (
	"encoding/json"
)

type Actor struct {
	Id          int64  `json:"id" DbType:"BIGINT" ColumnName:"id" PK:"LastId"`
	Name        string `json:"name" DbType:"VARCHAR" ColumnName:"name"`
	Career      string `json:"career"`
	Area        string `json:"area" DbType:"VARCHAR" ColumnName:"country"`
	Description string `json:"description" DbType:"VARCHAR" ColumnName:"description"`
	Img         string `json:"img" DbType:"VARCHAR" ColumnName:"image"`
	StaffType   int8   `json:"staff_type" DbType:"INT" ColumnName:"staff_type"`
}
type PlayContent struct {
	Id       int64  `json:"id" DbType:"BIGINT" ColumnName:"id" PK:"LastId"`
	PlayName string `json:"play_name" DbType:"VARCHAR" ColumnName:"title"`
	PlayNum  int    `json:"play_num"`
	VodId    int64  `json:"vod_id" DbType:"BIGINT" ColumnName:"vod_id"`
	PlayType string `json:"play_type" DbType:"VARCHAR" ColumnName:"play_type"`
	PlayUrl  string `json:"play_url" DbType:"VARCHAR" ColumnName:"play_url"`
}

type MovieMessage struct {
	Id            int64         `json:"id" DbType:"BIGINT" ColumnName:"id" PK:"LastId"`
	Name          string        `json:"name" DbType:"VARCHAR" ColumnName:"name"`
	Title         string        `json:"title" DbType:"VARCHAR" ColumnName:"title"`
	Url           string        `json:"url"`
	Img           string        `json:"img" DbType:"VARCHAR" ColumnName:"image"`
	EnName        string        `json:"en_name"`
	MovieType     string        `json:"type"`
	Area          string        `json:"area" DbType:"VARCHAR" ColumnName:"area"`
	Language      string        `json:"language" DbType:"VARCHAR" ColumnName:"language"`
	TotalEpisodes string        `json:"total_episodes" DbType:"VARCHAR" ColumnName:"total_episode"`
	Year          string        `json:"year" DbType:"VARCHAR" ColumnName:"show_year"`
	UpdateTime    string        `json:"update_time" DbType:"VARCHAR" ColumnName:"update_time"`
	Keywords      string        `json:"keywords" DbType:"VARCHAR" ColumnName:"keywords"`
	Description   string        `json:"description" DbType:"VARCHAR" ColumnName:"description"`
	Score         string        `json:"score" DbType:"VARCHAR" ColumnName:"score"`
	DoubanUrl     string        `json:"douban_url"`
	DoubanScore   string        `json:"douban_score" DbType:"VARCHAR" ColumnName:"douban_score"`
	Director      Actor         `json:"director"`
	Actors        []Actor       `json:"actors"`
	PlayUrls      []PlayContent `json:"play_urls"`
}

var _ IMessage = (*MovieMessage)(nil)

func (msg *MovieMessage) DoDecoder(body []byte) error {
	return json.Unmarshal(body, msg)
}

func (msg *MovieMessage) GetMsgCode() int {
	return 0x1000
}

func (msg *MovieMessage) GetDownField() string {
	return "Img"
}

func (msg *MovieMessage) GetTranslateField() string {
	return ""
}
func (msg *MovieMessage) GetSegmentField() string {
	return "Content"
}

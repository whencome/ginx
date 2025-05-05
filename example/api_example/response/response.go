package response

// TimeResponse 时间查询结果
type TimeResponse struct {
	TimeZone string `json:"time_zone" desc:"时区，默认Asia/Shanghai"`
	Time     string `json:"time" desc:"时间，当前时间"`
}

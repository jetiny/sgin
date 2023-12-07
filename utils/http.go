package utils

type PageInfo struct {
	Current   int    `json:"current,omitempty"`
	PageSize  int    `json:"pageSize,omitempty"`
	Total     int64  `json:"total,omitempty"`
	PageToken string `json:"token,omitempty"`
}

func (s PageInfo) Offset() int {
	if s.Current < 1 {
		return 0
	}
	return (s.Current - 1) * s.PageSize
}

type Data[T any] struct {
	Data     T         `json:"data,omitempty"`
	Code     string    `json:"code"`
	Message  string    `json:"message"`
	PageInfo *PageInfo `json:"pageInfo,omitempty"`
}

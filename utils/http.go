package utils

type PageInfo struct {
	PageNumber int    `json:"pageNumber,omitempty"`
	PageSize   int    `json:"pageSize,omitempty"`
	Total      int64  `json:"total,omitempty"`
	PageToken  string `json:"token,omitempty"`
}

func (s PageInfo) Offset() int {
	if s.PageNumber < 1 {
		return 0
	}
	return (s.PageNumber - 1) * s.PageSize
}

type Data[T any] struct {
	Data     T         `json:"data,omitempty"`
	Code     string    `json:"code"`
	Message  string    `json:"message"`
	PageInfo *PageInfo `json:"pageInfo,omitempty"`
}

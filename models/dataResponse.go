package models

type DataResponsePosts struct {
	Data DataPosts `json:"data"`
}

type DataPosts struct {
	Posts []Post `json:"posts"`
}

type DataResponseComments struct {
	Data DataComments `json:"data"`
}

type DataComments struct {
	Comments []Comment `json:"comments"`
}

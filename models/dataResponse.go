package models

type DataResponsePosts struct {
	Data DataPosts `json:"data"`
}

func NewDataResponsePosts(data DataPosts) *DataResponsePosts {
	return &DataResponsePosts{
		Data: data,
	}
}

type DataPosts struct {
	Posts []Post `json:"posts"`
}

func NewDataPosts(posts []Post) *DataPosts {
	return &DataPosts{
		Posts: posts,
	}
}

type DataResponseComments struct {
	Data DataComments `json:"data"`
}

func NewDataResponseComments(data DataComments) *DataResponseComments {
	return &DataResponseComments{
		Data: data,
	}
}

type DataComments struct {
	Comments []Comment `json:"comments"`
}

func NewDataComments(comments []Comment) *DataComments {
	return &DataComments{
		Comments: comments,
	}
}

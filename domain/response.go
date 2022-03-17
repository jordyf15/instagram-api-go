package domain

type Message struct {
	Message string `json:"message"`
}

func NewMessage(message string) *Message {
	return &Message{
		Message: message,
	}
}

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

func NewDataResponseComments(data *DataComments) *DataResponseComments {
	return &DataResponseComments{
		Data: *data,
	}
}

type DataComments struct {
	Comments []Comment `json:"comments"`
}

func NewDataComments(comments *[]Comment) *DataComments {
	return &DataComments{
		Comments: *comments,
	}
}

type DataResponseAuthentication struct {
	Message string             `json:"message"`
	Data    DataAuthentication `json:"data"`
}

func NewDataResponseAuthentication(message string, data DataAuthentication) *DataResponseAuthentication {
	return &DataResponseAuthentication{
		Message: message,
		Data:    data,
	}
}

type DataAuthentication struct {
	AccessToken string `json:"access_token"`
}

func NewDataAuthentication(token string) *DataAuthentication {
	return &DataAuthentication{
		AccessToken: token,
	}
}

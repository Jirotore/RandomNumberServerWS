package internal

type RequestModel[T any] struct {
	Request T `json:"request"`
}

type ResponseModel[T any] struct {
	Response T      `json:"response"`
	ErrorMsg string `json:"error"`
}

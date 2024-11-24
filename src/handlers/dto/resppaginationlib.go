package dto

type RespPaginationLib struct {
	Next  bool   `json:"next"`
	Songs []Song `json:"songs"`
}

package req

type AccountReq struct {
	Number string `json:"document_number" binding:"required"`
}

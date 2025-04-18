package req

type TransactionReq struct {
	AccountID       int     `json:"account_id" binding:"required"`
	OperationTypeId int     `json:"operation_type_id" binding:"required"`
	Amount          float64 `json:"amount" binding:"required"`
}

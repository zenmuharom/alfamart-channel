package dto

type Trx struct {
	UserName         string `json:"userName" binding:"required"`
	Signature        string `json:"signature" binding:"required"`
	ProductCode      string `json:"productCode" binding:"required"`
	Terminal         string `json:"terminal" binding:"required"`
	TerminalName     string `json:"terminalName"`
	TerminalLocation string `json:"terminalLocation"`
	TransactionType  string `json:"transactionType" binding:"required"`
	Channel          string `json:"channel"`
	BillNumber       string `json:"billNumber" binding:"required"`
	Amount           string `json:"amount" binding:"required"`
	FeeAmount        string `json:"feeAmount" binding:"required"`
	Bit32            string `json:"bit32"`
	Bit33            string `json:"bit33"`
	Bit61            string `json:"bit61" binding:"required"`
	TraxId           string `json:"traxId" binding:"required"`
	AddInfo1         string `json:"addInfo1"`
	AddInfo2         string `json:"addInfo2"`
	AddInfo3         string `json:"addInfo3"`
	TimeStamp        string `json:"timeStamp" binding:"required"`
}

type Response struct {
	Bit61           string `json:"bit61"`
	CustomerData    string `json:"customerData"`
	ResultCode      string `json:"resultCode"`
	ResultDesc      string `json:"resultDesc"`
	SysCode         string `json:"sysCode"`
	ProductCode     string `json:"productCode"`
	Terminal        string `json:"terminal"`
	TransactionType string `json:"transactionType"`
	Amount          string `json:"amount"`
	FeeAmount       string `json:"feeAmount"`
	Bit48           string `json:"bit48"`
	TraxId          string `json:"traxId"`
	Timestamp       string `json:"timeStamp"`
	TimestampServer string `json:"timeStampServer"`
}

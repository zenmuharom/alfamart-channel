package dto

type TsReq struct {
	UserName         string `json:"userName"`
	CaCode           string `json:"caCode"`
	SubcaCode        string `json:"subcaCode"`
	ProductCode      string `json:"productCode"`
	Channel          string `json:"channel"`
	Terminal         string `json:"terminal"`
	TerminalName     string `json:"terminalName"`
	TerminalLocation string `json:"terminalLocation"`
	TransactionType  string `json:"transactionType"`
	BillNumber       string `json:"billNumber"`
	Amount           string `json:"amount"`
	FeeAmount        string `json:"feeAmount"`
	Bit61            string `json:"bit61"`
	TraxId           string `json:"traxId"`
	TimeStamp        string `json:"timeStamp"`
	AddInfo1         string `json:"addInfo1"`
	AddInfo2         string `json:"addInfo2"`
	AddInfo3         string `json:"addInfo3"`
}

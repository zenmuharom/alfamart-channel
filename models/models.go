package models

type ErrorMsg struct {
	GoCode int64
	Err    error
}

type Variable struct {
	Field   string              `json:"field"`
	Type    string              `json:"type"`
	Value   any                 `json:"value"`
	As      string              `json:"as"`
	RouteId int64               `json:"route_id"`
	Childs  map[string]Variable `json:"childs"`
	Parent  string              `json:"parent"`
}

type Request struct {
	AgentID         string `json:"AgentID"`
	AgentPIN        string `json:"AgentPIN"`
	AgentTrxID      string `json:"AgentTrxID"`
	AgentStoreID    string `json:"AgentStoreID"`
	ProductID       string `json:"ProductID"`
	CustomerID      string `json:"CustomerID"`
	DateTimeRequest string `json:"DateTimeRequest"`
	Signature       string `json:"Signature"`
}

type InquiryReq struct {
	AgentID         string `json:"AgentID"`
	AgentPIN        string `json:"AgentPIN"`
	AgentTrxID      string `json:"AgentTrxID"`
	AgentStoreID    string `json:"AgentStoreID"`
	ProductID       string `json:"ProductID"`
	CustomerID      string `json:"CustomerID"`
	DateTimeRequest string `json:"DateTimeRequest"`
	Signature       string `json:"Signature"`
}

type PaymentReq struct {
	AgentID         string `json:"AgentID"`
	AgentPIN        string `json:"AgentPIN"`
	AgentTrxID      string `json:"AgentTrxID"`
	AgentStoreID    string `json:"AgentStoreID"`
	ProductID       string `json:"ProductID"`
	CustomerID      string `json:"CustomerID"`
	DateTimeRequest string `json:"DateTimeRequest"`
	PaymentPeriod   string `json:"PaymentPeriod"`
	Amount          string `json:"Amount"`
	Charge          string `json:"Charge"`
	Total           string `json:"Total"`
	AdminFee        string `json:"AdminFee"`
	Signature       string `json:"Signature"`
}

type CommitReq struct {
	AgentID         string `json:"AgentID"`
	AgentPIN        string `json:"AgentPIN"`
	AgentTrxID      string `json:"AgentTrxID"`
	AgentStoreID    string `json:"AgentStoreID"`
	ProductID       string `json:"ProductID"`
	CustomerID      string `json:"CustomerID"`
	DateTimeRequest string `json:"DateTimeRequest"`
	Signature       string `json:"Signature"`
}

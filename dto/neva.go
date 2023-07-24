package dto

import "encoding/xml"

type NevaField struct {
	Type  string `xml:"xsi:type,attr"`
	Value string `xml:",chardata"`
}

type NevaInputToken struct {
	Type     string    `xml:"xsi:type,attr"`
	Username NevaField `xml:"username"`
	Password NevaField `xml:"password"`
	MitraCo  NevaField `xml:"mitraCo"`
}

type NevaInputTransaction struct {
	Type            string    `xml:"xsi:type,attr"`
	Description     NevaField `xml:"description"`
	Dest1Acc        NevaField `xml:"dest1Acc"`
	Dest1Amount     NevaField `xml:"dest1Amount"`
	Dest2Acc        NevaField `xml:"dest2Acc"`
	Dest2Amount     NevaField `xml:"dest2Amount"`
	Notidesc        NevaField `xml:"notiDesc"`
	Notiphone       NevaField `xml:"notiPhone"`
	Phoneno         NevaField `xml:"phoneNo"`
	Sessionid       NevaField `xml:"sessionId"`
	Source1Acc      NevaField `xml:"source1Acc"`
	Source1Amount   NevaField `xml:"source1Amount"`
	Source2Acc      NevaField `xml:"source2Acc"`
	Source2Amount   NevaField `xml:"source2Amount"`
	TransactionType NevaField `xml:"transactionType"`
	Mitraco         NevaField `xml:"mitraCo"`
	Prodcode        NevaField `xml:"prodCode"`
	Billingno       NevaField `xml:"billingNo"`
	Sofcode         NevaField `xml:"sofCode"`
	Sofid           NevaField `xml:"sofId"`
	TraxId          NevaField `xml:"traxId"`
}

type NevaGetTokenReq struct {
	XMLName        xml.Name       `xml:"urn:getToken"`
	EncodingStyle  string         `xml:"soapenv:encodingStyle,attr"`
	NevaInputToken NevaInputToken `xml:"inputToken"`
}

type NevaDebetAccountV2Req struct {
	XMLName              xml.Name             `xml:"urn:debetAccountV2"`
	EncodingStyle        string               `xml:"soapenv:encodingStyle,attr"`
	NevaInputTransaction NevaInputTransaction `xml:"inputTransaction"`
}

type NevaBody struct {
	XMLName xml.Name `xml:"soapenv:Body"`
	Request interface{}
}

type NevaHeader struct {
	XMLName xml.Name `xml:"soapenv:Header"`
}

type NevaRoot struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	Xsi     string   `xml:"xmlns:xsi,attr"`
	Xsd     string   `xml:"xmlns:xsd,attr"`
	Env     string   `xml:"xmlns:soapenv,attr"`
	Urn     string   `xml:"xmlns:urn,attr"`
	Header  NevaHeader
	Body    NevaBody
}

type NevaGetTokenResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		GetTokenResponse struct {
			OutputToken struct {
				SessionId   string `xml:"sessionId"`
				ExpiredDate string `xml:"expiredDate"`
				Status      struct {
					HostRef    string `xml:"hostRef"`
					ResultCode string `xml:"resultCode"`
					ResultDesc string `xml:"resultDesc"`
				} `xml:"status"`
			} `xml:"outputToken"`
		} `xml:"getTokenResponse"`
	} `xml:"Body"`
}

type NevaDebetAccountV2Response struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		DebetAccountV2Response struct {
			OutputTransaction struct {
				PhoneNo   string `xml:"phoneNo"`
				SessionId string `xml:"sessionId"`
				Status    struct {
					HostRef    string `xml:"hostRef"`
					ResultCode string `xml:"resultCode"`
					ResultDesc string `xml:"resultDesc"`
				} `xml:"status"`
				TransactionType string `xml:"transactionType"`
				TraxId          string `xml:"traxId"`
			} `xml:"outputTransaction"`
		} `xml:"debetAccountV2Response"`
	} `xml:"Body"`
}

type NevaDebet struct {
	Username  string
	Password  string
	DestAcc   string
	Amount    string
	PhoneNo   string
	SourceAcc string
	MitraCo   string
	ProdCode  string
	BillingNo string
	TraxId    string
}

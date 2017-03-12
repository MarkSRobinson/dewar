package plaid

import (
	"encoding/json"

	"net/http"
	"strings"
	"io/ioutil"
	"github.com/op/go-logging"
	"time"
	"errors"
)

var log = logging.MustGetLogger("plaid")

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type PlaidCredentials struct {
	Public_key string         `json:"public_key"`
	Secret string         `json:"secret"`
	ClientId string         `json:"client_id"`
}

type Connection struct {
	Url       string
	publicKey string
	ClientId  string
	Secret    string
}

func CreateConnection(url string, publicKey string, clientId string, secret string) Connection {
	conn := Connection{Url:url, publicKey: publicKey, ClientId: clientId, Secret: secret}
	return conn
}

type InstitutionReq struct {
	Public_key string         `json:"public_key"`
	Query      string         `json:"query"`
	Products   []string         `json:"products"`
}

type InstitutionResp struct {
	Institutions []Institutions
	request_id   string         `json:"request_id"`
}

type Institutions struct {
	Has_mfa        bool  `json:"has_mfa"`
	Institution_id string  `json:"institution_id"`
	Mfa            []string   `json:"mfa"`
	Name           string    `json:"name"`
	Products       []string  `json:"products"`
}

type CreateItemReq struct {
	Creds            Credentials `json:"credentials"`
	Institution_id   string `json:"institution_id"`
	Options          []string        `json:"options"`
	Initial_products []string `json:"initial_products"`
	Public_key       string         `json:"public_key"`
}

type GetItemReq struct {
	Access_token string         `json:"access_token"`
	Secret string         `json:"secret"`
	ClientId string         `json:"client_id"`
}

type GetItemResp struct {
	Item Item         `json:"item"`
}

type DeleteItemReq struct {
	Access_token string         `json:"access_token"`
	Secret string         `json:"secret"`
	ClientId string         `json:"client_id"`
}

type Item struct {
	InstitutionId string         `json:"institution_id"`
	ItemId string         `json:"item_id"`
}

type CreateItemResp struct {
	Accounts     []Account        `json:"accounts"`
	Public_token string         `json:"public_token"`
}

type Balances struct {
	Available *int  `json:"available"`
	Current   *int         `json:"current"`
	Limit     *int  `json:"limit"`
}

type Account struct {
	Account_id    string `json:"account_id"`
	Mask          string        `json:"mask"`
	Name          string `json:"name"`
	Balance       Balances `json:"balances"`
	Official_name string         `json:"official_name"`
	Type          string         `json:"type"`
	Subtype       string         `json:"subtype"`
}

type MfaDevice struct {
	DeviceId string         `json:"device_id"`
	Mask string         `json:"mask"`
	Type string         `json:"type"`
}

type MfaSelection struct {
	Answers []string         `json:"answers"`
	Question string         `json:"question"`
}

type MfaRequest struct {
	Devices [] MfaDevice `json:"device_list"`
	Device string         `json:"device"`
	MfaType string         `json:"mfa_type"`
	PublicToken string         `json:"public_token"`
	Questions []string         `json:"questions"`
	Selections []MfaDevice         `json:"selections"`
}

type AccessToken struct {
	Access_token string         `json:"access_token"`
}

type InstitutionLink struct {
	Accounts    []Account
	AccessToken string
}

type InstitutionInfo struct {
	InstitutionId    string    `json:"institution_id"`
	ItemId string    `json:"item_id"`
}

type ExchangeTokenReq struct {
	Secret       string        `json:"secret"`
	Client_id    string `json:"client_id"`
	Public_token string         `json:"public_token"`
}

func convertToAccessToken(conn Connection, public_token string) string {

	createItem := ExchangeTokenReq{}
	createItem.Client_id = conn.ClientId
	createItem.Secret = conn.Secret
	createItem.Public_token = public_token

	bolB, _ := json.Marshal(createItem)

	resp, err := http.Post(conn.Url + "/item/public_token/exchange", "application/json", strings.NewReader(string(bolB)))

	if (err == nil) {
		log.Info("Success")
		body, _ := ioutil.ReadAll(resp.Body)
		s := string(body[:])
		log.Info(s)

		res := AccessToken{}

		json.Unmarshal([]byte(s), &res)

		return res.Access_token
	} else {
		log.Info("Failure")
		log.Info(err)
		return ""
	}
}

func DeleteFinancialLink(conn Connection,  accessToken string) error {

	deleteItemReq := DeleteItemReq{}
	deleteItemReq.Access_token = accessToken
	deleteItemReq.ClientId = conn.ClientId
	deleteItemReq.Secret = conn.Secret

	bolB, _ := json.Marshal(deleteItemReq)

	resp, err := http.Post(conn.Url + "/item/delete", "application/json", strings.NewReader(string(bolB)))

	body, _ := ioutil.ReadAll(resp.Body)
	log.Debug(string(body[:]))


	if resp.StatusCode == 200 {
		return nil
	}

	return err
}

func GetFinancialLink(conn Connection,  accessToken string) (*InstitutionInfo, error) {

	getItemReq := GetItemReq{}
	getItemReq.Access_token = accessToken
	getItemReq.ClientId = conn.ClientId
	getItemReq.Secret = conn.Secret

	bolB, _ := json.Marshal(getItemReq)

	resp, err := http.Post(conn.Url + "/item/get", "application/json", strings.NewReader(string(bolB)))

	if err != nil {
		log.Error("Failure")
		log.Error(err)
		return nil, err;
	}

	body, _ := ioutil.ReadAll(resp.Body)
	log.Debug(string(body[:]))

	log.Debugf("%#v", resp.StatusCode)

	if(resp.StatusCode == 200) {
		res := GetItemResp{}

		err = json.Unmarshal(body, &res)

		log.Debugf("%#v", res)
		log.Debugf("Error: %#v", err)

		return &InstitutionInfo{InstitutionId: res.Item.InstitutionId, ItemId: res.Item.ItemId}, nil
	} else {
		return nil, errors.New("Error!")
	}
}

func CreateFinancialLink(conn Connection, institutionId string, username string, password string) (*InstitutionLink, *MfaRequest) {
	createItem := CreateItemReq{}
	createItem.Institution_id = institutionId
	createItem.Public_key = conn.publicKey
	createItem.Initial_products = []string{"transactions"}
	createItem.Creds.Username = username
	createItem.Creds.Password = password

	bolB, _ := json.Marshal(createItem)

	log.Info(string(bolB[:]))

	resp, err := http.Post(conn.Url + "/link/item/create", "application/json", strings.NewReader(string(bolB)))

	if err != nil {
		log.Error("Failure")
		log.Error(err)
		return nil, nil;
	}

	body, _ := ioutil.ReadAll(resp.Body)
	log.Debug(string(body[:]))

	switch(resp.StatusCode) {
	case 200:
		res := CreateItemResp{}

		json.Unmarshal(body, &res)

		return &InstitutionLink{res.Accounts, convertToAccessToken(conn, res.Public_token)}, nil
	case 210:
		res := MfaRequest{}

		json.Unmarshal(body, &res)

		return nil, &res
		break
	default:
		break
	}

	return nil, nil
}

func SearchInstitutions(conn Connection, query string) *InstitutionResp {
	createItem := InstitutionReq{}
	createItem.Query = query
	createItem.Public_key = conn.publicKey
	createItem.Products = []string{"transactions"}

	bolB, _ := json.Marshal(createItem)

	resp, err := http.Post(conn.Url + "/institutions/search", "application/json", strings.NewReader(string(bolB)))

	if (err == nil) {
		log.Info("Success")
		body, _ := ioutil.ReadAll(resp.Body)
		s := string(body[:])

		res := InstitutionResp{}

		json.Unmarshal([]byte(s), &res)

		log.Debug(s)
		return &res
	} else {
		log.Info("Failure")
		log.Debug(err)
		return nil
	}
}

type TransactionReq struct {
	Secret       string        `json:"secret"`
	Client_id    string `json:"client_id"`
	Access_token string         `json:"access_token"`
	StartDate    string         `json:"start_date"`
	EndDate      string         `json:"end_date"`
}

type PaymentMetadata struct {
	ByOrderOf string        `json:"by_order_of"`
	Payee     string        `json:"payee"`
	Payer     string        `json:"payer"`
	Method    string        `json:"payment_method"`
	Processor string        `json:"payment_processor"`
	PpdId     string        `json:"ppd_id"`
	Reason    string        `json:"reason"`
	Reference string        `json:"reference_number"`
}

type TransactionItem struct {
	AccountId       string        `json:"account_id"`
	//Category string `json:"category"`
	Amount            json.Number         `json:"amount"`
	Date            string         `json:"date"`
	Name            string         `json:"name"`
	TransactionType string         `json:"transaction_type"`
	TransactionId   string         `json:"transaction_id"`
	Metadata        PaymentMetadata         `json:"payment_meta"`
}

type TransactionResp struct {
	Transactions     []TransactionItem         `json:"transactions"`
	TransactionCount int         `json:"total_transactions"`
	Accounts         []Account         `json:"accounts"`
}

func GetTransactions(conn Connection, accessToken string, startDate string, endDate string) *TransactionResp {
	createItem := TransactionReq{}
	createItem.Client_id = conn.ClientId
	createItem.Secret = conn.Secret
	createItem.Access_token = accessToken
	createItem.StartDate = startDate
	createItem.EndDate = endDate

	bolB, _ := json.Marshal(createItem)

	log.Info(string(bolB[:]))

	for true {

		resp, err := http.Post(conn.Url + "/transactions/get", "application/json", strings.NewReader(string(bolB)))

		if (err == nil) {
			log.Debugf("Status Code: %v", resp.Status)

			switch(resp.StatusCode) {
			case 200:
				body, _ := ioutil.ReadAll(resp.Body)
				s := string(body[:])

				res := TransactionResp{}

				json.Unmarshal([]byte(s), &res)

				log.Info(s)

				return &res
			case 400:
				time.Sleep(10 * time.Second/time.Nanosecond)
				break
			case 500:
				break
			default:
				log.Errorf("Unkown error code %v", resp.StatusCode)
			}

		} else {
			log.Info("Failure")
			log.Info(err)
		}
	}

	return nil
}
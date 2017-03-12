package plaid

import (
	"testing"
	"io/ioutil"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

var baseUrl string = "https://sandbox.plaid.com"

func getConnection() Connection {

	dat, _ := ioutil.ReadFile("testingcredentials.json")

	jsonResp := string(dat)

	res := PlaidCredentials{}

	json.Unmarshal([]byte(jsonResp), &res)

	return	CreateConnection(baseUrl, res.Public_key, res.ClientId, res.Secret)
}

func TestSearchInstitutions(t *testing.T) {

	conn := getConnection()

	institutionResp := SearchInstitutions(conn, "F")

	assert.NotNil(t, institutionResp)
	assert.Equal(t, len(institutionResp.Institutions), 10)
}

func TestCreateFinancialLink(t *testing.T) {

	conn := getConnection()

	institutionLink, mfaRequest := CreateFinancialLink(conn, "ins_100069", "user_good", "pass_good")

	assert.NotNil(t, institutionLink)
	assert.Nil(t, mfaRequest)
	assert.NotEmpty(t, institutionLink.AccessToken)
	assert.Equal(t, 4, len(institutionLink.Accounts))
}

func TestGetFinancialLink(t *testing.T) {

	conn := getConnection()

	institutionLink, mfaRequest := CreateFinancialLink(conn, "ins_100069", "user_good", "pass_good")

	assert.NotNil(t, institutionLink)
	assert.Nil(t, mfaRequest)

	institutionInfo,err := GetFinancialLink(conn, institutionLink.AccessToken)

	assert.Nil(t, err)
	assert.NotNil(t, institutionInfo)
	assert.NotEmpty(t, institutionInfo.ItemId)
	assert.NotEmpty(t, institutionInfo.InstitutionId)
}

func TestGetFinancialLink_InvalidAccessToken(t *testing.T) {

	conn := getConnection()

	institutionInfo,err := GetFinancialLink(conn, "access-sandbox-0ae4f3f5-8994-415c-bb18-b0fb3db9ef4f")

	assert.NotNil(t, err)
	assert.Nil(t, institutionInfo)
}

func TestDeleteFinancialLink(t *testing.T) {

	conn := getConnection()

	institutionLink, mfaRequest := CreateFinancialLink(conn, "ins_100069", "user_good", "pass_good")

	assert.NotNil(t, institutionLink)
	assert.Nil(t, mfaRequest)

	err := DeleteFinancialLink(conn, institutionLink.AccessToken)

	assert.Nil(t, err)

	institutionInfo,err := GetFinancialLink(conn, institutionLink.AccessToken)

	assert.Nil(t, institutionInfo)
	assert.NotNil(t, err)
}

func TestCreateFinancialLinkMfaDevice(t *testing.T) {

	conn := getConnection()

	institutionLink, mfaRequest := CreateFinancialLink(conn, "ins_100069", "user_good", "mfa_device")

	assert.Nil(t, institutionLink)
	assert.NotNil(t, mfaRequest)
	assert.NotNil(t, mfaRequest.PublicToken)
	assert.Equal(t, 4, len(mfaRequest.Devices))
	assert.Equal(t, "device_list", mfaRequest.MfaType)
}

func TestCreateFinancialLinkMfaQuestions(t *testing.T) {

	conn := getConnection()

	institutionLink, mfaRequest := CreateFinancialLink(conn, "ins_100069", "user_good", "mfa_questions_2_2")

	assert.Nil(t, institutionLink)
	assert.NotNil(t, mfaRequest)
	assert.NotNil(t, mfaRequest.PublicToken)
}

func TestCreateFinancialLinkMfaSelections(t *testing.T) {

	conn := getConnection()

	institutionLink, mfaRequest := CreateFinancialLink(conn, "ins_100069", "user_good", "mfa_selections")

	assert.Nil(t, institutionLink)
	assert.NotNil(t, mfaRequest)
	assert.NotNil(t, mfaRequest.PublicToken)
}

func TestGetTransactionsPath(t *testing.T) {

	conn := getConnection()

	institutionLink, mfaRequest := CreateFinancialLink(conn, "ins_100069", "user_good", "pass_good")

	assert.NotNil(t, institutionLink)
	assert.Nil(t, mfaRequest)
	assert.NotEmpty(t, institutionLink.AccessToken)
	assert.Equal(t, 4, len(institutionLink.Accounts))

	transactions := GetTransactions(conn, institutionLink.AccessToken,
		"2017-01-26", "2017-03-01")

	assert.NotNil(t, transactions)
}

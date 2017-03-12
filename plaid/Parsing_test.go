package plaid

import (
	"testing"
	"io/ioutil"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

func TestParseInstitutions(t *testing.T) {

	dat, _ := ioutil.ReadFile("searchInstitution.json")

	jsonResp := string(dat)

	res := InstitutionResp{}

	err := json.Unmarshal([]byte(jsonResp), &res)

	assert.Nil(t, err)

	assert.Equal(t, 10, len(res.Institutions), "they should be equal")
	for i := 0; i < len(res.Institutions); i++ {
		inst := res.Institutions[i]
		assert.NotEmpty(t, inst.Name, "Institution missing name")
		assert.NotEmpty(t, inst.Institution_id, "Institution missing id")
		assert.True(t, len(inst.Products) > 0, "Institution has no products")
	}
}

func TestParseItem(t *testing.T) {
	dat, _ := ioutil.ReadFile("createItem.json")

	jsonResp := string(dat)

	res := CreateItemResp{}

	err := json.Unmarshal([]byte(jsonResp), &res)

	assert.Nil(t, err)

	assert.Equal(t, 4, len(res.Accounts), "they should be equal")
	for i := 0; i < len(res.Accounts); i++ {
		inst := res.Accounts[i]
		assert.NotEmpty(t, inst.Name, "Account missing name")
		assert.NotEmpty(t, inst.Official_name, "Account missing official_name")
		assert.NotEmpty(t, inst.Mask, "Account has no mask")
		assert.NotEmpty(t, inst.Type, "Account has no Type")
		assert.NotEmpty(t, inst.Subtype, "Account has no Subtype")
		//
		assert.True(t, *inst.Balance.Current > 0, "Institution has no products")
		if(inst.Type == "depository") {
			if(inst.Subtype == "cd") {
				assert.Nil(t, inst.Balance.Available, "CD account has available balance")
			} else {
				assert.NotNil(t, inst.Balance.Available, "Depository account has no available balance")
			}
		} else if(inst.Type == "credit") {
			assert.True(t, *inst.Balance.Limit > 0, "Credit account has no limit")
			assert.Nil(t, inst.Balance.Available, "Credit account has available balance")
		}
	}
}

func TestParseTransactions(t *testing.T) {
	dat, _ := ioutil.ReadFile("gettransactions.json")

	jsonResp := string(dat)

	res := TransactionResp{}

	err := json.Unmarshal([]byte(jsonResp), &res)

	assert.Nil(t, err)

	assert.Equal(t, 4, len(res.Accounts), "they should be equal")
	assert.Equal(t, 21, len(res.Transactions), "they should be equal")
	assert.Equal(t, 21, res.TransactionCount, "they should be equal")
	for i := 0; i < len(res.Transactions); i++ {
		transaction := res.Transactions[i]
		assert.NotEmpty(t, transaction.AccountId)
		assert.NotEmpty(t, transaction.Amount)
		assert.NotEmpty(t, transaction.Date)
		assert.NotEmpty(t, transaction.Name)
		assert.NotEmpty(t, transaction.TransactionType)
		assert.NotEmpty(t, transaction.TransactionId)
	}
}

func TestParseMfaRequestDevice(t *testing.T) {
	dat, _ := ioutil.ReadFile("mfa_request_device.json")

	res := MfaRequest{}

	err := json.Unmarshal(dat, &res)

	assert.Nil(t, err)
	assert.Equal(t, res.MfaType, "device_list")
	assert.Equal(t, 4, len(res.Devices))
	assert.NotEmpty(t, res.MfaType)
	assert.NotEmpty(t, res.PublicToken)
}

func TestParseGetItem(t *testing.T) {
	dat, _ := ioutil.ReadFile("get_item.json")

	res := GetItemResp{}

	err := json.Unmarshal(dat, &res)

	assert.Nil(t, err)
	assert.NotEmpty(t, res.Item.ItemId)
}

func TestParseMfaRequestSelection(t *testing.T) {
	dat, _ := ioutil.ReadFile("mfa_request_selection.json")

	res := MfaRequest{}

	err := json.Unmarshal(dat, &res)

	assert.Nil(t, err)
	assert.Equal(t, res.MfaType, "selections")
	assert.Nil(t, res.Devices)
	assert.NotEmpty(t, res.PublicToken)
}

func TestParseMfaRequestQuestions(t *testing.T) {
	dat, _ := ioutil.ReadFile("mfa_request_questions.json")

	res := MfaRequest{}

	err := json.Unmarshal(dat, &res)

	assert.Nil(t, err)
	assert.Nil(t, res.Devices)
	assert.Equal(t, res.MfaType, "questions")
	assert.NotEmpty(t, res.PublicToken)
}
package caledonia

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"bytes"
	"io/ioutil"
	"github.com/sergi/go-diff/diffmatchpatch"
	"fmt"
	"strings"
)

func verifyInterfaceOfxWriter(ofxWriter OfxWriter) {

}

func TestImplementsApi(t *testing.T) {

	ofxWriter := OfxV1Writer{}

	verifyInterfaceOfxWriter(&ofxWriter)
	assert.True(t, true)
}

func TestWrite_CreditCard(t *testing.T) {

	ofxWriter := NewOfxV1Writer()

	ofxWriter.AddAccount(Account{AccountId:"12345678910241877", Mask:"12345678910241877", Name:"Credit Card",
	Type:"credit", Subtype:"credit"})

	ofxWriter.AddTransaction("12345678910241877", Transaction{   TransactionType: "CREDIT",
		Posted: "20170228170000.000",
		Amount: "294.02",
		Fitid: "FITID20170228294.02F7D86",
		Name: "PAYMENT - THANK YOU"})

	ofxWriter.AddTransaction("12345678910241877", Transaction{   TransactionType: "DEBIT",
		Posted: "20170301170000.000",
		Amount: "-10",
		Fitid: "FITID20170303-10.0FRED8",
		Name: "VISIT PATREON.COM/INFO 877-887-7"})

	ofxWriter.AddTransaction("12345678910241877", Transaction{   TransactionType: "DEBIT",
		Posted: "20170304170000.000",
		Amount: "-17.16",
		Fitid: "FITID20170305-17.16FEJR8",
		Name: "AMZ*TINGTX#123456 866-216-1072 W",
		Memo: "PEFT2IYOYK3"})

	ofxWriter.AddTransaction("12345678910241877", Transaction{   TransactionType: "DEBIT",
		Posted: "20170305170000.000",
		Amount: "-9.75",
		Fitid: "FITID20170306-9.75FDFYG",
		Name: "AMAZON MKTPLACE PMTS AMZN.COM/BI",
		Memo: "QYN7L8LIK18",
	})

	ofxWriter.AddTransaction("12345678910241877", Transaction{   TransactionType: "DEBIT",
		Posted: "20170307170000.000",
		Amount: "-60",
		Fitid: "FITID20170308-60.0DPZY4",
		Name: "CLIPPER SERVICE CONCORD CA",
	})


	buf := new(bytes.Buffer)
	ofxWriter.Write(buf)

	assert.NotEmpty(t, buf.String())

	dat, _ := ioutil.ReadFile("creditcard_transaction.ofx")

	ofxFile := string(dat)

	unpretty := strings.Replace(ofxFile, "    ", "", -1)
	assert.Equal(t, unpretty, buf.String())


	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(buf.String(), unpretty, false)
	assert.Equal(t, 1, len(diffs))
	assert.Equal(t, diffmatchpatch.DiffEqual, diffs[0].Type)

	if(len(diffs) > 0) {
		if(len(diffs) > 1 || diffs[0].Type != diffmatchpatch.DiffEqual) {
			fmt.Println(diffUncoloured(diffs))
		}
	}
}

func TestWrite_BankAccount(t *testing.T) {
	ofxWriter := NewOfxV1Writer()

	ofxWriter.AddAccount(Account{AccountId:"12345678910241877", Mask:"12345678910241877", Name:"Chequing",
		Type:"depository", Subtype:"checking"})


	ofxWriter.AddTransaction("12345678910241877", Transaction{   TransactionType: "DEBIT",
		Posted: "20170308000000",
		Amount: "-1350.0000",
		Fitid: "2017030806391731675000001000",
		Name: " ACH DEBIT- FID BKG SVC LLC",
		Memo: " ACH DEBIT- FID BKG SVC LLC DESC:MO",
	})


	ofxWriter.AddTransaction("12345678910241877", Transaction{   TransactionType: "CREDIT",
		Posted: "20170306000000",
		Amount: "1350.0000",
		Fitid: "2017030618200107073000136000",
		Name: " ATM Deposit PATELCO CREDIT  124",
		Memo: " ATM Deposit PATELCO CREDIT  124 2ND ST SAN FRANCISCO",
	})



	ofxWriter.AddTransaction("12345678910241877", Transaction{   TransactionType: "DEBIT",
		Posted: "20170217000000",
		Amount: "-1250.0000",
		Fitid: "2017021706572721405000001000",
		Name: " ACH DEBIT- FID BKG SVC LLC",
		Memo: " ACH DEBIT- FID BKG SVC LLC DESC:MO",
	})

	ofxWriter.AddTransaction("12345678910241877", Transaction{   TransactionType: "CREDIT",
		Posted: "20170206000000",
		Amount: "1250.0000",
		Fitid: "2017020609335306326000126000",
		Name: " Deposit",
	})

	ofxWriter.AddTransaction("12345678910241877", Transaction{   TransactionType: "DEBIT",
		Posted: "20170110000000",
		Amount: "-1350.0000",
		Fitid: "2017011006552810719000001000",
		Name: " ACH DEBIT- FID BKG SVC LLC",
		Memo: " ACH DEBIT- FID BKG SVC LLC DESC:MO",
	})

	ofxWriter.AddTransaction("12345678910241877", Transaction{   TransactionType: "CREDIT",
		Posted: "20170106000000",
		Amount: "1350.0000",
		Fitid: "2017010615324015847000136000",
		Name: " Deposit",
	})



	buf := new(bytes.Buffer)
	ofxWriter.Write(buf)

	assert.NotEmpty(t, buf.String())

	dat, _ := ioutil.ReadFile("bank_transaction.ofx")

	ofxFile := string(dat)

	unpretty := strings.Replace(ofxFile, "    ", "", -1)
	assert.Equal(t, unpretty, buf.String())


	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(buf.String(), unpretty, false)
	assert.Equal(t, 1, len(diffs))
	assert.Equal(t, diffmatchpatch.DiffEqual, diffs[0].Type)

	if(len(diffs) > 0) {
		if(len(diffs) > 1 || diffs[0].Type != diffmatchpatch.DiffEqual) {
			fmt.Println(diffUncoloured(diffs))
		}
	}
}

func diffUncoloured(diffs []diffmatchpatch.Diff) string {
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := diff.Text

		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			_, _ = buff.WriteString("\n+++++++++\n")
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("\n>>>>>>>\n")
		case diffmatchpatch.DiffDelete:
			_, _ = buff.WriteString("--------\n")
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("\n<<<<<<<<\n")
		case diffmatchpatch.DiffEqual:
			_, _ = buff.WriteString(text)
		}
	}

	return buff.String()
}
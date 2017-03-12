package caledonia

import (
	"io"
	"github.com/op/go-logging"
	"text/tabwriter"
	"os"
	"fmt"
	"strconv"
)

var log = logging.MustGetLogger("plaid")

type OfxV1Writer struct {
	transactions map[string][]Transaction
	accounts map[string]Account
}

func NewOfxV1Writer() OfxV1Writer {
	writer := OfxV1Writer{}

	writer.transactions = make(map[string][]Transaction, 0)
	writer.accounts = make(map[string]Account, 0)

	return writer
}

func (w *OfxV1Writer) AddTransaction(accountId string, transaction Transaction) {

	if(w.transactions[accountId] == nil) {
		w.transactions[accountId] = make([]Transaction, 0)
	}

	w.transactions[accountId] = append(w.transactions[accountId], transaction)
}

func (w *OfxV1Writer) AddAccount(account Account) {
	w.accounts[account.AccountId] = account
}

func (w *OfxV1Writer) writeDepository(writer io.Writer, accountId string) {
	writer.Write([]byte("<BANKMSGSRSV1>" + lineTerminator))
	writer.Write([]byte("<STMTTRNRS>" + lineTerminator))

	writer.Write([]byte("<TRNUID>0" + lineTerminator))
	writer.Write([]byte("<STATUS>" + lineTerminator))
	writer.Write([]byte("<CODE>0" + lineTerminator))
	writer.Write([]byte("<SEVERITY>INFO" + lineTerminator))
	writer.Write([]byte("</STATUS>" + lineTerminator))
	writer.Write([]byte("<STMTRS>" + lineTerminator))
	writer.Write([]byte("<CURDEF>USD" + lineTerminator))

	writer.Write([]byte("<BANKACCTFROM>" + lineTerminator))
	writer.Write([]byte("<ACCTID>" + w.accounts[accountId].Name +"-" + w.accounts[accountId].Mask + lineTerminator)) // TODO
	writer.Write([]byte("</BANKACCTFROM>" + lineTerminator))

	writer.Write([]byte("<BANKTRANLIST>" + lineTerminator))

	writer.Write([]byte("<DTSTART>20161214000000" + lineTerminator))
	writer.Write([]byte("<DTEND>20170314000000" + lineTerminator))

	for i := 0; i < len(w.transactions[accountId]); i++ {
		transaction := w.transactions[accountId][i]
		// do something with e.Value
		writer.Write([]byte("<STMTTRN>"+lineTerminator))

		f, _ := strconv.ParseFloat(transaction.Amount, 64)

		f = -f

		writer.Write([]byte("<TRNTYPE>CREDIT" + lineTerminator))
		writer.Write([]byte("<DTPOSTED>"+ transaction.Posted + lineTerminator))
		writer.Write([]byte("<TRNAMT>"+ strconv.FormatFloat(f, 'f', -1, 64) + lineTerminator))
		writer.Write([]byte("<FITID>"+ transaction.Fitid + lineTerminator))
		writer.Write([]byte("<NAME>"+ transaction.Name + lineTerminator))

		if transaction.Memo != "" {
			writer.Write([]byte("<MEMO>" + transaction.Memo + lineTerminator))
		}

		writer.Write([]byte("</STMTTRN>"+lineTerminator))
	}
	writer.Write([]byte("</BANKTRANLIST>" + lineTerminator))

	writer.Write([]byte("</STMTRS>" + lineTerminator))
	writer.Write([]byte("</STMTTRNRS>" + lineTerminator))
	writer.Write([]byte("</BANKMSGSRSV1>" + lineTerminator))
}

func (w *OfxV1Writer) writeCredit(writer io.Writer, accountId string) {

	writer.Write([]byte("<CREDITCARDMSGSRSV1>" + lineTerminator))
	writer.Write([]byte("<CCSTMTTRNRS>" + lineTerminator))

	writer.Write([]byte("<TRNUID>0" + lineTerminator))
	writer.Write([]byte("<STATUS>" + lineTerminator))
	writer.Write([]byte("<CODE>0" + lineTerminator))
	writer.Write([]byte("<SEVERITY>INFO" + lineTerminator))
	writer.Write([]byte("</STATUS>" + lineTerminator))
	writer.Write([]byte("<CCSTMTRS>" + lineTerminator))
	writer.Write([]byte("<CURDEF>USD" + lineTerminator))

	writer.Write([]byte("<CCACCTFROM>" + lineTerminator))
	writer.Write([]byte("<ACCTID>" + w.accounts[accountId].Name +"-" + w.accounts[accountId].Mask + lineTerminator)) // TODO
	writer.Write([]byte("</CCACCTFROM>" + lineTerminator))

	writer.Write([]byte("<BANKTRANLIST>" + lineTerminator))

	writer.Write([]byte("<DTSTART>20170308190000.000[-5:EST]" + lineTerminator))
	writer.Write([]byte("<DTEND>20170310120000.000[-5:EST]" + lineTerminator))

	for i := 0; i < len(w.transactions[accountId]); i++ {
		transaction := w.transactions[accountId][i]
		// do something with e.Value
		writer.Write([]byte("<STMTTRN>"+lineTerminator))

		f, _ := strconv.ParseFloat(transaction.Amount, 64)

		f = -f

		//writer.Write([]byte("<TRNTYPE>"+ transaction.TransactionType + lineTerminator))
		writer.Write([]byte("<TRNTYPE>CREDIT" + lineTerminator))
		writer.Write([]byte("<DTPOSTED>"+ transaction.Posted + lineTerminator))
		writer.Write([]byte("<TRNAMT>"+ strconv.FormatFloat(f, 'f', -1, 64) + lineTerminator))
		writer.Write([]byte("<FITID>"+ transaction.Fitid + lineTerminator))
		writer.Write([]byte("<NAME>"+ transaction.Name + lineTerminator))

		if transaction.Memo != "" {
			writer.Write([]byte("<MEMO>" + transaction.Memo + lineTerminator))
		}

		writer.Write([]byte("</STMTTRN>"+lineTerminator))
	}
	writer.Write([]byte("</BANKTRANLIST>" + lineTerminator))

	writer.Write([]byte("</CCSTMTRS>" + lineTerminator))
	writer.Write([]byte("</CCSTMTTRNRS>" + lineTerminator))
	writer.Write([]byte("</CREDITCARDMSGSRSV1>" + lineTerminator))

}

var lineTerminator = "\r\n"

func (w *OfxV1Writer) Debug() {
	writer := tabwriter.NewWriter(os.Stdout, 8, 1, 1, ' ', 0)

	writer.Write([]byte(fmt.Sprintf("\n\tName \t Amount\t Type \t SubType \n")))

	for k := range w.accounts {
		for i := 0; i < len(w.transactions[k]); i++ {

			inst := w.accounts[k]
			xact := w.transactions[k][i]

			writer.Write([]byte(fmt.Sprintf("\t%s \t %v \t %v \t %v\n", xact.Name, xact.Amount, inst.Type, inst.Subtype)))
		}
	}

	writer.Flush()
}

func (w *OfxV1Writer) Write(writer io.Writer) {

	writer.Write([]byte("OFXHEADER:100\r\n"))
	writer.Write([]byte("DATA:OFXSGML\r\n"))
	writer.Write([]byte("VERSION:102\r\n"))
	writer.Write([]byte("SECURITY:NONE\r\n"))
	writer.Write([]byte("ENCODING:USASCII\r\n"))
	writer.Write([]byte("CHARSET:1252\r\n"))
	writer.Write([]byte("COMPRESSION:NONE\r\n"))
	writer.Write([]byte("OLDFILEUID:NONE\r\n"))
	writer.Write([]byte("NEWFILEUID:NONE\r\n"))
	writer.Write([]byte("\r\n"))

	writer.Write([]byte("<OFX>" + lineTerminator))

	writer.Write([]byte("<SIGNONMSGSRSV1>" + lineTerminator))
	writer.Write([]byte("<SONRS>" + lineTerminator))
	writer.Write([]byte("<STATUS>" + lineTerminator))
	writer.Write([]byte("<CODE>0" + lineTerminator))
	writer.Write([]byte("<SEVERITY>INFO" + lineTerminator))
	writer.Write([]byte("</STATUS>" + lineTerminator))
	writer.Write([]byte("<LANGUAGE>ENG" + lineTerminator))

	writer.Write([]byte("</SONRS>" + lineTerminator))
	writer.Write([]byte("</SIGNONMSGSRSV1>" + lineTerminator))

	for k := range w.accounts {
		switch(w.accounts[k].Type) {
		case "credit":
			w.writeCredit(writer, k)
			break
		case "depository":
			w.writeDepository(writer, k)
			break
		default:
			log.Errorf("Unimplemented account type %s", w.accounts[k].Type)
		}
	}

	writer.Write([]byte("</OFX>"+lineTerminator))
}
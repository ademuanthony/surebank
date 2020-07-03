package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"merryworld/surebank/internal/account"
	"merryworld/surebank/internal/customer"
	"merryworld/surebank/internal/platform/notify"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/postgres/models"

	"github.com/gorilla/schema"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"
)

// BulkSMS represents the endpoint for sending SMS notification to customers
type BulkSMS struct {
	CustomerRepo *customer.Repository
	AccountRepo  *account.Repository
	NotifySMS    notify.SMS
	Renderer     web.Renderer
	Redis        *redis.Client
	DbConn       *sql.DB
}

type sendSMSRequest struct {
	Message      string
	AccountType  string
	PhoneNumbers string
	AccountNos   string
}

func (h BulkSMS) Index(ctx context.Context, w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	//
	data := make(map[string]interface{})
	req := new(sendSMSRequest)
	data["form"] = req
	f := func() (bool, error) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				return false, err
			}

			decoder := schema.NewDecoder()
			decoder.IgnoreUnknownKeys(true)

			if err := decoder.Decode(req, r.PostForm); err != nil {
				return false, err
			}

			var recipients = map[string]string{}
			if req.PhoneNumbers != "" {
				numbers := strings.Split(req.PhoneNumbers, ",")
				for _, number := range numbers {
					recipients[number] = "Customer"
				}
			}
			if req.AccountNos != "" {
				numbers := strings.Split(req.AccountNos, ",")
				accounts, err := models.Accounts(
					qm.Load(models.AccountRels.Customer),
					models.AccountWhere.Number.IN(numbers)).All(ctx, h.DbConn)
				if err != nil {
					if err.Error() != sql.ErrNoRows.Error() {
						return false, err
					}
				}
				for _, acc := range accounts {
					recipients[acc.R.Customer.PhoneNumber] = acc.R.Customer.Name
				}
			}
			if req.AccountType != "" {
				var queries = []qm.QueryMod{
					qm.Load(models.AccountRels.Customer),
				}
				if req.AccountType != "all" {
					queries = append(queries, models.AccountWhere.AccountType.EQ(req.AccountType))
				}
				accounts, err := models.Accounts(queries...).All(ctx, h.DbConn)
				if err != nil {
					if err.Error() != sql.ErrNoRows.Error() {
						return false, err
					}
				}
				for _, acc := range accounts {
					recipients[acc.R.Customer.PhoneNumber] = acc.R.Customer.Name
				}
			}
			go func() {
				for number, name := range recipients {
					var message = strings.ReplaceAll(req.Message, "@name", name)
					if err = h.NotifySMS.SendStr(ctx, number, message); err != nil {
						// TODO: log critical error. Send message to monitoring account
						fmt.Println(err)
					}
				}
			}()
			return true, nil
		}
		return false, nil
	}

	success, err := f()
	if err != nil {
		data["form"] = req
		data["error"] = err.Error()
	}

	if success {
		data["message"] = "Messages sent to server successfully"
	}

	data["accountTypes"] = customer.AccountTypes

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "sms-send.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

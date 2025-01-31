package main

import (
	"github.com/diki-haryado/go-payment-proxy/util/localconfig"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	config, err := localconfig.LoadConfig("example/server/config.yaml")
	if err != nil {
		panic(err)
	}

	secret, err := localconfig.LoadSecret("example/server/secret.yaml")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("example/server/gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	db.AutoMigrate(
		&midtrans.TransactionStatus{},
		&invoice.Invoice{},
		&invoice.Payment{},
		&invoice.CreditCardDetail{},
		&invoice.LineItem{},
		&invoice.BillingAddress{},
		&subscription.Subscription{},
		&subscription.Schedule{},
	)

	m := manage.NewManager(*config, secret.Payment)
	m.MustMidtransTransactionStatusRepository(dssql.NewMidtransTransactionStatusRepository(db))
	m.MustInvoiceRepository(dssql.NewInvoiceRepository(db))
	m.MustSubscriptionRepository(dssql.NewSubscriptionRepository(db))
	m.MustPaymentConfigReader(inmemory.NewPaymentConfigRepository("example/server/payment-methods.yaml"))

	srv := srv{
		Router: mux.NewRouter(),
	}

}

type srv struct {
	Router     *mux.Router
	PaymentSrv *server.Server
}

func (s srv) GetHandler() http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:     []string{"http://localhost:3000", "https://localhost:3000"},
		AllowedMethods:     []string{"POST", "GET", "PUT", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Mode"},
		MaxAge:             60, // 1 minutes
		AllowCredentials:   true,
		OptionsPassthrough: false,
		Debug:              false,
	})

	return c.Handler(s.Router)
}

func (s *srv) Healthcheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}
}

func (s srv) routes() {
	s.Router.HandleFunc("/payment/methods", s.paymentSrv.GetPaymentMethodsHandler()).Methods("GET")
	s.Router.HandleFunc("/payment/invoices", s.paymentSrv.CreateInvoiceHandler()).Methods("POST")
	s.Router.HandleFunc("/payment/midtrans/callback", s.paymentSrv.MidtransTransactionCallbackHandler()).Methods("POST")
	s.Router.HandleFunc("/payment/xendit/invoice/callback", s.paymentSrv.XenditInvoiceCallbackHandler()).Methods("POST")
	s.Router.HandleFunc("/payment/xendit/ovo/callback", s.paymentSrv.XenditOVOCallbackHandler()).Methods("POST")
	s.Router.HandleFunc("/payment/xendit/dana/callback", s.paymentSrv.XenditDanaCallbackHandler()).Methods("POST")
	s.Router.HandleFunc("/payment/xendit/linkaja/callback", s.paymentSrv.XenditLinkAjaCallbackHandler()).Methods("POST")
	s.Router.HandleFunc("/payment/xendit/ewallet/callback", s.paymentSrv.XenditEWalletCallbackHandler()).Methods("POST")
	s.Router.HandleFunc("/payment/subscriptions", s.paymentSrv.CreateSubscriptionHandler()).Methods("POST")
	s.Router.HandleFunc("/payment/subscriptions/{subscription_number}/pause", s.paymentSrv.PauseSubscriptionHandler()).Methods("POST", "PUT")
	s.Router.HandleFunc("/payment/subscriptions/{subscription_number}/stop", s.paymentSrv.StopSubscriptionHandler()).Methods("POST", "PUT")
	s.Router.HandleFunc("/payment/subscriptions/{subscription_number}/resume", s.paymentSrv.ResumeSubscriptionHandler()).Methods("POST", "PUT")
}

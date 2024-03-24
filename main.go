package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type AccountNotifier interface {
	NotifyAccountCreated(context.Context, Account) error
}

type Account struct {
	Username string
	Email    string
}

type AccountHandler struct {
	AccountNotifier AccountNotifier
}

func NewAccountHandler(notifier AccountNotifier) *AccountHandler {
	return &AccountHandler{
		AccountNotifier: notifier,
	}
}

type SimpleAccountNotifier struct{}

func (n SimpleAccountNotifier) NotifyAccountCreated(ctx context.Context, account Account) error {
	slog.Info("Account was created successfully", "Email", account.Email, "Username", account.Username)
	return nil
}

type NetflixAccountNotifier struct{}

func (n NetflixAccountNotifier) NotifyAccountCreated(ctx context.Context, account Account) error {
	slog.Info("Account was created successfully by netflix notifier", "Email", account.Email, "Username", account.Username)
	return nil
}

func (h AccountHandler) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var account Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		slog.Error("Failed to decode body of request", "err", err)
		return
	}

	// adding logic
	if err := h.AccountNotifier.NotifyAccountCreated(r.Context(), account); err != nil {
		slog.Error("Could not send notification to", "email", account.Email, "err", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

func main() {

	mux := http.NewServeMux()

	accountHandler := NewAccountHandler(NetflixAccountNotifier{})

	mux.HandleFunc("POST /add-account", accountHandler.handleCreateAccount)

	http.ListenAndServe(":3000", mux)

}

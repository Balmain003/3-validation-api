package handler

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/smtp"
	"strings"
	"study/config"
	"sync"

	"github.com/jordan-wright/email"
)

var (
	verificationStore = struct {
		sync.RWMutex
		m map[string]string
	}{m: make(map[string]string)}
)

func SendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	emailAddr := r.FormValue("email")
	if emailAddr == "" {
		http.Error(w, "Internal service error", http.StatusBadRequest)
		return
	}
	token, err := generateToken()
	if err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}
	verificationStore.Lock()
	verificationStore.m[token] = emailAddr
	verificationStore.Unlock()

	if err = sendVerificationEmail(emailAddr, token); err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Verification email sent")
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	hash := strings.TrimPrefix(r.URL.Path, "/verify/")
	if hash == "" {
		http.Error(w, "Invalid verification link", http.StatusBadRequest)
		return
	}
	verificationStore.Lock()
	emailAddr, exists := verificationStore.m[hash]
	if exists {
		delete(verificationStore.m, hash)
		verificationStore.Unlock()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Email %s verified successfully", emailAddr)
	} else {
		verificationStore.Unlock()
		http.Error(w, "Invalid verification link", http.StatusBadRequest)
	}
}

func generateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func sendVerificationEmail(to, token string) error {
	e := email.NewEmail()
	e.From = config.Cfg.Email
	e.To = []string{to}
	e.Subject = "Email Verification"
	e.Text = []byte(fmt.Sprintf("Click to verify: http://localhost:8080/verify/%s", token))
	host, _, err := net.SplitHostPort(config.Cfg.Address)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", config.Cfg.Email, config.Cfg.Password, host)

	// Для порта 465 используем SSL/TLS, для 587 - STARTTLS
	if strings.HasSuffix(config.Cfg.Address, ":465") {
		return e.SendWithTLS(config.Cfg.Address, auth, &tls.Config{
			ServerName: host,
		})
	}

	return e.Send(config.Cfg.Address, auth)
}

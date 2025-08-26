package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type DataOfEmail struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

func main() {
	MailConfig := &DataOfEmail{
		Email:    "r_mukhamedzhanov@bk.ru",
		Password: "2iSFCxX0B640WJHMzfKo",
		Address:  "smtp.mail.ru:587",
	}
	http.HandleFunc("/send", MailConfig.Send)
	fmt.Println("Сервер запущен")
	http.ListenAndServe(":8080", nil)
}

func generateHash(length int) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	result := make([]rune, length)

	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}

func (data *DataOfEmail) Send(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метода нет", http.StatusMethodNotAllowed)
		return
	}
	var requestData struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	fmt.Printf("Конфигурация SMTP: Email=%s, Password=%s, Address=%s\n",
		data.Email, data.Password, data.Address)
	Hash := generateHash(5)
	MsgWithHash := fmt.Sprintf("Пройдите по ссылке: http://localhost:8080/verify/{%v}", Hash)
	fmt.Println(MsgWithHash)
	result := SendEmail(data.Email, data.Password, data.Address, requestData.Email, MsgWithHash)
	if !result {
		fmt.Printf("FUNC SEND EMAIL RETURN %v", result)
	} else {
		response := map[string]interface{}{
			"message": "Email sent sucessfully",
			"to":      requestData.Email,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func SendEmail(Email, Password, Address, Reciver, MsgWithHash string) bool {
	e := email.NewEmail()
	e.From = Email
	e.To = []string{Reciver}
	e.Subject = "Подтвержение"
	e.Text = []byte(MsgWithHash)
	err := e.Send(Address, smtp.PlainAuth("", Email, Password, "smtp.mail.ru"))
	return err == nil
}

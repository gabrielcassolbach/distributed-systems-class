package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv" 
	"github.com/resend/resend-go/v2"
)

func SendEmail(toEmail string, produto string, detalhes string) error {
	err := godotenv.Load()

	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY não foi configurada no arquivo .env ou no sistema")
	}

	client := resend.NewClient(apiKey)
	htmlContent := fmt.Sprintf(`
		<div style="font-family: sans-serif; padding: 20px; border: 1px solid #eee; border-radius: 5px;">
			<h2 style="color: #ff4500;">🔥 Nova Promoção Detectada!</h2>
			<p>Olá! Uma nova oportunidade surgiu na categoria que você monitora:</p>
			<blockquote style="background: #f9f9f9; padding: 15px; border-left: 5px solid #ff4500;">
				<strong>Produto:</strong> %s <br>
				<strong>Detalhes:</strong> %s
			</blockquote>
			<p style="font-size: 12px; color: #999;">Você recebeu este e-mail porque está inscrito no nosso sistema de notificações.</p>
		</div>
	`, produto, detalhes)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev", 
		To:      []string{toEmail},
		Subject: "🔔 Alerta de Promoção: " + produto,
		Html:    htmlContent,
	}

	_, err = client.Emails.SendWithContext(context.Background(), params)
	if err != nil {
		return fmt.Errorf("falha ao enviar e-mail via Resend: %v", err)
	}

	log.Printf("[EMAIL] Notificação enviada com sucesso para %s", toEmail)
	return nil
}
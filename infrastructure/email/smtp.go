package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"myapp/infrastructure/config"
)

const mailtrapAPIURL = "https://send.api.mailtrap.io/api/send"

// Notifier sends transactional emails via Mailtrap API.
type Notifier struct {
	token string
	from  string
}

// New creates a Notifier from config.
func New(cfg config.Config) *Notifier {
	return &Notifier{
		token: cfg.MailtrapToken,
		from:  cfg.EmailFrom,
	}
}

var registrationTmpl = template.Must(template.New("registration").Parse(`<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Pendaftaran Berhasil</title>
</head>
<body style="margin:0;padding:0;background-color:#f5f0e8;font-family:'Helvetica Neue',Helvetica,Arial,sans-serif;">
  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f5f0e8;padding:40px 16px;">
    <tr>
      <td align="center">
        <table width="100%" cellpadding="0" cellspacing="0" style="max-width:560px;background-color:#fffef9;border-radius:24px;overflow:hidden;border:1px solid rgba(0,0,0,0.08);">
          <tr>
            <td style="background-color:#1a1a1a;padding:32px 40px;">
              <p style="margin:0;font-size:11px;font-family:monospace;color:#c8f55a;letter-spacing:0.1em;text-transform:uppercase;">Amigos Care Club</p>
              <h1 style="margin:8px 0 0;font-size:24px;color:#fffef9;font-weight:800;letter-spacing:-0.02em;">Pendaftaran Diterima!</h1>
            </td>
          </tr>
          <tr>
            <td style="padding:36px 40px;">
              <p style="margin:0 0 16px;font-size:15px;color:#1a1a1a;line-height:1.6;">Halo, <strong>{{.Name}}</strong>!</p>
              <p style="margin:0 0 24px;font-size:15px;color:#444;line-height:1.6;">
                Pendaftaran kamu untuk event berikut telah kami terima dan sedang <strong>dalam proses verifikasi</strong> oleh tim kami:
              </p>
              <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f5f0e8;border-radius:16px;margin-bottom:24px;">
                <tr>
                  <td style="padding:20px 24px;">
                    <p style="margin:0 0 4px;font-size:11px;font-family:monospace;color:#888;text-transform:uppercase;letter-spacing:0.08em;">Event</p>
                    <p style="margin:0 0 16px;font-size:16px;font-weight:700;color:#1a1a1a;">40% of Heart Rate Run — Vol.2</p>
                    <p style="margin:0 0 4px;font-size:11px;font-family:monospace;color:#888;text-transform:uppercase;letter-spacing:0.08em;">ID Registrasi</p>
                    <p style="margin:0;font-size:14px;font-weight:600;color:#1a1a1a;font-family:monospace;">{{.RegistrationID}}</p>
                  </td>
                </tr>
              </table>
              <p style="margin:0 0 16px;font-size:14px;color:#555;line-height:1.6;">
                Tim kami akan memverifikasi bukti pembayaran kamu. Setelah terverifikasi, kamu akan mendapatkan email konfirmasi beserta nomor tiket.
              </p>
              <p style="margin:0;font-size:14px;color:#555;line-height:1.6;">Sampai jumpa di garis start! 🏃</p>
            </td>
          </tr>
          <tr>
            <td style="padding:20px 40px;border-top:1px solid rgba(0,0,0,0.06);">
              <p style="margin:0;font-size:12px;font-family:monospace;color:#999;">
                Tim Amigos Care Club &mdash; Run With Fun
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`))

// SendRegistrationConfirmation sends a "pendaftaran diterima, sedang diproses" email.
func (n *Notifier) SendRegistrationConfirmation(ctx context.Context, email, name, registrationID string) error {
	var buf bytes.Buffer
	if err := registrationTmpl.Execute(&buf, struct {
		Name           string
		RegistrationID string
	}{Name: name, RegistrationID: registrationID}); err != nil {
		return fmt.Errorf("email template: %w", err)
	}
	subject := "Pendaftaran Kamu Sedang Diproses — 40% of Heart Rate Run"
	return n.send(ctx, email, subject, buf.String())
}

var verificationTmpl = template.Must(template.New("verification").Parse(`<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Registrasi Terverifikasi</title>
</head>
<body style="margin:0;padding:0;background-color:#f5f0e8;font-family:'Helvetica Neue',Helvetica,Arial,sans-serif;">
  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f5f0e8;padding:40px 16px;">
    <tr>
      <td align="center">
        <table width="100%" cellpadding="0" cellspacing="0" style="max-width:560px;background-color:#fffef9;border-radius:24px;overflow:hidden;border:1px solid rgba(0,0,0,0.08);">
          <tr>
            <td style="background-color:#1a1a1a;padding:32px 40px;">
              <p style="margin:0;font-size:11px;font-family:monospace;color:#c8f55a;letter-spacing:0.1em;text-transform:uppercase;">Amigos Care Club</p>
              <h1 style="margin:8px 0 0;font-size:24px;color:#fffef9;font-weight:800;letter-spacing:-0.02em;">Registrasi Terverifikasi! 🎉</h1>
            </td>
          </tr>
          <tr>
            <td style="padding:36px 40px;">
              <p style="margin:0 0 16px;font-size:15px;color:#1a1a1a;line-height:1.6;">Halo, <strong>{{.Name}}</strong>!</p>
              <p style="margin:0 0 24px;font-size:15px;color:#444;line-height:1.6;">
                Kabar gembira! Pembayaran kamu sudah kami verifikasi dan registrasi kamu untuk event berikut telah <strong>resmi diterima</strong>:
              </p>
              <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f5f0e8;border-radius:16px;margin-bottom:24px;">
                <tr>
                  <td style="padding:20px 24px;">
                    <p style="margin:0 0 4px;font-size:11px;font-family:monospace;color:#888;text-transform:uppercase;letter-spacing:0.08em;">Event</p>
                    <p style="margin:0 0 16px;font-size:16px;font-weight:700;color:#1a1a1a;">40% of Heart Rate Run — Vol.2</p>
                    <p style="margin:0 0 4px;font-size:11px;font-family:monospace;color:#888;text-transform:uppercase;letter-spacing:0.08em;">Nomor Tiket</p>
                    <p style="margin:0;font-size:18px;font-weight:800;color:#1a1a1a;font-family:monospace;letter-spacing:0.05em;">{{.TicketNumber}}</p>
                  </td>
                </tr>
              </table>
              <p style="margin:0 0 16px;font-size:14px;color:#555;line-height:1.6;">
                Simpan nomor tiket kamu dan pantau terus informasi lebih lanjut mengenai detail acara yang akan kami kirimkan mendekati hari H.
              </p>
              <p style="margin:0;font-size:14px;color:#555;line-height:1.6;">Sampai jumpa di garis start! 🏃</p>
            </td>
          </tr>
          <tr>
            <td style="padding:20px 40px;border-top:1px solid rgba(0,0,0,0.06);">
              <p style="margin:0;font-size:12px;font-family:monospace;color:#999;">
                Tim Amigos Care Club &mdash; Run With Fun
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`))

// SendVerificationConfirmation sends a confirmation email after admin verifies a registration.
func (n *Notifier) SendVerificationConfirmation(ctx context.Context, email, name, ticketNumber string) error {
	var buf bytes.Buffer
	if err := verificationTmpl.Execute(&buf, struct {
		Name         string
		TicketNumber string
	}{Name: name, TicketNumber: ticketNumber}); err != nil {
		return fmt.Errorf("email template: %w", err)
	}
	subject := "Registrasi Kamu Sudah Terverifikasi! — 40% of Heart Rate Run"
	return n.send(ctx, email, subject, buf.String())
}

type mailtrapAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type mailtrapPayload struct {
	From    mailtrapAddress   `json:"from"`
	To      []mailtrapAddress `json:"to"`
	Subject string            `json:"subject"`
	HTML    string            `json:"html"`
}

func (n *Notifier) send(ctx context.Context, to, subject, htmlBody string) error {
	fromEmail, fromName := parseAddress(n.from)

	payload := mailtrapPayload{
		From:    mailtrapAddress{Email: fromEmail, Name: fromName},
		To:      []mailtrapAddress{{Email: to}},
		Subject: subject,
		HTML:    htmlBody,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal email payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, mailtrapAPIURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+n.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("mailtrap api status %d", resp.StatusCode)
	}
	return nil
}

// parseAddress splits "Name <email@example.com>" into name and email parts.
func parseAddress(from string) (email, name string) {
	start := -1
	for i, c := range from {
		if c == '<' {
			start = i
		}
		if c == '>' && start != -1 {
			return from[start+1 : i], strings.TrimSpace(from[:start])
		}
	}
	return from, ""
}

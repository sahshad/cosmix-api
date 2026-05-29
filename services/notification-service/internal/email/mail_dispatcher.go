package email

const (
	TemplateAuthVerification    = "auth/verification.html"
	TemplateAuthWelcome         = "auth/welcome.html"
	TemplateAuthForgotPassword  = "auth/forgot_password.html"
	TemplateAuthPasswordChanged = "auth/password_changed.html"
)

type SendEmailDTO struct {
	To           string
	Subject      string
	TemplateName string
	TemplateData any
}

type MailDispatcher struct {
	transporter *Transport
	renderer    *Renderer
}

func NewMailDispatcher(host string, port int, from string, templateDir string) *MailDispatcher {
	return &MailDispatcher{
		transporter: NewTransport(
			host,
			port,
			from,
		),

		renderer: NewRenderer(
			templateDir,
		),
	}
}

func (m *MailDispatcher) Send(dto SendEmailDTO) error {
	html, err := m.renderer.Render(
		dto.TemplateName,
		dto.TemplateData,
	)
	if err != nil {
		return err
	}

	return m.transporter.Transport(
		Email{
			To:       dto.To,
			Subject:  dto.Subject,
			HTMLBody: html,
		},
	)
}

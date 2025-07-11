package mailer

import "embed"

const (
	FromName               = "GopherSocial"
	maxRetries             = 3
	UserInvitationTemplate = "user_invitation.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile string, username string, email string, data any, isSandbox bool) error
}

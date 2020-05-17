package Models

type MSG struct {
	To      string
	From    string
	Cc      string
	Subject string
	Body    string
}

func NewMessage(to, from, cc, subject, body string) *MSG {
	m := &MSG{
		To:      to,
		From:    from,
		Cc:      cc,
		Subject: subject,
		Body:    body,
	}
	return m
}

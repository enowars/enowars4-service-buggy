package Mailing

type Msg struct {
	To      string
	From    string
	Cc      string
	Subject string
	Body    string
}

func NewMessage(to, from, cc, subject, body string) *Msg {
	m := &Msg{
		To:      to,
		From:    from,
		Cc:      cc,
		Subject: subject,
		Body:    body,
	}
	return m
}

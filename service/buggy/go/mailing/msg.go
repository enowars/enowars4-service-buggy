package Mailing

type Msg struct {
	To      string
	From    string
	Cc      string
	Subject string
	Body    string
}

func NewMessage(to, from, subject, body, cc string) *Msg {

	m := &Msg{
		To:      to,
		From:    from,
		Subject: subject,
		Body:    body,
		Cc:      cc,
	}
	return m
}

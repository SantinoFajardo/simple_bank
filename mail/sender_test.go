package mail

import (
	"fmt"
	"testing"

	"github.com/santinofajardo/simpleBank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config, err := util.LoadConfig("..")
	require.NoError(t, err)
	fmt.Printf("NAME,ADDRESS,PASSWORD = %s, %s, %s", config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	subject := "A test email"
	content := `
	<h1>Herllo world</h1>
	<p>This is a test messages used with Golang</p>
	`
	to := []string{"fajardosantinoguillermo@gmail.com", "santinogfajardo@hotmail.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}

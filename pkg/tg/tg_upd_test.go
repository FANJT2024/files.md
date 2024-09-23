package tg

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/require"
)

func TestCmdNil(t *testing.T) {
	r := require.New(t)

	m := tgbotapi.Message{}
	m.Text = "j new journal record"
	rawUpdate := tgbotapi.Update{
		UpdateID: 0,
		Message:  &m,
	}

	u := NewTGUpd(rawUpdate)
	cmd := u.Cmd()

	r.Nil(cmd)
}

func TestCmdInTheBeginning(t *testing.T) {
	r := require.New(t)

	m := tgbotapi.Message{}
	m.Text = "/j new journal record"
	m.Entities = []tgbotapi.MessageEntity{{
		Type:   "bot_command",
		Offset: 0,
		Length: 2,
	}}
	rawUpdate := tgbotapi.Update{
		UpdateID: 0,
		Message:  &m,
	}

	u := NewTGUpd(rawUpdate)
	cmd := u.Cmd()

	r.NotNil(cmd)
	r.Equal("j", cmd.Name)
	r.Equal("New journal record", cmd.Params[0])
}

func TestCmdAtTheEnd(t *testing.T) {
	r := require.New(t)

	m := tgbotapi.Message{}
	m.Text = "new journal record /j"
	m.Entities = []tgbotapi.MessageEntity{{
		Type:   "bot_command",
		Offset: 19,
		Length: 2,
	}}
	rawUpdate := tgbotapi.Update{
		UpdateID: 0,
		Message:  &m,
	}

	u := NewTGUpd(rawUpdate)
	cmd := u.Cmd()

	r.NotNil(cmd)
	r.Equal("j", cmd.Name)
	r.Equal("New journal record", cmd.Params[0])
}

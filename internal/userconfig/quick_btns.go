package userconfig

import (
	"zakirullin/stuffbot/i18n"
	"zakirullin/stuffbot/internal/consts"
	"zakirullin/stuffbot/pkg/tg"
)

type QuickBtn struct {
	Cmd         string
	CmdType     string
	Emoji       string
	Description string
}

var AvailableQuickBtns = []tg.Btn{
	//tg.NewBtn("Later", tg.NewCmd(consts.)),
	NewQuickBtn(consts.CmdLater, tg.CmdTypeCallback, i18n.Emoji("Later"), "Later"),
	NewQuickBtn(consts.CmdInlineQuerySearchEveryWhere, tg.CmdTypeInlineQueryCurrentChat, i18n.Emoji("Search"), "Search"),
	NewQuickBtn(consts.CmdShowFiles, tg.CmdTypeCallback, i18n.Emoji("Files"), "Files"),
	NewQuickBtn(consts.CmdShowChecklists, tg.CmdTypeCallback, i18n.Emoji("Checklists"), "Checklists"),
	NewQuickBtn(consts.CmdShowPostpone, tg.CmdTypeCallback, i18n.Emoji("Postpone"), "Postpone"),
	NewQuickBtn(consts.CmdShowReadChecklist, tg.CmdTypeCallback, i18n.Emoji("Read"), "Read"),
	NewQuickBtn(consts.CmdShowWatchChecklist, tg.CmdTypeCallback, i18n.Emoji("Watch"), "Watch"),
	NewQuickBtn(consts.CmdShowShopChecklist, tg.CmdTypeCallback, i18n.Emoji("Shop"), "Shop"),
	NewQuickBtn(consts.CmdWebAppHabits, tg.CmdTypeWebApp, i18n.Emoji("Habits"), "Habits"),
}

var (
	QuickPanelAddButton = "➕"
	QuickPanelDelButton = "➖"
)

func NewQuickBtn(cmd, cmdType, emoji, description string) QuickBtn {
	return QuickBtn{cmd, cmdType, emoji, description}
}

func (c *Config) AddQuickBtn(button string) bool {
	// Does this button already exist?
	for _, curBtn := range c.raw.QuickCmds {
		if curBtn == button {
			return false
		}
	}
	c.raw.QuickCmds = append(c.raw.QuickCmds, button)
	return true
}

func (c *Config) QuickCmds() []string {
	return c.raw.QuickCmds
}

func (c *Config) HasQuickCmd(cmd string) bool {
	for _, pref := range c.raw.QuickCmds {
		if cmd == pref {
			return true
		}
	}
	return false
}

func (c *Config) DelQuickBtn(toDelete string) bool {
	var newButtons []string
	found := false // Was the target
	for _, curBtn := range c.raw.QuickCmds {
		if curBtn == toDelete {
			found = true
		} else {
			newButtons = append(newButtons, curBtn)
		}
	}
	c.raw.QuickCmds = newButtons
	return found
}

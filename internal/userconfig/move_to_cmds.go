package userconfig

import (
	"zakirullin/stuffbot/i18n"
	"zakirullin/stuffbot/internal/consts"
	"zakirullin/stuffbot/pkg/tg"
)

var AvailableMoveToBtns = []tg.Btn{
	tg.NewBtn(i18n.StrToTomorrow, tg.NewCmd(consts.CmdScheduleForTmrw, nil)),
	tg.NewBtn(i18n.StrMoveToLater, tg.NewCmd(consts.CmdMoveToLater, nil)),
	tg.NewBtn(i18n.StrToADay, tg.NewCmd(consts.CmdShowScheduleForDay, nil)),
	tg.NewBtn(i18n.StrToFile, tg.NewCmd(consts.CmdShowMoveToFile, nil)),
	tg.NewBtn(i18n.StrToJournal, tg.NewCmd(consts.CmdMoveToJournal, nil)),
	tg.NewBtn(i18n.StrToRead, tg.NewCmd(consts.CmdMoveToRead, nil)),
	tg.NewBtn(i18n.StrToWatch, tg.NewCmd(consts.CmdMoveToWatch, nil)),
	tg.NewBtn(i18n.StrToShop, tg.NewCmd(consts.CmdMoveToShop, nil)),
	tg.NewBtn(i18n.StrToChecklist, tg.NewCmd(consts.CmdShowMoveToChecklist, nil)),
}

func (c *Config) AddMoveToCmd(cmd string) bool {
	// Does this cmd already exist?
	for _, existingCmds := range c.raw.MoveToCmds {
		if existingCmds == cmd {
			return false
		}
	}
	c.raw.MoveToCmds = append(c.raw.MoveToCmds, cmd)

	return true
}

func (c *Config) MoveToCmds() []string {
	return c.raw.MoveToCmds
}

func (c *Config) HasMoveToCmd(cmd string) bool {
	for _, pref := range c.raw.MoveToCmds {
		if cmd == pref {
			return true
		}
	}
	return false
}

func (c *Config) DelMoveToCmd(cmd string) bool {
	var newCmds []string
	found := false // Was the target
	for _, existingCmd := range c.raw.MoveToCmds {
		if existingCmd == cmd {
			found = true
		} else {
			newCmds = append(newCmds, existingCmd)
		}
	}
	c.raw.MoveToCmds = newCmds

	return found
}

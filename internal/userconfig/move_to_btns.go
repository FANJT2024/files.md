package userconfig

//import (
//	"zakirullin/stuffbot/i18n"
//	"zakirullin/stuffbot/internal/constants"
//	"zakirullin/stuffbot/pkg/tg"
//)
//
//type QuickBtn struct {
//	Cmd         string
//	CmdType     string
//	Emoji       string
//	Description string
//}
//
//var AvailableQuickBtns = []QuickBtn{
//	NewQuickBtn(constants.CmdLater, tg.CmdTypeCallback, i18n.Emoji("Later"), "Later"),
//	NewQuickBtn(constants.CmdInlineQuerySearchEveryWhere, tg.CmdTypeInlineQueryCurrentChat, i18n.Emoji("Search"), "Search"),
//	NewQuickBtn(constants.CmdShowFiles, tg.CmdTypeCallback, i18n.Emoji("Files"), "Files"),
//	NewQuickBtn(constants.CmdShowChecklists, tg.CmdTypeCallback, i18n.Emoji("Checklists"), "Checklists"),
//	NewQuickBtn(constants.CmdShowPostpone, tg.CmdTypeCallback, i18n.Emoji("Postpone"), "Postpone"),
//	NewQuickBtn(constants.CmdShowReadChecklist, tg.CmdTypeCallback, i18n.Emoji("Read"), "Read"),
//	NewQuickBtn(constants.CmdShowWatchChecklist, tg.CmdTypeCallback, i18n.Emoji("Watch"), "Watch"),
//	NewQuickBtn(constants.CmdShowShopChecklist, tg.CmdTypeCallback, i18n.Emoji("Shop"), "Shop"),
//	NewQuickBtn(constants.CmdWebAppHabits, tg.CmdTypeWebApp, i18n.Emoji("Habits"), "Habits"),
//}
//
//var (
//	QuickPanelAddButton = "➕"
//	QuickPanelDelButton = "➖"
//)
//
//func NewQuickBtn(cmd, cmdType, emoji, description string) QuickBtn {
//	return QuickBtn{cmd, cmdType, emoji, description}
//}
//
//func (c *Config) AddQuickBtn(button string) bool {
//	// Does this button already exist?
//	for _, curBtn := range c.raw.QuickCmds {
//		if curBtn == button {
//			return false
//		}
//	}
//	c.raw.QuickCmds = append(c.raw.QuickCmds, button)
//	return true
//}
//
//func (c *Config) QuickCmds() []string {
//	return c.raw.QuickCmds
//}
//
//func (c *Config) HasQuickCmd(cmd string) bool {
//	for _, pref := range c.raw.QuickCmds {
//		if cmd == pref {
//			return true
//		}
//	}
//	return false
//}
//
//func (c *Config) DelQuickBtn(toDelete string) bool {
//	var newButtons []string
//	found := false // Was the target
//	for _, curBtn := range c.raw.QuickCmds {
//		if curBtn == toDelete {
//			found = true
//		} else {
//			newButtons = append(newButtons, curBtn)
//		}
//	}
//	c.raw.QuickCmds = newButtons
//	return found
//}

package journal

import (
	"fmt"
	"strings"
	"time"

	"zakirullin/stuffbot/internal/fs"
	"zakirullin/stuffbot/pkg/txt"
)

var now = time.Now

func AddRecord(userFS *fs.FS, noteFilename string) error {
	record, err := userFS.Read(fs.DirToday, noteFilename)
	if err != nil {
		return fmt.Errorf("failed to move to journal: can't get note content: %w", err)
	}

	record = strings.TrimSpace(record)
	if len(record) == 0 {
		record = fs.Title(noteFilename)
	}

	journalFilename := now().Format("2006 January.md")
	exists, err := userFS.Exists(fs.DirJournal, journalFilename)
	if err != nil {
		return err
	}

	var md string
	if exists {
		md, err = userFS.Read(fs.DirJournal, journalFilename)
		if err != nil {
			return err
		}
		md = txt.NormNewLines(md)
		md = strings.TrimSpace(md)
		if len(md) != 0 {
			md += "\n"
		}
	}

	header := fmt.Sprintf("#### %d, %s", now().Day(), now().Weekday())
	if !strings.Contains(md, header) {
		md += header + "\n"
	}

	record = fmt.Sprintf("%s %s\n", now().Format("`15:04`"), record)
	md += record

	return userFS.Write(fs.DirJournal, journalFilename, md)
}

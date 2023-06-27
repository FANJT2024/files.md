package worker

import (
	"fmt"
	"time"

	"github.com/spf13/afero"
	"golang.org/x/exp/slog"

	"zakirullin/dumpbot/internal/fs"
	"zakirullin/dumpbot/internal/sched"
	"zakirullin/dumpbot/internal/userconfig"
)

func MoveDueTasksToToday(config userconfig.Config, fsBackend afero.Fs) error {
	ids, err := fs.AllUserIDs()
	if err != nil {
		return fmt.Errorf("moveDueTasksForToday: %s\n", err)
	}

	for _, id := range ids {
		sch := config.Schedules()

		fsys, err := fs.NewFS(id, fsBackend)
		if err != nil {
			return fmt.Errorf("moveDueTasksForToday: can't create FS: %s", err)
		}
		for _, schedule := range sch {
			if time.Now().Unix() >= schedule.ScheduleAt {
				err = moveTaskToToday(schedule.Filename, fsys)
				if err != nil {
					slog.Error("moveDueTasksForToday: can't move", "err", err)
				}
				slog.Debug("Scheduled task moved to today", schedule.Filename, "filename")
				if len(schedule.Cron) != 0 {
					runAt := sched.Next(schedule.Cron)
					config.AddToSchedule(schedule.Filename, runAt, schedule.Cron)
					slog.Debug("Task was rescheduled", "filename", schedule.Filename, "schedule", schedule.Cron, "runAt", runAt)
					continue
				}

				config.DelFromSchedule(schedule.Filename)
			}
		}
	}
	return nil
}

func moveTaskToToday(filename string, fsys *fs.FS) error {
	dirsToLookFor := []string{fs.DirLater, fs.DirArchive}
	for _, dir := range dirsToLookFor {
		filenames, err := fsys.FilesAndDirs(dir)
		if err != nil {
			return fmt.Errorf("moveTaskForToday: %w", err)
		}

		for _, f := range filenames {
			if f.Name == filename {
				err = fsys.Rename(dir, filename, fs.DirToday, filename)
				if err != nil {
					return fmt.Errorf("moveTaskForToday: can't rename: %w", err)
				}
			}
		}
	}

	return nil
}

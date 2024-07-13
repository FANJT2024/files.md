package habits

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/rivo/uniseg"

	"zakirullin/stuffbot/internal/fs"
	"zakirullin/stuffbot/pkg/txt"
)

// [1 => false, <year day> => 0, 1...]
type Year map[int]int

const (
	habitSkipped            = "⚪️"
	habitCompleted          = "🟢"
	habitCompletedAtWeekend = "🟡"

	mood = "Mood"
)

var (
	moodEmojis            = []string{"⚪️", "🤕", "😔", "😐", "🙂", "😊"}
	errMalformedMonthLine = errors.New("malformed month line")
	now                   = time.Now
)

func Habits(userFS *fs.FS, year int) (map[string]Year, error) {
	filename := "%d Habits.md"
	habitsStr, err := userFS.Read(fs.DirInsights, fmt.Sprintf(filename, year))
	if err != nil {
		return nil, fmt.Errorf("read %s error: %w", filename, err)
	}

	habits := make(map[string]Year)
	month := time.January
	lines := strings.Split(txt.NormNewLines(habitsStr), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		// Parsing month line
		isMonthLine := strings.HasPrefix(line, "###")
		if isMonthLine {
			parts := strings.Split(line, " ")
			if len(parts) < 2 {
				return nil, fmt.Errorf("read habits: can't parse month line '%s': %w", line, errMalformedMonthLine)
			}

			date, err := time.Parse("January", parts[1])
			if err != nil {
				return nil, fmt.Errorf("read habits: can't parse month %s: %w", line, err)
			}
			month = date.Month()

			continue
		}

		// Tolerant reader, if we encounter gibberish,
		// we skip it. See ADRs in README.md for details for details
		// TODO preserve gibberish between parsing seesions
		daysAndHabit := strings.SplitN(line, " ", 2)
		if len(daysAndHabit) < 2 {
			continue
		}
		days, habit := daysAndHabit[0], daysAndHabit[1]

		firstDayOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		dayOfTheYear := firstDayOfMonth.YearDay()

		// Moods line
		moodsMarker := mood
		if strings.Contains(habit, moodsMarker) {
			gr := uniseg.NewGraphemes(days)
			dayOffset := 0
			for gr.Next() {
				power := slices.Index(moodEmojis, gr.Str())
				habits[mood][dayOfTheYear+dayOffset] = power
				dayOfTheYear++
			}
			continue
		}

		habitsMarker := fmt.Sprintf("%s%s%s", habitSkipped, habitCompletedAtWeekend, habitCompleted)
		if !strings.ContainsAny(days, habitsMarker) {
			continue
		}

		// At this point we are on habits line, which is
		// [⚪️🟢... Habit name] i.e. completion status
		// for every day of the above found month
		habitName := strings.TrimSpace(habit)
		if _, ok := habits[habitName]; !ok {
			habits[habitName] = make(Year)
		}


		// See README.md ADRs
		gr := uniseg.NewGraphemes(days)
		dayOffset := 0
		for gr.Next() {
			habits[habitName][dayOfTheYear+dayOffset] = 0
			if gr.Str() != habitSkipped {
				habits[habitName][dayOfTheYear+dayOffset] = 1
			}
			dayOfTheYear++
		}
	}

	return habits, nil
}

func LastWeekHabits(userFS *fs.FS) (map[string]Year, error) {
	habitsForYear, err := Habits(userFS, now().Year())
	if err != nil {
		return nil, fmt.Errorf("can't get habits for last week: %w", err)
	}

	var currentDay = now()
	for currentDay.Weekday() != time.Monday {
		currentDay = currentDay.Add(-24 * time.Hour)
	}

	habits := make(map[string]Year)
	for habit, statuses := range habitsForYear {
		habits[habit] = make(Year)
		for offset := range 7 {
			yearDay := currentDay.Add(time.Duration(offset) * 24 * time.Hour).YearDay()
			habits[habit][yearDay] = 0
			if status, ok := statuses[yearDay]; ok {
				habits[habit][yearDay] = status
			}
		}
	}

	return habits, nil
}

// func Write(userFS *fs.FS, habits []Habit) error {
// 	return nil
// }

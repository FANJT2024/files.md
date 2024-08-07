package worker

//func TestBot_togglePomodoro(t *testing.T) {
//	r := require.New(t)
//	userFS, err := fs.NewFS("/", afero.NewMemMapFs())
//	r.NoError(err)
//	tgram := fake.NewTG()
//	redis, err := miniredis.Run()
//	r.NoError(err)
//	defer redis.Close()
//	b2 := internal.NewBot(-1, tgram, userFS,db.NewDB(), &userconfig.DefaultConfig)
//	b := b2
//
//	pomodoroIn := func(dirName string) bool {
//		hasPomodoroInDir, err := b.fs.Exists(dirName, fs.FilePomodoro)
//		r.NoError(err)
//		return hasPomodoroInDir
//	}
//	r.False(pomodoroIn(fs.DirToday) || pomodoroIn(fs.DirArchive))
//
//	// Add pomodoro	to today
//	r.Nil(b.togglePomodoro(nil))
//	r.True(pomodoroIn(fs.DirToday) && !pomodoroIn(fs.DirArchive))
//	// and remove pomodoro from today
//	r.Nil(b.togglePomodoro(nil))
//	r.False(pomodoroIn(fs.DirToday) || pomodoroIn(fs.DirArchive))
//
//	// Add pomodoro	to today
//	r.Nil(b.togglePomodoro(nil))
//	r.True(pomodoroIn(fs.DirToday) && !pomodoroIn(fs.DirArchive))
//	// complete it
//	r.Nil(b.complete([]string{fs.DirToday, fs.FilePomodoro}))
//	r.True(!pomodoroIn(fs.DirToday) && pomodoroIn(fs.DirArchive))
//	// and remove pomodoro from trash
//	r.Nil(b.togglePomodoro(nil))
//	r.False(pomodoroIn(fs.DirToday) || pomodoroIn(fs.DirArchive))
//}
//
// func TestWorkerReturnsPomodoroBackToToday(t *testing.T) {
// 	r := require.New(t)

// 	fsBackend := afero.NewMemMapFs()
// 	userFS, err := fs.NewFS("/-1", fsBackend)
// 	r.NoError(err)
// 	err = userFS.CreateUserDirs()
// 	r.NoError(err)

// 	tgram := fake.NewTG()
// 	redis, err := miniredis.Run()
// 	r.NoError(err)
// 	defer redis.Close()

// 	b := NewBot(-1, tgram, userFS,db.NewDB(), &userconfig.DefaultConfig)

// 	currentBackend := fs.DefaultBackend
// 	fs.DefaultBackend = fsBackend
// 	defer func() {
// 		fs.DefaultBackend = currentBackend
// 	}()

// 	pomodoroIn := func(dirName string) bool {
// 		hasPomodoroInDir, err := b.fs.Exists(dirName, fs.FilePomodoro)
// 		r.NoError(err)
// 		return hasPomodoroInDir
// 	}
// 	r.False(pomodoroIn(fs.DirToday) || pomodoroIn(fs.DirArchive))

// 	// Add pomodoro	to today
// 	r.Nil(b.togglePomodoro(nil))
// 	r.True(pomodoroIn(fs.DirToday) && !pomodoroIn(fs.DirArchive))
// 	// set pomodoro duration to 1us
// 	r.NoError(b.conf.SetPomodoroDuration(time.Nanosecond))
// 	// complete it
// 	r.NoError(b.complete([]string{fs.DirToday, fs.FilePomodoro}))
// 	r.True(!pomodoroIn(fs.DirToday) && pomodoroIn(fs.DirArchive))
// 	// "wait" until it gets back to today
// 	r.NoError(worker.MoveDueTasksToToday("", "conf", fsBackend))
// 	r.True(pomodoroIn(fs.DirToday) && !pomodoroIn(fs.DirArchive))
// }

//func TestWorkerPomodoroIsNotReturnedUntilItIsDue(t *testing.T) {
//	r := require.New(t)
//	fsBackend := afero.NewMemMapFs()
//	userFS, err := fs.NewFS("/-1", fsBackend)
//	r.NoError(err)
//	tgram := fake.NewTG()
//	redis, err := miniredis.Run()
//	r.NoError(err)
//	defer redis.Close()
//	b := NewBot(-1, tgram, userFS,db.NewDB(), &userconfig.DefaultConfig)
//
//	currentBackend := fs.DefaultBackend
//	fs.DefaultBackend = fsBackend
//	defer func() {
//		fs.DefaultBackend = currentBackend
//	}()
//
//	pomodoroIn := func(dirName string) bool {
//		hasPomodoroInDir, err := b.fs.Exists(dirName, fs.FilePomodoro)
//		r.NoError(err)
//		return hasPomodoroInDir
//	}
//	r.False(pomodoroIn(fs.DirToday) || pomodoroIn(fs.DirArchive))
//
//	r.NoError(b.togglePomodoro(nil))
//	r.True(pomodoroIn(fs.DirToday) && !pomodoroIn(fs.DirArchive))
//	r.NoError(b.conf.SetPomodoroDuration(2 * time.Second))
//	r.NoError(b.complete([]string{fs.DirToday, fs.FilePomodoro}))
//	r.True(!pomodoroIn(fs.DirToday) && pomodoroIn(fs.DirArchive))
//	// trigger due tasks processing
//	r.NoError(worker.MoveDueTasksToToday("", "conf", fsBackend))
//	// pomodoro is not returned back to today
//	r.True(!pomodoroIn(fs.DirToday) && pomodoroIn(fs.DirArchive))
//}

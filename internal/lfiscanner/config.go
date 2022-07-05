package lfiscanner

type Config struct {
	Targets         map[string]string // absolute paths to files and strings they contain
	LevelUpAttempts int               // how many times scanner will go up in directory tree
}

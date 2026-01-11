package constants

import "fmt"

type ExerciseObj struct {
	ID            ExerciseID
	Name          string
	Url           string
	Type          string
	Accent        string
	RestInSeconds int
}

func (e *ExerciseObj) GetName() string {
	if e == nil {
		return ""
	}
	return e.Name
}

func (e *ExerciseObj) GetAccent() string {
	if e == nil {
		return ""
	}
	return e.Accent
}

func (e *ExerciseObj) GetHint() string {
	if e == nil {
		return ""
	}
	if e.Url == "" {
		return ""
	}
	return WrapYandexLink(e.Url)
}

func WrapYandexLink(url string) string {
	return fmt.Sprintf("\n<a href=\"%s\"><b>⚠️Техника выполнения:</b></a>", url)
}

type ExerciseID = int

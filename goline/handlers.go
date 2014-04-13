package goline

func Finish(l *GoLine) (bool, error) {
	return true, nil
}

func UserTerminated(l *GoLine) (bool, error) {
	return true, UserTerminatedError
}

func Backspace(l *GoLine) (bool, error) {
	if l.Len > 0 && l.Pos > 0 {
		l.CurLine = append(l.CurLine[:l.Pos-1], l.CurLine[l.Pos:]...)
		l.Len--
		l.Pos--
		l.CurLine[l.Len] = 0
	}
	return false, nil
}

func MoveLeft(l *GoLine) (bool, error) {
	if l.Pos > 0 {
		l.Pos--
	}
	return false, nil

}

func MoveRight(l *GoLine) (bool, error) {
	if l.Pos != l.Len {
		l.Pos++
	}
	return false, nil

}

func DeleteLine(l *GoLine) (bool, error) {
	l.CurLine = make([]rune, MAX_LINE)
	l.Pos = 0
	l.Len = 0
	return false, nil

}

func DeleteRestofLine(l *GoLine) (bool, error) {
	copy(l.CurLine, l.CurLine[:l.Pos])
	l.Len = l.Pos
	return false, nil

}

func MoveStartofLine(l *GoLine) (bool, error) {
	l.Pos = 0
	return false, nil

}

func MoveEndofLine(l *GoLine) (bool, error) {
	l.Pos = l.Len
	return false, nil

}

func ClearScreen(l *GoLine) (bool, error) {
	l.ClearScreen()
	return false, nil

}

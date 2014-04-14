package goline

func Finish(l *GoLine) (bool, error) {
	return true, nil
}

func UserTerminated(l *GoLine) (bool, error) {
	return true, UserTerminatedError
}

func Backspace(l *GoLine) (bool, error) {
	if l.Len > 0 && l.Position > 0 {
		l.CurLine = append(l.CurLine[:l.Position-1], l.CurLine[l.Position:]...)
		l.Len--
		l.Position--
		l.CurLine[l.Len] = 0
	}
	return false, nil
}

func MoveLeft(l *GoLine) (bool, error) {
	if l.Position > 0 {
		l.Position--
	}
	return false, nil
}

func MoveRight(l *GoLine) (bool, error) {
	if l.Position != l.Len {
		l.Position++
	}
	return false, nil
}

func DeleteLine(l *GoLine) (bool, error) {
	l.CurLine = make([]rune, MAX_LINE)
	l.Position = 0
	l.Len = 0
	return false, nil
}

func DeleteRestofLine(l *GoLine) (bool, error) {
	copy(l.CurLine, l.CurLine[:l.Position])
	l.Len = l.Position
	return false, nil
}

func DeleteLastWord(l *GoLine) (bool, error) {
	for i := l.Position - 1; i > 0; i-- {
		if l.CurLine[i-1] == ' ' {
			copy(l.CurLine, l.CurLine[:i])
			l.Len = i
			l.Position = i
			return false, nil
		}
	}
	return DeleteLine(l)
}

func MoveStartofLine(l *GoLine) (bool, error) {
	l.Position = 0
	return false, nil
}

func MoveEndofLine(l *GoLine) (bool, error) {
	l.Position = l.Len
	return false, nil
}

func ClearScreen(l *GoLine) (bool, error) {
	l.ClearScreen()
	return false, nil
}

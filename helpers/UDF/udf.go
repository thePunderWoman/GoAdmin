package UDF

func IntOrNil(num int) *int {
	if num > 0 {
		return &num
	} else {
		return nil
	}
}

func StrOrNil(str string) *string {
	if str != "" {
		return &str
	} else {
		return nil
	}
}

package user

import "strings"

func LinterLogin(login string) (bool, string) {
	if strings.Contains(login, " ") {
		return false, "логин не должен содержать пробелов"
	}
	if len(login) < 4 {
		return false, "логин слишком короткий. Минимум 4 символов"
	}

	return true, ""
}

func LinterPass(pass string) (bool, string) {
	if len(pass) < 128 {
		return false, "некорректный пароль. Смените метод шифрации указанный в документации."
	}
	return true, ""
}

func LinterEmail(email string) (bool, string) {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") || len(email) < 6 {
		return false, "некорректная электронная почта"
	}
	return true, ""
}

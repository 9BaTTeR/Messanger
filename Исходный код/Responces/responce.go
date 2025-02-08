package responces

import (
	"encoding/json"
)

const (
	operdone = "Operation done"
)

// Внутреняя ошибка
func (r Responce) InternalError(message string) Responce {
	r.Code = 500
	r.Description = message
	return r
}

func (r Responce) OK(message string) Responce {
	r.Code = 200
	r.Description = message
	return r
}

// Объект создан
func (r Responce) Created() Responce {
	r.Code = 201
	r.Description = operdone
	return r
}

// Запрос выполнен
func (r Responce) Accepted() Responce {
	r.Code = 202
	r.Description = operdone
	return r
}

// Некорректный запрос
func (r Responce) BadRequest(message string) Responce {
	r.Code = 400
	r.Description = message
	return r
}

// Неавторизованный пользователь
func (r Responce) Unauthorized(message string) Responce {
	r.Code = 401
	r.Description = message
	return r
}

// Доступ запрещён
func (r Responce) Forbidden(message string) Responce {
	r.Code = 403
	r.Description = message
	return r
}

func (vr VerifedResponce) Responcer() ([]byte, error) {
	answer, err := json.Marshal(vr)
	return answer, err
}

func (rr RegResponce) Responser() ([]byte, error) {
	answer, err := json.Marshal(rr)
	return answer, err
}

func (ar AuthResponce) Responser() ([]byte, error) {
	answer, err := json.Marshal(ar)
	return answer, err
}

func (rr AEResponce) Responser() ([]byte, error) {
	answer, err := json.Marshal(rr)
	return answer, err
}

func (rr MediaResponce) Responser() ([]byte, error) {
	answer, err := json.Marshal(rr)
	return answer, err
}

func (rr CertResponce) Responser() ([]byte, error) {
	answer, err := json.Marshal(rr)
	return answer, err
}

func (rr MesagesResponce) Responser() ([]byte, error) {
	answer, err := json.Marshal(rr)
	return answer, err
}

func (rr DialogResponce) Responser() ([]byte, error) {
	answer, err := json.Marshal(rr)
	return answer, err
}

func (cf CountFriend) Responser() ([]byte, error) {
	answer, err := json.Marshal(cf)
	return answer, err
}

func (gf GetFriends) Responser() ([]byte, error) {
	answer, err := json.Marshal(gf)
	return answer, err
}
func (uf UpdateFriend) Responser() ([]byte, error) {
	answer, err := json.Marshal(uf)
	return answer, err
}

func (ar AppendResponce) Responser() ([]byte, error) {
	answer, err := json.Marshal(ar)
	return answer, err
}
func (afr AddFriendResponce) Responser() ([]byte, error) {
	answer, err := json.Marshal(afr)
	return answer, err
}
func (rfr RemoveFriendResponce) Responser() ([]byte, error) {
	answer, err := json.Marshal(rfr)
	return answer, err
}
func (cr ComingResponce) Responser() ([]byte, error) {
	answer, err := json.Marshal(cr)
	return answer, err
}

func (blr BlackListResponce) Responser() ([]byte, error) {
	answer, err := json.Marshal(blr)
	return answer, err
}

func (blr ListBlocksResponce) Responser() ([]byte, error) {
	answer, err := json.Marshal(blr)
	return answer, err
}

func (gm GetMessages) Responser() ([]byte, error) {
	answer, err := json.Marshal(gm)
	return answer, err
}

func (gd GetDialogs) Responser() ([]byte, error) {
	answer, err := json.Marshal(gd)
	return answer, err
}

func (gd GetMedia) Responser() ([]byte, error) {
	answer, err := json.Marshal(gd)
	return answer, err
}

func (vtr ValidateTockenResponce) Responser() ([]byte, error) {
	answer, err := json.Marshal(vtr)
	return answer, err
}

func (sdcr SearchDialogsContainsResponce) Responser() ([]byte, error) {
	answer, err := json.Marshal(sdcr)
	return answer, err
}

func (smr SearchMessageResponce) Responser() ([]byte, error) {
	answer, err := json.Marshal(smr)
	return answer, err
}

func (mr MultiResponce) Responser() ([]byte, error) {
	answer, err := json.Marshal(mr)
	return answer, err
}

func (nar NonAddResponce) Responser() ([]byte, error) {
	answer, err := json.Marshal(nar)
	return answer, err
}

func (gu GetUser) Responser() ([]byte, error) {
	answer, err := json.Marshal(gu)
	return answer, err
}

func (gu GetAllMedia) Responser() ([]byte, error) {
	answer, err := json.Marshal(gu)
	return answer, err
}

func (gud GetUsersDialog) Responser() ([]byte, error) {
	answer, err := json.Marshal(gud)
	return answer, err
}

func (gem GetDelMessages) Responser() ([]byte, error) {
	answer, err := json.Marshal(gem)
	return answer, err
}
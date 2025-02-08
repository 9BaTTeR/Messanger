package responces

type Resp interface {
	Responser() ([]byte, error)
}

type Responce struct {
	Code        uint   `json:"code"`
	Description string `json:"message"`
}

type RegResponce struct {
	Status Responce
	Tocken string `json:"Tocken"`
}

type AuthResponce struct {
	Status Responce
	Tocken string
}

type AEResponce struct {
	Status Responce
}

type MediaResponce struct {
	Status Responce
	Key    string `json:"Key"`
}

type VerifedResponce struct {
	Responce Responce
}

type CertResponce struct {
	Status Responce
	Hash   string `json:"hash"`
}

type MesagesResponce struct {
	Status Responce
	MsgKey string `json:"msgkey"`
}

type DialogResponce struct {
	Status Responce
	Hash   string `json:"hash"`
}

type CountFriend struct {
	Status Responce
	Count  uint `json:"countfriend"`
}

type GetFriends struct {
	Status     Responce
	TakeFriend []Friend `json:"friends"`
}

type UpdateFriend struct {
	Status Responce
}

type Friend struct {
	ID          uint64 `json:"id"`
	Description string `json:"description"`
	DateAdd     string `json:"dateadd"`
	Photo       string `json:"photo"`
	Nickname    string `json:"nickname"`
}

type AppendResponce struct {
	Status Responce
}

type AddFriendResponce struct {
	Status Responce
}

type RemoveFriendResponce struct {
	Status Responce
}

type ComingResponce struct {
	Status  Responce
	Comings []Friend
}

type BlackListResponce struct {
	Status Responce
}

type ListBlocksResponce struct {
	Status Responce
	Blocks []Friend
}

type DMessages struct {
	DMKey        string `json:"msgKey" sqlite:"DMKey"`
	OMKey        string `json:"-" sqlite:"-"`
	IdUser       uint64 `json:"idUser" sqlite:"IdUser"`
	Date         string `json:"date" sqlite:"Date"`
	Content      string `json:"content" sqlite:"Content"`
	UpdateAt     string `json:"updateAt" sqlite:"UpdateAt"`
	DeleteAt     string `json:"deleteAt" sqlite:"DeletedAt"`
	ForwardedKey string `json:"forwardedKey" sqlite:"ForwardedKey"`
	Read         string `json:"read" sqlite:"Read"`
	Important    string `json:"important" sqlite:"Important"`
}

type Media struct {
	DMessages DMessages `json:"infoMessage" sqlite:"DMessages"`
	Hash      []string  `json:"hash" sqlite:"Hash"`
	Order     []uint64  `json:"order" sqlite:"Order"`
}

type GetMessages struct {
	Status   Responce
	Msg      []Media
	CountMsg uint64
}

type Dialog struct {
	Hash  string `json:"hash" sqlite:"Hash"`
	Name  string `json:"name" sqlite:"Name"`
	Photo string `json:"photo" sqlite:"Photo"`
}

type GetDialogs struct {
	Status      Responce
	Dialogs     []Dialog
	CountDialog uint64
}

type SearchMessageResponce struct {
	Status   Responce
	Messages []string
}

type SearchDialogsContainsResponce struct {
	Status  Responce
	Dialogs []string
}

type GetMedia struct {
	Status      Responce
	Hash        string
	BytesBase64 string
	Extension   string
}

type GetUser struct {
	Status   Responce
	Id       uint64
	Photo    string
	Nickname string
}

type GetDelMessages struct {
	Status Responce
	MsgKey []string
}

type GetAllMedia struct {
	Status Responce
	Hash   []string
}

type GetUsersDialog struct {
	Status  Responce
	UsersId []uint64
}

type ValidateTockenResponce struct {
	Status Responce
}

type MultiResponce struct {
	Status   Responce
	Messages []string
}

type NonAddResponce struct {
	Status Responce
}

package system

// import (
// 	"gin_chat/models"
// )

// 这个json要和前端一致
type User_Register struct {
	Name     string `json:"name" gorm:"unique;not null" validate:"required"`
	Password string `json:"password" validate:"required,min=2,max=20"`
	Identity string `json:"identity"`
	Salt     string `json:"salt"`
}

type User_Login struct {
	Name     string `json:"name" gorm:"unique;not null" validate:"required"`
	Password string `json:"password" validate:"required,min=2,max=20"`
}

type UpdateUserPayload struct{

}

type DeleteUserPayload struct{
	
}

type Frend struct{
	OwnerId int `json:"userid"`
	FrendId int `json:"users"`
}

type FrendsPayload struct {
	UserId string `json:"userid"`
	Users  interface{} `json:"users"`
}

type ChatPayload struct {
}

type CommunityPayload struct {

}

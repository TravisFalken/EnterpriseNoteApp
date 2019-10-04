package model

type User struct {
	UserName   string `json:"userName"`
	Password   string `json:"password"` //This will normally be encripted
	Email      string `json:"email"`
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
	SessionID  string `json:"sessionID"`
}

package user

type User struct {
	UserID     int    `json:"userID"`
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
	UserName   string `json:"userName"`
	Password   string `json:"password"` //This will normally be encripted
	SessionID  string `json:"sessionID"`
}

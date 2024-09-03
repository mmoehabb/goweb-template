package user

type Credentials struct{
    Username string `json:"username" xml:"username" form:"username"`
    Password string `json:"password" xml:"password" form:"password"`
}


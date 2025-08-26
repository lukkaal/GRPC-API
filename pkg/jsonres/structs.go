package jsonres

// response for gateway(gin)

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Error  string      `json:"error"`
}

// uerlogin:(userpbresp, token)
type TokenData struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

// tasklist
type DataList struct {
	Item  interface{} `json:"item"`
	Total int64       `json:"total"`
}

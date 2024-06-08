package types

type APIErrorMessage struct {
	ErrorMessage string `json:"error"`
}

type APISuccessMessage struct {
	SuccessMessage string `json:"success"`
	Data           any    `json:"data"`
}

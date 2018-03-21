package main


type ClientSafeError struct {
	Msg string // required string
	Code int // numeric code if wanted
}

func (c *ClientSafeError) Error() string {
	if c.Code > 0 {
		return "Error " + string(c.Code) + " " + c.Msg
	}
	return "Error" + c.Msg
}
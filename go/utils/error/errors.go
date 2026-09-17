package merr

var (
	/* password */

	ErrInvalidPassword  = NewErrorFactory(10001, "Invalid Password")
	ErrInvalidPwdParams = NewErrorFactory(10002, "Invalid Password Params")
	ErrInvalidPwdSalt   = NewErrorFactory(10003, "Invalid Password Salt")
	ErrInvalidPwdKey    = NewErrorFactory(10004, "Invalid Password Key")
)

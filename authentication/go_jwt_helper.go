package authentication

type JWTHelperImpl struct {
}

//GenerateJWT(id UserID) (*AuthenticatedUserJWT, error)
//ValidateJWT(jwt AuthenticatedUserJWT) (UserDetail, error)
//RenewJWT(jwt string) (*AuthenticatedUserJWT, error)

func (j *JWTHelperImpl) GenerateJWT(id UserID) (token *AuthenticatedUserJWT, err error) {
	panic("implement me")
}

func (j *JWTHelperImpl) RenewJWT(jwt string) (token *AuthenticatedUserJWT, err error) {
	panic("implement me")
}

func (j *JWTHelperImpl) ValidateJWT(jwt AuthenticatedUserJWT) (UserDetail, error) {
	panic("implement me")
}

func NewJWTHelper() *JWTHelperImpl {
	return &JWTHelperImpl{}
}

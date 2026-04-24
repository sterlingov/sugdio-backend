package domain

type UserShort struct {
	ID    int64
	Email string
	Role  string
}

type UserCredentials struct {
	ID           int64
	Email        string
	Role         string
	PasswordHash string
}

type UserCreate struct {
	Email    string
	Password string
	Role     *string
}

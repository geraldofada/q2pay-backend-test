package port

type AuthService interface {
	AuthorizeTransfer() (bool, error)
}

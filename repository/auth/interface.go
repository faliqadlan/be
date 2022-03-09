package auth

type Auth interface {
	Login(userName string, password string) (map[string]interface{}, error)
}

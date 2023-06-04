package constant

type Service string

const (
	SERVICE_USER  Service = "user"
	SERVICE_CACHE Service = "cache"
	SERVICE_LOBBY Service = "lobby"
)

// validate

const (
	MAX_ACCOUNT  = 30
	MIN_ACCOUNT  = 1
	MAX_PASSWORD = 30
	MIN_PASSWORD = 1
)

package piwigo

type PiwigoConfig struct {
	url string
	username string
	password string
}

type PiwigoCategory struct {
	id int
	name string
	key string
}
package piwigo

import "net/http/cookiejar"

type PiwigoContext struct {
	Url      string
	Username string
	Password string
	Cookies  *cookiejar.Jar
}

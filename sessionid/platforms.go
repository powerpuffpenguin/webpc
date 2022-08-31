package sessionid

var platforms = []string{
	`web`,
	`forward`,
	`socks5`,
	`android`,
	`ios`,
	`linux`,
	`windows`,
	`darwin`,
}

func Platforms() []string {
	return platforms
}

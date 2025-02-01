package sync

type config struct {
	clientID     string
	clientSecret string
	userName     string
	password     string
	debug        bool
	dir          string
}

func (c *config) ClientID() string {
	return c.clientID
}

func (c *config) ClientSecret() string {
	return c.clientSecret
}

func (c *config) UserName() string {
	return c.userName
}

func (c *config) Password() string {
	return c.password
}

func (c *config) Directory() string {
	return c.dir
}

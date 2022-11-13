package channel

type commonChannel struct {
	RequestChannel chan string
}

var CommonChannel commonChannel = commonChannel{make(chan string)}

func (c *commonChannel) Destroy() {
	close(c.RequestChannel)
}

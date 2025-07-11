package commands

func (c *commands) registerModules(modules ...commandModule) *commands {
	for _, module := range modules {
		c.Commands = append(c.Commands, module.Definitions()...)
		for name, handler := range module.Commands() {
			c.handlers[name] = handler
		}
	}
	return c
}

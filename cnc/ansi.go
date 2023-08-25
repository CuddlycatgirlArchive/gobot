package main

func getConsoleTitleAnsi(title string) string {
	return "\u001B]0;" + title + "\a"
}

func (c *client) setConsoleTitle(title string) {
	c.conn.Write([]byte(getConsoleTitleAnsi(title)))
}

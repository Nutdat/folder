package main

import (
	"github.com/Nutdat/folder"
	"github.com/Nutdat/logger"
)

func main() {
	defer logger.RecoverAndFlush()
	folder.CreateFolder("media")
	folder.CreateFolder("media/file")
	folder.DeleteFolder("media")
}

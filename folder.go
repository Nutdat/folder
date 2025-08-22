package folder

import (
	"github.com/Nutdat/folder/core"
	"github.com/Nutdat/logger"
)

func init() {
	logger.LogInit("Folder")
	core.CheckAndRestoreFolders()

}

func CreateFolder(folderPath string) {
	core.CreateFolder(folderPath)
}

func DeleteFolder(folderPath string) {
	core.RemoveFolder(folderPath)
}

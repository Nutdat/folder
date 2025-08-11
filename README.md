Folder Package 

This Go package provides easy and persistent folder management for your projects. It automatically creates, remembers, and restores folders inside a hidden `.Nutdat` base directory. Folder metadata is cached in a JSON file, ensuring your folder structure survives restarts. 
---

Features

- Persistent folder tracking: Creates folders under `./.Nutdat`, records them in a JSON cache (`created_folders.json`) inside a hidden `.nutburrow` folder.
- Auto restore: On start, checks all remembered folders and recreates missing ones automatically.
- Thread-safe: Uses mutex locks to ensure concurrent safety.
- Simple API: Create, check/restore, and remove folders with easy-to-use functions.


---

Installation


go get github.com/Nutdat/folder

---

Usage

Initialization

On app startup, load the remembered folders and ensure the structure is restored:

func init() {
    core.CheckAndRestoreFolders()
}

Create a new folder

core.CreateFolder("log")
core.CreateFolder("cache")

Remove a folder

core.RemoveFolder("temp")

---

Folder Structure

- All folders are created under the hidden `.Nutdat` directory in your project root.
- Metadata is stored at `.Nutdat/.nutburrow/created_folders.json`.
- Folder paths are stored relative to `.Nutdat` for portability.

---

Logger Integration

This package automatically logs:

- Folder creation, restoration, and removal events.
- Errors when folder operations fail.
- Duration of folder checks/restoration.

---

Thread Safety

All public functions use internal mutex locking to allow safe concurrent usage.

---

Notes

- This package does not return errors; all issues are logged instead.
- You can extend the hardcoded folder list by adding calls to `CreateFolder` on init.
- Ensure your `logger` package is imported and properly initialized for best experience.

---

License

MIT License — feel free to use and modify!

package internal

import "os"

func InitMessageFile() {
	err := os.WriteFile(plainTextFileName+defaultExtension, []byte("[Write your message here]"), 0777)
	if err != nil {
		Error("Initialize message file", err)
		return
	}

	Success("Initialize message file")
}

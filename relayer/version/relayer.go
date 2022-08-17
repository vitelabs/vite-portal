package version

import "fmt"

func PrintBuildVersion() {
	if PROJECT_COMMIT_VERSION != "" {
		fmt.Printf("version: %s, git commit: %s\n", PROJECT_BUILD_VERSION, PROJECT_COMMIT_VERSION)
	} else {
		fmt.Println("cannot read gitversion file please build with Make first")
	}
}
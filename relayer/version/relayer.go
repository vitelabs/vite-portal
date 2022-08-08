package version

import "fmt"

func PrintBuildVersion() {
	if RELAYER_COMMIT_VERSION != "" {
		fmt.Printf("version: %s, git commit: %s\n", RELAYER_BUILD_VERSION, RELAYER_COMMIT_VERSION)
	} else {
		fmt.Println("cannot read gitversion file please build with Make first")
	}
}
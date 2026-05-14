package utils

import (
	"fmt"
	"os/exec"
)


const accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiYXVkIjoiaHR0cDovLzIwLjI0NC41Ni4xNDQvZXZhbHVhdGlvbi1zZXJ2aWNlIiwiZW1haWwiOiJwYXJ2YWdnYXJ3YWwxMzBAZ21haWwuY29tIiwiZXhwIjoxNzc4NzYzOTM0LCJpYXQiOjE3Nzg3NjMwMzQsImlzcyI6IkFmZm9yZCBNZWRpY2FsIFRlY2hub2xvZ2llcyBQcml2YXRlIExpbWl0ZWQiLCJqdGkiOiI0OTIzMDRiNi0yNjYyLTRiOWUtOWI4NC04MzNmZTZiYTI0NDEiLCJsb2NhbGUiOiJlbi1JTiIsIm5hbWUiOiJwYXJ2IGFnZ2Fyd2FsIiwic3ViIjoiNzViOGI3NjQtMmM2MS00Mzg4LWFhOTYtZDNjMTI2ZWMxODA1In0sImVtYWlsIjoicGFydmFnZ2Fyd2FsMTMwQGdtYWlsLmNvbSIsIm5hbWUiOiJwYXJ2IGFnZ2Fyd2FsIiwicm9sbE5vIjoiMTIzMTA3MDAiLCJhY2Nlc3NDb2RlIjoiVFJ2WldxIiwiY2xpZW50SUQiOiI3NWI4Yjc2NC0yYzYxLTQzODgtYWE5Ni1kM2MxMjZlYzE4MDUiLCJjbGllbnRTZWNyZXQiOiJId3hSQVRlaFlydmpiVUROIn0.V51US7AXMsdzkafe3WREw--9vO_HebM7wJQsehX3jiU"


func Log(stack, level, pkg, message string) {
	script := fmt.Sprintf(`
		const { initLogger, Log } = require('../loging_middleware/index.js');
		initLogger('%s');
		Log('%s', '%s', '%s', '%s');
	`, accessToken, stack, level, pkg, message)

	cmd := exec.Command("node", "-e", script)


	go func() {
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("--- LOGGER BRIDGE ERROR ---\nNode JS Output: %s\nError: %v\n---------------------------\n", string(out), err)
		}
	}()
}
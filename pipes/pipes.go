package pipes

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/sys/windows"
)

type PipeInfo struct {
	PID          uint32 `json:"pid"`
	Name         string `json:"name"`
	Instances    uint32 `json:"instances"`
	MaxInstances uint32 `json:"max_instances"`
	Flags        string `json:"flags"`
}

func pipeFlagsToString(flags uint32) string {
	var flagString []string
	if flags&windows.PIPE_SERVER_END != 0 {
		flagString = append(flagString, "PIPE_SERVER_END")
	} else {
		flagString = append(flagString, "PIPE_CLIENT_END")
	}

	if flags&windows.PIPE_TYPE_MESSAGE != 0 {
		flagString = append(flagString, "PIPE_TYPE_MESSAGE")
	} else {
		flagString = append(flagString, "PIPE_TYPE_BYTE")
	}

	return strings.Join(flagString, ",")
}

func GenPipes() ([]PipeInfo, error) {
	var pipeSearch string = `\\.\pipe\*`
	var findFileData windows.Win32finddata

	findHandle, err := windows.FindFirstFile(windows.StringToUTF16Ptr(pipeSearch), &findFileData)
	if err != nil {
		return nil, fmt.Errorf("failed to find first pipe: %v", err)
	}
	defer windows.FindClose(findHandle)

	var pipes []PipeInfo
	for {
		pipeName := windows.UTF16ToString(findFileData.FileName[:])
		pipePath := `\\.\pipe\` + pipeName

		pipeInfo := PipeInfo{
			Name:         pipeName,
			PID:          0,
			Instances:    0,
			MaxInstances: 0,
			Flags:        "",
		}

		pipeHandle, err := windows.CreateFile(
			windows.StringToUTF16Ptr(pipePath),
			windows.GENERIC_READ,
			windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
			nil,
			windows.OPEN_EXISTING,
			0,
			0,
		)

		if err == nil {
			// Try to get server PID first, then client PID
			if err = windows.GetNamedPipeServerProcessId(pipeHandle, &pipeInfo.PID); err != nil {
				if err = windows.GetNamedPipeClientProcessId(pipeHandle, &pipeInfo.PID); err != nil {
					log.Printf("failed to get pipe process ID: %v", err)
				}
			}

			if err = windows.GetNamedPipeHandleState(
				pipeHandle,
				nil,
				&pipeInfo.Instances,
				nil,
				nil,
				nil,
				0,
			); err != nil {
				log.Printf("failed to get pipe handle state: %v", err)
			}

			var pipeFlags uint32
			if err = windows.GetNamedPipeInfo(
				pipeHandle,
				&pipeFlags,
				nil,
				nil,
				&pipeInfo.MaxInstances,
			); err != nil {
				log.Printf("failed to get pipe info: %v", err)
			} else {
				pipeInfo.Flags = pipeFlagsToString(pipeFlags)
			}
		} else {
			log.Printf("failed to open pipe %s: %v", pipePath, err)
		}

		pipes = append(pipes, pipeInfo)

		// Close the pipe handle
		windows.CloseHandle(pipeHandle)

		if err = windows.FindNextFile(findHandle, &findFileData); err != nil {
			if err == windows.ERROR_NO_MORE_FILES {
				break
			}
			return nil, fmt.Errorf("failed to find next file: %v", err)
		}
	}
	return pipes, nil
}

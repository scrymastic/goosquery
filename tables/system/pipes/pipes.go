package pipes

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/sys/windows"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

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

func GenPipes(ctx *sqlctx.Context) (*result.Results, error) {
	var pipeSearch string = `\\.\pipe\*`
	var findFileData windows.Win32finddata

	findHandle, err := windows.FindFirstFile(windows.StringToUTF16Ptr(pipeSearch), &findFileData)
	if err != nil {
		return nil, fmt.Errorf("failed to find first pipe: %v", err)
	}
	defer windows.FindClose(findHandle)

	pipes := result.NewQueryResult()
	for {
		pipeName := windows.UTF16ToString(findFileData.FileName[:])
		pipePath := `\\.\pipe\` + pipeName

		pipeInfo := result.NewResult(ctx, Schema)
		pipeInfo.Set("name", pipeName)

		pipeHandle, err := windows.CreateFile(
			windows.StringToUTF16Ptr(pipePath),
			windows.GENERIC_READ,
			windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
			nil,
			windows.OPEN_EXISTING,
			0,
			0,
		)

		if err != nil {
			log.Printf("failed to open pipe %s: %v", pipePath, err)
			goto NextPipe
		}

		if ctx.IsColumnUsed("pid") {
			// Try to get server PID first, then client PID
			var pid uint32
			if err = windows.GetNamedPipeServerProcessId(pipeHandle, &pid); err != nil {
				log.Printf("failed to get pipe process ID: %v", err)
			}
			pipeInfo.Set("pid", pid)
		}

		if ctx.IsColumnUsed("instances") {
			var instances uint32
			if err = windows.GetNamedPipeHandleState(
				pipeHandle,
				nil,
				&instances,
				nil,
				nil,
				nil,
				0,
			); err != nil {
				log.Printf("failed to get pipe handle state: %v", err)
			}
			pipeInfo.Set("instances", instances)
		}

		if ctx.IsAnyOfColumnsUsed([]string{"max_instances", "flags"}) {
			var pipeFlags uint32
			var maxInstances uint32
			if err = windows.GetNamedPipeInfo(
				pipeHandle,
				&pipeFlags,
				nil,
				nil,
				&maxInstances,
			); err != nil {
				log.Printf("failed to get pipe info: %v", err)
			} else {
				pipeInfo.Set("flags", pipeFlagsToString(pipeFlags))
				pipeInfo.Set("max_instances", maxInstances)
			}
		}

		pipes.AppendResult(*pipeInfo)

	NextPipe:
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

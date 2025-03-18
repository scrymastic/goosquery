package scheduled_tasks

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// The hidden field looks not correct. It still needs to be refined.
type ScheduledTask struct {
	Name           string `json:"name"`
	Action         string `json:"action"`
	Path           string `json:"path"`
	Enabled        bool   `json:"enabled"`
	State          string `json:"state"`
	Hidden         bool   `json:"hidden"`
	LastRunTime    int64  `json:"last_run_time"`
	NextRunTime    int64  `json:"next_run_time"`
	LastRunMessage string `json:"last_run_message"`
	LastRunCode    uint32 `json:"last_run_code"`
}

var taskStates = map[int]string{
	0: "unknown",
	1: "disabled",
	2: "queued",
	3: "ready",
	4: "running",
}

var (
	CLSID_TaskScheduler = ole.NewGUID("{0F87369F-A4E5-4CFC-BD3E-73E6154572DD}")
	IID_ITaskScheduler  = ole.NewGUID("{148BD520-A2AB-11CE-B11F-00AA00476E5D}")
)

func getTasksFromFolder(folder *ole.IDispatch) ([]ScheduledTask, error) {
	tasks := []ScheduledTask{}

	// Process subfolders
	subfoldersVariant, err := oleutil.CallMethod(folder, "GetFolders", 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get subfolders: %w", err)
	}
	folders := subfoldersVariant.ToIDispatch()
	defer folders.Release()

	count, err := oleutil.GetProperty(folders, "Count")
	if err != nil {
		return nil, fmt.Errorf("failed to get subfolder count: %w", err)
	}

	for i := 1; i <= int(count.Val); i++ {
		subfolderVariant, err := oleutil.GetProperty(folders, "Item", i)
		if err != nil {
			continue
		}

		subfolder := subfolderVariant.ToIDispatch()
		subtasks, err := getTasksFromFolder(subfolder)
		if err != nil {
			subfolder.Release()
			continue
		}
		tasks = append(tasks, subtasks...)
		subfolder.Release()
	}

	// Process tasks in current folder
	tasksVariant, err := oleutil.CallMethod(folder, "GetTasks", 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	taskCollection := tasksVariant.ToIDispatch()
	defer taskCollection.Release()

	count, err = oleutil.GetProperty(taskCollection, "Count")
	if err != nil {
		return nil, fmt.Errorf("failed to get task count: %w", err)
	}

	for i := 1; i <= int(count.Val); i++ {
		taskVariant, err := oleutil.GetProperty(taskCollection, "Item", i)
		if err != nil {
			continue
		}
		task := taskVariant.ToIDispatch()
		if parsedTask := parseTask(task); parsedTask != nil {
			tasks = append(tasks, *parsedTask)
		}
		task.Release()
	}

	return tasks, nil
}

func parseTask(taskObj *ole.IDispatch) *ScheduledTask {
	task := &ScheduledTask{}

	// Get basic properties
	if name, err := oleutil.GetProperty(taskObj, "Name"); err == nil {
		task.Name = name.ToString()
	}
	if path, err := oleutil.GetProperty(taskObj, "Path"); err == nil {
		task.Path = path.ToString()
	}
	if enabled, err := oleutil.GetProperty(taskObj, "Enabled"); err == nil {
		task.Enabled = enabled.Val != 0
	}
	if hidden, err := oleutil.GetProperty(taskObj, "Hidden"); err == nil {
		fmt.Printf("Hidden property value: %v\n", hidden.Val)
		task.Hidden = hidden.Val != 0
	}
	if state, err := oleutil.GetProperty(taskObj, "State"); err == nil {
		if s, ok := taskStates[int(state.Val)]; ok {
			task.State = s
		} else {
			task.State = fmt.Sprintf("unknown (%d)", int(state.Val))
		}
	}

	// Get time properties
	if lastRun, err := oleutil.GetProperty(taskObj, "LastRunTime"); err == nil && lastRun.VT == ole.VT_DATE {
		task.LastRunTime = lastRun.Value().(time.Time).Unix()
	}
	if nextRun, err := oleutil.GetProperty(taskObj, "NextRunTime"); err == nil && nextRun.VT == ole.VT_DATE {
		task.NextRunTime = nextRun.Value().(time.Time).Unix()
	} else {
		task.NextRunTime = -1 // Use -1 to indicate "N/A" for int64
	}

	// Get result properties
	if result, err := oleutil.GetProperty(taskObj, "LastTaskResult"); err == nil {
		task.LastRunCode = uint32(result.Val)
		if task.LastRunCode == 0 {
			task.LastRunMessage = "The operation completed successfully."
		} else if err := ole.NewError(uintptr(task.LastRunCode)); err != nil {
			task.LastRunMessage = err.Error()
		}
	}

	// Get actions
	definition, err := oleutil.GetProperty(taskObj, "Definition")
	if err != nil {
		return task
	}
	defObj := definition.ToIDispatch()
	defer defObj.Release()

	actions, err := oleutil.GetProperty(defObj, "Actions")
	if err != nil {
		return task
	}

	actionsObj := actions.ToIDispatch()
	defer actionsObj.Release()

	var actionStrings []string
	count, err := oleutil.GetProperty(actionsObj, "Count")
	if err != nil {
		return task
	}
	for i := 1; i <= int(count.Val); i++ {
		actionItem, err := oleutil.GetProperty(actionsObj, "Item", i)
		if err != nil {
			continue
		}
		action := actionItem.ToIDispatch()
		defer action.Release()

		actionType, err := oleutil.GetProperty(action, "Type")
		if err != nil || actionType.Val != 0 {
			continue
		}
		parts := []string{}
		for _, prop := range []string{"WorkingDirectory", "Path", "Arguments"} {
			if val, err := oleutil.GetProperty(action, prop); err == nil && val.ToString() != "" {
				parts = append(parts, strings.TrimSpace(val.ToString()))
			}
		}
		if len(parts) > 0 {
			actionStrings = append(actionStrings, strings.Join(parts, " "))
		}
	}
	if len(actionStrings) > 0 {
		task.Action = strings.Join(actionStrings, ",")
	}

	return task
}

func GenScheduledTasks() ([]ScheduledTask, error) {
	if err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED); err != nil {
		return nil, fmt.Errorf("failed to initialize COM: %w", err)
	}
	defer ole.CoUninitialize()

	unknown, err := ole.CreateInstance(CLSID_TaskScheduler, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Task Scheduler instance: %w", err)
	}
	defer unknown.Release()

	service, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, fmt.Errorf("failed to query interface: %w", err)
	}
	defer service.Release()

	if _, err := service.CallMethod("Connect"); err != nil {
		return nil, fmt.Errorf("failed to connect to Task Scheduler: %w", err)
	}

	rootFolder, err := oleutil.CallMethod(service, "GetFolder", "\\")
	if err != nil {
		return nil, fmt.Errorf("failed to get root folder: %w", err)
	}
	folder := rootFolder.ToIDispatch()
	defer folder.Release()

	return getTasksFromFolder(folder)
}

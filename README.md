<div align="center">

<h1 align="center" style="font-family: 'Segoe UI', sans-serif; font-size: 72px; font-style: italic; font-weight: bold; margin-bottom: 20px;">
  <span style="color: #0066cc;">GO</span><span style="color: #eeeeee;">OSQUERY</span>
</h1>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/scrymastic/goosquery">
    <img src="https://goreportcard.com/badge/github.com/scrymastic/goosquery" alt="Go Report Card" />
  </a>
  <a href="https://pkg.go.dev/github.com/scrymastic/goosquery">
    <img src="https://pkg.go.dev/badge/github.com/scrymastic/goosquery.svg" alt="Go Reference" />
  </a>
  <a href="https://github.com/scrymastic/goosquery/blob/main/LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="License" />
  </a>
  <a href="#">
    <img src="https://img.shields.io/badge/platform-Windows-blue" alt="Platform" />
  </a>
</p>

<p align="center">
  <b>Go osquery with JSON output, currently working on Windows.</b>
</p>

</div>

<p align="center">
  The goal is to provide a lightweight, portable, and easy-to-use version of osquery that can be integrated into other projects.
</p>

---

## ✨ Features

- 🚀 Lightweight and portable implementation
- 📊 JSON output for easy parsing
- 🖥️ Windows-specific system information
- 🧩 Modular table implementation
- 🔌 Easy integration with other Go projects

## 📦 Installation

```bash
go get github.com/scrymastic/goosquery
```

## 🚀 Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/scrymastic/goosquery/system/os_version"
)

func main() {
    // Get OS version information
    osInfo, err := os_version.GenOSVersion()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Print as JSON
    jsonData, err := json.MarshalIndent(osInfo, "", "  ")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(string(jsonData))
}
```

## 📋 Example Output

```json
{
  "name": "Microsoft Windows 10 Pro",
  "version": "10.0.19041",
  "major": 10,
  "minor": 0,
  "patch": 19041,
  "build": "19041",
  "platform": "windows",
  "platform_like": "windows",
  "codename": "Microsoft Windows 10 Pro",
  "arch": "64-bit",
  "install_date": 1596240000,
  "revision": 928
}
```

## 📊 Implementation Status

<table>
  <tr>
    <th>Status</th>
    <th>Icon</th>
    <th>Description</th>
  </tr>
  <tr>
    <td>Not Started</td>
    <td align="center">⏳</td>
    <td>Work has not begun on this table yet.</td>
  </tr>
  <tr>
    <td>In Progress</td>
    <td align="center">🛠️</td>
    <td>Actively being developed.</td>
  </tr>
  <tr>
    <td>Completed</td>
    <td align="center">✅</td>
    <td>Fully implemented and tested.</td>
  </tr>
  <tr>
    <td>Testing</td>
    <td align="center">🧪</td>
    <td>Development is done, but under testing for bugs or issues.</td>
  </tr>
  <tr>
    <td>Blocked</td>
    <td align="center">⛔</td>
    <td>Development is paused due to dependencies, blockers, or technical issues.</td>
  </tr>
  <tr>
    <td>Planned</td>
    <td align="center">🗓️</td>
    <td>Table is planned for future implementation but hasn't started yet.</td>
  </tr>
  <tr>
    <td>Deprecated</td>
    <td align="center">🗑️</td>
    <td>This table is no longer relevant or supported in this implementation.</td>
  </tr>
</table>

## 📑 Tables

<details>
<summary><b>Click to expand table list</b></summary>
<br>

| Table Name                       | Status  |
|----------------------------------|---------|
| appcompat_shims                  | 🧪      |
| arp_cache                        | 🧪      |
| authenticode                     | ✅      |
| autoexec                         | ⏳      |
| azure_instance_metadata          | ⏳      |
| azure_instance_tags              | ⏳      |
| background_activities_moderator  | 🧪      |
| battery                          | ⛔      |
| bitlocker_info                   | ✅      |
| carbon_black_info                | ⏳      |
| carves                           | ⏳      |
| certificates                     | 🛠️      |
| chassis_info                     | ✅      |
| chocolatey_packages              | ✅      |
| chrome_extension_content_scripts | ⏳      |
| chrome_extensions                | ⏳      |
| connectivity                     | ✅      |
| cpu_info                         | ✅      |
| cpuid                            | ⏳      |
| curl                             | ✅      |
| curl_certificate                 | ⏳      |
| default_environment              | ✅      |
| deviceguard_status               | ✅      |
| disk_info                        | ✅      |
| dns_cache                        | ✅      |
| drivers                          | ✅      |
| ec2_instance_metadata            | ⏳      |
| ec2_instance_tags                | ⏳      |
| etc_hosts                        | ✅      |
| etc_protocols                    | ✅      |
| etc_services                     | ✅      |
| file                             | ✅      |
| firefox_addons                   | ⏳      |
| groups                           | ✅      |
| hash                             | ✅      |
| ie_extensions                    | ⏳      |
| intel_me_info                    | ⏳      |
| interface_addresses              | ✅      |
| interface_details                | ✅      |
| kernel_info                      | ✅      |
| kva_speculative_info             | ✅      |
| listening_ports                  | ✅      |
| logged_in_users                  | ✅      |
| logical_drives                   | ✅      |
| logon_sessions                   | ✅      |
| memory_devices                   | 🧪      |
| npm_packages                     | ⏳      |
| ntdomains                        | ✅      |
| ntfs_acl_permissions             | ⏳      |
| ntfs_journal_events              | ⏳      |
| office_mru                       | ⏳      |
| os_version                       | ✅      |
| osquery_events                   | 🗑️      |
| osquery_extensions               | 🗑️      |
| osquery_flags                    | 🗑️      |
| osquery_info                     | 🗑️      |
| osquery_packs                    | 🗑️      |
| osquery_registry                 | 🗑️      |
| osquery_schedule                 | 🗑️      |
| patches                          | ✅      |
| physical_disk_performance        | ⏳      |
| pipes                            | 🧪      |
| platform_info                    | 🧪      |
| powershell_events                | ⏳      |
| prefetch                         | ⏳      |
| process_etw_events               | ⏳      |
| process_memory_map               | ✅      |
| process_open_sockets             | ✅      |
| processes                        | 🧪      |
| programs                         | 🧪      |
| python_packages                  | 🧪      |
| registry                         | 🧪      |
| routes                           | 🧪      |
| scheduled_tasks                  | 🧪      |
| secureboot                       | ⏳      |
| security_profile_info            | 🛠️      |
| services                         | ✅      |
| shared_resources                 | ✅      |
| shellbags                        | ⏳      |
| shimcache                        | ⏳      |
| ssh_configs                      | ⏳      |
| startup_items                    | 🛠️      |
| system_info                      | 🧪      |
| time                             | ✅      |
| tpm_info                         | ⏳      |
| uptime                           | ✅      |
| user_groups                      | 🧪      |
| user_ssh_keys                    | ⏳      |
| userassist                       | ⏳      |
| users                            | ✅      |
| video_info                       | ⏳      |
| vscode_extensions                | ⏳      |
| winbaseobj                       | 🧪      |
| windows_crashes                  | ⏳      |
| windows_eventlog                 | ⏳      |
| windows_events                   | ⏳      |
| windows_firewall_rules           | 🧪      |
| windows_optional_features        | ✅      |
| windows_search                   | ⛔      |
| windows_security_center          | 🧪      |
| windows_security_products        | 🛠️      |
| windows_update_history           | 🛠️      |
| wmi_bios_info                    | ⏳      |
| wmi_cli_event_consumers          | ⏳      |
| wmi_event_filters                | ⏳      |
| wmi_filter_consumer_binding      | ⏳      |
| wmi_script_event_consumers       | ⏳      |
| yara                             | ⛔      |
| yara_events                      | ⛔      |
| ycloud_instance_metadata         | ⛔      |
</details>

<div align="center">
  <p>
    <a href="https://github.com/scrymastic/goosquery/issues">Report Bug</a>
    ·
    <a href="https://github.com/scrymastic/goosquery/issues">Request Feature</a>
  </p>
</div>

<div align="center">

# <span style="font-weight: bold; font-size: 1.5em;">GOOSQUERY</span>
<img src="https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go Version" />
<img src="https://img.shields.io/badge/Platform-Windows-0078D6?style=flat-square&logo=windows&logoColor=white" alt="Platform" />
<img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square" alt="License" />

### Windows System Information Collector

</div>

Goosquery is a Go-based system information collection tool inspired by OSQuery. It provides a unified interface to collect various system information from Windows systems.

## Features

- Collects system information from various sources
- Organizes data into JSON format for easy analysis
- Provides benchmarking capabilities to measure performance
- Logs all operations for debugging and auditing
- Modular design with clean separation of concerns

## Usage

### Basic Usage

```bash
# Run with default output directory (reports)
goosquery

# Specify a custom output directory
goosquery /path/to/output
```

### Options

```bash
# Display help information
goosquery --help
```

## Output

Goosquery generates the following output:

- JSON files for each collector's data
- A compressed zip file containing all collected data
- Detailed logging information

## Development

### Adding a New Collector

1. Create a new package in the appropriate directory (networking, system, or utility)
2. Implement the collector function with the standard Gen* naming convention
3. Add the type and function to the corresponding main package file (networking.go, system.go, or utility.go)
4. Add the collector to the collectors list in main.go

### Implementation Status

<details>
<summary><strong>Status Legend</strong></summary>

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
</details>

<details>
<summary><strong>Tables Implementation Status</strong></summary>

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
| memory_devices                   | ✅      |
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

## Requirements

- Go 1.18 or higher
- Windows operating system

## License

MIT

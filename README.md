# <span style="font-weight: bold">goosquery </span>
![Go](https://img.shields.io/badge/made%20with-Go-00ADD8)
![Platform](https://img.shields.io/badge/platform-windows-0078D6)
![License](https://img.shields.io/badge/license-MIT-yellow)
</div>

Goosquery is a Go-based system information collection tool inspired by OSQuery. It provides a unified interface to collect various system information from Windows systems using SQL-like queries.

## Table of Contents

- [Notice](#notice)
- [Features](#features)
- [Usage](#usage)
  - [Interactive Mode](#interactive-mode)
  - [Command Line Mode](#command-line-mode)
  - [Output Formats](#output-formats)
- [Examples](#examples)
- [Development](#development)
- [Implementation Status](#implementation-status)
- [Adding New Tables](#adding-new-tables)

## Notice

> [!IMPORTANT]  
> This project is under development.

## Features

- ✅ SQL-like query interface for accessing system information
- ✅ Interactive mode with command history and autocompletion
- ✅ Multiple output formats (table and JSON)
- ✅ Modular design for easy extension with new tables
- ✅ Efficient data collection with column-based filtering
- ✅ Comprehensive table collection for Windows systems
- ✅ Support for WHERE clauses to filter results

## Usage

### Interactive Mode

```bash
# Start in interactive mode with default table output
goosquery -i

# Start in interactive mode with JSON output
goosquery -i -json
```

In interactive mode, you can use the following commands:

```
.quit        - Exit the program
.json        - Switch to JSON output mode
.table       - Switch to table output mode  
.mode        - Show current output mode
.help        - Show help message
```

### Command Line Mode

```bash
# Execute a query directly from the command line
goosquery -q "SELECT * FROM processes LIMIT 5"

# Execute a query and output as JSON
goosquery -q "SELECT name, pid FROM processes" -json
```

### Output Formats

GoOsquery supports two output formats:

1. **Table Format (Default)**: Displays results in a formatted ASCII table
   ```
   +------+-----+
   | name | pid |
   +------+-----+
   | cmd  | 123 |
   +------+-----+
   ```

2. **JSON Format**: Outputs results as JSON data
   ```json
   [
     {
       "name": "cmd",
       "pid": 123
     }
   ]
   ```

## Examples

Query processes:
```sql
SELECT name, pid, parent FROM processes WHERE name LIKE '%svc%';
```

Query network interfaces:
```sql
SELECT interface, address, mask FROM interface_addresses;
```

Check HTTP response from a website:
```sql
SELECT url, response_code, round_trip_time FROM curl WHERE url = 'https://example.com';
```

## Development

### Implementation Status

<details>
<summary><strong>Status Legend</strong></summary>

| Status | Icon | Description |
|--------|------|-------------|
| Not Started | ⏳ | Work has not begun on this table yet. |
| In Progress | 🛠️ | Actively being developed. |
| Completed | ✅ | Fully implemented and tested. |
| Testing | 🧪 | Development is done, but under testing for bugs or issues. |
| Blocked | ⛔ | Development is paused due to dependencies, blockers, or technical issues. |
| Planned | 🗓️ | Table is planned for future implementation but hasn't started yet. |
| Deprecated | 🗑️ | This table is no longer relevant or supported in this implementation. |
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
| platform_info                    | ✅      |
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

## Adding New Tables

To add a new table to Goosquery, follow these steps:

1. Create a new package for your table under the appropriate category (e.g., `tables/system/newtable/`)

2. Implement a data generator function with the following signature:
   ```go
   func GenNewTable(context context.Context) ([]map[string]interface{}, error)
   ```
   
   This function should:
   - Use the context.IsColumnUsed() method to check which columns are needed
   - Only fetch the data needed for the columns requested
   - Return data as a slice of maps where keys are column names

3. Define column types in a map to ensure consistent type handling:
   ```go
   var columnDefs = map[string]string{
       "column1": "string",
       "column2": "int64",
       "column3": "int32",
   }
   ```

4. Use the utility function to initialize columns with appropriate default values:
   ```go
   result := util.InitColumns(ctx, columnDefs)
   ```

5. Register the table in `sql/executor/executor.go` by adding a new case to the `GetExecutor` function:
   ```go
   case "newtable":
       return &tableExecutor{
           BaseExecutor: BaseExecutor{},
           tableName:    "newtable",
           dataFunc:     newtable.GenNewTable,
       }, nil
   ```

The `tableExecutor` pattern eliminates the need to write custom executor code for each table. It handles:
- Extracting selected columns 
- Filtering with WHERE clauses
- Projecting only requested columns

This design pattern reduces code duplication and ensures consistent behavior across all tables.

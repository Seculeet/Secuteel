# Secuteel
The *Secuteel* app is an auditing tool, with multi-platform support, written in go.
- [Installation](https://github.com/Seculeet/secuteel#installation)
	- [Compile](https://github.com/Seculeet/secuteel#compile)
		- [Git](https://github.com/Seculeet/secuteel#git)
		- [Go](https://github.com/Seculeet/secuteel#go)
- [Help](https://github.com/Seculeet/secuteel#help)
- [Documentation](https://github.com/Seculeet/secuteel#documentation)
	- [Command overview](https://github.com/Seculeet/secuteel#command-overview)
	- [Create a config file](https://github.com/Seculeet/secuteel#create-a-config-file)
	- [Note](https://github.com/Seculeet/secuteel#note)
	- [Additional JavaScript functions](https://github.com/Seculeet/secuteel#additional-javascript-functions)
	- [Supported Commands through wrapper](https://github.com/Seculeet/secuteel#supported-commands-through-wrapper)
	- [Start a scan](https://github.com/Seculeet/secuteel#start-a-scan)
- [Example usage](https://github.com/Seculeet/secuteel#example-usage)
	- [Config to get started](https://github.com/Seculeet/secuteel#config-to-get-started)
		- [Windows](https://github.com/Seculeet/secuteel#windows)
		- [Linux](https://github.com/Seculeet/secuteel#linux)
## Installation
- You can download our pre-compiled binaries [here](https://github.com/Seculeet/secuteel/tree/main/binaries)
- On Linux you may want to run `chmod +x secuteel`, before executing it
### Compile
- *Secuteel* depends on Go 1.16 or higher
#### Git
- You can get and compile it from scratch using the following command:
```bash
git clone https://github.com/seculeet/secuteel.git; cd secuteel; go get; go build
```
#### Go
- You can also get it using go:
```bash
go get -u github.com/seculeet/secuteel
```
## Help
- If you are stuck at getting a command to run, activate the `-debug` flag at startup
- You can also activate verbose mode using `-v`
## Documentation
- This will get you started on running your first audit. If you are running into problems check the `error.log`
### Command overview
```bash
-input, --input= 'expects path to config file (.json)'
-output, --output= 'specify your output file (.zip)'
-add, --add= 'add a list of allowed commands'
-v 'toggle verbose mode for console'
-s 'skip sanity check'
-h 'help'
-debug 'activate debug mode for log files'
-p 'set password to encrypt output zip folder'
```

### Create a config file
```json
{  
  "system":
  {
    "systemName": "The current operating system",
    "version": "The OS version",
    "shell": "A system shell (e.g. CMD, Powershell), optional",
    "argument": "The argument used to execute commands (e.g. /C), optional",
    "root": "Specify if the audit has to be run as root, optional (default false)"
  },
  "commands": [  
    {
      "name": "Example Audit Step",  
      "command": "You can script in JavaScript and add your audit command in here",  
      "dontSaveArtefact": true,  
      "blackenContent": "Regex pattern",  
      "typeExpected": "containsReg",  
      "expected": "Regex pattern",  
      "description": "This is just for taking notes of what is happening"  
    }
  ]
}
```
### Note
- `name` Has to be unique and without special characters (required)
- `command` You can use everything JavaScript has to offer and additionaly our own implemented functions (required)
- `dontSaveArtefact` Default is false, if set to true no artefacts for this Auditstep will be saved (optional)
- `blackenContent` Censor the saved artefacts with the given pattern. Can't be used if `dontSaveArtefact` is toggled true (optional)
- `typeExpected` Default is `==`, (optional)
	- Supported operators for Strings: `==, !=, contains, containsReg, nil`
	- Supported operators for integers: `==, !=, <, <=, >, >=, nil`
- `expected` Default is an empty string, it is compared with the ``command`` output using the chosen operator in `typeExpected` (optional)
- `description` Is just for taking notes of what is happening (optional)

### Additional JavaScript functions 
- `call('command')` Specify your command type (e.g. ls, grep), additionaly add arguments (e.g. -la, -e). You can also specify a file to directly save it to the artefacts (e.g. §file§pathToFile).
- `callCompare('command', 'string to compare to') bool` Can be used in a JavaScript if statement. If command output == string it return true
- `callContains('command', 'string') bool` Can be used in a JavaScript if statement. If command output contains string return true
- `shell('command without Wrapper')` The command gets directly executed on the shell, without further checking if its valid or harmful to the system. Only use if necessary, try `call()` Instead, because this will deliver better debugging output.
- `regQuery('registry','value')` This is our own function to query the registry on windows.
	- The format has to be: `regQuery('HKLM:\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion', 'ProductID')`
- `printToConsole('string')` Prints the string to your stdout (Good for debugging)
- `printToLog('string')` Prints the string to your log file (Good for taking notes)

### Supported Commands through wrapper
Here is a list of all supported Commands you can use in `Call()`,  `CallCompare()` and `CallContains()`. If the command you want to use is not included you might want to add it through `-add` instead of using `shell()`.

```bash
"Get-MpComputerStatus", "Get-ItemPropertyValue", "Get-NetFirewallProfile", "type", "echo", "reg", "findstr",
"dir", "ls", "cat", "grep", "find", "useradd", "stat", "mount", "systemctl", "egrep", "test", "call", "ps",
"Select-String", "%", "modprobe", "df", "rpm", "zypper", "crontab", "stat", "sysctl", "journalctl", "sestatus",
"apparmor_status", "timedatectl", "ss", "lsof", "iw", "ip", "lsmod", "firewall-cmd", "nmcli", "iptables",
"ip6tables", "nft", "auditctl", "sshd", "useradd", "rmmod", "awk", "xargs", "subscription-manager", "dnf", "authselect"
```

- **Keep in mind, that some commands only work on Windows and others only on Linux . Some are multi-platform but might behave differently.**

### Start a scan
- When starting the tool via command line you have to provide an input file. Everything else is optional.
```bash
./secuteel -input <path/to/config(.json)>
```
## Example usage
- Run an audit on Linux with an custom output file (.zip), verbose output and debugging enabled.
```bash
./secuteel -input <path/to/config(.json)> -output <name> -v -debug
```
- Add new commands for the current audit.
```bash
./secuteel -input <path/to/config(.json)> -add <custom,command,here>
```
- Run the help command on Windows.
```bash
secuteel.exe -h
```
### Config to get started
#### Windows
```json
{
  "system": 
  {
    "systemName": "Windows",
    "version": "10.0.19042"
  },
  "commands": [
    {
      "name": "get_bluetooth_status",
      "command": "regQuery('HKLM:\\SYSTEM\\CurrentControlSet\\Services\\bthserv', 'Start')",
      "typeExpected": "==",
      "expected": "3",
      "description": "Checking bluetooth enabled"
    }
  ]
}
```
#### Linux
```json
{
  "system": 
  {
    "systemName": "Linux",
    "version": "20.04.1-Ubuntu"
  },
  "commands": [
    {
      "name": "get_password_expiration",
      "command": "useradd -D | grep INACTIVE",
      "typeExpected": "==",
      "expected": "INACTIVE=-1",
      "description": "Check for the default time an account gets disabled after password expires"
    }
  ]
}
```

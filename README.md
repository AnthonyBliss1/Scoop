# Scoop
Scoop is a simple REST API client built for testing and discovery.
I built it as a personal tool that doesn't require accounts or logins.
It’s also been a practical way to build experience with TypeScript.

## How to Build
Scoop requires the Wails framework to build

### Set Up Wails 
- Install Wails V3 (currently in alpha)
```bash
go install github.com/wailsapp/wails/v3/cmd/wails@latest
```

- Check system dependencies
```bash
wails3 doctor
```

### Run in Dev Mode
```bash
wails3 dev
```

### Build with Wails
```bash
wails3 build
```

## How to Use
The application data is structured in the following way:
  - `Scoops` are made up of Request and Response objects
  - Each `Scoop` has a name
  - `Collections` are groups of `Scoops`
  - Each `Collection` has a name

### Keybindings
Scoop utilizes various keybindings to help with navigation
Keybindings can also be view by clicking the 'I' (Info) button in the bottom right corner


| Action                | Binding           |
|---------------------- | ----------------- |
| Open Command Palette  | Ctrl + Shift + P  |
| Go to Scoop           | Ctrl + Num        |
| Rename Scoop          | Ctrl + R          |
| Expand Response Body  | Ctrl + E          |


## Scoop Server
Scoop comes with a self-hosted sync server to ensure consistent app data across multiple devices

### Scoop Server Set Up
- Clone the repo
```bash 
git clone https://github.com/AnthonyBliss1/Scoop-Server.git
```

- Build the project
```bash
cd Scoop-Server && go build -o build/scoop-server .
```

- Run the executable
```bash
./build/scoop-server
```

### Scoop Server Flags
Scoop-Server has two flag options


| Flag        | Description                  |
|------------ | ---------------------------  |
| -port=XXXX  | Specify port (default 2767)  |
| -deploy     | Create a systemD service     |


> [!IMPORTANT]
> When running with the `-deploy` flag the executable must be run with `sudo`

### Examples
- Running Scoop Server using the default port `2767`
```bash
./scoop-server
```

- Run Scoop Server using port `8000`
```bash
./scoop-server -port=8000
```

- Deploy Scoop Server as a `systemD` service using port `8000`
```bash
sudo ./scoop-server -port=8000 -deploy
```

> [!NOTE]
>`-deploy` support for MacOS (launchD) coming soon

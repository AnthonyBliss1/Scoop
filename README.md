<p align="center">
  <img src="build/appicon.png" alt="Scoop Icon" width="128" />
</p>

<h1 align="center">Scoop</h1>

<p align="center">
  Scoop is a simple REST API client built for testing and discovery.
  I built it as a personal tool that doesn't require accounts or logins.
  It’s also been a practical way to build experience with Typescript.
</p>

## How to Build
Scoop requires Go and the Wails framework to build

### Set Up Wails 
- Install Wails V3 (currently in alpha)
```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
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
| Open Command Palette  | Shift + Ctrl + P  |
| Go to Scoop           | Ctrl + Num        |
| Rename Collection     | Shift + Ctrl + R  |
| Rename Scoop          | Ctrl + R          |
| Delete Collection     | Shift + Ctrl + D  |
| Delete Scoop          | Ctrl + D          |
| Expand Response Body  | Ctrl + E          |


## Scoop Server
Scoop comes with a self-hosted sync server to ensure consistent app data across multiple devices.
See the Scoop Server repo [here](https://github.com/AnthonyBliss1/Scoop-Server) for setup instructions.

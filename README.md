Leapp Daemon
=========

- Website: https://www.leapp.cloud/
- Roadmap: [Roadmap](https://github.com/Noovolari/leapp/projects/4)
- Wiki: [Wiki](https://github.com/Noovolari/leapp/wiki)
- Chat with us: [Slack](https://join.slack.com/t/noovolari/shared_invite/zt-noc0ju05-18_GRX~Zi6Jz8~95j5CySA)

![logo](.github/images/README-1.png)

Leapp-daemon is the core Business logic of the [Leapp project](https://github.com/Noovolari/leapp).

The daemon is the engine designed to manage and secure cloud access in multi-account environments.

> The project is in active development to replace the current core logic of [Leapp](https://github.com/Noovolari/leapp)

> If you wanted to download Leapp click [here to download the stable version](https://github.com/Noovolari/leapp)

# Contributing

Please read through our [contributing guidelines](.github/CONTRIBUTING.md) and [code of conduct](.github/CODE_OF_CONDUCT.md). We included directions
for opening issues, coding standards, and notes on development.

> 🤙 To jumpstart your contribution to the project we suggest you to come to our [Slack](https://join.slack.com/t/noovolari/shared_invite/zt-noc0ju05-18_GRX~Zi6Jz8~95j5CySA) and have a chat with us; we will point you in the right direction as fast as possible.

Editor preferences are available in the [editor config](.editorconfig) for easy use in common text editors. Read more and download plugins at [editorconfig.org](http://editorconfig.org).

# Developing
Development on leapp-daemon can be done on Mac, Windows, or Linux as long as you have Go installed. See the [go.mod](https://github.com/Noovolari/leapp-daemon/blob/master/go.mod) file located in the project root for the correct Go version.

## Quickstart
- Clone the repository with ```git clone https://github.com/Noovolari/leapp-daemon```
- Change directory into the project root
- Install dependencies with ```go get ./...```
- Check the [good first issues](https://github.com/Noovolari/leapp-daemon/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22)

## Basic functionality
Leapp-daemon is a set of REST APIs wrapped around an http client exposed on port 8080.

<details>
<summary>Main elements and entry-point</summary>
  
- The entry point is [main.go](https://github.com/Noovolari/leapp-daemon/blob/master/main.go) file located in the project root
- The [configuration](https://github.com/Noovolari/leapp-daemon/blob/616470d9e8d668dd067eb63cac2024a2b463f67a/core/configuration/configuration.go) represent the current state of the software
- The [http-engine](https://github.com/Noovolari/leapp-daemon/blob/616470d9e8d668dd067eb63cac2024a2b463f67a/api/engine/engine.go) to respond to API calls
- The websocket for enabling full-duplex communication against multiple consumers
- The [timer](https://github.com/Noovolari/leapp-daemon/blob/616470d9e8d668dd067eb63cac2024a2b463f67a/core/timer/timer.go) to auto-rotate credentials
</details>

<details>
<summary>Project Structure</summary>

- **api**     - interfaces for interacting with core logic
- **core**    - business logic
- **service** - middleware that serves as communication between api and core
</details>

## Testing
To test business logic you can use any API client like [Insomnia](https://insomnia.rest/) or [Postman](https://www.postman.com/).

# Logs


## Documentation
Here you can find our [documentation](https://github.com/Noovolari/leapp-daemon/wiki).

## Links
- [Roadmap](https://github.com/Noovolari/leapp/projects/4): view our next steps and stay up to date
- [Contributing](./.github/CONTRIBUTING.md): follow the guidelines if you'd like to contribute to the project

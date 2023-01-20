# todoassistant.api

backend api for getticked.com

## Getting Started

To get server up running locally follow these steps:

### Install go

- Install **go >=1.19** by following the steps [here](https://go.dev/doc/install)
- Preferably use [JetBrains GoLand](https://youtu.be/vetAfxQxyJE) and open this project as it simplifies this entire process
  <details>
    <summary>for vscode to solve go module error </summary>

  ```sh
  go work init
  go work use ./email-service ./todoassistant.api ./listener-srv/
  ```

  </details>

- Run `go get ./...` in `email-service/` `listener-srv/cmd/api/` && `todoassistant.api/` directories to install the required go modules
- make sure RabitMQ is up and running on your localhost
- Start the app by running `go run main.go` in `email-service/` `listener-srv/cmd/api/` && `todoassistant.api/` directories in separate terminal.
- Head over to `http://localhost:2022/api/v1` in your browser of choice and confirm that the page has loaded without any problems.
- Voilà! Happy coding! :sparkles:

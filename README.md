# CS Demo Analyzer

A CLI to analyze and export CS2/CS:GO demos.

## Usage

Ready-to-use binaries are available on the [releases page](https://github.com/akiver/cs-demo-analyzer/releases).

### Options

```
csda -help

Usage of csda:
  -demo-path string
        Demo file path (mandatory)
  -format string
        Export format, valid values: [csv,json,csdm] (default "csv")
  -minify
        Minify JSON file, it has effect only when -format is set to json
  -output string
        Output folder or file path, must be a folder when exporting to CSV (mandatory)
  -positions
        Include entities (players, grenades...) positions (default false)
  -source string
        Force demo's source, valid values: [challengermode,ebot,esea,esl,faceit,perfectworld,popflash,valve]
```

### Examples

Export a demo into CSV files in the current folder.

`csda -demo-path=myDemo.dem -output=.`

Export a demo in a specific folder into a minified JSON file including entities positions.

`csda -demo-path=/path/to/myDemo.dem -output=/path/to/folder -format=json -positions -minify`

## API

### GO API

This API exposes functions to analyze/export a demo using the Go language.

#### Analyze

This function analyzes the demo located at the given path and returns a `Match`.

```go
package main

import (
	"fmt"
	"os"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

func main() {
	match, err := api.AnalyzeDemo("./myDemo.dem", api.AnalyzeDemoOptions{
		IncludePositions: true,
		Source:           constants.DemoSourceValve,
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, kill := range match.Kills {
		fmt.Printf("(%d): %s killed %s with %s\n", kill.Tick, kill.KillerName, kill.VictimName, kill.WeaponName)
	}
}
```

#### Analyze and export

This function analyzes and exports a demo into the given output path.

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

func main() {
	exePath, _ := os.Executable()
	outputPath := filepath.Dir(exePath)
	err := api.AnalyzeAndExportDemo("./myDemo.dem", outputPath, api.AnalyzeAndExportDemoOptions{
		IncludePositions: false,
		Source:           constants.DemoSourceValve,
		Format:           constants.ExportFormatJSON,
		MinifyJSON:       true,
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("Demo analyzed and exported in " + outputPath)
}
```

### CLI

This API exposes the command-line interface.

```go
package main

import (
	"os"

	"github.com/akiver/cs-demo-analyzer/pkg/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
```

### Node.js API

A Node.js module called `@akiver/cs-demo-analyzer` is available on NPM.  
It exposes a function that under the hood is a wrapper around the Go CLI.  
The module also exports TypeScript types and constants.

```js
import { analyzeDemo, DemoSource, ExportFormat } from '@akiver/cs-demo-analyzer';

async function main() {
  await analyzeDemo({
    demoPath: './myDemo.dem',
    outputFolderPath: '.',
    format: ExportFormat.JSON,
    source: DemoSource.Valve,
    analyzePositions: false,
    minify: false,
    onStderr: console.error,
    onStdout: console.log,
    onStart: () => {
      console.log('Starting!');
    },
    onEnd: () => {
      console.log('Done!');
    },
  });
}

main();
```

## Developing

### Requirements

- [Go](https://golang.org/dl/)
- [Make](https://www.gnu.org/software/make/)

### Build

#### Windows

`make build-windows`

#### macOS

`make build-darwin` / `make build-darwin-arm64`

#### Linux

`make build-linux` / `make build-linux-arm64`

### Tests

1. `./download-demos.sh` it will download the demos used for the tests
2. `make test`

### VSCode debugger

1. Inside the `.vscode` folder, copy/paste the file `launch.template.json` and name it `launch.json`
2. Place a demo in the `debug` folder
3. Update the `-demo-path` argument to point to the demo you just placed
4. Adjust the other arguments as you wish
5. Start the debugger from VSCode

## Acknowledgements

This project uses the demo parser [demoinfocs-golang](https://github.com/markus-wa/demoinfocs-golang) created by [@markus-wa](https://github.com/markus-wa) and maintained by him and [@akiver](https://github.com/akiver).

## License

[MIT](https://github.com/akiver/cs-demo-analyzer/blob/main/LICENSE.md)

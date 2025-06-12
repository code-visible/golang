# gopher

gopher is a code visible Golang parser.

## Features

- Parses Go projects and outputs structured data.
- Supports custom module and type filtering.
- Excludes specified directories and files.
- Outputs results in JSON format for further processing.

## Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/code-visible/gopher.git
cd gopher
CGO_ENABLED=0 go build -o gopher .
```

## Usage

```bash
gopher \
  --project=kubernetes-1.12.0 \
  --directory=. \
  --dump=dist/kubernets.json \
  --module=k8s.io/kubernetes \
  --excludes=vendor,logo,.* ,tests \
  --types=*.json,*.yaml,README.md
```

### Arguments

- `--project` : Name of the project.
- `--directory` : Root directory to parse.
- `--dump` : Output file path (JSON).
- `--module` : Go module to analyze.
- `--excludes` : Comma-separated list of directories/files to exclude.
- `--types` : Comma-separated list of file types to include.

## Output

The output is a JSON file containing structured information about the parsed Go project, such as modules, types, and dependencies.

## Language Protocol

[code visible protocol definition](https://github.com/code-visible/protocol)

## Contributing

Contributions are welcome! Please open issues or submit pull requests.

## LICENSE

GPL-3.0

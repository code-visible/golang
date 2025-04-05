# gopher

gopher is a code visible Golang parser

## Usage

```bash
gopher \
  --project=kubernetes-1.12.0 \
  --directory=. \
  --dump=dist/kubernets.json \
  --module=k8s.io/kubernetes \
  --excludes=vendor,logo,.\*,tests \
  --types=\*.json,\*.yaml,README.md
```

## Language Protocol

[code visible protocol definition](https://github.com/code-visible/protocol)

## LICENSE

GPL-3.0

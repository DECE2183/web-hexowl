# hexowl web version

The web implementation of the hexowl calculator.

## Build

### Windows

```powershell
$env:GOOS='js'; $env:GOARCH='wasm'; go build -C hexowl -o ../wasm/hexowl.wasm
```

### Unix

```bash
GOOS=js GOARCH=wasm go build -C hexowl -o ../wasm/hexowl.wasm
```
# hexowl web version

The web implementation of the hexowl calculator.

## Build

### Windows

```powershell
cd hexwol
```
```powershell
$env:GOOS='js'; $env:GOARCH='wasm'; go build -o ../wasm/hexowl.wasm
```

### Unix

```bash
cd hexowl
```
```bash
GOOS=js GOARCH=wasm go build -o ../wasm/hexowl.wasm
```
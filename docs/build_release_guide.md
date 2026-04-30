# Media Collector — Build & Release Guide

## Versioning

Use semantic versioning: `vMAJOR.MINOR.PATCH`

## Development Build

```bash
wails dev
```

## Production Build

### Windows (amd64)

```powershell
wails build -platform windows/amd64
```

Output: `build/bin/Media Collector.exe`

### Using the build script

```powershell
.\build\build.ps1
```

## Output Location

All build artifacts are placed in:
```
build/bin/
```

## Testing Checklist

Before releasing, verify:

- [ ] App launches without errors
- [ ] Master folder selection dialog works
- [ ] Media scan discovers all supported file types
- [ ] Collections tree displays correctly
- [ ] Media can be moved between collections
- [ ] Player config creation works
- [ ] Slideshow plays with correct transition timing
- [ ] Video files play for full duration
- [ ] F11 toggles fullscreen
- [ ] ESC exits fullscreen player
- [ ] `go test ./... -cover` passes with ≥80% on core packages

## Packaging

The `wails build` command produces a standalone `.exe` file.
No additional installer is required — distribute the binary directly.

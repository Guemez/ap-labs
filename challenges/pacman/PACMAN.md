# REQUIREMENTS
- Go version 1.10 or upper
```
go get github.com/faiface/pixel
go get github.com/faiface/glhf
go get github.com/go-gl/glfw/v3.2/glfw 
```

* On macOS, you need Xcode or Command Line Tools for Xcode (xcode-select --install) for required headers and libraries.
* On Ubuntu/Debian-like Linux distributions, you need libgl1-mesa-dev and xorg-dev packages.
* On CentOS/Fedora-like Linux distributions, you need libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel packages.

**Pixel go library doesn’t support threads, that’s why we can observe some glitches

# TO RUN
```
go run pacman.go -g <# ghosts>
```

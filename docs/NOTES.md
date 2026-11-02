# NOTES

## Go Runtime & Module Init

### Cek/Instal Go

Cek versi yang terpasang:

```bash
brew install go
```

> Catatan: Pastikan $PATH mengarah ke /usr/local/bin/go (atau /opt/homebrew/bin/go pada Apple Silicon).
> Verifikasi lagi:

which go # → /usr/local/bin/go
go version

### Inisialisasi Go Module

Inisialisasi module dengan nama yang akan digunakan di GitHub:

```bash
go mod init github.com/AankTia/go-fintrack
```

Ini akan membuat file go.mod berisi:

```text
module github.com/AankTia/go-fintrack

go 1.26.4
```

Pastikan module sudah berfungsi, Buat file sederhana untuk test:

```bash
echo 'package main; import "fmt"; func main(){ fmt.Println("Hello Go!") }' > main.go
go run main.go   # Output: Hello Go!
```

Jika berhasil, module sudah siap.

### Verifikasi Build

Install `golangci-lint`:

```bash
brew install golangci-lint
golangci-lint run
```

Run go fmt for formatting

```bash
go fmt ./...
```

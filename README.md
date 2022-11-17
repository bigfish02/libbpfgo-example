# libbpfgo-example

Basic eBPF examples in Golang using [libbpfgo](https://github.com/aquasecurity/tracee/tree/main/libbpfgo). 
* Accompanying [slides from my talk at GOTOpia 2021](https://speakerdeck.com/lizrice/beginners-guide-to-ebpf-programming-with-go) called Beginner's Guide to eBPF Programming in Go
* See also my [original Python examples](https://github.com/lizrice/ebpf-beginners) from my [Beginner's Guide to eBPF talk](https://speakerdeck.com/lizrice/liz-rice-beginners-guide-to-ebpf)  

## Install packages

```sh
sudo apt-get update
sudo apt-get install make clang llvm
```

## Install libbpf

```sh
mkdir /build && \
    git clone --branch v1.0.1 --depth 1 https://github.com/libbpf/libbpf.git /build/libbpf && \
    make -j $(nproc) -C /build/libbpf/src BUILD_STATIC_ONLY=y LIBSUBDIR=lib install
```

You must keep your libbpf-go version same with libbpf.
If you use libbpfgo with version: `github.com/aquasecurity/libbpfgo@v0.4.4-libbpf-1.0.1`, you must install libbpf version with 1.0.1.

## Building and running hello

```sh
make all
sudo ./hello
```

This builds two things:
* dist/hello.bpf.o - an object file for the eBPF program
* hello - a Go executable

The Go executable reads in the object file at runtime. Take a look at the .o file with readelf if you want to see the sections defined in it.

## Docker

To avoid compatibility issues, you can use the `Dockerfile` provided in this repository.

Build it by your own:

```bash
docker build -t hello .
```

And the run it from the project directory to compile the program:

```bash
docker run --rm -v $(pwd)/:/app/:z hello
```

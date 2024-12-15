This needs to be able to find the z3 header and library files to link against them.

On Linux using Homebrew to install z3, you need to run:

```shell
export CGO_CFLAGS=-I/home/linuxbrew/.linuxbrew/include CGO_LDFLAGS=-L/home/linuxbrew/.linuxbrew/lib
export LD_LIBRARY_PATH=/home/linuxbrew/.linuxbrew/lib
```

On other operating systems with Homebrew, the Z3 prefix will be different. If z3 is installed under a more common location, this might not be necessary.
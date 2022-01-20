Install dependencies
```shell
go install github.com/google/pprof@latest
brew install graphviz
```
pprof index page
```shell
open http://localhost:7777/debug/pprof
```
inspect cpu profile
```shell
go tool pprof http://localhost:7777/debug/pprof/profile
```
# Development
## Requirements
- Go: go1.17.6
- Helm: v3.0
- Kind: v0.11.1
- Chart Releaser: (https://github.com/helm/chart-releaser)
- Chart Testing: (https://github.com/helm/chart-testing)
- Kubectl 
   ```txt
    Client Version: version.Info{Major:"1", Minor:"22", GitVersion:"v1.22.5" }
    Server Version: version.Info{Major:"1", Minor:"21", GitVersion:"v1.21.1" }
    ```
## Submitting Changes
Raise a PR with a clear break down of what has been done. The more tests the better. We can always use more test coverage. Please follow our coding conventions (below) and make sure all of your commits are atomic (one feature per commit).

Always write a clear log message for your commits. One-line messages are fine for small changes, but bigger changes should look like this:
```shell
git commit -m "A brief summary of the commit
> 
> A paragraph describing what changed and its impact."
```
## Coding Convention
Start reading our code, and you'll get the hang of it. We try to prioritise readability (not always possible), nevertheless:
- We indent with soft tabs
- Add a single blank line at the end of the file
- Try to split functionality into distinct, independent  functions to ease the logical structure of the project and make it easier to understand
- This is open source software. Consider the people who will read your code. The more readable the better.
- We advise that you perform code scanning with Snyk or Sonar Lint to help fix some code issues.

## Linting
### Go
We simply use `go fmt`. Simple is better and "60% of the time, it works every time." :)

### Helm Charts
We use both https://github.com/helm/chart-testing and `helm lint`. Again, simple is better.

## Testing
We have a number of Ginkgo integration and unit tests under `tests` split based on functionality. Please write tests,if possible, for any new code.

## Debugging
Install dependencies
```shell
go install github.com/google/pprof@latest
brew install graphviz
```
Run the service locally
```shell
ARDOQ_BASEURI=https://{custom_domain_here}/api/  ARDOQ_ORG={org_label_here} ARDOQ_WORKSPACE_ID={workspace_id_here} ARDOQ_APIKEY={key_here} ARDOQ_CLUSTER=local make run
```
pprof index page
```shell
open http://localhost:7777/debug/pprof
```
Inspect cpu profile
```shell
go tool pprof http://localhost:7777/debug/pprof/profile
```



# Code References
- https://github.com/AlexsJones/KubeOps/blob/main/DEVELOPMENT.md

Remember D.R.Y. and K.I.S.S

`#StaySmartlyLazy`
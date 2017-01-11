# git-sync

Allows to synchronize multiple GIT repositories remotes in parallel.
The GIT repositories must have "origin" set to your "fork" and "upstream" remote pointing
to the original GIT repository you forked your from.

## Install

```
go install github.com/mfojtik/git-sync
```

## Usage

### git sync add [directory]

Add the GIT repository to tracking.

### git sync list

List repositories we track.

### git sync run

Perform the synchronization:

```console
openshift/origin [====================================>--------------------------------------------------------------------------------------------------------------------------------------------------]  20.00% 0s
containers/image [=============================================================================================================>-------------------------------------------------------------------------]  60.00% 1s
...
```

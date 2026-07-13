Small wrapper program around git commands to easily read source code, see codes evolve commit by commit, from start to end.

Usage

```
commit2commit start
```

Functions

```
start
next
prev
goto <commit>
tree
show
diff
status
quit/exit
```

Phase 1:
- `start`
- `next/n` - go to next commit
- `prev/p` - go to previous commit
- `goto <commit>` - go to specific commit
- `tree` - show previous 5 and next 5 commits
- `show` - show the current commit's metadata
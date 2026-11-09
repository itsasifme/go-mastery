# Files and Folders

> 📝 **Planned.** No code yet. This is what I plan to add here.

Read, write, and walk through files with `os` and `path/filepath`.

## Examples to add

- `01-read-write`: `os.ReadFile` and `os.WriteFile`
- `02-open-close`: open a file and close it with `defer`
- `03-walk`: go through a folder with `filepath.WalkDir`
- `04-temp-files`: `os.CreateTemp` and `os.MkdirTemp`

## Things to remember

- Use `path/filepath` for file paths. It handles `\` on Windows.
- When writing a file, **check the error from `Close` too**.

## Go version notes

- `io/ioutil` is old since Go 1.16. Use `os` and `io` instead.

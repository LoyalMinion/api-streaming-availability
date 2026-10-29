---
name: shape-code
description: >-
  Shapes functions, names, checks, and errors when writing or reviewing code.
  Use when implementing a change, splitting a function, naming a helper,
  returning an error, or reviewing a pull request for readability.
---

# Shape code

Apply these rules when writing or changing code. They are about shape and errors, not about one feature. The examples show one case of each rule; apply the rule to any code with the same shape.

Prefer the simpler version. A shorter function, a direct check, and one name for one thing beat a helper, a branch, or an abstraction that does not remove a real step.

## Functions

- A function that sequences several steps only calls the next step. Give a step its own function when it returns a value or makes a check of its own, named after that value or check. A one-off that is a single call stays inline. For example: `parseQuestion`, `readHeader`.
- Define a helper immediately after the function that calls it.
- Define a private type that only supports the code above it, such as an interface or a key struct, at the bottom of the file.
- Do not pay for input the function does not use. Read up to the part this function needs and stop; loading or scanning the rest costs time and memory for nothing. For example: read one CSV record to get the author instead of `ReadAll` on the whole file.
- Keep the types for one feature together in one file. Add a type in the change that first needs it.
- Do not hand-write what the standard library does in one call. For example: `strings.CutPrefix` checks and strips a prefix, where a manual slice needs a length check and an index.

## Names

- Name a variable for what it holds, in one word where one word says it: `question`, `answer`. Try the word first. Fall back to a letter only when no single word fits and the name is read within a few lines of where it is set (`i` in a short loop, `r` for a reader used twice). A name used across a function or beyond always gets a word. A method receiver is one letter, the type's initial, however long the method is. Do not pad a name with its type or with context the reader already has. Name a function for what it returns (a noun, no `Get`) or what it does (a verb), and do not repeat the package, receiver, or parameter names in it. See the [Go style best practices on naming](https://google.github.io/styleguide/go/best-practices.html#naming).
- Name a value after what it means to the reader, not how it was computed. For example: `sourceLine` for the line a record started on.
- Name a dependency after its role so the call site reads like plain English. For example: `q.csvParser.ParseAuthor(contents)` says what happens; `q.csv.ParseAuthor(contents)` names a format, not a collaborator.
- Group related constants in their own `const` block. A value shared by several groups gets its own block.
- Avoid a package-level variable when the value can live closer to its use. A constant is still fine. For example: `var errInvalidCSV` becomes a new error at the check that failed.
- Prefer a constant when a value is needed in more than one place. A value used once stays inline, next to the comparison that uses it.
- Use a standard library or project constant when one matches. Write one by hand only when none fits.
- Export only what another package uses.
- Do not add code to control output that no consumer depends on. For example: a custom JSON marshaler only to fix key order.

## Checks

- Check a value as soon as it exists. Do not compute the next value when this one already fails. For example:

```go
question = strings.TrimSpace(row[0])
if question == "" {
    return "", "", errors.New("question is required")
}

answer = strings.TrimSpace(row[1])
```

- Exit early. Return on the failure as soon as it is known. Do not nest the rest of the function under a success branch, and do not carry a bad value further so a later check can reject it.
- When neither branch is a failure, return on the short one. Test for the branch that takes a line or two, return, and let the long branch run at the top level instead of inside the `if`. For example:

```go
if len(tt.Want.Error) == 0 {
    assert.Empty(t, errs)

    return
}

require.Len(t, errs, len(tt.Want.Error))

for i, want := range tt.Want.Error {
    require.EqualError(t, errs[i], want)
}
```

- Guard nil first, then take the value. Do not start from a zero value and overwrite it when the input is present. For example:

```go
if ip == nil || ip.FileName == nil {
    return "", errors.New("invalid filename")
}

filename := *ip.FileName
```

- One `if` per failure. Do not combine two different failures in one condition.
- Keep error messages unique. When two checks would return the same message, find out why: if one check is already covered by the other, it is not needed; if both are needed, give each its own message.
- Be consistent: the same rule uses the same comparison everywhere. For example: if a record must have exactly two fields, every record with that rule checks `!= 2`; do not check `< 2` in one place and `!= 2` in another.
- Replace several guards on the same value with one check that covers them.
- Delete a check that an earlier check makes impossible or a later call rejects anyway. For example: `time.ParseInLocation` fails on a date of the wrong length, so a length check before it adds nothing.
- Before adding or removing a check, confirm what the library already does. Run the call on the inputs in question when the documentation does not say.
- When a check is removed, the error for that input becomes the one the remaining call returns. Update the test to that error rather than keeping the old check for its message.

## Errors

- Wrap an error that crosses a boundary with what this function was doing: an error from another package, or one returned through an interface. Return an error this function creates unwrapped. For example:

```go
record, err := reader.Read()
if err != nil {
    return nil, errors.Wrap(err, "failed to read header")
}
```

- Keep the underlying error. Do not replace it with a new error that drops the cause.
- Create a separate error for each failure. Do not share one error value, such as a package-level variable, across different failures.

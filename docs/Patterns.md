# Patterns

This is a collection of ideas for handling patterns in the IR for optimization.

## Ternary Ops

There are a lot of Go shapes that could be converted into ternary operators
to simplify the efforts. Here are several patterns that should be detected
if possible and replaced with a ternary.

```Go
if Comp {
    return X
}
return Y
```

```Go
if Comp {
    return X
} else {
    return Y
}
```

```Go
V = X
if Comp {
    V = Y
}
```

```Go
if Comp {
    V = X
} else {
    V = Y
}
```

MultiError
==========

A Go language package with an abstraction for errors which represent multiple child errors.

MultiError is an interface which includes `error`, and adds two additional methods:

```go
Add(error) MultiError
Errors() []error
```

Think of `Add` like `append`, it can add another error into those being wrapped by a `MultiError` instance.

`Errors` gives back the list of contained errors. But the `MultiError` implementation offered by this package also implements `As` and `Is`, so for most things you might want to do with the contained errors, reach for `errors.Is` or `errors.As` first.

There is only one way to create instances of this package's `MultiError` implementation, the function `Wrap(...error) MultiError`.


Group
-----

Because the most common source of multiple errors arising is concurrent operations which can each produce errors, the multierror package also provides a concurrency group structure which supports starting goroutines and collecting the errors they produce.

The zero value of `Group` is usable immediately, and its method `(*Group) Go(func() error)` kicks off goroutines. The errors returned by the given functions will be grouped into a `MultiError` which is conveniently given back by the method `(*Group) Wait() MultiError`. The group also manages a `sync.WaitGroup` so the `Wait` method will block until all goroutines have completed.

As a convenience, `Group` also supports a method `(*Group) Error() MultiError` which produces the complete `MultiError` again. This is also safe for concurrent use, so you could use this method to pull a `MultiError` wrapping all the errors *so far* when only a subset of goroutines have completed.

Finally, by default a `Group` will capture panics by the goroutines it starts, turn them into errors, and include them in the `MultiError`. To change this behavior call `(*Group) RecoverPanics(bool)` and pass it `false` - before starting any goroutines.

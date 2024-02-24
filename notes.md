Maybe one more thing about it. The gocui.ErrUnknownView is probably the design decision here, where you can always run SetView and just check and decide accordingly what should be done.

1. `err == gocui.ErrUnknownView` - newly created View, you can populate it with the "default" content
2. `err == nil` - this is existing view and was updated with new dimension (or the same dimension), it may contain some content already
3. `err != nil` - something is wrong, maybe incorrectly specified dimension or name


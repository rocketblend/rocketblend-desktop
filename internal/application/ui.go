package application

type (
	UI interface {
		Emit(event string, data interface{})

		// On returns a function that removes the handler.
		On(event string, handler func(data ...interface{})) func()

		Once(event string, handler func(data ...interface{})) func()

		OpenFileDialog(opts OpenDialogOptions) (string, error)
		OpenDirectoryDialog(opts OpenDialogOptions) (string, error)
		SaveFileDialog(opts SaveDialogOptions) (string, error)

		// ShowWindow also unminimises the window.
		ShowWindow()
		Quit()
	}
)

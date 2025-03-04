package main

import (
	"log"
	"os"

	"github.com/gotk3/gotk3/gtk"
)

func main() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	// Create a new top-level window, set its title, and connect the "destroy" event to the gtk.MainQuit function.
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
		os.Exit(1)
	}
	win.SetTitle("GTK Example")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Set the default window size.
	win.SetDefaultSize(800, 600)

	// Create a new label widget to display a message.
	label, err := gtk.LabelNew("Hello, GTK!")
	if err != nil {
		log.Fatal("Unable to create label:", err)
		os.Exit(1)
	}

	// Add the label to the window.
	win.Add(label)

	// Make all widgets within the window visible.
	win.ShowAll()

	// Run the GTK main loop.
	gtk.Main()
}

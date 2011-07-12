package main

import (
	"github.com/mattn/go-xmpp"
	"github.com/mattn/go-gtk/gtk"
	"log"
	"os"
	"strings"
)

func main() {
	gtk.Init(&os.Args)

	window := gtk.Window(gtk.GTK_WINDOW_TOPLEVEL)
	window.SetTitle("GoTalk")
	window.Connect("destroy", func() {
		gtk.MainQuit()
	})
	vbox := gtk.VBox(false, 1)
	scrolledwin := gtk.ScrolledWindow(nil, nil)
	textview := gtk.TextView()
	textview.SetEditable(false)
	textview.SetCursorVisible(false)
	scrolledwin.Add(textview)
	vbox.Add(scrolledwin)

	buffer := textview.GetBuffer()

	entry := gtk.Entry()
	vbox.PackEnd(entry, false, false, 0)

	window.Add(vbox)
	window.SetSizeRequest(300, 400)
	window.ShowAll()

	dialog := gtk.Dialog()
	dialog.SetTitle(window.GetTitle())
	sgroup := gtk.SizeGroup(gtk.GTK_SIZE_GROUP_HORIZONTAL)

	hbox := gtk.HBox(false, 1)
	dialog.GetVBox().Add(hbox)
	label := gtk.Label("username:")
	sgroup.Add(label)
	hbox.Add(label)
	username := gtk.Entry()
	hbox.Add(username)

	hbox = gtk.HBox(false, 1)
	dialog.GetVBox().Add(hbox)
	label = gtk.Label("password:")
	sgroup.Add(label)
	hbox.Add(label)
	password := gtk.Entry()
	password.SetVisibility(false)
	hbox.Add(password)

	dialog.AddButton(gtk.GTK_STOCK_OK, int(gtk.GTK_RESPONSE_OK))
	dialog.AddButton(gtk.GTK_STOCK_CANCEL, int(gtk.GTK_RESPONSE_CANCEL))
	dialog.SetDefaultResponse(int(gtk.GTK_RESPONSE_OK))
	dialog.SetTransientFor(window)
	dialog.ShowAll()
	res := dialog.Run()
	username_ := username.GetText()
	password_ := password.GetText()
	dialog.Destroy()
	if res != int(gtk.GTK_RESPONSE_OK) {
		os.Exit(0)
	}

	talk, err := xmpp.NewClient("talk.google.com:443", username_, password_)
	if err != nil {
		log.Fatal(err)
	}

	entry.Connect("activate", func() {
		text := entry.GetText()
		tokens := strings.SplitN(text, " ", 2)
		if len(tokens) == 2 {
			func() {
				defer recover()
				talk.Send(xmpp.Chat{Remote: tokens[0], Type: "chat", Text: tokens[1]})
				entry.SetText("")
			}()
		}
	})

	go func() {
		for {
			func() {
				defer recover()
				chat, err := talk.Recv()
				if err != nil {
					log.Fatal(err)
				}

				var iter gtk.GtkTextIter
				buffer.GetStartIter(&iter)
				buffer.Insert(&iter, chat.Remote+": "+chat.Text+"\n")
			}()
		}
	}()

	gtk.Main()
}

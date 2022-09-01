// a gui youtube video downloader
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
    "github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)		

func main() {
	gtk.Init(nil)
	builder, err := gtk.BuilderNew()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	builder.AddFromFile("ui.glade")
	obj, err := builder.GetObject("window1")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	window := obj.(*gtk.Window)
	window.Connect("destroy", func() {
		gtk.MainQuit()
	})
	obj, err = builder.GetObject("button1")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	button := obj.(*gtk.Button)
	button.Connect("clicked", func() {
		obj, err = builder.GetObject("entry1")
		if err != nil {
			fmt.Println("Error: ", err)
		}
		entry := obj.(*gtk.Entry)
		url, _ := entry.GetText()
		if strings.Contains(url, "youtube.com") {
			downloadVideo(url)
		} else {
			fmt.Println("Invalid url")
		}
	})
	window.ShowAll()
	gtk.Main()
}	

func downloadVideo(url string) {
	// get the video id
	videoID := strings.Split(url, "v=")[1]
	// get the video title
	resp, err := http.Get("https://www.youtube.com/watch?v=" + videoID)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer resp.Body.Close()
	buf := make([]byte, 1024)
	var title string
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("Error: ", err)
		}
		if n == 0 {
			break
		}
		title += string(buf[:n])
	}
	title = strings.Split(strings.Split(title, "title")[1], "\"")[2]
	// download the video
	resp, err = http.Get("https://www.youtube.com/watch?v=" + videoID)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer resp.Body.Close()
	file, err := os.Create(title + ".mp4")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer file.Close()
	io.Copy(file, resp.Body)
	// open the video
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "start", title+".mp4")
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("xdg-open", title+".mp4")
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("open", title+".mp4")
	}
	cmd.Start()
	// move the video to the videos folder
	homeDir, err := glib.GetUserSpecialDir(glib.USER_DIRECTORY_VIDEOS)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	homeDir = filepath.Join(homeDir, title+".mp4")
	err = os.Rename(title+".mp4", homeDir)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}	

// Language: go
// Path: src\youtube downloader\ui.glade
<?xml version="1.0" encoding="UTF-8"?>
<!-- Generated with glade 3.22.2 -->
<interface>
  <requires lib="gtk+" version="3.20"/>
  <object class="GtkWindow" id="window1">
    <property name="can_focus">False</property>
	<property name="title" translatable="yes">Youtube Downloader</property>
	<property name="default_width">400</property>
	<property name="default_height">200</property>
	<child>
	  <object class="GtkBox" id="box1">
		<property name="can_focus">False</property>
		<property name="orientation">vertical</property>
		<property name="spacing">10</property>
		<child>
		  <object class="GtkEntry" id="entry1">
			<property name="can_focus">True</property>
			<property name="hexpand">True</property>
			<property name="placeholder_text" translatable="yes">Enter the url of the video</property>
		  </object>
		  <packing>
			<property name="expand">False</property>
			<property name="fill">True</property>
			<property name="position">0</property>
		  </packing>
		</child>
		<child>
		  <object class="GtkButton" id="button1">
			<property name="can_focus">True</property>
			<property name="label" translatable="yes">Download</property>
		  </object>
		  <packing>
			<property name="expand">False</property>
			<property name="fill">True</property>
			<property name="position">1</property>
		  </packing>
		</child>
	  </object>
	</child>
  </object>
</interface>


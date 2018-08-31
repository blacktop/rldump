package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"runtime"
	"strings"
	"time"

	"howett.net/plist"
)

// ReadingListItem com.apple.ReadingList Child ReadingList item
type ReadingListItem struct {
	DateAdded   time.Time `plist:"DateAdded"`
	PreviewText string    `plist:"PreviewText"`
}

// ReadingListNonSync struct
type ReadingListNonSync struct {
	SiteName                             string    `plist:"siteName"`
	FetchResult                          int       `plist:"FetchResult"`
	PreviewText                          string    `plist:"PreviewText"`
	DateLastFetched                      time.Time `plist:"DateLastFetched"`
	DidAttemptToFetchIconFromImageURLKey bool      `plist:"didAttemptToFetchIconFromImageUrlKey"`
	Title                                string    `plist:"Title"`
	NeverFetchMetadata                   bool      `plist:"neverFetchMetadata"`
}

// ReadingListChild com.apple.ReadingList Child
type ReadingListChild struct {
	WebBookmarkUUID    string             `plist:"WebBookmarkUUID"`
	WebBookmarkType    string             `plist:"WebBookmarkType"`
	URLString          string             `plist:"URLString"`
	ImageURL           string             `plist:"imageURL"`
	ReadingList        ReadingListItem    `plist:"ReadingList"`
	URIDictionary      map[string]string  `plist:"URIDictionary"`
	ReadingListNonSync ReadingListNonSync `plist:"ReadingListNonSync"`
}

// BookmarkChild Bookmark.plist Child
type BookmarkChild struct {
	WebBookmarkUUID       string             `plist:"WebBookmarkUUID"`
	WebBookmarkType       string             `plist:"WebBookmarkType"`
	Children              []ReadingListChild `plist:"Children"`
	Title                 string             `plist:"Title"`
	URLString             string             `plist:"URLString"`
	WebBookmarkIdentifier string             `plist:"WebBookmarkIdentifier"`
}

// BookmarksPlist Bookmarks.plist structure
type BookmarksPlist struct {
	WebBookmarkUUID        string          `plist:"WebBookmarkUUID"`
	WebBookmarkFileVersion int             `plist:"WebBookmarkFileVersion"`
	Children               []BookmarkChild `plist:"Children"`
	Title                  string          `plist:"Title"`
	WebBookmarkType        string          `plist:"WebBookmarkType"`
}

func main() {

	var bplist BookmarksPlist

	if runtime.GOOS != "darwin" {
		log.Fatal("only runs on macOS")
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	bookmarks, err := os.Open(path.Join(usr.HomeDir, "Library/Safari/Bookmarks.plist"))
	if err != nil {
		log.Fatal(err)
	}
	defer bookmarks.Close()

	outputMarkdown, err := os.Create("Bookmarks.md")
	if err != nil {
		log.Fatal(err)
	}
	defer outputMarkdown.Close()

	w := bufio.NewWriter(outputMarkdown)
	defer w.Flush()

	decoder := plist.NewDecoder(bookmarks)
	err = decoder.Decode(&bplist)
	if err != nil {
		log.Fatal(err)
	}

	// create bookmark backup markdown file
	w.WriteString(fmt.Sprintf("# Bookmarks\n\n"))
	for _, child := range bplist.Children {
		if strings.Contains(child.Title, "com.apple.ReadingList") {
			for _, rChild := range child.Children {
				w.WriteString(fmt.Sprintf("## %s\n\n", rChild.URIDictionary["title"]))
				w.WriteString(fmt.Sprintf("### %s\n\n", rChild.URLString))
				if len(rChild.ReadingList.PreviewText) > 0 {
					w.WriteString(fmt.Sprintf("> %s\n\n", rChild.ReadingList.PreviewText))
				}
			}
		}
	}
}

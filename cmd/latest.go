package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/jaroddev/katana"
	"github.com/jaroddev/katana/chapter"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// latestCmd represents the latest command
var latestCmd = &cobra.Command{
	Use:   "latest",
	Short: "Return the latest updates from mangakatana",
	Run: func(cmd *cobra.Command, args []string) {

		ptr := katana.GetUpdates(katana.Url(1))

		menu := createUpdateSelectionMenu(ptr.Updates)
		index, err := selectUpdate(menu)

		if err != nil {
			os.Exit(1)
		}

		choice := ptr.Updates[index]

		manga := katana.GetManga(choice.Url)

		chapterSelectionMenu := createChapterSelectionMenu(*manga)
		index, err = selectChapter(chapterSelectionMenu)

		if err != nil {
			os.Exit(1)
		}

		choosenChapter := manga.Chapters[index]

		chapterScraper := chapter.DScraper{
			Chapter: &chapter.Chapter{
				Images: make([]string, 0),
			},
		}

		chapterScraper.GetChapters(choosenChapter.Url)

		appFolder, _ := cmd.Flags().GetString("path")

		for _, image := range chapterScraper.Chapter.Images {

			data, err := downloadImage(image)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			mangaPath := MangaBasePath(appFolder, choice.Title, choosenChapter.Name)

			os.MkdirAll(mangaPath, os.ModePerm)

			imgPath := fmt.Sprintf("./%s/%s", mangaPath, path.Base(image))

			err = os.WriteFile(imgPath, data, 0666)

			if err != nil {
				fmt.Printf("image could not be saved %d", err)
				os.Exit(1)
			}

		}

	},
}

func createUpdateSelectionMenu(items []katana.Update) promptui.Select {

	titles := make([]string, 0)

	for _, item := range items {
		titles = append(titles, item.Title)
	}

	// titles = append(titles, "next page")

	return promptui.Select{
		Label: "Which manga would you like to scrap chapters of ?",
		Items: titles,
	}
}

func selectUpdate(menu promptui.Select) (id int, err error) {
	id, title, err := menu.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %d %q\n", id, title)
	return
}

func createChapterSelectionMenu(item katana.Manga) promptui.Select {
	titles := make([]string, 0)

	for _, chapter := range item.Chapters {
		titles = append(titles, chapter.Name)
	}

	return promptui.Select{
		Label: "Which chapter would you like to scrap ?",
		Items: titles,
	}
}

func selectChapter(menu promptui.Select) (id int, err error) {
	id, title, err := menu.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %d %q\n", id, title)
	return
}

func downloadImage(url string) (data []byte, err error) {
	res, err := http.DefaultClient.Get(url)

	if err != nil {
		return
	}

	data, err = ioutil.ReadAll(res.Body)

	if err != nil {
		return
	}

	res.Body.Close()

	return
}

func MangaBasePath(appFolder, manga, chapter string) string {
	return fmt.Sprintf("%s/%s/%s", appFolder, manga, chapter)
}

func init() {
	rootCmd.AddCommand(latestCmd)

	latestCmd.Flags().StringP("path", "p", "pages", "Path where data is saved")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// latestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// latestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

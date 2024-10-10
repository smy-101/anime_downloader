package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/smy-101/anime_downloader/download"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download anime using magnet link or torrent file",
	Run: func(cmd *cobra.Command, args []string) {
		magnet, _ := cmd.Flags().GetString("magnet")
		torrent, _ := cmd.Flags().GetString("torrent")
		outputDir, _ := cmd.Flags().GetString("outputDir")

		if magnet == "" && torrent == "" && outputDir == "" {
			downloadTorrent()
		} else if magnet != "" && torrent != "" {
			fmt.Println("Error: You must specify either a magnet link or a torrent file, but not both.")
			cmd.Usage()
			os.Exit(1)
		}

		fmt.Printf("magnet: %v\n", magnet)
		fmt.Printf("torrent: %v\n", torrent)
		fmt.Printf("outputDir: %v\n", outputDir)

		//todo 下载操作
		if torrent != "" {
			download.AnimeDownloader(torrent, outputDir)
		} else if magnet != "" {
			download.AnimeDownloader(magnet, outputDir)
		}
	},
}

func init() {
	downloadCmd.Flags().StringP("magnet", "m", "", "Magnet link")
	downloadCmd.Flags().StringP("torrent", "t", "", "Path to torrent file")
	downloadCmd.Flags().StringP("outputDir", "o", ".", "Output directory")
}

func downloadTorrent() {
	//获取上次使用的下载目录
	downloadDir := viper.GetString("download_dir")

	//提示用户输入种子文件路径或磁力链接
	prompt := promptui.Prompt{
		Label: "Enter the path to the torrent file or magnet link",
	}

	input, _ := prompt.Run()

	//用户下载路径
	prompt = promptui.Prompt{
		Label:   "Enter the download directory",
		Default: downloadDir,
	}
	downloadDir, _ = prompt.Run()

	//保存下载路径
	viper.Set("download_dir", downloadDir)
	viper.WriteConfig()

	//开始执行下载操作
	fmt.Printf("Downloading from: %s\nOutput directory: %s\n", input, downloadDir)
	download.AnimeDownloader(input, downloadDir)
}

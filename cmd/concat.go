package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

const tempListFile = ".list.tmp"

// concatCmd represents the concat command
var concatCmd = &cobra.Command{
	Use:   "concat",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(2),
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("concat Run() called")
		fmt.Println("args: " + strings.Join(args, " "))

		// create list.txt file
		f, err := os.Create(tempListFile)
		if err != nil {
			panic(err)
		}

		defer f.Close()
		for _, filename := range args {
			_, err := os.Stat(filename)
			if err != nil {
				panic(err)
			}
			f.WriteString("file '" + filename + "'\n")
		}
		f.Sync()
	},
	Run: func(cmd *cobra.Command, args []string) {

		// exec ffmpeg
		ffmpeg, err := exec.LookPath("ffmpeg")
		if err != nil {
			panic(err)
		}
		ffmpegArgs := []string{"ffmpeg", "-f", "concat", "-safe", "0", "-i", tempListFile, "-c", "copy", "out.mp4"}
		env := os.Environ()

		execErr := syscall.Exec(ffmpeg, ffmpegArgs, env)
		if execErr != nil {
			panic(execErr)
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("concat PostRun() called")

		// clean up file
		os.Remove(tempListFile)
	},
}

func init() {
	rootCmd.AddCommand(concatCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// concatCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// concatCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

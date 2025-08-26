package cli

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"bscli/pkg/brightsign"
	"github.com/spf13/cobra"
)

func addFileCommands() {
	fileCmd := &cobra.Command{
		Use:     "file",
		Aliases: []string{"files"},
		Short:   "File management commands",
		Long:    "Commands for managing files on the BrightSign player",
	}

	// List files command
	listCmd := &cobra.Command{
		Use:   "list [path]",
		Aliases: []string{"ls"},
		Short: "List files and directories",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			path := "/storage/sd/"
			if len(args) > 0 {
				path = args[0]
			}

			raw, _ := cmd.Flags().GetBool("raw")
			options := &brightsign.ListOptions{Raw: raw}

			files, err := client.Storage.ListFiles(path, options)
			if err != nil {
				handleError(err)
			}

			if jsonOutput {
				outputJSON(files)
				return
			}

			if len(files) == 0 {
				fmt.Println("No files found")
				return
			}

			// Print in table format
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "TYPE\tNAME\tSIZE\tMODIFIED")
			fmt.Fprintln(w, "----\t----\t----\t--------")
			
			for _, file := range files {
				fileType := "file"
				if file.Type == "directory" {
					fileType = "dir"
				}
				size := formatSize(file.Size)
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", fileType, file.Name, size, file.Modified)
			}
			w.Flush()
		},
	}
	listCmd.Flags().Bool("raw", false, "Return raw directory listing")

	// Upload command
	uploadCmd := &cobra.Command{
		Use:   "upload [local-file] [remote-path]",
		Aliases: []string{"put", "cp"},
		Short: "Upload file to player",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			localPath := args[0]
			remotePath := args[1]

			// Ensure remote path is absolute
			if !strings.HasPrefix(remotePath, "/") {
				remotePath = "/storage/sd/" + remotePath
			}

			// Check if local file exists
			if _, err := os.Stat(localPath); err != nil {
				handleError(fmt.Errorf("local file not found: %s", localPath))
			}

			if !jsonOutput {
				fmt.Printf("Uploading %s to %s...\n", localPath, remotePath)
			}
			
			err = client.Storage.UploadFile(localPath, remotePath)
			if err != nil {
				handleError(err)
			}

			if jsonOutput {
				outputJSON(map[string]interface{}{
					"success": true,
					"action":  "upload",
					"source":  localPath,
					"destination": remotePath,
				})
			} else {
				fmt.Println("Upload complete")
			}
		},
	}

	// Download command
	downloadCmd := &cobra.Command{
		Use:   "download [remote-path] [local-file]",
		Aliases: []string{"get"},
		Short: "Download file from player",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			remotePath := args[0]
			localPath := args[1]

			// Ensure remote path is absolute
			if !strings.HasPrefix(remotePath, "/") {
				remotePath = "/storage/sd/" + remotePath
			}

			if !jsonOutput {
				fmt.Printf("Downloading %s to %s...\n", remotePath, localPath)
			}
			
			err = client.Storage.DownloadFile(remotePath, localPath)
			if err != nil {
				handleError(err)
			}

			if jsonOutput {
				outputJSON(map[string]interface{}{
					"success": true,
					"action":  "download",
					"source":  remotePath,
					"destination": localPath,
				})
			} else {
				fmt.Println("Download complete")
			}
		},
	}

	// Delete command
	deleteCmd := &cobra.Command{
		Use:   "delete [path]",
		Aliases: []string{"rm"},
		Short: "Delete file or directory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			path := args[0]

			// Ensure path is absolute
			if !strings.HasPrefix(path, "/") {
				path = "/storage/sd/" + path
			}

			force, _ := cmd.Flags().GetBool("force")
			if !force {
				fmt.Printf("Delete %s? (y/N): ", path)
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					fmt.Println("Cancelled")
					return
				}
			}

			err = client.Storage.DeleteFile(path)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Deleted %s\n", path)
		},
	}
	deleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation")

	// Rename command
	renameCmd := &cobra.Command{
		Use:   "rename [old-path] [new-name]",
		Aliases: []string{"mv"},
		Short: "Rename a file",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			oldPath := args[0]
			newName := args[1]

			// Ensure path is absolute
			if !strings.HasPrefix(oldPath, "/") {
				oldPath = "/storage/sd/" + oldPath
			}

			err = client.Storage.RenameFile(oldPath, newName)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Renamed to %s\n", newName)
		},
	}

	// Create directory command
	mkdirCmd := &cobra.Command{
		Use:   "mkdir [path]",
		Short: "Create directory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			path := args[0]

			// Ensure path is absolute
			if !strings.HasPrefix(path, "/") {
				path = "/storage/sd/" + path
			}

			err = client.Storage.CreateDirectory(path)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Created directory %s\n", path)
		},
	}

	// Format storage command
	formatCmd := &cobra.Command{
		Use:   "format [device]",
		Short: "Format storage device (requires autorun disabled)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			device := args[0]

			force, _ := cmd.Flags().GetBool("force")
			if !force {
				fmt.Printf("WARNING: This will format %s and delete all data. Continue? (y/N): ", device)
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					fmt.Println("Cancelled")
					return
				}
			}

			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			err = client.Storage.FormatStorage(device)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Formatted %s\n", device)
		},
	}
	formatCmd.Flags().BoolP("force", "f", false, "Skip confirmation")

	fileCmd.AddCommand(listCmd, uploadCmd, downloadCmd, deleteCmd, renameCmd, mkdirCmd, formatCmd)
	rootCmd.AddCommand(fileCmd)
}

// formatSize formats bytes into human-readable size
func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
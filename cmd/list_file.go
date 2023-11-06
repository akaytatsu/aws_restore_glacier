package cmd

import (
	"app/aws"
	"errors"

	"github.com/spf13/cobra"
)

// codeCmd represents the code command
var listCmd = &cobra.Command{
	Use:                   "list",
	Short:                 "Restore milions files of a bucket in Glacier",
	Long:                  `This command will help to restore milions of files of a bucket in Glacier.`,
	DisableFlagsInUseLine: false,
	//Args:                  cobra.ExactArgs(1),
	Args: func(cmd *cobra.Command, args []string) error {

		if BucketName == "" {
			return errors.New("bucket name is required")
		}

		if Region == "" {
			return errors.New("region name is required")
		}

		if AccessKey == "" || SecretKey == "" {
			return errors.New("access key and secret key are required")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		proccess := aws.AWSOperations{
			BucketName: BucketName,
			Region:     Region,
			PartialKey: PartialPatch,
		}

		proccess.Init(AccessKey, SecretKey)
		proccess.ListObjects(FilePath)
	},
}

func init_list() {
	listCmd.PersistentFlags().StringVar(&BucketName, "bucket", "", "Bucket name")
	listCmd.PersistentFlags().StringVar(&Region, "region", "", "Region name")
	listCmd.PersistentFlags().StringVar(&FilePath, "file", "", "File path")
	listCmd.PersistentFlags().StringVar(&PartialPatch, "partial", "", "Partial patch")
	listCmd.PersistentFlags().StringVar(&AccessKey, "access_key", "", "Access key")
	listCmd.PersistentFlags().StringVar(&SecretKey, "secret_key", "", "Secret key")

	RootCmd.AddCommand(listCmd)
}

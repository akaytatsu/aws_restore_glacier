package cmd

import (
	"app/aws"
	"errors"

	"github.com/spf13/cobra"
)

// codeCmd represents the code command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore milions files of a bucket in Glacier",
	Long:  `This command will help to restore milions of files of a bucket in Glacier.`,
	// DisableFlagsInUseLine: true,
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
		objects, _ := proccess.ListObjects(FilePath)

		proccess.RestoreObjects(&objects)
	},
}

func init_restore() {
	restoreCmd.PersistentFlags().StringVar(&BucketName, "bucket", "", "Bucket name")
	restoreCmd.PersistentFlags().StringVar(&Region, "region", "", "Region name")
	restoreCmd.PersistentFlags().StringVar(&PartialPatch, "partial", "", "Partial patch")
	restoreCmd.PersistentFlags().StringVar(&AccessKey, "access_key", "", "Access key")
	restoreCmd.PersistentFlags().StringVar(&SecretKey, "secret_key", "", "Secret key")

	RootCmd.AddCommand(restoreCmd)
}

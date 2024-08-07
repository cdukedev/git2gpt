// cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/chand1012/git2gpt/prompt"
	"github.com/spf13/cobra"
)

var repoPath string
var preambleFile string
var outputFile string
var estimateTokens bool
var ignoreFilePath string
var selectFilePath string
var ignoreGitignore bool
var outputJSON bool
var debug bool
var scrubComments bool

var rootCmd = &cobra.Command{
	Use:   "git2gpt [flags] /path/to/git/repository",
	Short: "git2gpt is a utility to convert a Git repository to a text file for input into GPT-4",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoPath = args[0]
		ignoreList := prompt.GenerateIgnoreList(repoPath, ignoreFilePath, !ignoreGitignore)
		selectList := prompt.GenerateSelectList(repoPath, selectFilePath)
		repo, err := prompt.ProcessGitRepo(repoPath, ignoreList, selectList)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		if outputJSON {
			output, err := prompt.MarshalRepo(repo, scrubComments)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)
			}
			if outputFile != "" {
				// if output file exists, throw error
				if _, err := os.Stat(outputFile); err == nil {
					fmt.Printf("Error: output file %s already exists\n", outputFile)
					os.Exit(1)
				}
				err = os.WriteFile(outputFile, []byte(output), 0644)
				if err != nil {
					fmt.Printf("Error: could not write to output file %s\n", outputFile)
					os.Exit(1)
				}
			} else {
				if !debug {
					fmt.Println(string(output))
				}
			}
			return
		}
		output, err := prompt.OutputGitRepo(repo, preambleFile, scrubComments)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		if outputFile != "" {
			// if output file exists, throw error
			if _, err := os.Stat(outputFile); err == nil {
				fmt.Printf("Error: output file %s already exists\n", outputFile)
				os.Exit(1)
			}
			err = os.WriteFile(outputFile, []byte(output), 0644)
			if err != nil {
				fmt.Printf("Error: could not write to output file %s\n", outputFile)
				os.Exit(1)
			}
		} else {
			if !debug {
				fmt.Println(output)
			}
		}
		if estimateTokens {
			fmt.Printf("Estimated number of tokens: %d\n", prompt.EstimateTokens(output))
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&preambleFile, "preamble", "p", "", "path to preamble text file")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "path to output file")
	rootCmd.Flags().BoolVarP(&estimateTokens, "estimate", "e", false, "estimate the number of tokens in the output")
	rootCmd.Flags().StringVarP(&ignoreFilePath, "ignore", "i", "", "path to .gptignore file")
	rootCmd.Flags().StringVarP(&selectFilePath, "select", "s", "", "path to .gptselect file")
	rootCmd.Flags().BoolVarP(&ignoreGitignore, "ignore-gitignore", "g", false, "ignore .gitignore file")
	rootCmd.Flags().BoolVarP(&outputJSON, "json", "j", false, "output JSON")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "debug mode. Do not output to standard output")
	rootCmd.Flags().BoolVarP(&scrubComments, "scrub-comments", "c", false, "scrub comments from the output. Decreases token count")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

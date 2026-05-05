package cmd

import (
	"fmt"
	"strings"

	appconfig "github.com/igorrochap/commitgen/internal/config"
	"github.com/igorrochap/commitgen/internal/prompts"
	"github.com/spf13/cobra"
)

var (
	configLanguage string
	configModel    string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage commitgen defaults",
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set default language or model",
	RunE:  runConfigSet,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current defaults",
	RunE:  runConfigShow,
}

func init() {
	configSetCmd.Flags().StringVar(&configLanguage, "language", "", "Default commit language")
	configSetCmd.Flags().StringVar(&configModel, "model", "", "Default Ollama model")

	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	languageChanged := cmd.Flags().Changed("language")
	modelChanged := cmd.Flags().Changed("model")
	if !languageChanged && !modelChanged {
		return fmt.Errorf("provide --language, --model, or both")
	}

	path, err := configPath()
	if err != nil {
		return err
	}
	cfg, err := appconfig.LoadFile(path)
	if err != nil {
		return err
	}

	if languageChanged {
		if strings.TrimSpace(configLanguage) == "" {
			return fmt.Errorf("language cannot be empty")
		}
		if !prompts.IsSupported(configLanguage) {
			return fmt.Errorf("language %s not supported", configLanguage)
		}
		cfg.Language = configLanguage
	}

	if modelChanged {
		configModel = strings.TrimSpace(configModel)
		if configModel == "" {
			return fmt.Errorf("model cannot be empty")
		}
		cfg.Model = configModel
	}

	if err := appconfig.SaveFile(path, cfg); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Defaults updated: language=%s model=%s\n", cfg.Language, cfg.Model)
	return nil
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	cfg, err := appconfig.LoadFile(path)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "language=%s\nmodel=%s\nconfig=%s\n", cfg.Language, cfg.Model, path)
	return nil
}

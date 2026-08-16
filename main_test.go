package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

// TestLoggerLevels tests that log levels are properly filtered
func TestLoggerLevels(t *testing.T) {
	tests := []struct {
		name          string
		logLevel      string
		expectedLevel LogLevel
		message       string
		logFunc       func(*Logger, string, ...interface{})
		shouldLog     bool
	}{
		// Debug level tests
		{"DebugLevel_DebugMessage", "debug", LogLevelDebug, "debug message", (*Logger).Debug, true},
		{"DebugLevel_InfoMessage", "debug", LogLevelDebug, "info message", (*Logger).Info, true},
		{"DebugLevel_WarnMessage", "debug", LogLevelDebug, "warn message", (*Logger).Warn, true},
		{"DebugLevel_ErrorMessage", "debug", LogLevelDebug, "error message", (*Logger).Error, true},

		// Info level tests
		{"InfoLevel_DebugMessage", "info", LogLevelInfo, "debug message", (*Logger).Debug, false},
		{"InfoLevel_InfoMessage", "info", LogLevelInfo, "info message", (*Logger).Info, true},
		{"InfoLevel_WarnMessage", "info", LogLevelInfo, "warn message", (*Logger).Warn, true},
		{"InfoLevel_ErrorMessage", "info", LogLevelInfo, "error message", (*Logger).Error, true},

		// Warn level tests
		{"WarnLevel_DebugMessage", "warn", LogLevelWarn, "debug message", (*Logger).Debug, false},
		{"WarnLevel_InfoMessage", "warn", LogLevelWarn, "info message", (*Logger).Info, false},
		{"WarnLevel_WarnMessage", "warn", LogLevelWarn, "warn message", (*Logger).Warn, true},
		{"WarnLevel_ErrorMessage", "warn", LogLevelWarn, "error message", (*Logger).Error, true},

		// Error level tests
		{"ErrorLevel_DebugMessage", "error", LogLevelError, "debug message", (*Logger).Debug, false},
		{"ErrorLevel_InfoMessage", "error", LogLevelError, "info message", (*Logger).Info, false},
		{"ErrorLevel_WarnMessage", "error", LogLevelError, "warn message", (*Logger).Warn, false},
		{"ErrorLevel_ErrorMessage", "error", LogLevelError, "error message", (*Logger).Error, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			logger := NewLogger(tt.logLevel)

			// Verify the logger's level was set correctly
			if logger.level != tt.expectedLevel {
				t.Errorf("Expected log level %v, got %v", tt.expectedLevel, logger.level)
			}

			// Call the logging function
			tt.logFunc(logger, tt.message)

			output := buf.String()

			// Check if the message was logged
			if tt.shouldLog {
				if !strings.Contains(output, tt.message) {
					t.Errorf("Expected log output to contain '%s', but it didn't. Output: %s", tt.message, output)
				}
			} else {
				if strings.Contains(output, tt.message) {
					t.Errorf("Expected log output to NOT contain '%s', but it did. Output: %s", tt.message, output)
				}
			}
		})
	}
}

// TestNewLoggerDefaultLevel tests that the logger defaults to info level for invalid input
func TestNewLoggerDefaultLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected LogLevel
	}{
		{"debug", LogLevelDebug},
		{"info", LogLevelInfo},
		{"warn", LogLevelWarn},
		{"warning", LogLevelWarn},
		{"error", LogLevelError},
		{"DEBUG", LogLevelDebug},
		{"INFO", LogLevelInfo},
		{"invalid", LogLevelInfo}, // Should default to info
		{"", LogLevelInfo},        // Should default to info
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			logger := NewLogger(tt.input)
			if logger.level != tt.expected {
				t.Errorf("NewLogger(%q) = %v, want %v", tt.input, logger.level, tt.expected)
			}
		})
	}
}

// TestIsValidRepoName tests the repository name validation
func TestIsValidRepoName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Valid_AlphanumericHyphen", "my-awesome-repo", true},
		{"Valid_AlphanumericUnderscore", "my_awesome_repo", true},
		{"Valid_AlphanumericDot", "my.awesome.repo", true},
		{"Valid_Mixed", "My-Repo_2.0", true},
		{"Invalid_Space", "my repo", false},
		{"Invalid_SpecialChar", "my@repo", false},
		{"Invalid_Empty", "", false},
		{"Invalid_TooLong", strings.Repeat("a", 101), false},
		{"Valid_MaxLength", strings.Repeat("a", 100), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidRepoName(tt.input)
			if got != tt.want {
				t.Errorf("isValidRepoName(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestExtractViewValues tests the extraction of values from view submissions
func TestExtractViewValues(t *testing.T) {
	tests := []struct {
		name     string
		payload  ViewSubmissionPayload
		expected map[string]string
	}{
		{
			name: "TextInputOnly",
			payload: ViewSubmissionPayload{
				View: struct {
					CallbackID string `json:"callback_id"`
					State      struct {
						Values map[string]map[string]struct {
							Type            string `json:"type"`
							Value           string `json:"value"`
							SelectedOptions []struct {
								Value string `json:"value"`
							} `json:"selected_options"`
						} `json:"values"`
					} `json:"state"`
				}{
					State: struct {
						Values map[string]map[string]struct {
							Type            string `json:"type"`
							Value           string `json:"value"`
							SelectedOptions []struct {
								Value string `json:"value"`
							} `json:"selected_options"`
						} `json:"values"`
					}{
						Values: map[string]map[string]struct {
							Type            string `json:"type"`
							Value           string `json:"value"`
							SelectedOptions []struct {
								Value string `json:"value"`
							} `json:"selected_options"`
						}{
							"repo-name": {
								"repo_name_input": {
									Type:  "plain_text_input",
									Value: "test-repo",
								},
							},
						},
					},
				},
			},
			expected: map[string]string{
				"repo-name": "test-repo",
			},
		},
		{
			name: "CheckboxSelected",
			payload: ViewSubmissionPayload{
				View: struct {
					CallbackID string `json:"callback_id"`
					State      struct {
						Values map[string]map[string]struct {
							Type            string `json:"type"`
							Value           string `json:"value"`
							SelectedOptions []struct {
								Value string `json:"value"`
							} `json:"selected_options"`
						} `json:"values"`
					} `json:"state"`
				}{
					State: struct {
						Values map[string]map[string]struct {
							Type            string `json:"type"`
							Value           string `json:"value"`
							SelectedOptions []struct {
								Value string `json:"value"`
							} `json:"selected_options"`
						} `json:"values"`
					}{
						Values: map[string]map[string]struct {
							Type            string `json:"type"`
							Value           string `json:"value"`
							SelectedOptions []struct {
								Value string `json:"value"`
							} `json:"selected_options"`
						}{
							"vibeops-config": {
								"vibeops_config_checkbox": {
									Type: "checkboxes",
									SelectedOptions: []struct {
										Value string `json:"value"`
									}{
										{Value: "create_vibeops_config"},
									},
								},
							},
						},
					},
				},
			},
			expected: map[string]string{
				"vibeops-config": "true",
			},
		},
		{
			name: "CheckboxNotSelected",
			payload: ViewSubmissionPayload{
				View: struct {
					CallbackID string `json:"callback_id"`
					State      struct {
						Values map[string]map[string]struct {
							Type            string `json:"type"`
							Value           string `json:"value"`
							SelectedOptions []struct {
								Value string `json:"value"`
							} `json:"selected_options"`
						} `json:"values"`
					} `json:"state"`
				}{
					State: struct {
						Values map[string]map[string]struct {
							Type            string `json:"type"`
							Value           string `json:"value"`
							SelectedOptions []struct {
								Value string `json:"value"`
							} `json:"selected_options"`
						} `json:"values"`
					}{
						Values: map[string]map[string]struct {
							Type            string `json:"type"`
							Value           string `json:"value"`
							SelectedOptions []struct {
								Value string `json:"value"`
							} `json:"selected_options"`
						}{
							"vibeops-config": {
								"vibeops_config_checkbox": {
									Type: "checkboxes",
									SelectedOptions: []struct {
										Value string `json:"value"`
									}{},
								},
							},
						},
					},
				},
			},
			expected: map[string]string{
				"vibeops-config": "false",
			},
		},
		{
			name: "MixedTextAndCheckbox",
			payload: ViewSubmissionPayload{
				View: struct {
					CallbackID string `json:"callback_id"`
					State      struct {
						Values map[string]map[string]struct {
							Type            string `json:"type"`
							Value           string `json:"value"`
							SelectedOptions []struct {
								Value string `json:"value"`
							} `json:"selected_options"`
						} `json:"values"`
					} `json:"state"`
				}{
					State: struct {
						Values map[string]map[string]struct {
							Type            string `json:"type"`
							Value           string `json:"value"`
							SelectedOptions []struct {
								Value string `json:"value"`
							} `json:"selected_options"`
						} `json:"values"`
					}{
						Values: map[string]map[string]struct {
							Type            string `json:"type"`
							Value           string `json:"value"`
							SelectedOptions []struct {
								Value string `json:"value"`
							} `json:"selected_options"`
						}{
							"repo-name": {
								"repo_name_input": {
									Type:  "plain_text_input",
									Value: "my-repo",
								},
							},
							"vibeops-config": {
								"vibeops_config_checkbox": {
									Type: "checkboxes",
									SelectedOptions: []struct {
										Value string `json:"value"`
									}{
										{Value: "create_vibeops_config"},
									},
								},
							},
						},
					},
				},
			},
			expected: map[string]string{
				"repo-name":      "my-repo",
				"vibeops-config": "true",
			},
		},
		{
			name: "WithAIPrompt",
			payload: ViewSubmissionPayload{
				View: struct {
					CallbackID string `json:"callback_id"`
					State      struct {
						Values map[string]map[string]struct {
							Type            string `json:"type"`
							Value           string `json:"value"`
							SelectedOptions []struct {
								Value string `json:"value"`
							} `json:"selected_options"`
						} `json:"values"`
					} `json:"state"`
				}{
					State: struct {
						Values map[string]map[string]struct {
							Type            string `json:"type"`
							Value           string `json:"value"`
							SelectedOptions []struct {
								Value string `json:"value"`
							} `json:"selected_options"`
						} `json:"values"`
					}{
						Values: map[string]map[string]struct {
							Type            string `json:"type"`
							Value           string `json:"value"`
							SelectedOptions []struct {
								Value string `json:"value"`
							} `json:"selected_options"`
						}{
							"repo-name": {
								"repo_name_input": {
									Type:  "plain_text_input",
									Value: "my-repo",
								},
							},
							"ai-prompt": {
								"ai_prompt_input": {
									Type:  "plain_text_input",
									Value: "Test AI Prompt",
								},
							},
						},
					},
				},
			},
			expected: map[string]string{
				"repo-name": "my-repo",
				"ai-prompt": "Test AI Prompt",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractViewValues(tt.payload)
			if len(result) != len(tt.expected) {
				t.Errorf("extractViewValues() returned %d values, expected %d", len(result), len(tt.expected))
			}
			for key, expectedValue := range tt.expected {
				if gotValue, ok := result[key]; !ok {
					t.Errorf("extractViewValues() missing key %q", key)
				} else if gotValue != expectedValue {
					t.Errorf("extractViewValues()[%q] = %q, want %q", key, gotValue, expectedValue)
				}
			}
		})
	}
}

// TestLoadConfigWithGithubPrivateWorkingDir tests that the GITHUB_PRIVATE_WORKING_DIR is correctly loaded
func TestLoadConfigWithGithubPrivateWorkingDir(t *testing.T) {
	// Set up test environment
	os.Setenv("SLACK_BOT_TOKEN", "test-token")
	os.Setenv("GITHUB_ORG", "test-org")
	os.Setenv("GITHUB_PRIVATE_WORKING_DIR", "/test/github-private")
	defer func() {
		os.Unsetenv("SLACK_BOT_TOKEN")
		os.Unsetenv("GITHUB_ORG")
		os.Unsetenv("GITHUB_PRIVATE_WORKING_DIR")
	}()

	config, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() failed: %v", err)
	}

	if config.GithubPrivateWorkingDir != "/test/github-private" {
		t.Errorf("GithubPrivateWorkingDir = %q, want %q", config.GithubPrivateWorkingDir, "/test/github-private")
	}
}

// TestLoadConfigWithoutGithubPrivateWorkingDir tests that the config loads with empty GITHUB_PRIVATE_WORKING_DIR
func TestLoadConfigWithoutGithubPrivateWorkingDir(t *testing.T) {
	// Set up test environment
	os.Setenv("SLACK_BOT_TOKEN", "test-token")
	os.Setenv("GITHUB_ORG", "test-org")
	defer func() {
		os.Unsetenv("SLACK_BOT_TOKEN")
		os.Unsetenv("GITHUB_ORG")
	}()

	config, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() failed: %v", err)
	}

	if config.GithubPrivateWorkingDir != "" {
		t.Errorf("GithubPrivateWorkingDir = %q, want empty string", config.GithubPrivateWorkingDir)
	}
}

// TestLoadFileConfig tests the JSON config file parsing
func TestLoadFileConfig(t *testing.T) {
	t.Run("ValidConfigFile", func(t *testing.T) {
		content := `{
			"redis": {"addr": "redis:6380", "channel": "my-channel", "viewSubmissionChannel": "my-vsub", "poppitList": "my-poppit", "slackLinerList": "my-liner"},
			"slack": {"channelNewRepo": "#my-new-repo"},
			"github": {"org": "my-org"},
			"paths": {"workingDir": "/work", "vibeopsWorkingDir": "/vibeops", "githubPrivateWorkingDir": "/private"},
			"logging": {"level": "debug"},
			"issueCreationDelay": 10
		}`
		f, err := os.CreateTemp("", "config-*.json")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		if _, err := f.WriteString(content); err != nil {
			t.Fatal(err)
		}
		f.Close()

		os.Setenv("CONFIG_FILE", f.Name())
		defer os.Unsetenv("CONFIG_FILE")

		fc, path, err := loadFileConfig()
		if err != nil {
			t.Fatalf("loadFileConfig() error: %v", err)
		}
		if path != f.Name() {
			t.Errorf("path = %q, want %q", path, f.Name())
		}
		if fc.Redis.Addr != "redis:6380" {
			t.Errorf("Redis.Addr = %q, want %q", fc.Redis.Addr, "redis:6380")
		}
		if fc.Github.Org != "my-org" {
			t.Errorf("Github.Org = %q, want %q", fc.Github.Org, "my-org")
		}
		if fc.IssueCreationDelay != 10 {
			t.Errorf("IssueCreationDelay = %d, want 10", fc.IssueCreationDelay)
		}
		if fc.Logging.Level != "debug" {
			t.Errorf("Logging.Level = %q, want debug", fc.Logging.Level)
		}
	})

	t.Run("MissingConfigFile", func(t *testing.T) {
		os.Setenv("CONFIG_FILE", "/nonexistent/config.json")
		defer os.Unsetenv("CONFIG_FILE")

		fc, _, err := loadFileConfig()
		if err != nil {
			t.Fatalf("loadFileConfig() should not error for missing file, got: %v", err)
		}
		if fc.Redis.Addr != "" {
			t.Errorf("expected zero-value fileConfig for missing file")
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		f, err := os.CreateTemp("", "config-*.json")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		f.WriteString("{not valid json}")
		f.Close()

		os.Setenv("CONFIG_FILE", f.Name())
		defer os.Unsetenv("CONFIG_FILE")

		_, _, err = loadFileConfig()
		if err == nil {
			t.Error("loadFileConfig() should return error for invalid JSON")
		}
	})
}

// TestLoadConfigFromFile tests that loadConfig reads values from config.json
func TestLoadConfigFromFile(t *testing.T) {
	content := `{
		"redis": {"addr": "redis-from-file:6379"},
		"github": {"org": "file-org"},
		"slack": {"channelNewRepo": "#from-file"},
		"logging": {"level": "warn"},
		"issueCreationDelay": 7
	}`
	f, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString(content)
	f.Close()

	os.Setenv("CONFIG_FILE", f.Name())
	os.Setenv("SLACK_BOT_TOKEN", "test-token")
	defer func() {
		os.Unsetenv("CONFIG_FILE")
		os.Unsetenv("SLACK_BOT_TOKEN")
	}()

	config, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() failed: %v", err)
	}
	if config.RedisAddr != "redis-from-file:6379" {
		t.Errorf("RedisAddr = %q, want redis-from-file:6379", config.RedisAddr)
	}
	if config.GithubOrg != "file-org" {
		t.Errorf("GithubOrg = %q, want file-org", config.GithubOrg)
	}
	if config.LogLevel != "warn" {
		t.Errorf("LogLevel = %q, want warn", config.LogLevel)
	}
	if config.IssueCreationDelay != 7 {
		t.Errorf("IssueCreationDelay = %d, want 7", config.IssueCreationDelay)
	}
}

// TestLoadConfigEnvOverridesFile tests that env vars override config file values
func TestLoadConfigEnvOverridesFile(t *testing.T) {
	content := `{"redis": {"addr": "file-redis:6379"}, "github": {"org": "file-org"}}`
	f, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString(content)
	f.Close()

	os.Setenv("CONFIG_FILE", f.Name())
	os.Setenv("SLACK_BOT_TOKEN", "test-token")
	os.Setenv("REDIS_ADDR", "env-redis:6380")
	defer func() {
		os.Unsetenv("CONFIG_FILE")
		os.Unsetenv("SLACK_BOT_TOKEN")
		os.Unsetenv("REDIS_ADDR")
	}()

	config, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() failed: %v", err)
	}
	if config.RedisAddr != "env-redis:6380" {
		t.Errorf("RedisAddr = %q, want env-redis:6380 (env should override file)", config.RedisAddr)
	}
}

// TestLoadConfigWithIssueCreationDelay tests that the ISSUE_CREATION_DELAY is correctly loaded
func TestLoadConfigWithIssueCreationDelay(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected int
	}{
		{"ValidDelay", "10", 10},
		{"DefaultDelay", "", 5},
		{"InvalidDelay", "abc", 5},
		{"ZeroDelay", "0", 0},
		{"NegativeDelay", "-1", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up test environment
			os.Setenv("SLACK_BOT_TOKEN", "test-token")
			os.Setenv("GITHUB_ORG", "test-org")
			if tt.envValue != "" {
				os.Setenv("ISSUE_CREATION_DELAY", tt.envValue)
			} else {
				os.Unsetenv("ISSUE_CREATION_DELAY")
			}
			defer func() {
				os.Unsetenv("SLACK_BOT_TOKEN")
				os.Unsetenv("GITHUB_ORG")
				os.Unsetenv("ISSUE_CREATION_DELAY")
			}()

			config, err := loadConfig()
			if err != nil {
				t.Fatalf("loadConfig() failed: %v", err)
			}

			if config.IssueCreationDelay != tt.expected {
				t.Errorf("IssueCreationDelay = %d, want %d", config.IssueCreationDelay, tt.expected)
			}
		})
	}
}

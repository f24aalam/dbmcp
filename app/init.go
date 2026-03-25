package app

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/f24aalam/dbmcp/database"
	"github.com/f24aalam/dbmcp/internal/mcpinstall"
	"github.com/f24aalam/dbmcp/internal/projectscan"
	"github.com/f24aalam/dbmcp/storage"
	stepflow "github.com/f24aalam/stepflow/pkg"
)

type initInput struct {
	Name       string
	DBType     string
	URL        string
	Host       string
	Port       string
	Database   string
	User       string
	Password   string
	SQLitePath string
	Evidence   []string
	Source     string
}

func runStepflow(steps ...stepflow.Step) (stepflow.Result, error) {
	return stepflow.New().
		WithAltScreen(false).
		WithDoneScreen(false).
		WithTheme(stepflow.DefaultTheme()).
		WithSteps(steps...).
		Run()
}

func InitProject() {
	root, err := os.Getwd()
	if err != nil {
		showErrorFlow("Initialization failed", fmt.Sprintf("Error getting current directory: %v", err))
		return
	}

	if !projectscan.IsLikelyProjectDir(root) {
		res, err := runStepflow(
			stepflow.Confirm("proceed", "No common project markers found. Continue in manual mode?").
				Default("No"),
		)

		if err != nil {
			if errors.Is(err, stepflow.ErrCancelled) {
				showErrorFlow("Initialization cancelled", "No changes were made.")
				return
			}

			showErrorFlow("Initialization failed", err.Error())
			return
		}

		if !res.Bool("proceed") {
			showErrorFlow("Initialization cancelled", "No changes were made.")
			return
		}
	}

	creds, err := storage.ListCredentials()
	if err != nil {
		showErrorFlow("Initialization failed", fmt.Sprintf("Error loading existing connections: %v", err))
		return
	}

	selected, err := runInitSourceWizard(root, creds)
	if err != nil {
		if errors.Is(err, stepflow.ErrCancelled) {
			showErrorFlow("Initialization cancelled", "No changes were made.")
			return
		}

		showErrorFlow("Initialization failed", err.Error())
		return
	}

	if selected.Name == "" {
		selected.Name = defaultConnectionName(selected)
	}

	connURL, err := buildConnectionURL(selected)
	if err != nil {
		showErrorFlow("Invalid connection details", err.Error())
		return
	}

	_, err = runStepflow(
		stepflow.Info("review", "Connection details").Body(buildConnectionReviewDetails(selected)),
	)

	if err != nil {
		if errors.Is(err, stepflow.ErrCancelled) {
			showErrorFlow("Initialization cancelled", "No changes were made.")
			return
		}

		showErrorFlow("Initialization failed", err.Error())
		return
	}

	conn := database.Connection{
		Database:      selected.DBType,
		ConnectionURL: connURL,
	}

	if err := conn.Open(); err != nil {
		showErrorFlow("Database connection failed", err.Error())
		return
	}
	defer conn.Close()

	connectionID := ""
	if strings.HasPrefix(selected.Source, "existing:") {
		connectionID = strings.TrimPrefix(selected.Source, "existing:")
	} else {
		saveRes, saveFlowErr := runStepflow(
			stepflow.Confirm("save", "Save this connection to dbmcp?").Default("Yes"),
		)
		if saveFlowErr != nil {
			if errors.Is(saveFlowErr, stepflow.ErrCancelled) {
				showErrorFlow("Initialization cancelled", "No changes were made.")
				return
			}

			showErrorFlow("Initialization failed", saveFlowErr.Error())
			return
		}

		if !saveRes.Bool("save") {
			showSuccessFlow("Init completed", "Connection verified.\nNot saved to dbmcp.\nMCP installation skipped.")
			return
		}

		id, saveErr := storage.SaveCredential(nil, selected.Name, selected.DBType, connURL)
		if saveErr != nil {
			showErrorFlow("Failed to save connection", saveErr.Error())
			return
		}

		connectionID = id
	}

	var agentItems []stepflow.ListItem
	for _, t := range mcpinstall.SupportedTargets() {
		agentItems = append(agentItems, stepflow.Item(t.ID, t.Name+" · "+t.Path))
	}

	// Confirm MCP install in its own wizard so declining "No" does not show agent selection.
	mcpRes, err := runStepflow(
		stepflow.Confirm("install_mcp", "Install MCP config for coding agents in this project?").Default("No"),
	)

	if err != nil {
		if errors.Is(err, stepflow.ErrCancelled) {
			showSuccessFlow("Init completed", "Connection validated.\nMCP installation skipped.")
			return
		}

		showErrorFlow("Initialization failed", err.Error())
		return
	}

	if !mcpRes.Bool("install_mcp") {
		showSuccessFlow("Init completed", "Connection validated.\nMCP installation skipped.")
		return
	}

	agentsRes, err := runStepflow(
		stepflow.List("agents", "Select agents to configure (space toggles, enter confirms)").
			Items(agentItems...).
			MultiSelect(true).
			VisibleRows(8),
	)

	if err != nil {
		if errors.Is(err, stepflow.ErrCancelled) {
			showSuccessFlow("Init completed", "Connection validated.\nMCP installation skipped.")
			return
		}

		showErrorFlow("Initialization failed", err.Error())
		return
	}

	chosen := splitCommaAnswers(agentsRes.Get("agents"))
	if len(chosen) == 0 {
		showSuccessFlow("Init completed", "Connection validated.\nNo agents selected; MCP installation skipped.")
		return
	}

	results := mcpinstall.InstallForAgents(root, chosen, connectionID)
	var lines []string
	lines = append(lines, "Connection validated.")
	if !strings.HasPrefix(selected.Source, "existing:") {
		lines = append(lines, "Connection saved with ID: "+connectionID)
	}

	lines = append(lines, "MCP installation results:")
	hasErr := false
	for _, r := range results {
		if r.Err != nil {
			hasErr = true
			lines = append(lines, fmt.Sprintf("- %s: failed (%v)", r.TargetID, r.Err))

			continue
		}

		lines = append(lines, fmt.Sprintf("- %s: written to %s", r.TargetID, r.Path))
	}

	if hasErr {
		showErrorFlow("Init completed with errors", strings.Join(lines, "\n"))
		return
	}

	showSuccessFlow("Init completed successfully", strings.Join(lines, "\n"))
}

func runInitSourceWizard(root string, creds []storage.Credential) (initInput, error) {
	modeItems := []stepflow.ListItem{
		stepflow.Item("existing", "Use existing connections"),
		stepflow.Item("scan", "Scan this project for database envs"),
		stepflow.Item("new", "Add new connection"),
	}

	res, err := runStepflow(
		stepflow.List("init_mode", "How do you want to initialize dbmcp?").
			Items(modeItems...).
			MultiSelect(false).
			VisibleRows(5),
	)
	if err != nil {
		return initInput{}, err
	}

	mode := strings.TrimSpace(res.Get("init_mode"))
	switch mode {
	case "new":
		return runNewConnectionWizard()
	case "existing":
		existingItems := []stepflow.ListItem{}
		labelToID := map[string]string{}

		for i, c := range creds {
			_ = i
			label := fmt.Sprintf("%s (%s) · %s", c.Name, c.Database, c.ID)
			meta := "existing connection"
			labelToID[label] = c.ID
			existingItems = append(existingItems, stepflow.Item(label, meta))
		}

		existingRes, existingErr := runStepflow(
			stepflow.List("existing_pick", "Select an existing connection").
				Items(existingItems...).
				MultiSelect(false).
				VisibleRows(8),
		)
		if existingErr != nil {
			return initInput{}, existingErr
		}

		choice := strings.TrimSpace(existingRes.Get("existing_pick"))
		if len(creds) == 0 {
			return initInput{DBType: "mysql", Port: "3306"}, nil
		}

		id, ok := labelToID[choice]
		if !ok {
			return initInput{}, fmt.Errorf("invalid existing connection selection")
		}

		for _, cred := range creds {
			if cred.ID == id {
				return existingCredentialToInput(cred), nil
			}
		}

		return initInput{}, fmt.Errorf("selected connection not found")
	default:
		scanResult, scanErr := projectscan.ScanProject(root)
		if scanErr != nil {
			return initInput{}, scanErr
		}

		candidates := scanResult.Candidates
		if len(candidates) == 0 {
			return runNewConnectionWizard()
		}

		scanItems := []stepflow.ListItem{
			stepflow.Item("manual", "Enter connection manually"),
		}

		for i, c := range candidates {
			label := strconv.Itoa(i)
			meta := fmt.Sprintf("%s — %s [%s] score %d", candidateSummary(c), c.SourceFile, c.Parser, c.Confidence)
			scanItems = append(scanItems, stepflow.Item(label, meta))
		}

		scanRes, scanRunErr := runStepflow(
			stepflow.List("scan_pick", "Select discovered database config").
				Items(scanItems...).
				MultiSelect(false).
				VisibleRows(8),
		)

		if scanRunErr != nil {
			return initInput{}, scanRunErr
		}

		choice := strings.TrimSpace(scanRes.Get("scan_pick"))
		if choice == "manual" {
			return runNewConnectionWizard()
		}

		idx, convErr := strconv.Atoi(choice)
		if convErr != nil || idx < 0 || idx >= len(candidates) {
			return initInput{}, fmt.Errorf("invalid scanned connection selection")
		}

		return candidateToInput(candidates[idx]), nil
	}
}

func runNewConnectionWizard() (initInput, error) {
	res, err := runStepflow(
		stepflow.Text("name", "Connection name"),
		stepflow.List("db_type", "Database type").
			Items(stepflow.Item("mysql"), stepflow.Item("postgres"), stepflow.Item("sqlite")).
			MultiSelect(false).
			VisibleRows(4),
		stepflow.Text("url", "Database URL / path"),
	)

	if err != nil {
		return initInput{}, err
	}

	in := initInput{
		Name:   strings.TrimSpace(res.Get("name")),
		DBType: strings.TrimSpace(res.Get("db_type")),
		URL:    strings.TrimSpace(res.Get("url")),
		Source: "manual",
	}

	if in.DBType == "sqlite" {
		in.SQLitePath = in.URL
	}

	return in, nil
}

func splitCommaAnswers(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	parts := strings.Split(s, ", ")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}

	return out
}

func candidateToInput(c projectscan.Candidate) initInput {
	in := initInput{
		DBType:     c.DBType,
		URL:        c.URL,
		Host:       c.Host,
		Port:       c.Port,
		Database:   c.Database,
		User:       c.User,
		Password:   c.Password,
		SQLitePath: c.SQLitePath,
		Evidence:   append([]string(nil), c.Evidence...),
		Source:     "scan:" + c.SourceFile,
	}

	if in.Name == "" {
		in.Name = defaultConnectionName(in)
	}

	return in
}

func existingCredentialToInput(c storage.Credential) initInput {
	_, connURL, err := storage.GetCredentialById(c.ID)
	if err != nil {
		return initInput{
			Name:   c.Name,
			DBType: c.Database,
		}
	}

	in := initInput{
		Name:   c.Name,
		DBType: c.Database,
		URL:    connURL,
		Source: "existing:" + c.ID,
	}

	if c.Database == "sqlite" {
		in.SQLitePath = connURL
	}

	return in
}

func buildConnectionURL(in initInput) (string, error) {
	if strings.TrimSpace(in.URL) != "" {
		if in.DBType == "sqlite" && !strings.Contains(in.URL, "://") {
			return in.URL, nil
		}

		u, err := url.Parse(in.URL)
		if err == nil && u.Scheme != "" {
			switch in.DBType {
			case "mysql":
				if strings.EqualFold(u.Scheme, "mysql") {
					return canonicalMySQLDSNFromURL(in.URL)
				}
			case "postgres":
				if u.Scheme == "postgres" || u.Scheme == "postgresql" {
					return sanitizePostgresURLForPQ(in.URL)
				}
			}
		}

		if _, err := url.Parse(in.URL); err == nil {
			return in.URL, nil
		}
	}

	switch in.DBType {
	case "sqlite":
		if strings.TrimSpace(in.SQLitePath) == "" {
			if strings.TrimSpace(in.Database) == "" {
				return "", fmt.Errorf("sqlite path is required")
			}

			return in.Database, nil
		}

		return in.SQLitePath, nil
	case "postgres":
		host := fallback(in.Host, "localhost")
		port := fallback(in.Port, "5432")
		db := fallback(in.Database, "postgres")
		user := url.QueryEscape(in.User)
		pass := url.QueryEscape(in.Password)

		if in.User == "" {
			return fmt.Sprintf("postgres://%s:%s/%s?sslmode=disable", host, port, db), nil
		}

		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, db), nil
	case "mysql":
		host := fallback(in.Host, "127.0.0.1")
		port := fallback(in.Port, "3306")
		db := in.Database

		if db == "" {
			return "", fmt.Errorf("database name is required for mysql")
		}

		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", in.User, in.Password, host, port, db), nil
	default:
		return "", fmt.Errorf("unsupported database type: %s", in.DBType)
	}
}

func candidateSummary(c projectscan.Candidate) string {
	if c.URL != "" {
		return fmt.Sprintf("%s (url)", c.DBType)
	}

	if c.DBType == "sqlite" {
		return fmt.Sprintf("sqlite (%s)", c.SQLitePath)
	}

	return fmt.Sprintf("%s %s:%s/%s", c.DBType, c.Host, c.Port, c.Database)
}

func defaultConnectionName(in initInput) string {
	if in.DBType == "sqlite" {
		if in.SQLitePath != "" {
			return "sqlite-local"
		}

		return "sqlite-db"
	}

	host := in.Host
	if host == "" {
		host = "localhost"
	}

	db := in.Database
	if db == "" {
		db = in.DBType
	}

	return fmt.Sprintf("%s-%s", host, db)
}

func obfuscateSource(in initInput) string {
	if in.URL != "" {
		return "URL (hidden credentials)"
	}

	if in.DBType == "sqlite" {
		return in.SQLitePath
	}

	if in.Password != "" {
		return fmt.Sprintf("%s@%s:%s/%s", in.User, in.Host, in.Port, in.Database)
	}

	return fmt.Sprintf("%s:%s/%s", in.Host, in.Port, in.Database)
}

func buildConnectionReviewDetails(in initInput) string {
	lines := []string{
		"Name: " + in.Name,
		"Type: " + in.DBType,
		"Source: " + in.Source,
		"Target: " + obfuscateSource(in),
	}

	if len(in.Evidence) > 0 {
		lines = append(lines, "Evidence: "+strings.Join(in.Evidence, "; "))
	}

	return strings.Join(lines, "\n")
}

func fallback(s, defaultValue string) string {
	if strings.TrimSpace(s) == "" {
		return defaultValue
	}

	return s
}

func showSuccessFlow(title, body string) {
	_, _ = runStepflow(stepflow.Success("success", title).Body(body))
}

func showErrorFlow(title, body string) {
	_, _ = runStepflow(stepflow.Error("error", title).Body(body))
}

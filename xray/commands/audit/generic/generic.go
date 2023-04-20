package audit

import (
	"os"

	"github.com/jfrog/jfrog-cli-core/v2/utils/coreutils"
	cmdUtils "github.com/jfrog/jfrog-cli-core/v2/xray/commands/utils"
	xrutils "github.com/jfrog/jfrog-cli-core/v2/xray/utils"
	"github.com/jfrog/jfrog-client-go/xray/services"
)

type GenericAuditCommand struct {
	watches                []string
	workingDirs            []string
	projectKey             string
	targetRepoPath         string
	IncludeVulnerabilities bool
	IncludeLicenses        bool
	Fail                   bool
	PrintExtendedTable     bool
	*cmdUtils.GraphBasicParams
}

func NewGenericAuditCommand() *GenericAuditCommand {
	return &GenericAuditCommand{}
}

func (auditCmd *GenericAuditCommand) SetWatches(watches []string) *GenericAuditCommand {
	auditCmd.watches = watches
	return auditCmd
}

func (auditCmd *GenericAuditCommand) SetWorkingDirs(dirs []string) *GenericAuditCommand {
	auditCmd.workingDirs = dirs
	return auditCmd
}

func (auditCmd *GenericAuditCommand) SetProject(project string) *GenericAuditCommand {
	auditCmd.projectKey = project
	return auditCmd
}

func (auditCmd *GenericAuditCommand) SetTargetRepoPath(repoPath string) *GenericAuditCommand {
	auditCmd.targetRepoPath = repoPath
	return auditCmd
}

func (auditCmd *GenericAuditCommand) SetIncludeVulnerabilities(include bool) *GenericAuditCommand {
	auditCmd.IncludeVulnerabilities = include
	return auditCmd
}

func (auditCmd *GenericAuditCommand) SetIncludeLicenses(include bool) *GenericAuditCommand {
	auditCmd.IncludeLicenses = include
	return auditCmd
}

func (auditCmd *GenericAuditCommand) SetFail(fail bool) *GenericAuditCommand {
	auditCmd.Fail = fail
	return auditCmd
}

func (auditCmd *GenericAuditCommand) SetPrintExtendedTable(printExtendedTable bool) *GenericAuditCommand {
	auditCmd.PrintExtendedTable = printExtendedTable
	return auditCmd
}

func (auditCmd *GenericAuditCommand) CreateXrayGraphScanParams() services.XrayGraphScanParams {
	params := services.XrayGraphScanParams{
		RepoPath: auditCmd.targetRepoPath,
		Watches:  auditCmd.watches,
		ScanType: services.Dependency,
	}
	if auditCmd.projectKey == "" {
		params.ProjectKey = os.Getenv(coreutils.Project)
	} else {
		params.ProjectKey = auditCmd.projectKey
	}
	params.IncludeVulnerabilities = auditCmd.IncludeVulnerabilities
	params.IncludeLicenses = auditCmd.IncludeLicenses
	return params
}

func (auditCmd *GenericAuditCommand) Run() (err error) {
	if err != nil {
		return
	}
	auditParams := NewAuditParams().
		SetXrayGraphScanParams(auditCmd.CreateXrayGraphScanParams())
	auditParams.GraphBasicParams = auditCmd.GraphBasicParams
	results, isMultipleRootProject, auditErr := GenericAudit(auditParams)

	if auditCmd.Progress != nil {
		err = auditCmd.Progress.Quit()
		if err != nil {
			return
		}
	}
	// Print Scan results on all cases except if errors accrued on Generic Audit command and no security/license issues found.
	printScanResults := !(auditErr != nil && xrutils.IsEmptyScanResponse(results))
	if printScanResults {
		err = xrutils.PrintScanResults(results,
			nil,
			auditCmd.OutputFormat,
			auditCmd.IncludeVulnerabilities,
			auditCmd.IncludeLicenses,
			isMultipleRootProject,
			auditCmd.PrintExtendedTable, false,
		)
		if err != nil {
			return
		}
	}
	if auditErr != nil {
		err = auditErr
		return
	}

	// Only in case Xray's context was given (!auditCmd.IncludeVulnerabilities) and the user asked to fail the build accordingly, do so.
	if auditCmd.Fail && !auditCmd.IncludeVulnerabilities && xrutils.CheckIfFailBuild(results) {
		err = xrutils.NewFailBuildError()
	}
	return
}

func (auditCmd *GenericAuditCommand) CommandName() string {
	return "generic_audit"
}

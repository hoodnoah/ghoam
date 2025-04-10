package services

import (
	"context"

	"github.com/hoodnoah/ghoam/internal/accounting"
)

type ChartOfAccountsService struct {
	AccountRepo      accounting.AccountRepository
	AccountGroupRepo accounting.AccountGroupRepository
}

// Produces a ChartOfAccounts tree from the AccountRepo and AccountGroup repo
func (s *ChartOfAccountsService) GetChartOfAccounts(ctx context.Context) (*accounting.ChartOfAccountsNode, error) {
	return accounting.BuildChartOfAccountsTree(ctx, s.AccountGroupRepo, s.AccountRepo)
}

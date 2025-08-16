package orgsvc

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

func (s *service) ListTeamMembers(ctx core.Context, teamId string, current, pageSize int) (list interface{}, total int64, err error) {
	return []interface{}{}, 0, nil
}

func (s *service) AddTeamMember(ctx core.Context, teamId string, accountId string, roleType string) (id string, err error) {
	return "", nil
}

func (s *service) RemoveTeamMember(ctx core.Context, teamId string, accountId string) (err error) {
	return nil
}

func (s *service) UpdateTeamMemberRole(ctx core.Context, teamId string, accountId string, roleType string) (err error) {
	return nil
}

func (s *service) ListUnassignedAccounts(ctx core.Context, current, pageSize int, keyword string) (list interface{}, total int64, err error) {
	return []interface{}{}, 0, nil
}

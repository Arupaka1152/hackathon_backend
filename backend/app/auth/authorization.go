package auth

import (
	"backend/app/dao"
	"backend/app/model"
	"errors"
)

func UserAuth(workspaceId string, accountId string) (userId string, role string, err error) {
	targetUser := model.User{}
	err = dao.FetchUserById(&targetUser, workspaceId, accountId).Error
	if err != nil {
		return "", "", err
	}

	userId = targetUser.Id
	role = targetUser.Role
	return userId, role, nil
}

func ContributionAuth(contributionId string, userId string) (err error) {
	targetContribution := model.Contribution{}
	if err = dao.FetchContribution(&targetContribution, contributionId).Error; err != nil {
		return err
	}

	if targetContribution.From != userId {
		return errors.New("not permitted")
	}

	return nil
}

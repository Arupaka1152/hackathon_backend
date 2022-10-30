package auth

import (
	"backend/app/dao"
	"backend/app/model"
	"errors"
)

func UserAuth(workspaceId string, accountId string) (userId string, role string, err error) {
	targetUser := model.User{}
	err = dao.FindUserById(&targetUser, workspaceId, accountId).Error
	if err != nil {
		return "", "", err
	}

	if targetUser.Id == "" {
		return "", "", errors.New("not permitted")
	}

	userId = targetUser.Id
	role = targetUser.Role
	return userId, role, nil
}

func ContributionAuth(contributionId string, userId string) (err error) {
	targetContribution := model.Contribution{}
	if err = dao.FindContribution(&targetContribution, contributionId).Error; err != nil {
		return err
	}

	if targetContribution.From != userId {
		return errors.New("not permitted")
	}

	return nil
}

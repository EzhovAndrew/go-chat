package dto

import (
	"github.com/go-chat/users/internal/domain"
	usersv1 "github.com/go-chat/users/pkg/api/users/v1"
)

// ToProtoProfile converts domain.UserProfile to proto UserProfile
func ToProtoProfile(profile *domain.UserProfile) *usersv1.UserProfile {
	return &usersv1.UserProfile{
		UserId:    profile.UserID.String(),
		Nickname:  profile.Nickname,
		Bio:       profile.Bio,
		AvatarUrl: profile.AvatarURL,
	}
}

// ToProtoProfiles converts a slice of domain.UserProfile to proto UserProfile slice
func ToProtoProfiles(profiles []*domain.UserProfile) []*usersv1.UserProfile {
	result := make([]*usersv1.UserProfile, len(profiles))
	for i, profile := range profiles {
		result[i] = ToProtoProfile(profile)
	}
	return result
}


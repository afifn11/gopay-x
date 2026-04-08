package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/afifn11/gopay-x/services/user-service/internal/domain"
)

var (
	ErrProfileNotFound    = errors.New("user profile not found")
	ErrProfileExists      = errors.New("user profile already exists")
	ErrKYCAlreadyVerified = errors.New("KYC already verified")
)

type UserUsecase interface {
	CreateProfile(ctx context.Context, req *CreateProfileRequest) (*domain.UserProfile, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.UserProfile, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *UpdateProfileRequest) (*domain.UserProfile, error)
	SubmitKYC(ctx context.Context, userID uuid.UUID, req *SubmitKYCRequest) error
	UpdateKYCStatus(ctx context.Context, userID uuid.UUID, req *UpdateKYCStatusRequest) error
	DeleteProfile(ctx context.Context, userID uuid.UUID) error
}

type CreateProfileRequest struct {
	UserID   string `json:"user_id" validate:"required,uuid"`
	FullName string `json:"full_name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
}

type UpdateProfileRequest struct {
	FullName    string  `json:"full_name" validate:"omitempty,min=2,max=100"`
	Phone       string  `json:"phone" validate:"omitempty"`
	Gender      string  `json:"gender" validate:"omitempty,oneof=male female other"`
	DateOfBirth string  `json:"date_of_birth" validate:"omitempty"`
	Address     string  `json:"address" validate:"omitempty,max=255"`
	City        string  `json:"city" validate:"omitempty,max=100"`
	AvatarURL   string  `json:"avatar_url" validate:"omitempty,url"`
}

type SubmitKYCRequest struct {
	DocumentType string `json:"document_type" validate:"required,oneof=ktp passport sim"`
	DocumentURL  string `json:"document_url" validate:"required,url"`
}

type UpdateKYCStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=verified rejected"`
	Note   string `json:"note" validate:"omitempty,max=255"`
}

type userUsecase struct {
	profileRepo domain.UserProfileRepository
	kycRepo     domain.KYCDocumentRepository
}

func NewUserUsecase(
	profileRepo domain.UserProfileRepository,
	kycRepo domain.KYCDocumentRepository,
) UserUsecase {
	return &userUsecase{
		profileRepo: profileRepo,
		kycRepo:     kycRepo,
	}
}

func (uc *userUsecase) CreateProfile(ctx context.Context, req *CreateProfileRequest) (*domain.UserProfile, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, err
	}

	existing, err := uc.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrProfileExists
	}

	profile := &domain.UserProfile{
		UserID:   userID,
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	if err := uc.profileRepo.Create(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func (uc *userUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.UserProfile, error) {
	profile, err := uc.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, ErrProfileNotFound
	}
	return profile, nil
}

func (uc *userUsecase) UpdateProfile(ctx context.Context, userID uuid.UUID, req *UpdateProfileRequest) (*domain.UserProfile, error) {
	profile, err := uc.profileRepo.FindByUserID(ctx, userID)
	if err != nil || profile == nil {
		return nil, ErrProfileNotFound
	}

	if req.FullName != "" {
		profile.FullName = req.FullName
	}
	if req.Phone != "" {
		profile.Phone = req.Phone
	}
	if req.Gender != "" {
		g := domain.Gender(req.Gender)
		profile.Gender = &g
	}
	if req.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err == nil {
			profile.DateOfBirth = &dob
		}
	}
	if req.Address != "" {
		profile.Address = req.Address
	}
	if req.City != "" {
		profile.City = req.City
	}
	if req.AvatarURL != "" {
		profile.AvatarURL = req.AvatarURL
	}

	if err := uc.profileRepo.Update(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func (uc *userUsecase) SubmitKYC(ctx context.Context, userID uuid.UUID, req *SubmitKYCRequest) error {
	profile, err := uc.profileRepo.FindByUserID(ctx, userID)
	if err != nil || profile == nil {
		return ErrProfileNotFound
	}

	if profile.KYCStatus == domain.KYCVerified {
		return ErrKYCAlreadyVerified
	}

	doc := &domain.KYCDocument{
		UserID:       userID,
		DocumentType: req.DocumentType,
		DocumentURL:  req.DocumentURL,
		SubmittedAt:  time.Now(),
	}

	if err := uc.kycRepo.Create(ctx, doc); err != nil {
		return err
	}

	return uc.profileRepo.UpdateKYCStatus(ctx, userID, domain.KYCPending, "Document submitted, under review")
}

func (uc *userUsecase) UpdateKYCStatus(ctx context.Context, userID uuid.UUID, req *UpdateKYCStatusRequest) error {
	profile, err := uc.profileRepo.FindByUserID(ctx, userID)
	if err != nil || profile == nil {
		return ErrProfileNotFound
	}

	return uc.profileRepo.UpdateKYCStatus(ctx, userID, domain.KYCStatus(req.Status), req.Note)
}

func (uc *userUsecase) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	return uc.profileRepo.Delete(ctx, userID)
}
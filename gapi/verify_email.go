package gapi

import (
	"context"

	db "github.com/santinofajardo/simpleBank/db/sqlc"
	"github.com/santinofajardo/simpleBank/pb"
	"github.com/santinofajardo/simpleBank/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmailRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.store.VerifyEmailTransaction(ctx, db.VerifyEmailTransactionParams{
		SecretCode: req.GetSecretCode(),
		EmailId:    req.GetEmailId(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.VerifyEmailResponse{
		IsVerified: user.User.IsEmailVerified,
	}, nil
}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, fieldViolation("email_id", err))
	}

	if err := validation.ValidateEmailId(req.GetEmailId()); err != nil {
		violations = append(violations, fieldViolation("secret_code", err))
	}

	return violations
}

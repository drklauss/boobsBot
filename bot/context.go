package bot

import (
	"context"
	"errors"

	"github.com/jinzhu/gorm"
)

type botString string

// SetDB sets gorm.DB into context.
func SetDB(ctx context.Context, conn *gorm.DB) context.Context {
	return context.WithValue(ctx, "db", conn)
}

// GetDB returns gorm.DB from context.
func GetDB(ctx context.Context) (*gorm.DB, error) {
	conn, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not get db connection from context")
	}
	return conn, nil
}

// SetCategory sets category into context.
func SetCategory(ctx context.Context, category *string) context.Context {
	return context.WithValue(ctx, "category", category)
}

// GetCategory returns category from context.
func GetCategory(ctx context.Context) (*string, error) {
	cat, ok := ctx.Value("category").(*string)
	if !ok {
		return nil, errors.New("could not get category from context")
	}
	return cat, nil
}

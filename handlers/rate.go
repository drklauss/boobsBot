package handlers

import (
	"context"
	"fmt"

	"github.com/drklauss/boobsBot/telegram"
)

func Rate(ctx context.Context, u *telegram.Update) {
	fmt.Println("rate handler")
}

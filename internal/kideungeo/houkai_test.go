package kideungeo

import (
	"context"
	"fmt"
	"testing"
)

func Test_getCouponNumber(t *testing.T) {
	num, _ := getCouponNumber(context.Background())
	fmt.Println(num)
}

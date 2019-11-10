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

func Test_getValkyrieList(t *testing.T) {
	valkyries, _ := getValkyrieList(context.Background())
	fmt.Println(valkyries)
}

func Test_getValkyrieInfo(t *testing.T) {
	valkyrie, _ := getValkyrieInfo(context.Background(), "이치의 율자")
	fmt.Println(valkyrie)
}

func Test_getClosestMatchingValkyre(t *testing.T) {
	valkyries := []ValkyrieSimple{
		ValkyrieSimple{"키아나 카스리나", "투여복 백련"},
		ValkyrieSimple{"키아나 카스리나", "발키리 레인저"},
		ValkyrieSimple{"키아나 카스리나", "성녀의 기도"},
		ValkyrieSimple{"키아나 카스리나", "백기사 월광"},
		ValkyrieSimple{"키아나 카스리나", "공간의 율자"},
		ValkyrieSimple{"테레사 아포칼립스", "[증폭] 처형복 반혼초"},
		ValkyrieSimple{"테레사 아포칼립스", "처형복 반혼초"},
		ValkyrieSimple{"브로냐 자이칙", "이치의 율자"},
	}
	test1 := getClosestMatchingValkyre(valkyries, "율자")
	if test1 != "공간의 율자" {
		fmt.Errorf("Closest match failed: %v", test1)
	}
	test2 := getClosestMatchingValkyre(valkyries, "성녀")
	if test2 != "성녀의 기도" {
		fmt.Errorf("Closest match failed: %v", test2)
	}
	test3 := getClosestMatchingValkyre(valkyries, "기도")
	if test3 != "성녀의 기도" {
		fmt.Errorf("Closest match failed: %v", test3)
	}
	test4 := getClosestMatchingValkyre(valkyries, "증폭 반혼초")
	if test4 != "[증폭] 처형복 반혼초" {
		fmt.Errorf("Closest match failed: %v", test4)
	}
}

package services

import (
	"testing"

	"github.com/google/uuid"
)

func TestDirectKeyIsOrderIndependent(t *testing.T) {
	first := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	second := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	if directKey(first, second) != directKey(second, first) {
		t.Fatal("direct key must not depend on participant order")
	}
}

func TestDirectKeyIsUniquePerPair(t *testing.T) {
	first := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	second := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	third := uuid.MustParse("33333333-3333-3333-3333-333333333333")

	if directKey(first, second) == directKey(first, third) {
		t.Fatal("different participant pairs must have different direct keys")
	}
}

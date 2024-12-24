package main

import "github.com/golang-jwt/jwt/v5"

const (
	RealmAccess = "realm_access"
	Roles       = "roles"
)

// RolesClaims структура контейнер с методом ToMap
type RolesClaims struct {
	Claims jwt.MapClaims
}

type mapClaims map[string]interface{}

type sliceClaims []interface{}

var emptyMap = map[string]interface{}{}

// ToMap метод возвращает тип mapClaims из Claims по ключу RealmAccess
//
//goland:noinspection GoExportedFuncWithUnexportedType
func (r RolesClaims) ToMap() mapClaims {
	if ra, ok1 := r.Claims[RealmAccess]; !ok1 {
		return emptyMap
	} else if realmAccess, ok2 := ra.(map[string]interface{}); !ok2 {
		return emptyMap
	} else {
		return realmAccess
	}
}

// ToSlice метод возвращает тип sliceClaims из mapClaims по ключу Roles
func (mc mapClaims) ToSlice() sliceClaims {
	if a, ok1 := mc[Roles]; !ok1 {
		return []interface{}{}
	} else if roles, ok2 := a.([]interface{}); !ok2 {
		return []interface{}{}
	} else {
		return roles
	}
}

// ToRoles метод возвращает множество ролей
func (sc sliceClaims) ToRoles() map[string]struct{} {
	result := make(map[string]struct{})
	for _, value := range sc {
		if role, ok := value.(string); ok {
			result[role] = struct{}{}
		}
	}
	return result
}

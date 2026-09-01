package models

type AuthRole string 

const(
	UserRole AuthRole = "user"
	AdminRole AuthRole = "admin"
)
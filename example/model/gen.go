package model_test

//go:generate tools gen model -f=./ddl/orders.sql

//go:generate tools gen model -f=./ddl/test.txt

//go:generate tools kiple dao -i UserProfilesDao -e "../entity/user_profiles.go"

//go:generate tools gen method -n UserProfiles

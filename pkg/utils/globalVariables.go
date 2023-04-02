package utils

// WorkEnvironment ...
// OPTINAL VALUES: "LOCAL_DEV", "LOCAL_PROD", "DEV", "PROD"
// 'LOCAL' differs only the output format of the logger etc. Nothing else.
// var WorkEnvironment string = GetEnvValue(WorkEnvironmentKey)
var WorkEnvironment string = GetEnvValue()


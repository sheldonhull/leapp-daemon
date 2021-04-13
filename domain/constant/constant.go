package domain

const CredentialsFilePath = `.aws/credentials`


const SessionTypePlain = "PLAIN"
const SessionTypeFederated = "FEDERATED"

const SessionTokenDurationInSeconds int64 = 3600
const RotationIntervalInSeconds int64 = 15

package constant

const AwsCredentialsFilePath = `.aws/credentials`
const AlibabaCredentialsFilePath = `.aliyun/config.json`

const SessionTypePlain = "PLAIN"
const SessionTypeFederated = "FEDERATED"

const SessionTokenDurationInSeconds int64 = 3600
const RotationIntervalInSeconds int64 = 15

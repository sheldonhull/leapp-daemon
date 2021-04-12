swagger:
	swagger generate spec -o ./swagger.yaml --scan-models

swagger-serve:
	swagger serve ./swagger.yaml

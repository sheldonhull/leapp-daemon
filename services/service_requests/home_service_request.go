package service_requests

type HomeServiceRequest struct {
	Name string
}

/*func (serviceRequest *HomeServiceRequest) Build(m map[string]interface{}) error {
	name := m["Name"]

	if name == nil {
		return errors.New("service request build error")
	}

	castedName, ok := name.(string)

	if !ok {
		return errors.New("service request build error")
	}

	serviceRequest.Name = castedName
	return nil
}*/

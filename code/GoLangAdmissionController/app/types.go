package main

//Admission Review Request
type AdmissionReview struct {
	Request Request `json:"request"`
}

type Request struct {
	Object Object `json:"object"`
	UID    string `json:"uid"`
}

type Object struct {
	Spec Spec `json:"spec"`
}

type Spec struct {
	Containers []Container `json:"containers"`
}

type Container struct {
	Image string `json:"image"`
}

//Admission Review Response
type AdmissionReviewRes struct {
	UID     string         `json:"uid"`
	Allowed bool           `json:"allowed"`
	Status  ResponseStatus `json:"status"`
}

type ResponseStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

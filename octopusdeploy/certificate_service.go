package octopusdeploy

import (
	"fmt"

	"github.com/dghubble/sling"
)

// certificateService handles communication with Certificate-related methods of the Octopus API.
type certificateService struct {
	canDeleteService
}

// newCertificateService returns an certificateService with a preconfigured client.
func newCertificateService(sling *sling.Sling, uriTemplate string) *certificateService {
	certificateService := &certificateService{}
	certificateService.service = newService(ServiceCertificateService, sling, uriTemplate)

	return certificateService
}

func (s certificateService) getPagedResponse(path string) ([]*CertificateResource, error) {
	resources := []*CertificateResource{}
	loadNextPage := true

	for loadNextPage {
		resp, err := apiGet(s.getClient(), new(CertificateResources), path)
		if err != nil {
			return resources, err
		}

		responseList := resp.(*CertificateResources)
		resources = append(resources, responseList.Items...)
		path, loadNextPage = LoadNextPage(responseList.PagedResults)
	}

	return resources, nil
}

// Add creates a new certificate.
func (s certificateService) Add(resource *CertificateResource) (*CertificateResource, error) {
	path, err := getAddPath(s, resource)
	if err != nil {
		return nil, err
	}

	resp, err := apiAdd(s.getClient(), resource, new(CertificateResource), path)
	if err != nil {
		return nil, err
	}

	return resp.(*CertificateResource), nil
}

// GetAll returns all certificates. If none can be found or an error occurs, it
// returns an empty collection.
func (s certificateService) GetAll() ([]*CertificateResource, error) {
	items := []*CertificateResource{}
	path, err := getAllPath(s)
	if err != nil {
		return items, err
	}

	_, err = apiGet(s.getClient(), &items, path)
	return items, err
}

// GetByID returns the certificate that matches the input ID. If one cannot be
// found, it returns nil and an error.
func (s certificateService) GetByID(id string) (*CertificateResource, error) {
	path, err := getByIDPath(s, id)
	if err != nil {
		return nil, err
	}

	resp, err := apiGet(s.getClient(), new(CertificateResource), path)
	if err != nil {
		return nil, createResourceNotFoundError(s.getName(), "ID", id)
	}

	return resp.(*CertificateResource), nil
}

// GetByPartialName performs a lookup and returns instances of a Certificate with a matching partial name.
func (s certificateService) GetByPartialName(name string) ([]*CertificateResource, error) {
	path, err := getByPartialNamePath(s, name)
	if err != nil {
		return []*CertificateResource{}, err
	}

	return s.getPagedResponse(path)
}

// Update modifies a Certificate based on the one provided as input.
func (s certificateService) Update(resource CertificateResource) (*CertificateResource, error) {
	path, err := getUpdatePath(s, &resource)
	if err != nil {
		return nil, err
	}

	resp, err := apiUpdate(s.getClient(), resource, new(CertificateResource), path)
	if err != nil {
		return nil, err
	}

	return resp.(*CertificateResource), nil
}

func (s certificateService) Replace(certificateID string, replacementCertificate *ReplacementCertificate) (*CertificateResource, error) {
	if isEmpty(certificateID) {
		return nil, createInvalidParameterError(OperationReplace, ParameterCertificateID)
	}

	if replacementCertificate == nil {
		return nil, createInvalidParameterError(OperationReplace, ParameterReplacementCertificate)
	}

	err := validateInternalState(s)
	if err != nil {
		return nil, err
	}

	path := trimTemplate(s.getPath())
	path = fmt.Sprintf(path+"/%s/replace", certificateID)

	_, err = apiPost(s.getClient(), replacementCertificate, new(CertificateResource), path)
	if err != nil {
		return nil, err
	}

	//The API endpoint /certificates/id/replace returns the old cert, we need to re-query to get the updated one.
	return s.GetByID(certificateID)
}

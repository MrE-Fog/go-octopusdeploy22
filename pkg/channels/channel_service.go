package channels

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/internal"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/constants"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/services"
	"github.com/dghubble/sling"
)

type ChannelService struct {
	versionRuleTestPath string

	services.CanDeleteService
}

func NewChannelService(sling *sling.Sling, uriTemplate string, versionRuleTestPath string) *ChannelService {
	return &ChannelService{
		versionRuleTestPath: versionRuleTestPath,
		CanDeleteService: services.CanDeleteService{
			Service: services.NewService(constants.ServiceChannelService, sling, uriTemplate),
		},
	}
}

// Add creates a new channel.
func (s *ChannelService) Add(channel *Channel) (*Channel, error) {
	if IsNil(channel) {
		return nil, internal.CreateInvalidParameterError(constants.OperationAdd, constants.ParameterChannel)
	}

	if err := channel.Validate(); err != nil {
		return nil, internal.CreateValidationFailureError(constants.OperationAdd, err)
	}

	path, err := services.GetAddPath(s, channel)
	if err != nil {
		return nil, err
	}

	response, err := services.ApiAdd(s.GetClient(), channel, new(Channel), path)
	if err != nil {
		return nil, err
	}

	return response.(*Channel), nil
}

// Get returns a collection of channels based on the criteria defined by its
// input query parameter. If an error occurs, an empty collection is returned
// along with the associated error.
func (s *ChannelService) Get(channelsQuery Query) (*Channels, error) {
	path, err := s.GetURITemplate().Expand(channelsQuery)
	if err != nil {
		return &Channels{}, err
	}

	response, err := services.ApiGet(s.GetClient(), new(Channels), path)
	if err != nil {
		return &Channels{}, err
	}

	return response.(*Channels), nil
}

// GetAll returns all channels. If none can be found or an error occurs, it
// returns an empty collection.
func (s *ChannelService) GetAll() ([]*Channel, error) {
	items := []*Channel{}
	path, err := services.GetAllPath(s)
	if err != nil {
		return items, err
	}

	_, err = services.ApiGet(s.GetClient(), &items, path)
	return items, err
}

// GetByID returns the channel that matches the input ID. If one cannot be
// found, it returns nil and an error.
func (s *ChannelService) GetByID(id string) (*Channel, error) {
	if internal.IsEmpty(id) {
		return nil, internal.CreateInvalidParameterError(constants.OperationGetByID, "id")
	}

	path := s.BasePath + "/" + id
	resp, err := services.ApiGet(s.GetClient(), new(Channel), path)
	if err != nil {
		return nil, err
	}

	return resp.(*Channel), nil
}

// Update modifies a channel based on the one provided as input.
func (s *ChannelService) Update(channel *Channel) (*Channel, error) {
	path, err := services.GetUpdatePath(s, channel)
	if err != nil {
		return nil, err
	}

	resp, err := services.ApiUpdate(s.GetClient(), channel, new(Channel), path)
	if err != nil {
		return nil, err
	}

	return resp.(*Channel), nil
}
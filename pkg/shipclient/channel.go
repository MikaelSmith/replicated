package shipclient

import (
	"fmt"
	"github.com/pkg/errors"
	channels "github.com/replicatedhq/replicated/gen/go/v1"
	"github.com/replicatedhq/replicated/pkg/graphql"
	"github.com/replicatedhq/replicated/pkg/types"
)

type GraphQLResponseListChannels struct {
	Data   *ShipChannelData   `json:"data,omitempty"`
	Errors []graphql.GQLError `json:"errors,omitempty"`
}

type GraphQLResponseGetChannel struct {
	Data   *ShipGetChannelData `json:"data,omitempty"`
	Errors []graphql.GQLError  `json:"errors,omitempty"`
}

type GraphQLResponseCreateChannel struct {
	Data   *ShipCreateChannelData `json:"data,omitempty"`
	Errors []graphql.GQLError  `json:"errors,omitempty"`
}

type ShipCreateChannelData struct {
	ShipChannel *ShipChannel `json:"createChannel"`
}

type ShipGetChannelData struct {
	ShipChannel *ShipChannel `json:"getChannel"`
}

type ShipChannelData struct {
	ShipChannels []*ShipChannel `json:"getAppChannels"`
}

type ShipChannel struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	CurrentSequence int64  `json:"currentSequence"`
	CurrentVersion  string `json:"currentVersion"`
}

const listChannelsQuery = `
query getAppChannels($appId: ID!) {
  getAppChannels(appId: $appId) {
    id,
    appId,
    name,
    currentVersion,
    currentReleaseDate,
    currentSpec,
    releaseId,
    numReleases,
    description,
    channelIcon
    created,
    updated,
    isDefault,
    isArchived,
    adoptionRate {
      releaseId
      semver
      count
      percent
      totalOnChannel
    }
    releases {
      id
      semver
    }
  }
}`

func (c *GraphQLClient) ListChannels(appID string) ([]types.Channel, error) {
	response := GraphQLResponseListChannels{}

	request := graphql.Request{
		Query: listChannelsQuery,

		Variables: map[string]interface{}{
			"appId": appID,
		},
	}

	if err := c.ExecuteRequest(request, &response); err != nil {
		return nil, err
	}

	channels := make([]types.Channel, 0, 0)
	for _, shipChannel := range response.Data.ShipChannels {
		channel := types.Channel{
			ID:              shipChannel.ID,
			Name:            shipChannel.Name,
			Description:     shipChannel.Description,
			ReleaseSequence: shipChannel.CurrentSequence,
			ReleaseLabel:    shipChannel.CurrentVersion,
		}

		channels = append(channels, channel)
	}

	return channels, nil
}

const createChannelQuery = `
mutation createChannel($appId: String!, $channelName: String!, $description: String) {
  createChannel(appId: $appId, channelName: $channelName, description: $description) {
    id,
    name,
    description,
    channelIcon,
    currentVersion,
    currentReleaseDate,
    numReleases,
    created,
    updated,
    isDefault,
    isArchived
  }
}`

func (c *GraphQLClient) CreateChannel(appID string, name string, description string) (*types.Channel, error) {
	response := GraphQLResponseCreateChannel{}

	request := graphql.Request{
		Query: createChannelQuery,
		Variables: map[string]interface{}{
			"appId":       appID,
			"channelName": name,
			"description": description,
		},
	}

	if err := c.ExecuteRequest(request, &response); err != nil {
		return nil, err
	}
	return &types.Channel{
		ID:              response.Data.ShipChannel.ID,
		Name:            response.Data.ShipChannel.Name,
		Description:     response.Data.ShipChannel.Description,
		ReleaseSequence: response.Data.ShipChannel.CurrentSequence,
		ReleaseLabel:    response.Data.ShipChannel.CurrentVersion,
	}, nil
}

func (c *GraphQLClient) GetChannelByName(appID string, name string, description string, create bool) (*types.Channel, error) {
	allChannels, err := c.ListChannels(appID)
	if err != nil {
		return nil, err
	}

	matchingChannels := make([]*types.Channel, 0)
	for _, channel := range allChannels {
		if channel.ID == name || channel.Name == name {
			matchingChannels = append(matchingChannels, &types.Channel{
				ID:              channel.ID,
				Name:            channel.Name,
				Description:     channel.Description,
				ReleaseSequence: channel.ReleaseSequence,
				ReleaseLabel:    channel.ReleaseLabel,
			})
		}
	}

	if len(matchingChannels) == 0 {
		if create {
			channel, err := c.CreateChannel(appID, name, description)
			if err != nil {
				return nil, errors.Wrapf(err, "create channel %q ", name)
			}
			return channel, nil
		}

		return nil, fmt.Errorf("could not find channel %q", name)
	}

	if len(matchingChannels) > 1 {
		return nil, fmt.Errorf("channel %q is ambiguous, please use channel ID", name)
	}
	return matchingChannels[0], nil
}

var getShipChannel = `
  query getChannel($channelId: ID!) {
    getChannel(channelId: $channelId) {
      id
      appId
      name
      description
      channelIcon
      currentVersion
      currentReleaseDate
      installInstructions
      currentSpec
      numReleases
      adoptionRate {
        releaseId
        semver
        count
        percent
        totalOnChannel
      }
      customers {
        id
        name
        avatar
        actions {
          shipApplyDocker
        }
        installationId
        shipInstallStatus {
          status
          updatedAt
        }
      }
      githubRef {
        owner
        repoFullName
        branch
        path
      }
      extraLintRules
      created
      updated
      isDefault
      isArchived
      releases {
        id
        semver
        spec
        releaseNotes
        created
        updated
        sequence
      }
    }
  }
`

func (c *GraphQLClient) GetChannel(appID string, channelID string) (*channels.AppChannel, []channels.ChannelRelease, error) {
	response := GraphQLResponseGetChannel{}

	request := graphql.Request{
		Query: getShipChannel,
		Variables: map[string]interface{}{
			"appID":     appID,
			"channelId": channelID,
		},
	}
	if err := c.ExecuteRequest(request, &response); err != nil {
		return nil, nil, err
	}

	channelDetail := channels.AppChannel{
		Id:           response.Data.ShipChannel.ID,
		Name:         response.Data.ShipChannel.Name,
		Description:  response.Data.ShipChannel.Description,
		ReleaseLabel: response.Data.ShipChannel.CurrentVersion,
	}
	return &channelDetail, nil, nil
}

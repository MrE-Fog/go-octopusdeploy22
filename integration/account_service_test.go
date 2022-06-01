package integration

import (
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CreateTestAmazonWebServicesAccount(t *testing.T, client *octopusdeploy.Client) octopusdeploy.IAccount {
	if client == nil {
		client = getOctopusClient()
	}
	require.NotNil(t, client)

	accessKey := getRandomName()
	name := getRandomName()
	secretKey := octopusdeploy.NewSensitiveValue(getRandomName())

	account, err := octopusdeploy.NewAmazonWebServicesAccount(name, accessKey, secretKey)
	require.NotNil(t, account)
	require.NoError(t, err)
	require.NoError(t, account.Validate())

	createdAccount, err := client.Accounts.Add(account)
	require.NoError(t, err)
	require.NotNil(t, createdAccount)
	require.NotEmpty(t, createdAccount.GetID())
	require.Equal(t, octopusdeploy.AccountTypeAmazonWebServicesAccount, createdAccount.GetAccountType())
	require.Equal(t, name, createdAccount.GetName())

	return createdAccount
}

func CreateTestAzureServicePrincipalAccount(t *testing.T, client *octopusdeploy.Client) octopusdeploy.IAccount {
	if client == nil {
		client = getOctopusClient()
	}
	require.NotNil(t, client)

	applicationID := uuid.New()
	applicationPassword := octopusdeploy.NewSensitiveValue(getRandomName())
	azureEnvironment := getRandomAzureEnvironment()
	name := getRandomName()
	subscriptionID := uuid.New()
	tenantID := uuid.New()

	account, err := octopusdeploy.NewAzureServicePrincipalAccount(name, subscriptionID, tenantID, applicationID, applicationPassword)

	require.NotNil(t, account)
	require.NoError(t, err)
	require.NoError(t, account.Validate())

	// set Azure environment fields
	if !isEmpty(azureEnvironment.Name) {
		account.AzureEnvironment = azureEnvironment.Name
		account.AuthenticationEndpoint = azureEnvironment.AuthenticationEndpoint
		account.ResourceManagerEndpoint = azureEnvironment.ResourceManagerEndpoint
	}

	require.NoError(t, account.Validate())

	createdAccount, err := client.Accounts.Add(account)
	require.NoError(t, err)
	require.NotNil(t, createdAccount)
	require.NotEmpty(t, createdAccount.GetID())
	require.Equal(t, octopusdeploy.AccountTypeAzureServicePrincipal, createdAccount.GetAccountType())
	require.Equal(t, name, createdAccount.GetName())

	return createdAccount
}

func CreateTestAzureSubscriptionAccount(t *testing.T, client *octopusdeploy.Client) octopusdeploy.IAccount {
	if client == nil {
		client = getOctopusClient()
	}
	require.NotNil(t, client)

	azureEnvironment := getRandomAzureEnvironment()
	name := getRandomName()
	subscriptionID := uuid.New()

	account, err := octopusdeploy.NewAzureSubscriptionAccount(name, subscriptionID)

	require.NotNil(t, account)
	require.NoError(t, err)
	require.NoError(t, account.Validate())

	// set Azure environment fields
	if !isEmpty(azureEnvironment.Name) {
		account.AzureEnvironment = azureEnvironment.Name
		account.ManagementEndpoint = azureEnvironment.ManagementEndpoint
		account.StorageEndpointSuffix = azureEnvironment.StorageEndpointSuffix
	}

	require.NoError(t, account.Validate())

	resource, err := client.Accounts.Add(account)
	require.NoError(t, err)
	require.NotNil(t, resource)

	return resource
}

func CreateTestSSHKeyAccount(t *testing.T, client *octopusdeploy.Client) octopusdeploy.IAccount {
	if client == nil {
		client = getOctopusClient()
	}
	require.NotNil(t, client)

	name := getRandomName()
	username := getRandomName()
	privateKeyFile := octopusdeploy.NewSensitiveValue(getRandomName())

	account, err := octopusdeploy.NewSSHKeyAccount(name, username, privateKeyFile)

	require.NotNil(t, account)
	require.NoError(t, err)
	require.NoError(t, account.Validate())

	resource, err := client.Accounts.Add(account)
	require.NoError(t, err)
	require.NotNil(t, resource)

	return resource
}

func CreateTestTokenAccount(t *testing.T, client *octopusdeploy.Client) octopusdeploy.IAccount {
	if client == nil {
		client = getOctopusClient()
	}
	require.NotNil(t, client)

	name := getRandomName()
	token := octopusdeploy.NewSensitiveValue(getRandomName())

	account, err := octopusdeploy.NewTokenAccount(name, token)

	require.NotNil(t, account)
	require.NoError(t, err)
	require.NoError(t, account.Validate())

	resource, err := client.Accounts.Add(account)
	require.NoError(t, err)
	require.NotNil(t, resource)

	return resource
}

func CreateTestUsernamePasswordAccount(t *testing.T, client *octopusdeploy.Client) octopusdeploy.IAccount {
	if client == nil {
		client = getOctopusClient()
	}
	require.NotNil(t, client)

	name := getRandomName()

	account, err := octopusdeploy.NewUsernamePasswordAccount(name)
	require.NotNil(t, account)
	require.NoError(t, err)
	require.NoError(t, account.Validate())

	resource, err := client.Accounts.Add(account)
	require.NoError(t, err)
	require.NotNil(t, resource)

	return resource
}

func DeleteTestAccount(t *testing.T, client *octopusdeploy.Client, account octopusdeploy.IAccount) {
	if client == nil {
		client = getOctopusClient()
	}
	require.NotNil(t, client)

	err := client.Accounts.DeleteByID(account.GetID())
	assert.NoError(t, err)

	// verify the delete operation was successful
	deletedAccount, err := client.Accounts.GetByID(account.GetID())
	assert.Error(t, err)
	assert.Nil(t, deletedAccount)
}

func CompareAccounts(t *testing.T, expected octopusdeploy.IAccount, actual octopusdeploy.IAccount) {
	// equality cannot be determined through a direct comparison (below)
	// because APIs like GetByPartialName do not include the fields,
	// LastModifiedBy and LastModifiedOn
	//
	// assert.EqualValues(expected, actual)
	//
	// this statement (above) is expected to succeed, but it fails due to these
	// missing fields

	// type check
	assert.IsType(t, expected, actual)

	// IResource
	assert.Equal(t, expected.GetID(), actual.GetID())
	assert.True(t, IsEqualLinks(expected.GetLinks(), actual.GetLinks()))

	// IHasName
	assert.Equal(t, expected.GetName(), actual.GetName())

	// IHasSpace
	assert.Equal(t, expected.GetSpaceID(), actual.GetSpaceID())

	// IAccount
	require.Equal(t, expected.GetAccountType(), actual.GetAccountType())
	assert.Equal(t, expected.GetDescription(), actual.GetDescription())
	assert.Equal(t, expected.GetEnvironmentIDs(), actual.GetEnvironmentIDs())
	assert.Equal(t, expected.GetTenantedDeploymentMode(), actual.GetTenantedDeploymentMode())
	assert.Equal(t, expected.GetTenantIDs(), actual.GetTenantIDs())
	assert.Equal(t, expected.GetTenantTags(), actual.GetTenantTags())

	// the account types are equal -- therefore, it is assumed safe to perform a type assertion

	switch actual.GetAccountType() {
	case octopusdeploy.AccountTypeAmazonWebServicesAccount:
		expectedToCompare := expected.(*octopusdeploy.AmazonWebServicesAccount)
		actualToCompare := actual.(*octopusdeploy.AmazonWebServicesAccount)
		assert.Equal(t, expectedToCompare.AccessKey, actualToCompare.AccessKey)
		assert.Equal(t, expectedToCompare.SecretKey, actualToCompare.SecretKey)
	case octopusdeploy.AccountTypeAzureServicePrincipal:
		expectedToCompare := expected.(*octopusdeploy.AzureServicePrincipalAccount)
		actualToCompare := actual.(*octopusdeploy.AzureServicePrincipalAccount)
		assert.Equal(t, expectedToCompare.ApplicationID, actualToCompare.ApplicationID)
		assert.Equal(t, expectedToCompare.ApplicationPassword, actualToCompare.ApplicationPassword)
		assert.Equal(t, expectedToCompare.AuthenticationEndpoint, actualToCompare.AuthenticationEndpoint)
		assert.Equal(t, expectedToCompare.AzureEnvironment, actualToCompare.AzureEnvironment)
		assert.Equal(t, expectedToCompare.ResourceManagerEndpoint, actualToCompare.ResourceManagerEndpoint)
		assert.Equal(t, expectedToCompare.SubscriptionID, actualToCompare.SubscriptionID)
		assert.Equal(t, expectedToCompare.TenantID, actualToCompare.TenantID)
	case octopusdeploy.AccountTypeAzureSubscription:
		expectedToCompare := expected.(*octopusdeploy.AzureSubscriptionAccount)
		actualToCompare := actual.(*octopusdeploy.AzureSubscriptionAccount)
		assert.Equal(t, expectedToCompare.AzureEnvironment, actualToCompare.AzureEnvironment)
		assert.Equal(t, expectedToCompare.CertificateBytes, actualToCompare.CertificateBytes)
		assert.Equal(t, expectedToCompare.CertificateThumbprint, actualToCompare.CertificateThumbprint)
		assert.Equal(t, expectedToCompare.ManagementEndpoint, actualToCompare.ManagementEndpoint)
		assert.Equal(t, expectedToCompare.StorageEndpointSuffix, actualToCompare.StorageEndpointSuffix)
		assert.Equal(t, expectedToCompare.SubscriptionID, actualToCompare.SubscriptionID)
	case octopusdeploy.AccountTypeGoogleCloudPlatformAccount:
		expectedToCompare := expected.(*octopusdeploy.GoogleCloudPlatformAccount)
		actualToCompare := actual.(*octopusdeploy.GoogleCloudPlatformAccount)
		assert.Equal(t, expectedToCompare.JsonKey, actualToCompare.JsonKey)
	case octopusdeploy.AccountTypeSSHKeyPair:
		expectedToCompare := expected.(*octopusdeploy.SSHKeyAccount)
		actualToCompare := actual.(*octopusdeploy.SSHKeyAccount)
		assert.Equal(t, expectedToCompare.PrivateKeyFile, actualToCompare.PrivateKeyFile)
		assert.Equal(t, expectedToCompare.PrivateKeyPassphrase, actualToCompare.PrivateKeyPassphrase)
		assert.Equal(t, expectedToCompare.Username, actualToCompare.Username)
	case octopusdeploy.AccountTypeToken:
		expectedToCompare := expected.(*octopusdeploy.TokenAccount)
		actualToCompare := actual.(*octopusdeploy.TokenAccount)
		assert.Equal(t, expectedToCompare.Token, actualToCompare.Token)
	case octopusdeploy.AccountTypeUsernamePassword:
		expectedToCompare := expected.(*octopusdeploy.UsernamePasswordAccount)
		actualToCompare := actual.(*octopusdeploy.UsernamePasswordAccount)
		assert.Equal(t, expectedToCompare.Password, actualToCompare.Password)
		assert.Equal(t, expectedToCompare.Username, actualToCompare.Username)
	}
}

func UpdateAccount(t *testing.T, client *octopusdeploy.Client, account octopusdeploy.IAccount) octopusdeploy.IAccount {
	if client == nil {
		client = getOctopusClient()
	}
	require.NotNil(t, client)

	updatedAccount, err := client.Accounts.Update(account)
	require.NoError(t, err)
	require.NotNil(t, updatedAccount)

	return updatedAccount
}

func ValidateAccount(t *testing.T, account octopusdeploy.IAccount) {
	require.NoError(t, account.Validate())
	require.NotEmpty(t, account.GetLinks())
	require.NotEmpty(t, account.GetID())
	require.NotEmpty(t, account.GetModifiedBy())
	require.NotEmpty(t, account.GetModifiedOn())
	require.NotEmpty(t, account.GetSpaceID())

	// TODO: validate other fields/methods
}

func TestAccountServiceDeleteAll(t *testing.T) {
	client := getOctopusClient()
	require.NotNil(t, client)

	accounts, err := client.Accounts.GetAll()
	require.NoError(t, err)
	require.NotNil(t, accounts)

	for _, account := range accounts {
		defer DeleteTestAccount(t, client, account)
	}
}

func TestAccountServiceAddGetDelete(t *testing.T) {
	client := getOctopusClient()
	require.NotNil(t, client)

	amazonWebServicesAccount := CreateTestAmazonWebServicesAccount(t, client)
	ValidateAccount(t, amazonWebServicesAccount)
	defer DeleteTestAccount(t, client, amazonWebServicesAccount)

	azureServicePrincipalAccount := CreateTestAzureServicePrincipalAccount(t, client)
	ValidateAccount(t, azureServicePrincipalAccount)
	defer DeleteTestAccount(t, client, azureServicePrincipalAccount)

	azureSubscriptionAccount := CreateTestAzureSubscriptionAccount(t, client)
	ValidateAccount(t, azureSubscriptionAccount)
	defer DeleteTestAccount(t, client, azureSubscriptionAccount)

	sshKeyAccount := CreateTestSSHKeyAccount(t, client)
	ValidateAccount(t, sshKeyAccount)
	defer DeleteTestAccount(t, client, sshKeyAccount)

	tokenAccount := CreateTestTokenAccount(t, client)
	ValidateAccount(t, tokenAccount)
	defer DeleteTestAccount(t, client, tokenAccount)

	usernamePasswordAccount := CreateTestUsernamePasswordAccount(t, client)
	ValidateAccount(t, usernamePasswordAccount)
	defer DeleteTestAccount(t, client, usernamePasswordAccount)

	accounts, err := client.Accounts.GetAll()
	require.NoError(t, err)
	require.NotNil(t, accounts)

	for _, account := range accounts {
		name := account.GetName()
		query := octopusdeploy.AccountsQuery{
			PartialName: name,
			Take:        1,
		}

		namedAccounts, err := client.Accounts.Get(query)
		require.NoError(t, err)
		require.NotNil(t, namedAccounts)

		for _, namedAccount := range namedAccounts.Items {
			accountToCompare, err := client.Accounts.GetByID(namedAccount.GetID())
			require.NoError(t, err)
			require.NotNil(t, accountToCompare)
			CompareAccounts(t, namedAccount, accountToCompare)
		}

		accountToCompare, err := client.Accounts.GetByID(account.GetID())
		require.NoError(t, err)
		require.NotNil(t, accountToCompare)

		for _, namedAccount := range namedAccounts.Items {
			accountToCompare, err := client.Accounts.GetByID(namedAccount.GetID())
			require.NoError(t, err)
			require.NotNil(t, accountToCompare)
			CompareAccounts(t, namedAccount, accountToCompare)
		}

		accountUsages, err := client.Accounts.GetUsages(account)
		require.NoError(t, err)
		require.NotNil(t, accountUsages)
	}

	accountTypes := []octopusdeploy.AccountType{
		octopusdeploy.AccountTypeNone,
		octopusdeploy.AccountTypeUsernamePassword,
		octopusdeploy.AccountTypeSSHKeyPair,
		octopusdeploy.AccountTypeAzureSubscription,
		octopusdeploy.AccountTypeAzureServicePrincipal,
		octopusdeploy.AccountTypeAmazonWebServicesAccount,
		octopusdeploy.AccountTypeToken,
	}

	for _, accountType := range accountTypes {
		query := octopusdeploy.AccountsQuery{AccountType: accountType}
		accounts, err := client.Accounts.Get(query)
		require.NoError(t, err)

		for _, account := range accounts.Items {
			accountToCompare, err := client.Accounts.GetByID(account.GetID())
			require.NoError(t, err)
			CompareAccounts(t, account, accountToCompare)
		}
	}
}

func TestAccountServiceGetByID(t *testing.T) {
	client := getOctopusClient()
	require.NotNil(t, client)

	id := getRandomName()
	resource, err := client.Accounts.GetByID(id)
	require.Error(t, err)
	require.Nil(t, resource)

	apiError := err.(*octopusdeploy.APIError)
	assert.Equal(t, 404, apiError.StatusCode)

	accounts, err := client.Accounts.GetAll()
	require.NoError(t, err)
	require.NotNil(t, accounts)

	for _, account := range accounts {
		accountToCompare, err := client.Accounts.GetByID(account.GetID())
		require.NoError(t, err)
		CompareAccounts(t, account, accountToCompare)
	}
}

func TestAccountServiceGetByIDs(t *testing.T) {
	client := getOctopusClient()
	require.NotNil(t, client)

	accounts, err := client.Accounts.GetAll()
	require.NoError(t, err)
	require.NotNil(t, accounts)

	ids := []string{}
	for _, account := range accounts {
		ids = append(ids, account.GetID())
	}

	query := octopusdeploy.AccountsQuery{IDs: ids}

	accountsByIDs, err := client.Accounts.Get(query)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(accounts), len(accountsByIDs.Items))
}

func TestAccountServiceTokenAccounts(t *testing.T) {
	client := getOctopusClient()
	require.NotNil(t, client)

	for i := 0; i < 10; i++ {
		tokenAccount := CreateTestTokenAccount(t, client)
		ValidateAccount(t, tokenAccount)
		defer DeleteTestAccount(t, client, tokenAccount)
	}

	accounts, err := client.Accounts.GetAll()
	require.NoError(t, err)
	require.NotNil(t, accounts)

	for _, account := range accounts {
		accountToCompare, err := client.Accounts.GetByID(account.GetID())
		require.NoError(t, err)
		CompareAccounts(t, account, accountToCompare)
	}
}

func TestAccountServiceUpdate(t *testing.T) {
	client := getOctopusClient()
	require.NotNil(t, client)

	expected := CreateTestAzureServicePrincipalAccount(t, client)
	actual := UpdateAccount(t, client, expected)
	CompareAccounts(t, expected, actual)
	ValidateAccount(t, actual)
	defer DeleteTestAccount(t, client, expected)

	expected = CreateTestAzureSubscriptionAccount(t, client)
	expected.SetName(getRandomName())
	actual = UpdateAccount(t, client, expected)
	CompareAccounts(t, expected, actual)
	ValidateAccount(t, actual)
	defer DeleteTestAccount(t, client, expected)

	expected = CreateTestSSHKeyAccount(t, client)
	expected.SetName(getRandomName())
	actual = UpdateAccount(t, client, expected)
	CompareAccounts(t, expected, actual)
	ValidateAccount(t, actual)
	defer DeleteTestAccount(t, client, expected)

	expected = CreateTestTokenAccount(t, client)
	expected.SetName(getRandomName())
	actual = UpdateAccount(t, client, expected)
	CompareAccounts(t, expected, actual)
	ValidateAccount(t, actual)
	defer DeleteTestAccount(t, client, expected)

	expected = CreateTestUsernamePasswordAccount(t, client)
	expected.SetName(getRandomName())
	actual = UpdateAccount(t, client, expected)
	CompareAccounts(t, expected, actual)
	ValidateAccount(t, actual)
	defer DeleteTestAccount(t, client, expected)
}

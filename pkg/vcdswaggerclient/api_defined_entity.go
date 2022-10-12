/*
 * VMware Cloud Director OpenAPI
 *
 * VMware Cloud Director OpenAPI is a new API that is defined using the OpenAPI standards.<br/> This ReSTful API borrows some elements of the legacy VMware Cloud Director API and establishes new patterns for use as described below. <h4>Authentication</h4> Authentication and Authorization schemes are the same as those for the legacy APIs. You can authenticate using the JWT token via the <code>Authorization</code> header or specifying a session using <code>x-vcloud-authorization</code> (The latter form is deprecated). <h4>Operation Patterns</h4> This API follows the following general guidelines to establish a consistent CRUD pattern: <table> <tr>   <th>Operation</th><th>Description</th><th>Response Code</th><th>Response Content</th> </tr><tr>   <td>GET /items<td>Returns a paginated list of items<td>200<td>Response will include Navigational links to the items in the list. </tr><tr>   <td>POST /items<td>Returns newly created item<td>201<td>Content-Location header links to the newly created item </tr><tr>   <td>GET /items/urn<td>Returns an individual item<td>200<td>A single item using same data type as that included in list above </tr><tr>   <td>PUT /items/urn<td>Updates an individual item<td>200<td>Updated view of the item is returned </tr><tr>   <td>DELETE /items/urn<td>Deletes the item<td>204<td>No content is returned. </tr> </table> <h5>Asynchronous operations</h5> Asynchronous operations are determined by the server. In those cases, instead of responding as described above, the server responds with an HTTP Response code 202 and an empty body. The tracking task (which is the same task as all legacy API operations use) is linked via the URI provided in the <code>Location</code> header.<br/> All API calls can choose to service a request asynchronously or synchronously as determined by the server upon interpreting the request. Operations that choose to exhibit this dual behavior will have both options documented by specifying both response code(s) below. The caller must be prepared to handle responses to such API calls by inspecting the HTTP Response code. <h5>Error Conditions</h5> <b>All</b> operations report errors using the following error reporting rules: <ul>   <li>400: Bad Request - In event of bad request due to incorrect data or other user error</li>   <li>401: Bad Request - If user is unauthenticated or their session has expired</li>   <li>403: Forbidden - If the user is not authorized or the entity does not exist</li> </ul> <h4>OpenAPI Design Concepts and Principles</h4> <ul>   <li>IDs are full Uniform Resource Names (URNs).</li>   <li>OpenAPI's <code>Content-Type</code> is always <code>application/json</code></li>   <li>REST links are in the Link header.</li>   <ul>     <li>Multiple relationships for any link are represented by multiple values in a space-separated list.</li>     <li>Links have a custom VMware Cloud Director-specific &quot;model&quot; attribute that hints at the applicable data         type for the links.</li>     <li>title + rel + model attributes evaluates to a unique link.</li>     <li>Links follow Hypermedia as the Engine of Application State (HATEOAS) principles. Links are present if         certain operations are present and permitted for the user&quot;s current role and the state of the         referred entities.</li>   </ul>   <li>APIs follow a flat structure relying on cross-referencing other entities instead of the navigational style       used by the legacy VMware Cloud Director APIs.</li>   <li>Most endpoints that return a list support filtering and sorting similar to the query service in the legacy       VMware Cloud Director APIs.</li>   <li>Accept header must be included to specify the API version for the request similar to calls to existing legacy       VMware Cloud Director APIs.</li>   <li>Each feature has a version in the path element present in its URL.<br/>       <b>Note</b> API URL's without a version in their paths must be considered experimental.</li> </ul>
 *
 * API version: 36.0
 * Contact: https://code.vmware.com/support
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Linger please
var (
	_ context.Context
)

type DefinedEntityApiService service

/*
DefinedEntityApiService Creates a defined entity based on the entity type (URN).
Creates a defined entity based on the entity type (URN).
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param entity
 * @param id
 * @param optional nil or *DefinedEntityApiCreateDefinedEntityOpts - Optional Parameters:
     * @param "InvokeHooks" (optional.Interface of interface{}) -  Only users with Admin FullControl access to the Entity Type can pass this parameter. The default value is &#39;true&#39;.


*/

type DefinedEntityApiCreateDefinedEntityOpts struct {
	InvokeHooks optional.Interface
}

func (a *DefinedEntityApiService) CreateDefinedEntity(ctx context.Context, entity DefinedEntity, id string, orgID string, localVarOptionals *DefinedEntityApiCreateDefinedEntityOpts) (*http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Post")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte

	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/1.0.0/entityTypes/{id}"
	localVarPath = strings.Replace(localVarPath, "{"+"id"+"}", fmt.Sprintf("%v", id), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarOptionals != nil && localVarOptionals.InvokeHooks.IsSet() {
		localVarQueryParams.Add("invokeHooks", parameterToString(localVarOptionals.InvokeHooks.Value(), ""))
	}
	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"*_/_*;version=36.0"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	// body params
	localVarPostBody = &entity
	if ctx != nil {
		// API Key Authentication
		if auth, ok := ctx.Value(ContextAPIKey).(APIKey); ok {
			var key string
			if auth.Prefix != "" {
				key = auth.Prefix + " " + auth.Key
			} else {
				key = auth.Key
			}
			localVarHeaderParams["Authorization"] = key

		}
	}
	if orgID != "" {
		localVarHeaderParams[TenantContextHeader] = orgID
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarHttpResponse, err
	}


	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}

		if localVarHttpResponse.StatusCode == 400 {
			var v ModelError
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
			if err != nil {
				newErr.error = err.Error()
				return localVarHttpResponse, newErr
			}
			newErr.model = v
			return localVarHttpResponse, newErr
		}

		return localVarHttpResponse, newErr
	}

	return localVarHttpResponse, nil
}

/*
DefinedEntityApiService Deletes the defined entity with the unique identifier (URN)
Deletes the defined entity with the unique identifier (URN)
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param id
 * @param optional nil or *DefinedEntityApiDeleteDefinedEntityOpts - Optional Parameters:
     * @param "InvokeHooks" (optional.Interface of interface{}) -  Only users with Admin FullControl access to the Entity Type can pass this parameter. The default value is &#39;true&#39;.


*/

type DefinedEntityApiDeleteDefinedEntityOpts struct {
	InvokeHooks optional.Interface
}

func (a *DefinedEntityApiService) DeleteDefinedEntity(ctx context.Context, id string, orgID string, localVarOptionals *DefinedEntityApiDeleteDefinedEntityOpts) (*http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Delete")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/1.0.0/entities/{id}"
	localVarPath = strings.Replace(localVarPath, "{"+"id"+"}", fmt.Sprintf("%v", id), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarOptionals != nil && localVarOptionals.InvokeHooks.IsSet() {
		localVarQueryParams.Add("invokeHooks", parameterToString(localVarOptionals.InvokeHooks.Value(), ""))
	}
	// to determine the Content-Type header
	localVarHttpContentTypes := []string{}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"*_/_*;version=36.0"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	if ctx != nil {
		// API Key Authentication
		if auth, ok := ctx.Value(ContextAPIKey).(APIKey); ok {
			var key string
			if auth.Prefix != "" {
				key = auth.Prefix + " " + auth.Key
			} else {
				key = auth.Key
			}
			localVarHeaderParams["Authorization"] = key

		}
	}
	if orgID != "" {
		localVarHeaderParams[TenantContextHeader] = orgID
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarHttpResponse, err
	}


	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}

		return localVarHttpResponse, newErr
	}

	return localVarHttpResponse, nil
}

/*
DefinedEntityApiService Gets the collection of defined entities for the vCD-defined type with the specified vendor, nss and version.
Gets the collection of defined entities for the vCD-defined type with the specified vendor, nss and version. The version can act as a wildcard. If only &#39;1&#39; is specified as the version, all entity types with a major version of &#39;1&#39; will be matched (e.g. 1.0.0, 1.1.2). If &#39;1.0&#39; is specified, all entity types with a major version of &#39;1&#39; and a minor version of &#39;0&#39; will be included (e.g. 1.0.0, 1.0.1). If the full semver is specified, then no search will be performed.
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param vendor
 * @param nss
 * @param version
 * @param page Page to fetch, zero offset.
 * @param pageSize Results per page to fetch.
 * @param optional nil or *DefinedEntityApiGetDefinedEntitiesByEntityTypeOpts - Optional Parameters:
     * @param "Filter" (optional.String) -  Filter for a query.  FIQL format.
     * @param "SortAsc" (optional.String) -  Field to use for ascending sort
     * @param "SortDesc" (optional.String) -  Field to use for descending sort

@return DefinedEntities
*/

type DefinedEntityApiGetDefinedEntitiesByEntityTypeOpts struct {
	Filter optional.String
	SortAsc optional.String
	SortDesc optional.String
}

func (a *DefinedEntityApiService) GetDefinedEntitiesByEntityType(ctx context.Context, vendor string, nss string, version string, orgID string, page int32, pageSize int32, localVarOptionals *DefinedEntityApiGetDefinedEntitiesByEntityTypeOpts) (DefinedEntities, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue DefinedEntities
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/1.0.0/entities/types/{vendor}/{nss}/{version}"
	localVarPath = strings.Replace(localVarPath, "{"+"vendor"+"}", fmt.Sprintf("%v", vendor), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"nss"+"}", fmt.Sprintf("%v", nss), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"version"+"}", fmt.Sprintf("%v", version), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if strlen(vendor) < 1 {
		return localVarReturnValue, nil, reportError("vendor must have at least 1 elements")
	}
	if strlen(nss) < 1 {
		return localVarReturnValue, nil, reportError("nss must have at least 1 elements")
	}
	if strlen(version) < 1 {
		return localVarReturnValue, nil, reportError("version must have at least 1 elements")
	}
	if page < 1 {
		return localVarReturnValue, nil, reportError("page must be greater than 1")
	}
	if pageSize < 0 {
		return localVarReturnValue, nil, reportError("pageSize must be greater than 0")
	}
	if pageSize > 128 {
		return localVarReturnValue, nil, reportError("pageSize must be less than 128")
	}

	if localVarOptionals != nil && localVarOptionals.Filter.IsSet() {
		localVarQueryParams.Add("filter", parameterToString(localVarOptionals.Filter.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.SortAsc.IsSet() {
		localVarQueryParams.Add("sortAsc", parameterToString(localVarOptionals.SortAsc.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.SortDesc.IsSet() {
		localVarQueryParams.Add("sortDesc", parameterToString(localVarOptionals.SortDesc.Value(), ""))
	}
	localVarQueryParams.Add("page", parameterToString(page, ""))
	localVarQueryParams.Add("pageSize", parameterToString(pageSize, ""))
	// to determine the Content-Type header
	localVarHttpContentTypes := []string{}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"application/json;version=36.0"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	if ctx != nil {
		// API Key Authentication
		if auth, ok := ctx.Value(ContextAPIKey).(APIKey); ok {
			var key string
			if auth.Prefix != "" {
				key = auth.Prefix + " " + auth.Key
			} else {
				key = auth.Key
			}
			localVarHeaderParams["Authorization"] = key

		}
	}
	if orgID != "" {
		localVarHeaderParams[TenantContextHeader] = orgID
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}

		if localVarHttpResponse.StatusCode == 200 {
			var v DefinedEntities
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHttpResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHttpResponse, newErr
		}

		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}

/*
DefinedEntityApiService Gets the collection of defined entities for the vCD-defined interface with the specified vendor, nss and version
Gets the collection of defined entities for the vCD-defined interface with the specified vendor, nss and version. The version can act as a wildcard. If only &#39;1&#39; is specified as the version, all entity types with a major version of &#39;1&#39; will be matched (e.g. 1.0.0, 1.1.2). If &#39;1.0&#39; is specified, all entity types with a major version of &#39;1&#39; and a minor version of &#39;0&#39; will be included (e.g. 1.0.0, 1.0.1). If the full semver is specified, then no search will be performed.
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param vendor
 * @param nss
 * @param version
 * @param page Page to fetch, zero offset.
 * @param pageSize Results per page to fetch.
 * @param optional nil or *DefinedEntityApiGetDefinedEntitiesByInterfaceOpts - Optional Parameters:
     * @param "Filter" (optional.String) -  Filter for a query.  FIQL format.
     * @param "SortAsc" (optional.String) -  Field to use for ascending sort
     * @param "SortDesc" (optional.String) -  Field to use for descending sort

@return DefinedEntities
*/

type DefinedEntityApiGetDefinedEntitiesByInterfaceOpts struct {
	Filter optional.String
	SortAsc optional.String
	SortDesc optional.String
}

func (a *DefinedEntityApiService) GetDefinedEntitiesByInterface(ctx context.Context, vendor string, nss string, version string, orgID string, page int32, pageSize int32, localVarOptionals *DefinedEntityApiGetDefinedEntitiesByInterfaceOpts) (DefinedEntities, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue DefinedEntities
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/1.0.0/entities/interfaces/{vendor}/{nss}/{version}"
	localVarPath = strings.Replace(localVarPath, "{"+"vendor"+"}", fmt.Sprintf("%v", vendor), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"nss"+"}", fmt.Sprintf("%v", nss), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"version"+"}", fmt.Sprintf("%v", version), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if strlen(vendor) < 1 {
		return localVarReturnValue, nil, reportError("vendor must have at least 1 elements")
	}
	if strlen(nss) < 1 {
		return localVarReturnValue, nil, reportError("nss must have at least 1 elements")
	}
	if strlen(version) < 1 {
		return localVarReturnValue, nil, reportError("version must have at least 1 elements")
	}
	if page < 1 {
		return localVarReturnValue, nil, reportError("page must be greater than 1")
	}
	if pageSize < 0 {
		return localVarReturnValue, nil, reportError("pageSize must be greater than 0")
	}
	if pageSize > 128 {
		return localVarReturnValue, nil, reportError("pageSize must be less than 128")
	}

	if localVarOptionals != nil && localVarOptionals.Filter.IsSet() {
		localVarQueryParams.Add("filter", parameterToString(localVarOptionals.Filter.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.SortAsc.IsSet() {
		localVarQueryParams.Add("sortAsc", parameterToString(localVarOptionals.SortAsc.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.SortDesc.IsSet() {
		localVarQueryParams.Add("sortDesc", parameterToString(localVarOptionals.SortDesc.Value(), ""))
	}
	localVarQueryParams.Add("page", parameterToString(page, ""))
	localVarQueryParams.Add("pageSize", parameterToString(pageSize, ""))
	// to determine the Content-Type header
	localVarHttpContentTypes := []string{}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"application/json;version=36.0"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	if ctx != nil {
		// API Key Authentication
		if auth, ok := ctx.Value(ContextAPIKey).(APIKey); ok {
			var key string
			if auth.Prefix != "" {
				key = auth.Prefix + " " + auth.Key
			} else {
				key = auth.Key
			}
			localVarHeaderParams["Authorization"] = key

		}
	}
	if orgID != "" {
		localVarHeaderParams[TenantContextHeader] = orgID
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}

		if localVarHttpResponse.StatusCode == 200 {
			var v DefinedEntities
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHttpResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHttpResponse, newErr
		}

		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}

/*
DefinedEntityApiService Gets the defined entity with the unique identifier (URN)
Gets the defined entity with the unique identifier (URN)
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param id

@return DefinedEntity
*/
func (a *DefinedEntityApiService) GetDefinedEntity(ctx context.Context, id string, orgID string) (DefinedEntity, *http.Response, string, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue DefinedEntity
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/1.0.0/entities/{id}"
	localVarPath = strings.Replace(localVarPath, "{"+"id"+"}", fmt.Sprintf("%v", id), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"application/json;version=36.0"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	if ctx != nil {
		// API Key Authentication
		if auth, ok := ctx.Value(ContextAPIKey).(APIKey); ok {
			var key string
			if auth.Prefix != "" {
				key = auth.Prefix + " " + auth.Key
			} else {
				key = auth.Key
			}
			localVarHeaderParams["Authorization"] = key

		}
	}
	if orgID != "" {
		localVarHeaderParams[TenantContextHeader] = orgID
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, "", err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, "", err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, "", err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		etag := localVarHttpResponse.Header.Get("Etag")
		return localVarReturnValue, localVarHttpResponse, etag, err
	} else {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}

		return localVarReturnValue, localVarHttpResponse, "", newErr
	}
}

/*
DefinedEntityApiService Validates the defined entity against the entity type schema.
Validates the defined entity against the entity type schema. If the validation is successful, the entity will transition to a \&quot;RESOLVED\&quot; state. Otherwise, it will transition to an \&quot;ERROR\&quot; state.
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param id

@return EntityState
*/
func (a *DefinedEntityApiService) ResolveDefinedEntity(ctx context.Context, id string, orgID string) (EntityState, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Post")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue EntityState
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/1.0.0/entities/{id}/resolve"
	localVarPath = strings.Replace(localVarPath, "{"+"id"+"}", fmt.Sprintf("%v", id), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"application/json;version=36.0"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	if ctx != nil {
		// API Key Authentication
		if auth, ok := ctx.Value(ContextAPIKey).(APIKey); ok {
			var key string
			if auth.Prefix != "" {
				key = auth.Prefix + " " + auth.Key
			} else {
				key = auth.Key
			}
			localVarHeaderParams["Authorization"] = key

		}
	}
	if orgID != "" {
		localVarHeaderParams[TenantContextHeader] = orgID
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}

		if localVarHttpResponse.StatusCode == 200 {
			var v EntityState
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHttpResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHttpResponse, newErr
		}

		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}

/*
DefinedEntityApiService Updates the defined entity with the unique identifier (URN)
Update the defined entity with the unique identifier (URN)
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param entity
 * @param id
 * @param optional nil or *DefinedEntityApiUpdateDefinedEntityOpts - Optional Parameters:
     * @param "InvokeHooks" (optional.Interface of interface{}) -  Only users with Admin FullControl access to the Entity Type can pass this parameter. The default value is &#39;true&#39;.

@return DefinedEntity
*/

type DefinedEntityApiUpdateDefinedEntityOpts struct {
	InvokeHooks optional.Interface
}

func (a *DefinedEntityApiService) UpdateDefinedEntity(ctx context.Context, entity DefinedEntity, etag string, id string, orgID string, localVarOptionals *DefinedEntityApiUpdateDefinedEntityOpts) (DefinedEntity, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Put")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue DefinedEntity
	)
	if etag == "" {
		return localVarReturnValue, nil, fmt.Errorf("etag is empty when updating defined entity")
	}

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/1.0.0/entities/{id}"
	localVarPath = strings.Replace(localVarPath, "{"+"id"+"}", fmt.Sprintf("%v", id), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	localVarHeaderParams["If-Match"] = etag

	if localVarOptionals != nil && localVarOptionals.InvokeHooks.IsSet() {
		localVarQueryParams.Add("invokeHooks", parameterToString(localVarOptionals.InvokeHooks.Value(), ""))
	}
	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"application/json;version=36.0"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	// body params
	localVarPostBody = &entity
	if ctx != nil {
		// API Key Authentication
		if auth, ok := ctx.Value(ContextAPIKey).(APIKey); ok {
			var key string
			if auth.Prefix != "" {
				key = auth.Prefix + " " + auth.Key
			} else {
				key = auth.Key
			}
			localVarHeaderParams["Authorization"] = key

		}
	}
	if orgID != "" {
		localVarHeaderParams[TenantContextHeader] = orgID
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}

		if localVarHttpResponse.StatusCode == 400 {
			var v ModelError
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHttpResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHttpResponse, newErr
		}

		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}

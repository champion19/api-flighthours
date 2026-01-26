package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type HATEOASResource struct {
	Links []Link `json:"_links"`
}

func GetBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host
}

func SetLocationHeader(c *gin.Context, baseURL, resource, resourceID string) {
	locationURL := BuildResourceURL(baseURL, resource, resourceID)
	c.Header("Location", locationURL)
}

func BuildResourceURL(baseURL, resource, resourceID string) string {
	return fmt.Sprintf("%s/flighthours/api/v1/%s/%s", baseURL, resource, resourceID)
}

func BuildCollectionURL(baseURL, resource string) string {
	return fmt.Sprintf("%s/flighthours/api/v1/%s", baseURL, resource)
}

func BuildResourceLinks(baseURL, resource, resourceID string) []Link {
	resourceURL := BuildResourceURL(baseURL, resource, resourceID)
	collectionURL := BuildCollectionURL(baseURL, resource)

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL,
			Rel:    "update",
			Method: "PUT",
		},
		{
			Href:   resourceURL,
			Rel:    "delete",
			Method: "DELETE",
		},
		{
			Href:   collectionURL,
			Rel:    "collection",
			Method: "GET",
		},
	}
}

func BuildAccountLinks(baseURL string, accountID string) []Link {
	return BuildResourceLinks(baseURL, "accounts", accountID)
}

func BuildMessageLinks(baseURL string, messageID string) []Link {
	return BuildResourceLinks(baseURL, "messages", messageID)
}

func BuildMessageCreatedLinks(baseURL string, messageID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "messages", messageID)
	collectionURL := BuildCollectionURL(baseURL, "messages")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL,
			Rel:    "update",
			Method: "PUT",
		},
		{
			Href:   resourceURL,
			Rel:    "delete",
			Method: "DELETE",
		},
		{
			Href:   collectionURL,
			Rel:    "list",
			Method: "GET",
		},
	}
}


func BuildMessageUpdatedLinks(baseURL string, messageID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "messages", messageID)
	collectionURL := BuildCollectionURL(baseURL, "messages")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL,
			Rel:    "delete",
			Method: "DELETE",
		},
		{
			Href:   collectionURL,
			Rel:    "list",
			Method: "GET",
		},
	}
}

func BuildMessageListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "messages")

	return []Link{
		{
			Href:   collectionURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   collectionURL,
			Rel:    "create",
			Method: "POST",
		},
	}
}






func BuildEmployeeLinks(baseURL string, employeeID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "employees", employeeID)
	collectionURL := BuildCollectionURL(baseURL, "employees")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL,
			Rel:    "update",
			Method: "PUT",
		},
		{
			Href:   resourceURL,
			Rel:    "delete",
			Method: "DELETE",
		},
		{
			Href:   collectionURL,
			Rel:    "collection",
			Method: "GET",
		},
	}
}


func BuildEmployeeMeLinks(baseURL string) []Link {
	meURL := baseURL + "/flighthours/api/v1/employees/me"

	return []Link{
		{
			Href:   meURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   meURL,
			Rel:    "update",
			Method: "PUT",
		},
		{
			Href:   meURL,
			Rel:    "delete",
			Method: "DELETE",
		},
	}
}



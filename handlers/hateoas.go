package handlers

import (
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
)

const (
	pathActivate   = "/activate"
	pathDeactivate = "/deactivate"

	resourceDailyLogbooks    = "daily-logbooks"
	resourceTailNumbers      = "tail-numbers"
	resourceAircraftModels   = "aircraft-models"
	resourceAirlineEmployees = "airline-employees"
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

// SetLocationHeader sets the Location header using only the path portion
// to avoid open redirect vulnerabilities from user-controlled Host headers.
func SetLocationHeader(c *gin.Context, baseURL, resource, resourceID string) {
	locationURL := BuildResourceURL(baseURL, resource, resourceID)

	// Validate the URL to prevent injection
	parsed, err := url.Parse(locationURL)
	if err != nil {
		return
	}

	// Use only the path to prevent open redirects via Host header manipulation
	c.Header("Location", parsed.RequestURI())
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

func BuildAccountLinks(baseURL, accountID string) []Link {
	return BuildResourceLinks(baseURL, "accounts", accountID)
}

func BuildMessageLinks(baseURL, messageID string) []Link {
	return BuildResourceLinks(baseURL, "messages", messageID)
}

func BuildMessageCreatedLinks(baseURL, messageID string) []Link {
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

func BuildMessageUpdatedLinks(baseURL, messageID string) []Link {
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

// ============================================================================
// EMPLOYEE HATEOAS LINKS
// ============================================================================

// BuildEmployeeLinks construye links HATEOAS para un empleado específico
func BuildEmployeeLinks(baseURL, employeeID string) []Link {
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
	meURL := baseURL + "/flighthours/api/v1/employees"

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

// ============================================================================
// AIRLINE HATEOAS LINKS
// ============================================================================

// BuildAirlineLinks construye links HATEOAS para una aerolínea específica
func BuildAirlineLinks(baseURL, airlineID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "airlines", airlineID)
	collectionURL := BuildCollectionURL(baseURL, "airlines")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL + pathActivate,
			Rel:    "activate",
			Method: "PATCH",
		},
		{
			Href:   resourceURL + pathDeactivate,
			Rel:    "deactivate",
			Method: "PATCH",
		},
		{
			Href:   collectionURL,
			Rel:    "collection",
			Method: "GET",
		},
	}
}

// BuildAirlineListLinks construye links para la lista de aerolíneas
func BuildAirlineListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "airlines")

	return []Link{
		{
			Href:   collectionURL,
			Rel:    "self",
			Method: "GET",
		},
	}
}

// BuildAirlineStatusLinks construye links para respuesta de cambio de status
func BuildAirlineStatusLinks(baseURL, airlineID string, isActive bool) []Link {
	resourceURL := BuildResourceURL(baseURL, "airlines", airlineID)
	collectionURL := BuildCollectionURL(baseURL, "airlines")

	links := []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
	}

	// Si está activo, mostrar link para desactivar y viceversa
	if isActive {
		links = append(links, Link{
			Href:   resourceURL + pathDeactivate,
			Rel:    "deactivate",
			Method: "PATCH",
		})
	} else {
		links = append(links, Link{
			Href:   resourceURL + pathActivate,
			Rel:    "activate",
			Method: "PATCH",
		})
	}

	links = append(links, Link{
		Href:   collectionURL,
		Rel:    "collection",
		Method: "GET",
	})

	return links
}

// ============================================================================
// AIRPORT HATEOAS LINKS
// ============================================================================

// BuildAirportLinks construye links HATEOAS para un aeropuerto específico
func BuildAirportLinks(baseURL, airportID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "airports", airportID)
	collectionURL := BuildCollectionURL(baseURL, "airports")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL + pathActivate,
			Rel:    "activate",
			Method: "PATCH",
		},
		{
			Href:   resourceURL + pathDeactivate,
			Rel:    "deactivate",
			Method: "PATCH",
		},
		{
			Href:   collectionURL,
			Rel:    "collection",
			Method: "GET",
		},
	}
}

// BuildAirportListLinks construye links para la lista de aeropuertos
func BuildAirportListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "airports")

	return []Link{
		{
			Href:   collectionURL,
			Rel:    "self",
			Method: "GET",
		},
	}
}

// BuildAirportStatusLinks construye links para respuesta de cambio de status
func BuildAirportStatusLinks(baseURL, airportID string, isActive bool) []Link {
	resourceURL := BuildResourceURL(baseURL, "airports", airportID)
	collectionURL := BuildCollectionURL(baseURL, "airports")

	links := []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
	}

	// Si está activo, mostrar link para desactivar y viceversa
	if isActive {
		links = append(links, Link{
			Href:   resourceURL + pathDeactivate,
			Rel:    "deactivate",
			Method: "PATCH",
		})
	} else {
		links = append(links, Link{
			Href:   resourceURL + pathActivate,
			Rel:    "activate",
			Method: "PATCH",
		})
	}

	links = append(links, Link{
		Href:   collectionURL,
		Rel:    "collection",
		Method: "GET",
	})

	return links
}

// ============================================================================
// DAILY LOGBOOK HATEOAS LINKS
// ============================================================================

// BuildDailyLogbookLinks construye links HATEOAS para una bitácora diaria específica
func BuildDailyLogbookLinks(baseURL, logbookID string) []Link {
	resourceURL := BuildResourceURL(baseURL, resourceDailyLogbooks, logbookID)
	collectionURL := BuildCollectionURL(baseURL, resourceDailyLogbooks)

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
			Href:   resourceURL + pathActivate,
			Rel:    "activate",
			Method: "PATCH",
		},
		{
			Href:   resourceURL + pathDeactivate,
			Rel:    "deactivate",
			Method: "PATCH",
		},
		{
			Href:   collectionURL,
			Rel:    "collection",
			Method: "GET",
		},
	}
}

// BuildDailyLogbookListLinks construye links para la lista de bitácoras diarias
func BuildDailyLogbookListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, resourceDailyLogbooks)

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

// BuildDailyLogbookStatusLinks construye links para respuesta de cambio de status
func BuildDailyLogbookStatusLinks(baseURL, logbookID string, isActive bool) []Link {
	resourceURL := BuildResourceURL(baseURL, resourceDailyLogbooks, logbookID)
	collectionURL := BuildCollectionURL(baseURL, resourceDailyLogbooks)

	links := []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
	}

	// Si está activo, mostrar link para desactivar y viceversa
	if isActive {
		links = append(links, Link{
			Href:   resourceURL + pathDeactivate,
			Rel:    "deactivate",
			Method: "PATCH",
		})
	} else {
		links = append(links, Link{
			Href:   resourceURL + pathActivate,
			Rel:    "activate",
			Method: "PATCH",
		})
	}

	links = append(links, Link{
		Href:   collectionURL,
		Rel:    "collection",
		Method: "GET",
	})

	return links
}

// BuildDailyLogbookCreatedLinks construye links para una bitácora recién creada
func BuildDailyLogbookCreatedLinks(baseURL, logbookID string) []Link {
	resourceURL := BuildResourceURL(baseURL, resourceDailyLogbooks, logbookID)
	collectionURL := BuildCollectionURL(baseURL, resourceDailyLogbooks)

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

// BuildDailyLogbookDeletedLinks construye links para respuesta de eliminación
func BuildDailyLogbookDeletedLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, resourceDailyLogbooks)

	return []Link{
		{
			Href:   collectionURL,
			Rel:    "list",
			Method: "GET",
		},
		{
			Href:   collectionURL,
			Rel:    "create",
			Method: "POST",
		},
	}
}

// ============================================================================
// TAIL NUMBER HATEOAS LINKS
// ============================================================================

// BuildTailNumberLinks construye links HATEOAS para una matrícula específica
func BuildTailNumberLinks(baseURL, registrationID string) []Link {
	resourceURL := BuildResourceURL(baseURL, resourceTailNumbers, registrationID)
	collectionURL := BuildCollectionURL(baseURL, resourceTailNumbers)

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
			Href:   collectionURL,
			Rel:    "collection",
			Method: "GET",
		},
	}
}

// BuildTailNumberListLinks construye links para la lista de matrículas
func BuildTailNumberListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "tail-numbers")

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

// BuildTailNumberCreatedLinks construye links para una matrícula recién creada
func BuildTailNumberCreatedLinks(baseURL, registrationID string) []Link {
	resourceURL := BuildResourceURL(baseURL, resourceTailNumbers, registrationID)
	collectionURL := BuildCollectionURL(baseURL, resourceTailNumbers)

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
			Href:   collectionURL,
			Rel:    "list",
			Method: "GET",
		},
	}
}

// ============================================================================
// AIRCRAFT MODEL HATEOAS LINKS
// ============================================================================

// BuildAircraftModelLinks construye links HATEOAS para un modelo de aeronave específico
func BuildAircraftModelLinks(baseURL, modelID string) []Link {
	resourceURL := BuildResourceURL(baseURL, resourceAircraftModels, modelID)
	collectionURL := BuildCollectionURL(baseURL, resourceAircraftModels)

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL + pathActivate,
			Rel:    "activate",
			Method: "PATCH",
		},
		{
			Href:   resourceURL + pathDeactivate,
			Rel:    "deactivate",
			Method: "PATCH",
		},
		{
			Href:   collectionURL,
			Rel:    "collection",
			Method: "GET",
		},
	}
}

// BuildAircraftModelListLinks construye links para la lista de modelos de aeronave
func BuildAircraftModelListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "aircraft-models")

	return []Link{
		{
			Href:   collectionURL,
			Rel:    "self",
			Method: "GET",
		},
	}
}

// BuildAircraftFamilyListLinks construye links para la lista de familias de aeronave
func BuildAircraftFamilyListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "aircraft-families")

	return []Link{
		{
			Href:   collectionURL,
			Rel:    "self",
			Method: "GET",
		},
	}
}

// BuildAircraftModelStatusLinks construye links para respuesta de cambio de status (HU41, HU42)
func BuildAircraftModelStatusLinks(baseURL, modelID string, isActive bool) []Link {
	resourceURL := BuildResourceURL(baseURL, resourceAircraftModels, modelID)
	collectionURL := BuildCollectionURL(baseURL, resourceAircraftModels)

	links := []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
	}

	// Si está activo, mostrar link para desactivar y viceversa
	if isActive {
		links = append(links, Link{
			Href:   resourceURL + pathDeactivate,
			Rel:    "deactivate",
			Method: "PATCH",
		})
	} else {
		links = append(links, Link{
			Href:   resourceURL + pathActivate,
			Rel:    "activate",
			Method: "PATCH",
		})
	}

	links = append(links, Link{
		Href:   collectionURL,
		Rel:    "collection",
		Method: "GET",
	})

	return links
}

// ============================================================================
// DAILY LOGBOOK DETAIL HATEOAS LINKS
// ============================================================================

// BuildDailyLogbookDetailLinks construye links HATEOAS para un detalle de bitácora
// Retorna un array de Link para mantener consistencia con el resto del sistema
func BuildDailyLogbookDetailLinks(c *gin.Context, detailID string) []Link {
	baseURL := GetBaseURL(c)
	resourceURL := BuildResourceURL(baseURL, "daily-logbook-details", detailID)

	return []Link{
		{Href: resourceURL, Rel: "self", Method: "GET"},
		{Href: resourceURL, Rel: "update", Method: "PUT"},
	}
}

// BuildDailyLogbookDetailLinksArray construye links HATEOAS como array
func BuildDailyLogbookDetailLinksArray(baseURL, detailID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "daily-logbook-details", detailID)

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
	}
}

// BuildDailyLogbookDetailListLinks construye links para la lista de detalles de una bitácora
func BuildDailyLogbookDetailListLinks(baseURL, logbookID string) []Link {
	logbookURL := BuildResourceURL(baseURL, "daily-logbooks", logbookID)
	detailsURL := logbookURL + "/details"

	return []Link{
		{
			Href:   detailsURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   detailsURL,
			Rel:    "create",
			Method: "POST",
		},
		{
			Href:   logbookURL,
			Rel:    "logbook",
			Method: "GET",
		},
	}
}

// ============================================================================
// ENGINE HATEOAS LINKS
// ============================================================================

// BuildEngineLinks construye links HATEOAS para un motor específico
func BuildEngineLinks(baseURL, engineID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "engines", engineID)
	collectionURL := BuildCollectionURL(baseURL, "engines")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   collectionURL,
			Rel:    "collection",
			Method: "GET",
		},
	}
}

// BuildEngineListLinks construye links para la lista de motores
func BuildEngineListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "engines")

	return []Link{
		{
			Href:   collectionURL,
			Rel:    "self",
			Method: "GET",
		},
	}
}

// ============================================================================
// MANUFACTURER HATEOAS LINKS
// ============================================================================

// BuildManufacturerLinks construye links HATEOAS para un fabricante específico
func BuildManufacturerLinks(baseURL, manufacturerID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "manufacturers", manufacturerID)
	collectionURL := BuildCollectionURL(baseURL, "manufacturers")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   collectionURL,
			Rel:    "collection",
			Method: "GET",
		},
	}
}

// BuildManufacturerListLinks construye links para la lista de fabricantes
func BuildManufacturerListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "manufacturers")

	return []Link{
		{
			Href:   collectionURL,
			Rel:    "self",
			Method: "GET",
		},
	}
}

// ============================================================================
// AIRLINE EMPLOYEE HATEOAS LINKS
// ============================================================================

// BuildAirlineEmployeeLinks construye links HATEOAS para un empleado de aerolínea específico
func BuildAirlineEmployeeLinks(baseURL, employeeID string) []Link {
	resourceURL := BuildResourceURL(baseURL, resourceAirlineEmployees, employeeID)
	collectionURL := BuildCollectionURL(baseURL, resourceAirlineEmployees)

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
			Href:   resourceURL + pathActivate,
			Rel:    "activate",
			Method: "PATCH",
		},
		{
			Href:   resourceURL + pathDeactivate,
			Rel:    "deactivate",
			Method: "PATCH",
		},
		{
			Href:   collectionURL,
			Rel:    "collection",
			Method: "GET",
		},
	}
}

// BuildAirlineEmployeeStatusLinks construye links para respuesta de cambio de status
func BuildAirlineEmployeeStatusLinks(baseURL, employeeID string, isActive bool) []Link {
	resourceURL := BuildResourceURL(baseURL, resourceAirlineEmployees, employeeID)
	collectionURL := BuildCollectionURL(baseURL, resourceAirlineEmployees)

	links := []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
	}

	// Si está activo, mostrar link para desactivar y viceversa
	if isActive {
		links = append(links, Link{
			Href:   resourceURL + pathDeactivate,
			Rel:    "deactivate",
			Method: "PATCH",
		})
	} else {
		links = append(links, Link{
			Href:   resourceURL + pathActivate,
			Rel:    "activate",
			Method: "PATCH",
		})
	}

	links = append(links, Link{
		Href:   collectionURL,
		Rel:    "collection",
		Method: "GET",
	})

	return links
}

// BuildAirlineEmployeeCreatedLinks construye links para respuesta de creación de empleado
func BuildAirlineEmployeeCreatedLinks(baseURL, employeeID string) []Link {
	resourceURL := BuildResourceURL(baseURL, resourceAirlineEmployees, employeeID)
	collectionURL := BuildCollectionURL(baseURL, resourceAirlineEmployees)

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
			Href:   collectionURL,
			Rel:    "list",
			Method: "GET",
		},
	}
}

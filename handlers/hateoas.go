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

// ============================================================================
// EMPLOYEE HATEOAS LINKS
// ============================================================================

// BuildEmployeeLinks construye links HATEOAS para un empleado específico
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
func BuildAirlineLinks(baseURL string, airlineID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "airlines", airlineID)
	collectionURL := BuildCollectionURL(baseURL, "airlines")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL + "/activate",
			Rel:    "activate",
			Method: "PATCH",
		},
		{
			Href:   resourceURL + "/deactivate",
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
func BuildAirlineStatusLinks(baseURL string, airlineID string, isActive bool) []Link {
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
			Href:   resourceURL + "/deactivate",
			Rel:    "deactivate",
			Method: "PATCH",
		})
	} else {
		links = append(links, Link{
			Href:   resourceURL + "/activate",
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
func BuildAirportLinks(baseURL string, airportID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "airports", airportID)
	collectionURL := BuildCollectionURL(baseURL, "airports")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL + "/activate",
			Rel:    "activate",
			Method: "PATCH",
		},
		{
			Href:   resourceURL + "/deactivate",
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
func BuildAirportStatusLinks(baseURL string, airportID string, isActive bool) []Link {
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
			Href:   resourceURL + "/deactivate",
			Rel:    "deactivate",
			Method: "PATCH",
		})
	} else {
		links = append(links, Link{
			Href:   resourceURL + "/activate",
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
func BuildDailyLogbookLinks(baseURL string, logbookID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "daily-logbooks", logbookID)
	collectionURL := BuildCollectionURL(baseURL, "daily-logbooks")

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

// BuildDailyLogbookListLinks construye links para la lista de bitácoras diarias
func BuildDailyLogbookListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "daily-logbooks")

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

// BuildDailyLogbookCreatedLinks construye links para una bitácora recién creada
func BuildDailyLogbookCreatedLinks(baseURL string, logbookID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "daily-logbooks", logbookID)
	collectionURL := BuildCollectionURL(baseURL, "daily-logbooks")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   collectionURL,
			Rel:    "list",
			Method: "GET",
		},
	}
}

// ============================================================================
// LICENSE PLATE HATEOAS LINKS
// ============================================================================

// BuildLicensePlateLinks construye links HATEOAS para una matrícula específica
func BuildLicensePlateLinks(baseURL string, registrationID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "license-plates", registrationID)
	collectionURL := BuildCollectionURL(baseURL, "license-plates")

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

// BuildLicensePlateListLinks construye links para la lista de matrículas
func BuildLicensePlateListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "license-plates")

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

// BuildLicensePlateCreatedLinks construye links para una matrícula recién creada
func BuildLicensePlateCreatedLinks(baseURL string, registrationID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "license-plates", registrationID)
	collectionURL := BuildCollectionURL(baseURL, "license-plates")

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
func BuildAircraftModelLinks(baseURL string, modelID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "aircraft-models", modelID)
	collectionURL := BuildCollectionURL(baseURL, "aircraft-models")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL + "/activate",
			Rel:    "activate",
			Method: "PATCH",
		},
		{
			Href:   resourceURL + "/deactivate",
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
func BuildAircraftModelStatusLinks(baseURL string, modelID string, isActive bool) []Link {
	resourceURL := BuildResourceURL(baseURL, "aircraft-models", modelID)
	collectionURL := BuildCollectionURL(baseURL, "aircraft-models")

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
			Href:   resourceURL + "/deactivate",
			Rel:    "deactivate",
			Method: "PATCH",
		})
	} else {
		links = append(links, Link{
			Href:   resourceURL + "/activate",
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
// ROUTE HATEOAS LINKS
// ============================================================================

// BuildRouteLinks construye links HATEOAS para una ruta específica
func BuildRouteLinks(baseURL string, routeID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "routes", routeID)
	collectionURL := BuildCollectionURL(baseURL, "routes")

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

// BuildRouteListLinks construye links para la lista de rutas
func BuildRouteListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "routes")

	return []Link{
		{
			Href:   collectionURL,
			Rel:    "self",
			Method: "GET",
		},
	}
}

// ============================================================================
// AIRLINE ROUTE HATEOAS LINKS
// ============================================================================

// BuildAirlineRouteLinks construye links HATEOAS para una ruta de aerolínea específica
func BuildAirlineRouteLinks(baseURL string, airlineRouteID string, isActive bool) []Link {
	resourceURL := BuildResourceURL(baseURL, "airline-routes", airlineRouteID)
	collectionURL := BuildCollectionURL(baseURL, "airline-routes")

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
			Href:   resourceURL + "/deactivate",
			Rel:    "deactivate",
			Method: "PATCH",
		})
	} else {
		links = append(links, Link{
			Href:   resourceURL + "/activate",
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

// BuildAirlineRouteListLinks construye links para la lista de rutas aerolínea
func BuildAirlineRouteListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "airline-routes")

	return []Link{
		{
			Href:   collectionURL,
			Rel:    "self",
			Method: "GET",
		},
	}
}

// BuildAirlineRouteStatusLinks construye links para respuesta de cambio de status
func BuildAirlineRouteStatusLinks(baseURL string, airlineRouteID string, isActive bool) []Link {
	resourceURL := BuildResourceURL(baseURL, "airline-routes", airlineRouteID)
	collectionURL := BuildCollectionURL(baseURL, "airline-routes")

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
			Href:   resourceURL + "/deactivate",
			Rel:    "deactivate",
			Method: "PATCH",
		})
	} else {
		links = append(links, Link{
			Href:   resourceURL + "/activate",
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
// Retorna un mapa para usar directamente en el DTO response
func BuildDailyLogbookDetailLinks(c *gin.Context, detailID string) map[string]string {
	baseURL := GetBaseURL(c)
	resourceURL := BuildResourceURL(baseURL, "daily-logbook-details", detailID)

	return map[string]string{
		"self":   resourceURL,
		"update": resourceURL,
		"delete": resourceURL,
	}
}

// BuildDailyLogbookDetailLinksArray construye links HATEOAS como array
func BuildDailyLogbookDetailLinksArray(baseURL string, detailID string) []Link {
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
func BuildDailyLogbookDetailListLinks(baseURL string, logbookID string) []Link {
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
func BuildEngineLinks(baseURL string, engineID string) []Link {
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
func BuildManufacturerLinks(baseURL string, manufacturerID string) []Link {
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
func BuildAirlineEmployeeLinks(baseURL string, employeeID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "airline-employees", employeeID)
	collectionURL := BuildCollectionURL(baseURL, "airline-employees")

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
			Href:   resourceURL + "/activate",
			Rel:    "activate",
			Method: "PATCH",
		},
		{
			Href:   resourceURL + "/deactivate",
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
func BuildAirlineEmployeeStatusLinks(baseURL string, employeeID string, isActive bool) []Link {
	resourceURL := BuildResourceURL(baseURL, "airline-employees", employeeID)
	collectionURL := BuildCollectionURL(baseURL, "airline-employees")

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
			Href:   resourceURL + "/deactivate",
			Rel:    "deactivate",
			Method: "PATCH",
		})
	} else {
		links = append(links, Link{
			Href:   resourceURL + "/activate",
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
func BuildAirlineEmployeeCreatedLinks(baseURL string, employeeID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "airline-employees", employeeID)
	collectionURL := BuildCollectionURL(baseURL, "airline-employees")

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

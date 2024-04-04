// Package endpoints provides the API endpoints for handling events.
package endpoints

import (
	"database/sql"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	utils "github.com/Sea-Shell/gogear-api/pkg/utils"
	models "github.com/bateau84/pppoe-api/pkg/models"
	"github.com/gin-contrib/requestid"

	gin "github.com/gin-gonic/gin"
	zap "go.uber.org/zap"

	// needed for sqlite3 support
	_ "github.com/mattn/go-sqlite3"
)

// ListEvents will list all events in the database
//
//	@Summary		List events
//	@Description	Get a list of event items
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	false	"name search for event"		default(event name)
//	@Param			page	query		int		false	"Page number"				default(1)
//	@Param			limit	query		int		false	"Number of items per page"	default(30)
//	@Success		200		{object}	models.ResponsePayload{items=[]models.EventItem}
//	@Failure		default	{object}	models.Error
//	@Router			/event/list [get]
func ListEvents(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	currentQueryParameters := c.Request.URL.Query()

	page := c.Query("page")
	limit := c.Query("limit")
	eventname := c.Query("name")
	log := c.MustGet("logger").(*zap.SugaredLogger)
	db := c.MustGet("db").(*sql.DB)
	conf := c.MustGet("conf").(*models.General)

	log.Debugf("Request parameters: %#v", c.Request.URL.Query())

	if limit == "" {
		limit = "30"
	}

	if page == "" || page == "0" {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		log.Errorf("Error setting page to int: %#v", err)
		c.IndentedJSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		log.Errorf("Error setting limit to int: %#v", err)
		c.IndentedJSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	if !(utf8.RuneCountInString(eventname) >= conf.Event.SearchMinimumCharacters || eventname == "") {
		msg := fmt.Sprintf("minimum characters in event name search is %d", conf.Event.SearchMinimumCharacters)
		log.Errorf("Error event name (%#v) is less than %d characters", eventname, conf.Event.SearchMinimumCharacters)
		c.IndentedJSON(http.StatusBadRequest, models.Error{Error: msg})
		return
	}

	if pageInt <= 0 {
		log.Errorf("Error page is less than 0: %#v", err)
		c.IndentedJSON(http.StatusBadRequest, models.Error{Error: "Invalid page number"})
		return
	}

	if limitInt <= 0 {
		log.Errorf("Error limit is less than 0: %#v", err)
		c.IndentedJSON(http.StatusBadRequest, models.Error{Error: "Invalid limit number"})
		return
	}

	conditions := []string{}
	if eventname != "" {
		eName := fmt.Sprintf("eventName LIKE '%%%s%%'", eventname)
		conditions = append(conditions, eName)
	}
	// if category != "" {
	// 	cat := fmt.Sprintf("event.eventCategoryId = %s", category)
	// 	conditions = append(conditions, cat)
	// }
	// if manufacturer != "" {
	// 	cat := fmt.Sprintf("event.eventManufactureId = %s", manufacturer)
	// 	conditions = append(conditions, cat)
	// }

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	baseCountQuery := "SELECT COUNT(*) FROM event"
	countQuery := baseCountQuery + whereClause

	var totalCount int
	err = db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		log.Errorf("Error getting rating database: %#v", err)
		c.IndentedJSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	start := strconv.Itoa((pageInt - 1) * limitInt)
	totalPages := int(math.Ceil(float64(totalCount) / float64(limitInt)))

	startInt, err := strconv.Atoi(start)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	var paramEvent models.EventItem
	fields := utils.GetDBFieldNames(reflect.TypeOf(paramEvent))

	baseQuery := fmt.Sprintf(`SELECT %s FROM event`,
		strings.Join(fields, ", "))

	queryLimit := fmt.Sprintf(" LIMIT %v, %v", startInt, limitInt)

	query := baseQuery + whereClause + queryLimit

	log.Debugf("Query: %s", query)

	rows, err := db.Query(query)
	if err != nil {
		log.Errorf("Query error: %#v", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	dest, err := utils.GetScanFields(paramEvent)
	if err != nil {
		log.Errorf("Error getting destination arguments: %#v", err)
		c.IndentedJSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	var eventList []models.EventItem

	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Errorf("No events found")
				c.IndentedJSON(http.StatusNotFound, models.Error{Error: "No results"})
				return
			}
			log.Errorf("Scan error: %#v", err)
			c.IndentedJSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
			return
		}

		for i := 0; i < reflect.TypeOf(paramEvent).NumField(); i++ {
			reflect.ValueOf(&paramEvent).Elem().Field(i).Set(reflect.ValueOf(dest[i]).Elem())
		}

		eventList = append(eventList, paramEvent)
	}

	payload := models.ResponsePayload{
		TotalItemCount: totalCount,
		CurrentPage:    pageInt,
		ItemLimit:      limitInt,
		TotalPages:     totalPages,
		Items:          eventList,
	}

	if pageInt < totalPages {
		currentQueryParameters.Set("page", strconv.Itoa(pageInt+1))
		nextPage := url.URL{
			Path:     c.Request.URL.Path,
			RawQuery: currentQueryParameters.Encode(),
		}
		payload.NextPage = new(string)
		*payload.NextPage = nextPage.String()
	}

	if pageInt > 1 {
		currentQueryParameters.Set("page", strconv.Itoa(pageInt-1))
		prevPage := url.URL{
			Path:     c.Request.URL.Path,
			RawQuery: currentQueryParameters.Encode(),
		}
		payload.PrevPage = new(string)
		*payload.PrevPage = prevPage.String()
	}

	log.Infof("successfully fetched event with id: %s, eventName: %s", paramEvent.EventID, paramEvent.EventName)
	c.IndentedJSON(http.StatusOK, payload)
}

// GetEvent Gets spessific event with ID
//
// @Summary		Get event with ID
// @Description	Get event spessific to ID
// @Tags			Event
// @Accept			json
// @Produce		json
// @Param			event	path		int				true	"Unique ID of event you want to get"
// @Success		200		{object}	models.EventItem	"desc"
// @Failure		default	{object}	models.Error
// @Router			/event/{event}/get [get]
func GetEvent(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	log := c.MustGet("logger").(*zap.SugaredLogger)
	db := c.MustGet("db").(*sql.DB)
	function := "event"

	urlParameter, err := strconv.Atoi(c.Param(function))
	if err != nil {
		log.Errorf("urlParamter is of wrong type: %#v", err)
		c.IndentedJSON(http.StatusBadRequest, models.Error{Error: err.Error()})
	}

	var extraSQL []string
	// extraSQL = append(extraSQL, " LEFT JOIN manufacture ON event.eventManufactureId = manufacture.manufactureId ")
	// extraSQL = append(extraSQL, " LEFT JOIN event_top_category ON event.eventTopCategoryId = event_top_category.topCategoryId ")
	// extraSQL = append(extraSQL, "  LEFT JOIN event_category ON event.eventCategoryId = event_category.categoryId ")

	results, err := utils.GenericGet[models.EventItem]("event", urlParameter, extraSQL, db)
	if err != nil {
		log.Errorf("Unable to get %s with id: %s. Error: %#v", function, urlParameter, err)
		c.IndentedJSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	log.Infof("Successfully fetched %s with ID %s", function, urlParameter)
	c.IndentedJSON(http.StatusOK, results)
}

// InsertEvent will insert a new event to the db
//
//	@Summary		Insert new event
//	@Description	Insert new event with corresponding values
//	@Security		OAuth2Application[write]
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.EventItemNoID	true	"query params"	test
//	@Success		200		{object}	models.Status	"status: success when all goes well"
//	@Failure		default	{object}	models.Error
//	@Router			/event/insert [put]
func InsertEvent(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	log := c.MustGet("logger").(*zap.SugaredLogger)
	db := c.MustGet("db").(*sql.DB)

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		log.Error(err.Error())
		return
	}

	log.Infow("PayLoad", "Content", data, "requestID", requestid.Get(c))
	c.Set("logger", log)

	err = utils.GenericInsert[models.EventItem]("event", data, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		log.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

//	@Summary		Update event with ID
//	@Description	Update event identified by ID
//	@Security		OAuth2Application[write]
//	@Tags			event
//	@Accept			json
//	@Produce		json
//	@Param			event	path		int				true	"Unique ID of event you want to get"
//	@Param			request	body		models.EventItem		true	"query params"	test
//	@Success		200		{object}	models.Status	"status: success when all goes well"
//	@Failure		default	{object}	models.Error
//	@Router			/event/{event}/update [post]
// func Updateevent(c *gin.Context) {
// 	c.Header("Content-Type", "application/json")

// 	log := c.MustGet("logger").(*zap.SugaredLogger)
// 	db := c.MustGet("db").(*sql.DB)

// 	data, err := io.ReadAll(c.Request.Body)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
// 		log.Error(err.Error())
// 		return
// 	}

// 	err = utils.GenericUpdate[models.EventItem]("event", data, db)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
// 		log.Error(err.Error())
// 		return
// 	}

// 	c.JSON(http.StatusOK, map[string]string{"status": "success"})
// }

//	@Summary		Delete event with ID
//	@Description	Delete event with corresponding ID value
//	@Security		OAuth2Application[write]
//	@Tags			event
//	@Accept			json
//	@Produce		json
//	@Param			event	path		int				true	"Unique ID of event you want to delete"
//	@Success		200		{object}	models.Status	"status: success when all goes well"
//	@Failure		default	{object}	models.Error
//	@Router			/event/{event}/delete [delete]
// func Deleteevent(c *gin.Context) {
// 	c.Header("Content-Type", "application/json")

// 	log := c.MustGet("logger").(*zap.SugaredLogger)
// 	db := c.MustGet("db").(*sql.DB)
// 	function := "event"

// 	urlParameter, err := strconv.Atoi(c.Param(function))
// 	if err != nil {
// 		log.Errorf("urlParamter is of wrong type: %#v", err)
// 		c.IndentedJSON(http.StatusBadRequest, models.Error{Error: err.Error()})
// 		return
// 	}

// 	result, err := utils.GenericDelete[models.EventItem]("event", urlParameter, db)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
// 		log.Error(err.Error())
// 		return
// 	}

// 	log.Infof("success! event with event_id %v and event_name %s was deleted", result.eventId, result.eventName)
// 	c.JSON(http.StatusOK, map[string]string{
// 		"status": fmt.Sprintf("success! event with event_id %v and event_name %s was deleted", result.eventId, result.eventName),
// 	})
// }

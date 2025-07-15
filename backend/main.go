package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/slack-go/slack"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type TokenStore struct {
	mu     sync.RWMutex
	tokens map[string]string // userID -> token
}

var tokenStore = &TokenStore{
	tokens: make(map[string]string),
}

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/health/google", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/api/calendar", getCalendarEvents)
	r.POST("/api/token", saveToken)
	r.POST("/api/slack/notify", sendSlackNotification)
	r.POST("/api/google-meet/create", createGoogleMeetURL)
	r.POST("/api/zoom/create", createZoomURL)

	fmt.Println("Server running on :8080")
	r.Run(":8080")
}

func getCalendarEvents(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		return
	}

	ctx := context.Background()

	oauth2Token := &oauth2.Token{
		AccessToken: token,
		TokenType:   "Bearer",
	}

	config := &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{calendar.CalendarReadonlyScope},
		Endpoint:     google.Endpoint,
	}

	client := config.Client(ctx, oauth2Token)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Unable to retrieve Calendar client: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create calendar client"})
		return
	}

	events, err := srv.Events.List("primary").MaxResults(10).SingleEvents(true).OrderBy("startTime").Do()
	if err != nil {
		log.Printf("Unable to retrieve next ten of the user's events: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve calendar events"})
		return
	}

	var eventList []map[string]interface{}
	for _, item := range events.Items {
		eventData := map[string]interface{}{
			"id":      item.Id,
			"summary": item.Summary,
			"start":   item.Start,
			"end":     item.End,
		}
		eventList = append(eventList, eventData)
	}

	c.JSON(http.StatusOK, gin.H{
		"events": eventList,
		"count":  len(eventList),
	})
}

func saveToken(c *gin.Context) {
	var request struct {
		UserID string `json:"userId" binding:"required"`
		Token  string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenStore.mu.Lock()
	tokenStore.tokens[request.UserID] = request.Token
	tokenStore.mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Token saved successfully"})
}

func sendSlackNotification(c *gin.Context) {
	var request struct {
		UserID  string `json:"userId" binding:"required"`
		Channel string `json:"channel" binding:"required"`
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenStore.mu.RLock()
	token, exists := tokenStore.tokens[request.UserID]
	tokenStore.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token not found for user"})
		return
	}

	api := slack.New(token)
	_, _, err := api.PostMessage(request.Channel, slack.MsgOptionText(request.Message, false))
	if err != nil {
		log.Printf("Error sending Slack message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send Slack message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Slack notification sent successfully"})
}


func createGoogleMeetURL(c *gin.Context) {
	fmt.Println("createGoogleMeetURL")
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		return
	}

	var request struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description,omitempty"`
		StartTime   string `json:"startTime" binding:"required"`
		EndTime     string `json:"endTime" binding:"required"`
		Attendees   []string `json:"attendees,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	oauth2Token := &oauth2.Token{
		AccessToken: token,
		TokenType:   "Bearer",  
	}

	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			calendar.CalendarScope,
		},
		Endpoint: google.Endpoint,
	}

	client := config.Client(ctx, oauth2Token)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Unable to retrieve Calendar client: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create calendar client"})
		return
	}

	startTime, err := time.Parse(time.RFC3339, request.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format. Use RFC3339 format."})
		return
	}

	endTime, err := time.Parse(time.RFC3339, request.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time format. Use RFC3339 format."})
		return
	}

	var attendeesList []*calendar.EventAttendee
	for _, email := range request.Attendees {
		attendeesList = append(attendeesList, &calendar.EventAttendee{
			Email: email,
		})
	}

	event := &calendar.Event{
		Summary:     request.Title,
		Description: request.Description,
		Start: &calendar.EventDateTime{
			DateTime: startTime.Format(time.RFC3339),
			TimeZone: "UTC",
		},
		End: &calendar.EventDateTime{
			DateTime: endTime.Format(time.RFC3339),
			TimeZone: "UTC",
		},
		Attendees: attendeesList,
		ConferenceData: &calendar.ConferenceData{
			CreateRequest: &calendar.CreateConferenceRequest{
				RequestId: fmt.Sprintf("meet-%d", time.Now().Unix()),
				ConferenceSolutionKey: &calendar.ConferenceSolutionKey{
					Type: "hangoutsMeet",
				},
			},
		},
	}

	createdEvent, err := srv.Events.Insert("primary", event).ConferenceDataVersion(1).Do()
	if err != nil {
		log.Printf("Unable to create event: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create calendar event"})
		return
	}

	var meetURL string
	if createdEvent.ConferenceData != nil && len(createdEvent.ConferenceData.EntryPoints) > 0 {
		for _, entryPoint := range createdEvent.ConferenceData.EntryPoints {
			if entryPoint.EntryPointType == "video" {
				meetURL = entryPoint.Uri
				break
			}
		}
	}

	response := gin.H{
		"eventId":     createdEvent.Id,
		"meetURL":     meetURL,
		"title":       createdEvent.Summary,
		"description": createdEvent.Description,
		"startTime":   createdEvent.Start.DateTime,
		"endTime":     createdEvent.End.DateTime,
		"attendees":   request.Attendees,
		"createdAt":   time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}

type ZoomMeetingRequest struct {
	Topic      string `json:"topic"`
	Type       int    `json:"type"`
	StartTime  string `json:"start_time"`
	Duration   int    `json:"duration"`
	Timezone   string `json:"timezone"`
	Settings   struct {
		HostVideo        bool `json:"host_video"`
		ParticipantVideo bool `json:"participant_video"`
		JoinBeforeHost   bool `json:"join_before_host"`
		MuteUponEntry    bool `json:"mute_upon_entry"`
		Watermark        bool `json:"watermark"`
		UsePmi           bool `json:"use_pmi"`
		ApprovalType     int  `json:"approval_type"`
		Audio            string `json:"audio"`
		AutoRecording    string `json:"auto_recording"`
	} `json:"settings"`
}

type ZoomMeetingResponse struct {
	ID           int64  `json:"id"`
	Topic        string `json:"topic"`
	Type         int    `json:"type"`
	Status       string `json:"status"`
	StartTime    string `json:"start_time"`
	Duration     int    `json:"duration"`
	Timezone     string `json:"timezone"`
	JoinURL      string `json:"join_url"`
	StartURL     string `json:"start_url"`
	Password     string `json:"password"`
	HostID       string `json:"host_id"`
	UUID         string `json:"uuid"`
	CreatedAt    string `json:"created_at"`
}

func generateZoomJWT(apiKey, apiSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": apiKey,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
		"aud": "zoom",
		"appKey": apiKey,
		"tokenExp": time.Now().Add(time.Hour * 24).Unix(),
		"alg": "HS256",
	})

	tokenString, err := token.SignedString([]byte(apiSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func createZoomURL(c *gin.Context) {
	var request struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description,omitempty"`
		StartTime   string `json:"startTime" binding:"required"`
		Duration    int    `json:"duration" binding:"required"`
		Timezone    string `json:"timezone,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiKey := os.Getenv("ZOOM_API_KEY")
	apiSecret := os.Getenv("ZOOM_API_SECRET")
	
	if apiKey == "" || apiSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Zoom API credentials not configured"})
		return
	}

	jwtToken, err := generateZoomJWT(apiKey, apiSecret)
	if err != nil {
		log.Printf("Error generating JWT token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication token"})
		return
	}

	startTime, err := time.Parse(time.RFC3339, request.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format. Use RFC3339 format."})
		return
	}

	timezone := request.Timezone
	if timezone == "" {
		timezone = "UTC"
	}

	zoomRequest := ZoomMeetingRequest{
		Topic:     request.Title,
		Type:      2, // Scheduled meeting
		StartTime: startTime.Format("2006-01-02T15:04:05Z"),
		Duration:  request.Duration,
		Timezone:  timezone,
		Settings: struct {
			HostVideo        bool `json:"host_video"`
			ParticipantVideo bool `json:"participant_video"`
			JoinBeforeHost   bool `json:"join_before_host"`
			MuteUponEntry    bool `json:"mute_upon_entry"`
			Watermark        bool `json:"watermark"`
			UsePmi           bool `json:"use_pmi"`
			ApprovalType     int  `json:"approval_type"`
			Audio            string `json:"audio"`
			AutoRecording    string `json:"auto_recording"`
		}{
			HostVideo:        true,
			ParticipantVideo: true,
			JoinBeforeHost:   false,
			MuteUponEntry:    true,
			Watermark:        false,
			UsePmi:           false,
			ApprovalType:     2,
			Audio:            "both",
			AutoRecording:    "none",
		},
	}

	jsonData, err := json.Marshal(zoomRequest)
	if err != nil {
		log.Printf("Error marshaling request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req, err := http.NewRequest("POST", "https://api.zoom.us/v2/users/me/meetings", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to Zoom API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Zoom meeting"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("Zoom API returned status: %d", resp.StatusCode)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Zoom meeting"})
		return
	}

	var zoomResponse ZoomMeetingResponse
	if err := json.NewDecoder(resp.Body).Decode(&zoomResponse); err != nil {
		log.Printf("Error decoding Zoom response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Zoom response"})
		return
	}

	response := gin.H{
		"meetingId":   strconv.FormatInt(zoomResponse.ID, 10),
		"joinURL":     zoomResponse.JoinURL,
		"startURL":    zoomResponse.StartURL,
		"password":    zoomResponse.Password,
		"title":       zoomResponse.Topic,
		"description": request.Description,
		"startTime":   zoomResponse.StartTime,
		"duration":    zoomResponse.Duration,
		"timezone":    zoomResponse.Timezone,
		"createdAt":   time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}
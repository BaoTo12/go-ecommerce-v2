package domain
}
	}
		c.Metrics.ROI = (float64(c.Metrics.Revenue) - float64(c.Metrics.Cost)) / float64(c.Metrics.Cost) * 100
	if c.Metrics.Cost > 0 {
func (c *Campaign) calculateROI() {

}
	c.calculateROI()
	c.Metrics.Revenue += revenue
	c.Metrics.Converted++
func (c *Campaign) RecordConversion(revenue int64) {
// RecordConversion records a conversion

}
	c.Metrics.Clicked++
func (c *Campaign) RecordClick() {
// RecordClick records a click

}
	c.Metrics.Sent++
func (c *Campaign) RecordImpression() {
// RecordImpression records a campaign impression

}
	c.UpdatedAt = time.Now()
	c.Status = CampaignPaused
func (c *Campaign) Pause() {
// Pause pauses the campaign

}
	c.UpdatedAt = time.Now()
	c.Status = CampaignActive
func (c *Campaign) Activate() {
// Activate activates the campaign

}
	}
		UpdatedAt:   now,
		CreatedAt:   now,
		Metrics:     CampaignMetrics{},
		Status:      CampaignDraft,
		Type:        campaignType,
		Description: description,
		Name:        name,
		CampaignID:  uuid.New().String(),
	return &Campaign{
	now := time.Now()
func NewCampaign(name, description string, campaignType CampaignType) *Campaign {
// NewCampaign creates a new campaign

}
	ROI         float64
	Cost        int64
	Revenue     int64
	Converted   int
	Clicked     int
	Opened      int
	Delivered   int
	Sent        int
type CampaignMetrics struct {
// CampaignMetrics tracks campaign performance

}
	Timezone  string
	EndDate   time.Time
	StartDate time.Time
type CampaignSchedule struct {
// CampaignSchedule defines when the campaign runs

}
	Personalize bool
	CTALink     string
	CTAText     string
	ImageURL    string
	Body        string
	Subject     string
type CampaignContent struct {
// CampaignContent represents campaign creative content

}
	MaxOrderValue  int64
	MinOrderValue  int64
	Locations      []string
	MaxAge         int
	MinAge         int
	UserTags       []string
	SegmentIDs     []string
type TargetAudience struct {
// TargetAudience defines campaign targeting

}
	UpdatedAt      time.Time
	CreatedAt      time.Time
	Metrics        CampaignMetrics
	Schedule       CampaignSchedule
	Content        CampaignContent
	TargetAudience TargetAudience
	Status         CampaignStatus
	Type           CampaignType
	Description    string
	Name           string
	CampaignID     string
type Campaign struct {
// Campaign represents a marketing campaign

)
	CampaignSMS      CampaignType = "SMS"
	CampaignPopup    CampaignType = "POPUP"
	CampaignBanner   CampaignType = "BANNER"
	CampaignPush     CampaignType = "PUSH"
	CampaignEmail    CampaignType = "EMAIL"
const (

type CampaignType string
// CampaignType represents the type of campaign

)
	CampaignEnded     CampaignStatus = "ENDED"
	CampaignPaused    CampaignStatus = "PAUSED"
	CampaignActive    CampaignStatus = "ACTIVE"
	CampaignScheduled CampaignStatus = "SCHEDULED"
	CampaignDraft     CampaignStatus = "DRAFT"
const (

type CampaignStatus string
// CampaignStatus represents the status of a marketing campaign

)
	"github.com/google/uuid"

	"time"
import (



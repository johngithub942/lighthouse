package model

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"github.com/lbryio/lighthouse/app/es/index"

	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/null"
	"github.com/lbryio/lbry.go/v2/extras/util"

	"github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v6"
)

// Claim is the document type specified as a struct stored in elasticsearch
type Claim struct {
	ID                  uint64                 `json:"id,omitempty"`
	Name                string                 `json:"name,omitempty"`
	StrippedName        string                 `json:"stripped_name,omitempty"`
	ClaimID             string                 `json:"claimId,omitempty"`
	Channel             *null.String           `json:"channel,omitempty"`
	ChannelClaimID      *null.String           `json:"channel_claim_id,omitempty"`
	BidState            string                 `json:"bid_state,omitempty"`
	EffectiveAmount     uint64                 `json:"effective_amount,omitempty"`
	TransactionTimeUnix null.Uint64            `json:"-"` //Could be null in mempool
	TransactionTime     *null.Time             `json:"transaction_time,omitempty"`
	CertificateAmount   uint64                 `json:"certificate_amount,omitempty"`
	JSONValue           null.String            `json:"-"`
	Value               map[string]interface{} `json:"value,omitempty"`
	Title               *null.String           `json:"title,omitempty"`
	Description         *null.String           `json:"description,omitempty"`
	ReleaseTimeUnix     null.Uint64            `json:"-"`
	ReleaseTime         *null.Time             `json:"release_time,omitempty"`
	ContentType         *null.String           `json:"content_type,omitempty"`
	CertValid           bool                   `json:"cert_valid,omitempty"`
	ClaimType           *null.String           `json:"claim_type,omitempty"`
	FrameWidth          *null.Uint64           `json:"frame_width,omitempty"`
	FrameHeight         *null.Uint64           `json:"frame_height,omitempty"`
	Duration            *null.Uint64           `json:"duration,omitempty"`
	NSFW                bool                   `json:"nsfw,omitempty"`
	ViewCnt             *null.Uint64           `json:"view_cnt,omitempty"`
	SubCnt              *null.Uint64           `json:"sub_cnt,omitempty"`
	ThumbnailURL        *null.String           `json:"thumbnail_url,omitempty"`
	Fee                 *null.Float64          `json:"fee,omitempty"`
	TagsStr             *null.String           `json:"-"`
	Tags                []string               `json:"tags,omitempty"`
	ClaimCount          uint64                 `json:"claim_cnt,omitempty"`
	EffectiveSum        uint64                 `json:"effective_sum,omitempty"`
	ChannelEffectiveSum uint64                 `json:"channel_effective_sum,omitempty"`
	IOSFiltered         bool                   `json:"ios_filtered,omitempty"`
	Language            *null.String           `json:"language,omitempty"`
}

// NewClaim creates an instance of Claim with default values for pointers.
func NewClaim() Claim {
	return Claim{
		Channel:         util.PtrToNullString(""),
		ChannelClaimID:  util.PtrToNullString(""),
		TransactionTime: util.PtrToNullTime(time.Time{}),
		Title:           util.PtrToNullString(""),
		Description:     util.PtrToNullString(""),
		ReleaseTime:     util.PtrToNullTime(time.Time{}),
		ContentType:     util.PtrToNullString(""),
		ClaimType:       util.PtrToNullString(""),
		FrameWidth:      util.PtrToNullUint64(0),
		FrameHeight:     util.PtrToNullUint64(0),
		Duration:        util.PtrToNullUint64(0),
		ViewCnt:         util.PtrToNullUint64(0),
		SubCnt:          util.PtrToNullUint64(0),
		ThumbnailURL:    util.PtrToNullString(""),
		Fee:             util.PtrToNullFloat64(0),
		TagsStr:         util.PtrToNullString(""),
		Language:        util.PtrToNullString(""),
	}
}

// GetClaimsFromDBRows returns the claims from Chainquery DB.
func GetClaimsFromDBRows(rows *sql.Rows) ([]Claim, int, error) {
	claims := make([]Claim, 0)
	var lastID int
	for rows.Next() {
		claim := NewClaim()
		err := claim.PopulateFromDB(rows)
		if err != nil {
			return nil, 0, errors.Prefix("Failed to populate from db result:  ", err)
		}
		value := map[string]interface{}{}
		if !claim.JSONValue.IsNull() {
			err = json.Unmarshal([]byte(claim.JSONValue.String), &value)
			if err != nil {
				return nil, 0, errors.Prefix("could not parse json for value: ", err)
			}
		}
		claim.Name = getNormalizedName(claim.Name)
		claim.Value = value
		claim.StrippedName = getStrippedName(claim.Name)
		lastID = int(claim.ID)
		claims = append(claims, claim)
	}
	return claims, lastID, nil
}

// PopulateFromDB populates the data from the rows into claim objects
func (c *Claim) PopulateFromDB(rows *sql.Rows) error {
	if rows == nil {
		return errors.Err("DB rows do not exist")
	}
	err := rows.Scan(
		&c.ID,
		&c.Name,
		c.Channel,
		&c.ClaimCount,
		c.ChannelClaimID,
		&c.BidState,
		&c.EffectiveAmount,
		&c.EffectiveSum,
		&c.TransactionTimeUnix,
		&c.CertificateAmount,
		&c.ChannelEffectiveSum,
		&c.ClaimID,
		&c.JSONValue,
		c.Title,
		c.Description,
		&c.ReleaseTimeUnix,
		c.ContentType,
		&c.CertValid,
		c.ClaimType,
		c.FrameWidth,
		c.FrameHeight,
		c.Duration,
		&c.NSFW,
		c.ThumbnailURL,
		c.Fee,
		c.TagsStr,
		c.Language)
	if err != nil {
		err = errors.Prefix("Scan Err:", err)
	}
	return err
}

// ProcessFilters sets the filter booleans that apply for the data
func (c *Claim) ProcessFilters() {
	c.IOSFiltered = c.isIOSFilter()

}

var iosFilter = []string{"covid", "corona", "sars"}

func (c Claim) isIOSFilter() bool {
	for _, f := range iosFilter {
		if contains(c.Name, f) ||
			contains(c.Name, f) ||
			contains(c.Channel.String, f) ||
			contains(c.Title.String, f) ||
			contains(c.Description.String, f) ||
			c.tagsContain(f) {
			return true
		}
	}
	return false
}

func contains(field, value string) bool {
	return strings.Contains(strings.ToLower(field), strings.ToLower(value))
}

func (c Claim) tagsContain(s string) bool {
	for _, t := range c.Tags {
		if contains(t, s) {
			return true
		}
	}
	return false
}

// Add Inserts the claim as a document via the bulk processor into elasticsearch
func (c Claim) Add(p *elastic.BulkProcessor) {
	r := elastic.NewBulkIndexRequest().Index(index.Claims).Type(index.ClaimType).Id(c.ClaimID).Doc(c)
	p.Add(r)
}

// Delete removes the claim via the bulk processor from elasticsearch
func (c Claim) Delete(p *elastic.BulkProcessor) {
	r := elastic.NewBulkDeleteRequest().Index(index.Claims).Type(index.ClaimType).Id(c.ClaimID)
	p.Add(r)
}

// Update updates just the fields modified or with default values via the bulk processor in elasticsearch
func (c Claim) Update(p *elastic.BulkProcessor) {
	r := elastic.NewBulkUpdateRequest().Index(index.Claims).Type(index.ClaimType).Id(c.ClaimID).Doc(c)
	p.Add(r)
}

// AsJSON converts the object into a json string
func (c Claim) AsJSON() string {
	data, err := json.Marshal(&c)
	if err != nil {
		logrus.Error(errors.Err(err))
		return ""
	}
	return string(data)

}

var normalizedReplacement = map[string]string{
	".": "-",
	"&": "-",
}

func getNormalizedName(name string) string {
	var replacements []string
	for k, v := range strippedReplacement {
		replacements = append(replacements, k, v)
	}
	replacer := strings.NewReplacer(replacements...)
	return replacer.Replace(name)
}

var strippedReplacement = map[string]string{
	"-":   "",
	"_":   "",
	"The": "",
	"the": "",
	"&":   "",
}

func getStrippedName(name string) string {
	var replacements []string
	for k, v := range strippedReplacement {
		replacements = append(replacements, k, v)
	}
	replacer := strings.NewReplacer(replacements...)
	return replacer.Replace(name)
}
